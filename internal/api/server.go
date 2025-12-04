package api

import (
    "context"
    "log"
    "net/http"
    "os"

    "di_web/internal/app/handler"
    "di_web/internal/app/repository"
    "di_web/internal/db"
    minioStorage "di_web/internal/storage/minio"

    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
)

func Run() {
    // --- PostgreSQL ---
    gormDB, err := db.New(db.DSNFromEnv())
    if err != nil {
        log.Fatalf("db connect: %v", err)
    }

    repo := repository.New(gormDB)

    // --- MinIO ---
    minioClient := minioStorage.NewMinioClient()
    minioClient.EnsureBucket()

    // --- Handler ---
    h := handler.New(repo, minioClient.Client)

    // --- Gin ---
    r := gin.Default()
    r.LoadHTMLGlob("templates/*")
    r.Static("/static", "./static")

    // Прокси для изображений из MinIO
    r.GET("/images/:name", func(c *gin.Context) {
        name := c.Param("name")
        objectName := "images/" + name

        obj, err := minioClient.Client.GetObject(
            context.Background(),
            minioClient.BucketName,
            objectName,
            minio.GetObjectOptions{},
        )
        if err != nil {
            c.String(http.StatusNotFound, "Image not found")
            return
        }
        defer obj.Close()

        objInfo, err := obj.Stat()
        if err != nil {
            c.String(http.StatusNotFound, "Image not found")
            return
        }

        c.Header("Content-Type", objInfo.ContentType)
        c.Header("Cache-Control", "public, max-age=31536000")
        c.DataFromReader(http.StatusOK, objInfo.Size, objInfo.ContentType, obj, nil)
    })

    // --- HTML Routes ---
    r.GET("/", h.GetMonths)
    r.GET("/month/:id", h.GetMonthDetail)
    r.POST("/add-to-cart", h.AddToCart)
    r.POST("/remove-from-cart/:id", h.RemoveFromCart)
    r.GET("/deposit-request/:id", h.GetDepositRequest)
    // новая точка: сохранить сумму и даты по заявке
    r.POST("/deposit-request/:id/save", h.SaveRequest)

    // --- API ---
    api := r.Group("/api")
    {
        api.GET("/deposit-offers", h.ListOffers)
        api.GET("/deposit-app/draft", h.GetDraft)
        api.GET("/deposit-app/:id", h.GetAppByID)
        api.POST("/deposit-app/add-offer", h.AddOffer)
        api.POST("/deposit-app/:id/delete", h.DeleteApp)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "3001"
    }

    log.Printf("Server running on :%s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("server run: %v", err)
    }
}
