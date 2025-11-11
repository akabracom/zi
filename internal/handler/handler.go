// internal/handler/handler.go
package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"backend/internal/users"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/storage/minio"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo  *repository.Repository
	Minio *minio.MinioClient // Публичное поле — с большой буквы
}

func New(repo *repository.Repository, minioClient *minio.MinioClient) *Handler {
	return &Handler{
		repo:  repo,
		Minio: minioClient, // Правильное имя поля
	}
}

// GET /api/services?query=
func (h *Handler) ListServices(c *gin.Context) {
	query := c.Query("query")
	services, err := h.repo.GetServices(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

// POST /api/services
func (h *Handler) CreateService(c *gin.Context) {
	var input struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service := models.Service{Name: input.Name, Price: input.Price, IsActive: true}
	if err := h.repo.CreateService(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, service)
}

// PUT /api/services/:id
func (h *Handler) UpdateService(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var input struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.repo.UpdateService(id, input.Name, input.Price); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// DELETE /api/services/:id
func (h *Handler) DeleteService(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	service, _ := h.repo.GetServiceByID(id)
	if service != nil && service.ImageURL != "" {
		_ = h.Minio.RemoveImage(service.ImageURL) // Правильно: h.Minio
	}
	if err := h.repo.DeleteService(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// POST /api/services/:id/image
func (h *Handler) UploadServiceImage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image required"})
		return
	}
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("images/service_%d%s", id, ext)

	tmpPath := fmt.Sprintf("/tmp/service_%d%s", id, ext)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
		return
	}
	defer os.Remove(tmpPath)

	if err := h.Minio.UploadImage(objectName, tmpPath); err != nil { // Правильно: h.Minio
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	h.repo.UpdateServiceImage(id, objectName)
	c.JSON(http.StatusOK, gin.H{"image_url": "/images/" + objectName})
}

// GET /api/requests
func (h *Handler) ListRequests(c *gin.Context) {
	email := users.GetCurrentUser()
	requests, err := h.repo.ListRequests(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// GET /api/requests/draft
func (h *Handler) GetDraftRequest(c *gin.Context) {
	email := users.GetCurrentUser()
	req, err := h.repo.LoadDraftWithItems(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
		return
	}
	c.JSON(http.StatusOK, req)
}

// GET /api/requests/:id
func (h *Handler) GetRequestByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	req, err := h.repo.GetRequestByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, req)
}

// POST /api/requests/add-service
func (h *Handler) AddServiceToDraft(c *gin.Context) {
	var body struct {
		ServiceID int `json:"service_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	email := users.GetCurrentUser()
	if err := h.repo.AddServiceToDraft(email, body.ServiceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "added"})
}

// DELETE /api/requests/remove-service
func (h *Handler) RemoveFromRequest(c *gin.Context) {
	var body struct {
		RequestID int `json:"request_id" binding:"required"`
		ServiceID int `json:"service_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if err := h.repo.RemoveServiceFromRequest(body.RequestID, body.ServiceID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "removed"})
}

// PUT /api/requests/:id/submit
func (h *Handler) SubmitRequest(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.SubmitRequest(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}

// PUT /api/requests/:id/complete
func (h *Handler) CompleteRequest(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.CompleteRequest(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

// POST /api/auth/register
func (h *Handler) Register(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"email": input.Email})
}

// POST /api/auth/login
func (h *Handler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": "fake-jwt-token"})
}

// POST /api/requests/:id/delete — логическое удаление
func (h *Handler) DeleteRequest(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.DeleteRequestLogical(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}