package api

import (
    "context"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"

    "backend/internal/handler"
    "backend/internal/middleware"

    "github.com/gin-gonic/gin"
    minioClient "github.com/minio/minio-go/v7"
    
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

var router *gin.Engine
var h *handler.Handler

func Setup(handler *handler.Handler) {
    h = handler
    router = gin.Default()

    // Swagger UI
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Прокси для MinIO изображений
    router.GET("/images/:name", func(c *gin.Context) {
        name := c.Param("name")
        // Поддерживаем оба формата: с префиксом img/ и без
        objectPath := name
        if !strings.HasPrefix(name, "img/") {
            objectPath = "img/" + name
        }
        log.Printf("Requesting image: bucket=%s, path=%s", h.Minio.BucketName, objectPath)
        
        // Сначала получаем информацию об объекте
        objInfo, err := h.Minio.Client.StatObject(context.Background(), h.Minio.BucketName, objectPath, minioClient.StatObjectOptions{})
        if err != nil {
            log.Printf("Image not found in MinIO: %v, path=%s", err, objectPath)
            // CORS заголовки для ошибки
            c.Header("Access-Control-Allow-Origin", "*")
            c.String(http.StatusNotFound, "Image not found: "+name)
            return
        }
        
        log.Printf("Object info: size=%d, content-type=%s, etag=%s", objInfo.Size, objInfo.ContentType, objInfo.ETag)
        
        // Теперь получаем сам объект
        obj, err := h.Minio.Client.GetObject(context.Background(), h.Minio.BucketName, objectPath, minioClient.GetObjectOptions{})
        if err != nil {
            log.Printf("Failed to get image object: %v", err)
            // CORS заголовки для ошибки
            c.Header("Access-Control-Allow-Origin", "*")
            c.String(http.StatusNotFound, "Image not found: "+name)
            return
        }
        defer obj.Close()
        
        // Устанавливаем правильный Content-Type
        contentType := objInfo.ContentType
        if contentType == "" || contentType == "application/octet-stream" {
            // Определяем по расширению
            if len(name) > 4 && name[len(name)-4:] == ".png" {
                contentType = "image/png"
            } else if len(name) > 4 && name[len(name)-4:] == ".jpg" {
                contentType = "image/jpeg"
            } else {
                contentType = "image/png"
            }
        }
        
        log.Printf("Serving image: %s, size=%d, content-type=%s", name, objInfo.Size, contentType)
        
        // Читаем все данные из объекта
        data, err := io.ReadAll(obj)
        if err != nil {
            log.Printf("Failed to read image: %v", err)
            // CORS заголовки для ошибки
            c.Header("Access-Control-Allow-Origin", "*")
            c.String(http.StatusInternalServerError, "Failed to read image")
            return
        }
        
        if int64(len(data)) != objInfo.Size {
            log.Printf("Warning: read %d bytes, expected %d", len(data), objInfo.Size)
        }
        
        log.Printf("Successfully read %d bytes for image %s", len(data), name)
        
        // Проверяем, что это валидный PNG (должен начинаться с PNG signature)
        if len(data) >= 8 {
            pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
            actualSignature := data[:8]
            if string(actualSignature) != string(pngSignature) {
                log.Printf("ERROR: Image %s does not have valid PNG signature. First 8 bytes (hex): %x, expected: %x", name, actualSignature, pngSignature)
                previewLen := 100
                if len(data) < previewLen {
                    previewLen = len(data)
                }
                log.Printf("First %d bytes (hex): %x", previewLen, data[:previewLen])
                log.Printf("First %d bytes (ASCII): %s", previewLen, string(data[:previewLen]))
                // CORS заголовки для ошибки
                c.Header("Access-Control-Allow-Origin", "*")
                c.String(http.StatusInternalServerError, fmt.Sprintf("Invalid image format: %s is not a valid PNG file", name))
                return
            } else {
                log.Printf("Image %s has valid PNG signature", name)
            }
        }
        
        // Устанавливаем все заголовки ПЕРЕД отправкой данных
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type")
        c.Header("Content-Type", contentType)
        c.Header("Cache-Control", "public, max-age=31536000")
        c.Header("Content-Length", fmt.Sprintf("%d", len(data)))
        c.Header("Accept-Ranges", "bytes")
        c.Header("X-Content-Type-Options", "nosniff")
        
        // Используем стандартный метод Gin для отправки данных
        c.Data(http.StatusOK, contentType, data)
    })
    
    // OPTIONS для CORS preflight
    router.OPTIONS("/images/:name", func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type")
        c.Status(http.StatusNoContent)
    })

    // API
    api := router.Group("/api")
    {
        // --- АУТЕНТИФИКАЦИЯ (публичные роуты) ---
        auth := api.Group("/auth")
        {
            auth.POST("/register", h.Register)
            auth.POST("/login", h.Login)
        }

        // --- УСЛУГИ ---
        services := api.Group("/services")
        {
            // GET - доступен всем (включая гостей)
            services.GET("", middleware.AuthMiddleware(true), h.ListServices)
            
            // POST, PUT, DELETE - только для модераторов
            services.POST("", middleware.AuthMiddleware(false), middleware.RequireRole("moderator"), h.CreateService)
            services.PUT("/:id", middleware.AuthMiddleware(false), middleware.RequireRole("moderator"), h.UpdateService)
            services.DELETE("/:id", middleware.AuthMiddleware(false), middleware.RequireRole("moderator"), h.DeleteService)
            services.POST("/:id/image", middleware.AuthMiddleware(false), middleware.RequireRole("moderator"), h.UploadServiceImage)
        }

        // --- ЗАЯВКИ (требуется авторизация) ---
        requests := api.Group("/requests")
        requests.Use(middleware.AuthMiddleware(false)) // Все роуты требуют авторизации
        {
            // Доступны user и moderator
            requests.GET("", h.ListRequests)
            requests.GET("/draft", middleware.RequireRole("user", "moderator"), h.GetDraftRequest)
            requests.GET("/:id", h.GetRequestByID)
            
            // Только для user (управление своими заявками)
            requests.POST("/add-service", middleware.RequireRole("user", "moderator"), h.AddServiceToDraft)
            requests.DELETE("/remove-service", middleware.RequireRole("user", "moderator"), h.RemoveFromRequest)
            requests.PUT("/:id/services/:service_id/quantity", middleware.RequireRole("user", "moderator"), h.UpdateQuantity)
            requests.PUT("/:id/submit", middleware.RequireRole("user", "moderator"), h.SubmitRequest)
            requests.POST("/:id/delete", middleware.RequireRole("user", "moderator"), h.DeleteRequest)
            
            // Только для moderator (завершение заявок)
            requests.PUT("/:id/complete", middleware.RequireRole("moderator"), h.CompleteRequest)
        }
    }
}

func Run() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "3001"
    }
    log.Printf("API running on :%s", port)
    log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("server run: %v", err)
    }
}
