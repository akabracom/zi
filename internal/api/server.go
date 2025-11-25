package api

import (
    "context"
    "log"
    "net/http"
    "os"

    "backend/internal/handler"
    "backend/internal/middleware"

    "github.com/gin-gonic/gin"
    minioClient "github.com/minio/minio-go/v7"
    
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

var router *gin.Engine
var h *handler.Handler

// @title Deposits API
// @version 1.0
// @description API для управления заявками и услугами
// @host localhost:3001
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func Setup(handler *handler.Handler) {
    h = handler
    router = gin.Default()

    // Swagger UI
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Прокси для MinIO изображений
    router.GET("/images/:name", func(c *gin.Context) {
        name := c.Param("name")
        obj, err := h.Minio.Client.GetObject(context.Background(), h.Minio.BucketName, "images/"+name, minioClient.GetObjectOptions{})
        if err != nil {
            c.String(http.StatusNotFound, "Image not found")
            return
        }
        defer obj.Close()
        objInfo, _ := obj.Stat()
        c.Header("Content-Type", objInfo.ContentType)
        c.DataFromReader(http.StatusOK, objInfo.Size, objInfo.ContentType, obj, nil)
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
