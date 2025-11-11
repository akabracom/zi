package repository

import (
	"errors"
	"strings"
	"time"

	"backend/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
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
		return nil, errors.New("service not found")
	}
	return &s, err
}

func (r *Repository) CreateService(s *models.Service) error {
	s.CreatedAt = time.Now()
	return r.db.Create(s).Error
}

func (r *Repository) UpdateService(id int64, name string, price float64) error {
	s, err := r.GetServiceByID(id)
	if err != nil {
		return err
	}
	s.Name = name
	s.Price = price
	return r.db.Save(s).Error
}

func (r *Repository) DeleteService(id int64) error {
	s, err := r.GetServiceByID(id)
	if err != nil {
		return err
	}
	s.IsActive = false
	return r.db.Save(s).Error
}

func (r *Repository) UpdateServiceImage(id int64, imageURL string) error {
	s, err := r.GetServiceByID(id)
	if err != nil {
		return err
	}
	s.ImageURL = imageURL
	return r.db.Save(s).Error
}

// ===================== REQUEST =====================

func (r *Repository) ensureUser(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		u.Email = email
		u.Name = "Student"
		if err = r.db.Create(&u).Error; err != nil {
			return nil, err
		}
		return &u, nil
	}
	return &u, err
}

func (r *Repository) GetOrCreateDraft(email string) (*models.Request, error) {
	u, err := r.ensureUser(email)
	if err != nil {
		return nil, err
	}
	var req models.Request
	err = r.db.Where("user_id = ? AND status = 'draft'", u.ID).First(&req).Error
	if err == gorm.ErrRecordNotFound {
		req = models.Request{UserID: u.ID, Status: "draft", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		err = r.db.Create(&req).Error
	}
	return &req, err
}

func (r *Repository) AddServiceToDraft(email string, serviceID int) error {
	req, err := r.GetOrCreateDraft(email)
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

func (r *Repository) LoadDraftWithItems(email string) (*models.Request, error) {
	req, err := r.GetOrCreateDraft(email)
	if err != nil {
		return nil, err
	}
	err = r.db.Preload("Services").Preload("Services.Service").First(req, req.ID).Error
	return req, err
}

func (r *Repository) GetRequestByID(id int64) (*models.Request, error) {
	var req models.Request
	err := r.db.Preload("Services").Preload("Services.Service").First(&req, id).Error
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
	req.Status = "submitted"
	req.UpdatedAt = time.Now()
	return r.db.Save(&req).Error
}

func (r *Repository) CompleteRequest(id int64) error {
	var req models.Request
	err := r.db.First(&req, id).Error
	if err != nil {
		return err
	}
	req.Status = "completed"
	req.UpdatedAt = time.Now()
	return r.db.Save(&req).Error
}

// В конец файла repository.go
func (r *Repository) ListRequests(email string) ([]models.Request, error) {
	u, err := r.ensureUser(email)
	if err != nil {
		return nil, err
	}
	var reqs []models.Request
	err = r.db.Where("user_id = ? AND status != 'deleted'", u.ID).
		Preload("Services").
		Preload("Services.Service").
		Find(&reqs).Error
	return reqs, err
}

// internal/repository/repository.go
// ... (остальной код)

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