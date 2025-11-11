// internal/api/server.go
package api

import (
	"context"
	"log"
	"net/http"
	"os"

	"backend/internal/handler"

	"github.com/gin-gonic/gin"
	minioClient "github.com/minio/minio-go/v7"
)

var router *gin.Engine
var h *handler.Handler

func Setup(handler *handler.Handler) {
	h = handler
	router = gin.Default()

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
		// --- УСЛУГИ ---
		api.GET("/services", h.ListServices)
		api.POST("/services", h.CreateService)
		api.PUT("/services/:id", h.UpdateService)
		api.DELETE("/services/:id", h.DeleteService)
		api.POST("/services/:id/image", h.UploadServiceImage)

		// --- ЗАЯВКИ ---
		api.GET("/requests", h.ListRequests)
		api.GET("/requests/draft", h.GetDraftRequest)
		api.GET("/requests/:id", h.GetRequestByID)
		api.POST("/requests/add-service", h.AddServiceToDraft)
		api.DELETE("/requests/remove-service", h.RemoveFromRequest)
		api.PUT("/requests/:id/submit", h.SubmitRequest)
		api.PUT("/requests/:id/complete", h.CompleteRequest)
		api.POST("/requests/:id/delete", h.DeleteRequest) // логическое удаление

		// --- АУТЕНТИФИКАЦИЯ ---
		api.POST("/auth/register", h.Register)
		api.POST("/auth/login", h.Login)
	}
}

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("API running on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server run: %v", err)
	}
}