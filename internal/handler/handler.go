package handler

import (
    "fmt"
    "log"
    "net/http"
    "path/filepath"
    "strconv"
    "time"

    "backend/internal/auth"
    "backend/internal/cache"
    "backend/internal/models"
    "backend/internal/repository"
    "backend/internal/storage/minio"

    "github.com/gin-gonic/gin"
)

type Handler struct {
    repo  *repository.Repository
    Minio *minio.MinioClient
}

func New(repo *repository.Repository, minioClient *minio.MinioClient) *Handler {
    return &Handler{
        repo:  repo,
        Minio: minioClient,
    }
}

// ===================== AUTH =====================

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
    var input models.RegisterRequest

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.repo.CreateUser(input.Email, input.Password, input.Name)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Генерируем JWT токен
    token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
        return
    }

    // Сохраняем токен и сессию в Redis
    if err := cache.SaveToken(user.Email, token); err != nil {
        log.Printf("Failed to save token to Redis: %v", err)
    }
    if err := cache.SaveUserSession(user.ID, user.Email, user.Role); err != nil {
        log.Printf("Failed to save session to Redis: %v", err)
    }

    c.JSON(http.StatusCreated, models.AuthResponse{
        Token: token,
        User:  *user,
    })
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
    var input models.LoginRequest

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.repo.GetUserByEmail(input.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    if !h.repo.ValidatePassword(user, input.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // Генерируем JWT токен
    token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
        return
    }

    // Сохраняем токен и сессию в Redis
    if err := cache.SaveToken(user.Email, token); err != nil {
        log.Printf("Failed to save token to Redis: %v", err)
    }
    if err := cache.SaveUserSession(user.ID, user.Email, user.Role); err != nil {
        log.Printf("Failed to save session to Redis: %v", err)
    }

    c.JSON(http.StatusOK, models.AuthResponse{
        Token: token,
        User:  *user,
    })
}

// ===================== SERVICES =====================

// ListServices godoc
// @Summary Get list of services
// @Description Get all active services with optional search
// @Tags services
// @Produce json
// @Param query query string false "Search query"
// @Success 200 {array} models.Service
// @Router /services [get]
// @Security BearerAuth
func (h *Handler) ListServices(c *gin.Context) {
    query := c.Query("query")
    services, err := h.repo.GetServices(query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, services)
}

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service (moderator only)
// @Tags services
// @Accept json
// @Produce json
// @Param input body object{name=string,price=number} true "Service data"
// @Success 201 {object} models.Service
// @Failure 403 {object} map[string]string
// @Router /services [post]
// @Security BearerAuth
func (h *Handler) CreateService(c *gin.Context) {
    var input struct {
        Name  string  `json:"name" binding:"required"`
        Price float64 `json:"price" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    service := models.Service{
        Name:      input.Name,
        Price:     input.Price,
        IsActive:  true,
        CreatedAt: time.Now(),
    }
    if err := h.repo.CreateService(&service); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, service)
}

// UpdateService godoc
// @Summary Update service
// @Description Update service details (moderator only)
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param input body object{name=string,price=number} true "Service data"
// @Success 200 {object} map[string]string
// @Router /services/{id} [put]
// @Security BearerAuth
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

// DeleteService godoc
// @Summary Delete service
// @Description Delete service (moderator only)
// @Tags services
// @Param id path int true "Service ID"
// @Success 200 {object} map[string]string
// @Router /services/{id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteService(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    service, _ := h.repo.GetServiceByID(id)
    if service != nil && service.ImageURL != "" {
        _ = h.Minio.RemoveImage(service.ImageURL)
    }
    if err := h.repo.DeleteService(id); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// UploadServiceImage godoc
// @Summary Upload service image
// @Description Upload image for service (moderator only)
// @Tags services
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Service ID"
// @Param image formData file true "Image file"
// @Success 200 {object} map[string]string
// @Router /services/{id}/image [post]
// @Security BearerAuth
func (h *Handler) UploadServiceImage(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "image required", "details": err.Error()})
        return
    }

    ext := filepath.Ext(file.Filename)
    if ext == "" {
        ext = ".png"
    }
    objectName := fmt.Sprintf("images/service_%d%s", id, ext)

    fileHandle, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
        return
    }
    defer fileHandle.Close()

    if err := h.Minio.UploadImageDirect(objectName, fileHandle, file.Size); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed: " + err.Error()})
        return
    }

    if err := h.repo.UpdateServiceImage(id, objectName); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"image_url": "/images/" + objectName})
}

// ===================== REQUESTS =====================

// ListRequests godoc
// @Summary Get list of requests
// @Description Get requests (user sees only own, moderator sees all)
// @Tags requests
// @Produce json
// @Param date_from query string false "Filter by date from"
// @Param date_to query string false "Filter by date to"
// @Success 200 {array} models.Request
// @Router /requests [get]
// @Security BearerAuth
func (h *Handler) ListRequests(c *gin.Context) {
    userID, _ := c.Get("user_id")
    role, _ := c.Get("role")

    dateFrom := c.Query("date_from")
    dateTo := c.Query("date_to")

    isModerator := role.(string) == "moderator"
    requests, err := h.repo.ListRequests(userID.(int64), isModerator, dateFrom, dateTo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, requests)
}

// GetDraftRequest godoc
// @Summary Get draft request
// @Description Get or create draft request for current user
// @Tags requests
// @Produce json
// @Success 200 {object} models.Request
// @Router /requests/draft [get]
// @Security BearerAuth
func (h *Handler) GetDraftRequest(c *gin.Context) {
    userID, _ := c.Get("user_id")
    req, err := h.repo.LoadDraftWithItems(userID.(int64))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "draft not found"})
        return
    }
    c.JSON(http.StatusOK, req)
}

// GetRequestByID godoc
// @Summary Get request by ID
// @Description Get specific request (user can only see own, moderator sees all)
// @Tags requests
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} models.Request
// @Router /requests/{id} [get]
// @Security BearerAuth
func (h *Handler) GetRequestByID(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    userID, _ := c.Get("user_id")
    role, _ := c.Get("role")

    var req *models.Request
    var err error

    if role.(string) == "moderator" {
        req, err = h.repo.GetRequestByID(id)
    } else {
        req, err = h.repo.GetRequestByIDAndUser(id, userID.(int64))
    }

    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, req)
}

// AddServiceToDraft godoc
// @Summary Add service to draft
// @Description Add service to user's draft request
// @Tags requests
// @Accept json
// @Produce json
// @Param input body object{service_id=int} true "Service ID"
// @Success 200 {object} map[string]string
// @Router /requests/add-service [post]
// @Security BearerAuth
func (h *Handler) AddServiceToDraft(c *gin.Context) {
    var body struct {
        ServiceID int `json:"service_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }
    userID, _ := c.Get("user_id")
    if err := h.repo.AddServiceToDraft(userID.(int64), body.ServiceID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "added"})
}

// RemoveFromRequest godoc
// @Summary Remove service from request
// @Description Remove service from request
// @Tags requests
// @Accept json
// @Produce json
// @Param input body object{request_id=int,service_id=int} true "Request and Service IDs"
// @Success 200 {object} map[string]string
// @Router /requests/remove-service [delete]
// @Security BearerAuth
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

// UpdateQuantity godoc
// @Summary Update service quantity in request
// @Description Update quantity of service in request
// @Tags requests
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param service_id path int true "Service ID"
// @Param input body object{quantity=int} true "New quantity"
// @Success 200 {object} map[string]string
// @Router /requests/{id}/services/{service_id}/quantity [put]
// @Security BearerAuth
func (h *Handler) UpdateQuantity(c *gin.Context) {
    requestID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    serviceID, _ := strconv.Atoi(c.Param("service_id"))

    var input struct {
        Quantity int `json:"quantity" binding:"required,min=1"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.repo.UpdateQuantity(int(requestID), serviceID, input.Quantity); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "quantity updated"})
}

// SubmitRequest godoc
// @Summary Submit request
// @Description Submit draft request (change status to formed)
// @Tags requests
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]string
// @Router /requests/{id}/submit [put]
// @Security BearerAuth
func (h *Handler) SubmitRequest(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    if err := h.repo.SubmitRequest(id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}

// CompleteRequest godoc
// @Summary Complete request
// @Description Complete request (moderator only)
// @Tags requests
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]string
// @Router /requests/{id}/complete [put]
// @Security BearerAuth
func (h *Handler) CompleteRequest(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    if err := h.repo.CompleteRequest(id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

// DeleteRequest godoc
// @Summary Delete request
// @Description Logically delete request
// @Tags requests
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]string
// @Router /requests/{id}/delete [post]
// @Security BearerAuth
func (h *Handler) DeleteRequest(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    if err := h.repo.DeleteRequestLogical(id); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
