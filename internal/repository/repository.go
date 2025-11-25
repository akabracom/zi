package repository

import (
    "errors"
    "strings"
    "time"

    "backend/internal/models"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type Repository struct {
    db *gorm.DB
}

func New(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

// ===================== USER & AUTH =====================

// CreateUser создает нового пользователя с хешированным паролем
func (r *Repository) CreateUser(email, password, name string) (*models.User, error) {
    // Хешируем пароль
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := models.User{
        Email:     email,
        Name:      name,
        Password:  string(hashedPassword),
        Role:      "user", // По умолчанию роль "user"
        CreatedAt: time.Now(),
    }

    if err := r.db.Create(&user).Error; err != nil {
        return nil, err
    }

    return &user, nil
}

// GetUserByEmail находит пользователя по email
func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err == gorm.ErrRecordNotFound {
        return nil, errors.New("user not found")
    }
    return &user, err
}

// GetUserByID находит пользователя по ID
func (r *Repository) GetUserByID(id int64) (*models.User, error) {
    var user models.User
    err := r.db.First(&user, id).Error
    if err == gorm.ErrRecordNotFound {
        return nil, errors.New("user not found")
    }
    return &user, err
}

// ValidatePassword проверяет пароль пользователя
func (r *Repository) ValidatePassword(user *models.User, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    return err == nil
}

// ===================== SERVICE =====================

func (r *Repository) GetServices(query string) ([]models.Service, error) {
    var services []models.Service
    q := r.db.Model(&models.Service{}).Where("is_active = true")
    if query != "" {
        like := "%" + strings.ToLower(query) + "%"
        q = q.Where("lower(name) LIKE ? OR lower(description) LIKE ?", like, like)
    }
    err := q.Order("name ASC").Find(&services).Error
    return services, err
}

func (r *Repository) GetServiceByID(id int64) (*models.Service, error) {
    var s models.Service
    err := r.db.First(&s, id).Error
    if err == gorm.ErrRecordNotFound {
        return nil, errors.New("not found")
    }
    return &s, err
}

func (r *Repository) CreateService(s *models.Service) error {
    return r.db.Create(s).Error
}

func (r *Repository) UpdateService(id int64, name string, price float64) error {
    return r.db.Model(&models.Service{}).
        Where("id = ?", id).
        Updates(map[string]any{
            "name":  name,
            "price": price,
        }).Error
}

func (r *Repository) DeleteService(id int64) error {
    return r.db.Delete(&models.Service{}, id).Error
}

func (r *Repository) UpdateServiceImage(id int64, url string) error {
    return r.db.Model(&models.Service{}).
        Where("id = ?", id).
        Update("image_url", url).Error
}

// ===================== REQUEST =====================

// GetOrCreateDraft получает или создает черновик заявки для пользователя
func (r *Repository) GetOrCreateDraft(userID int64) (*models.Request, error) {
    var req models.Request
    err := r.db.Where("user_id = ? AND status = 'draft'", userID).First(&req).Error
    if err == gorm.ErrRecordNotFound {
        req = models.Request{
            UserID:    userID,
            Status:    "draft",
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        err = r.db.Create(&req).Error
    }
    return &req, err
}

// AddServiceToDraft добавляет услугу в черновик
func (r *Repository) AddServiceToDraft(userID int64, serviceID int) error {
    req, err := r.GetOrCreateDraft(userID)
    if err != nil {
        return err
    }
    var existing models.RequestService
    err = r.db.Where("request_id = ? AND service_id = ?", req.ID, serviceID).First(&existing).Error
    if err == gorm.ErrRecordNotFound {
        rs := models.RequestService{RequestID: req.ID, ServiceID: serviceID, Quantity: 1}
        return r.db.Create(&rs).Error
    } else if err != nil {
        return err
    }
    existing.Quantity++
    return r.db.Save(&existing).Error
}

func (r *Repository) RemoveServiceFromRequest(requestID, serviceID int) error {
    res := r.db.Delete(&models.RequestService{}, "request_id = ? AND service_id = ?", requestID, serviceID)
    if res.RowsAffected == 0 {
        return errors.New("not found")
    }
    return res.Error
}

func (r *Repository) UpdateQuantity(requestID, serviceID, quantity int) error {
    res := r.db.Model(&models.RequestService{}).
        Where("request_id = ? AND service_id = ?", requestID, serviceID).
        Update("quantity", quantity)
    if res.RowsAffected == 0 {
        return errors.New("not found")
    }
    return res.Error
}

// LoadDraftWithItems загружает черновик со всеми услугами
func (r *Repository) LoadDraftWithItems(userID int64) (*models.Request, error) {
    req, err := r.GetOrCreateDraft(userID)
    if err != nil {
        return nil, err
    }
    err = r.db.Preload("Services").Preload("Services.Service").First(req, req.ID).Error
    return req, err
}

// GetRequestByID получает заявку по ID
func (r *Repository) GetRequestByID(id int64) (*models.Request, error) {
    var req models.Request
    err := r.db.Preload("Services").Preload("Services.Service").First(&req, id).Error
    if err == gorm.ErrRecordNotFound {
        return nil, errors.New("not found")
    }
    return &req, err
}

// GetRequestByIDAndUser получает заявку по ID и проверяет владельца
func (r *Repository) GetRequestByIDAndUser(id, userID int64) (*models.Request, error) {
    var req models.Request
    err := r.db.Preload("Services").Preload("Services.Service").
        Where("id = ? AND user_id = ?", id, userID).
        First(&req).Error
    if err == gorm.ErrRecordNotFound {
        return nil, errors.New("not found")
    }
    return &req, err
}

func (r *Repository) SubmitRequest(id int64) error {
    var req models.Request
    err := r.db.First(&req, id).Error
    if err != nil {
        return err
    }
    req.Status = "formed"
    req.UpdatedAt = time.Now()
    return r.db.Save(&req).Error
}

func (r *Repository) CompleteRequest(id int64) error {
    var req models.Request
    err := r.db.First(&req, id).Error
    if err != nil {
        return err
    }
    if req.Status != "formed" {
        return errors.New("заявка должна быть сформирована перед завершением")
    }
    req.Status = "completed"
    req.UpdatedAt = time.Now()
    return r.db.Save(&req).Error
}

// ListRequests получает список заявок (для пользователя - только свои, для модератора - все)
func (r *Repository) ListRequests(userID int64, isModerator bool, dateFrom, dateTo string) ([]models.Request, error) {
    var reqs []models.Request
    q := r.db.Where("status != 'deleted'")

    // Если не модератор - показываем только свои заявки
    if !isModerator {
        q = q.Where("user_id = ?", userID)
    }

    // Опциональный фильтр по дате
    if dateFrom != "" {
        q = q.Where("created_at >= ?", dateFrom)
    }
    if dateTo != "" {
        q = q.Where("created_at <= ?", dateTo)
    }

    err := q.Preload("Services").
        Preload("Services.Service").
        Find(&reqs).Error
    return reqs, err
}

func (r *Repository) DeleteRequestLogical(id int64) error {
    var req models.Request
    err := r.db.First(&req, id).Error
    if err != nil {
        return err
    }
    req.Status = "deleted"
    req.UpdatedAt = time.Now()
    return r.db.Save(&req).Error
}
