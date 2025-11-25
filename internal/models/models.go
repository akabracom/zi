package models

import "time"

type User struct {
    ID        int64     `gorm:"primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex;not null" json:"email"`
    Name      string    `json:"name"`
    Password  string    `gorm:"not null" json:"-"` // НЕ возвращаем в JSON
    Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"` // user, moderator, guest
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Service struct {
    ID          int       `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;not null" json:"name"`
    Description string    `json:"description"`
    Price       float64   `gorm:"type:numeric(12,2);not null" json:"price"`
    IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
    ImageURL    string    `json:"image_url"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Request struct {
    ID        int64            `gorm:"primaryKey" json:"id"`
    UserID    int64            `gorm:"not null;index" json:"user_id"`
    Status    string           `gorm:"type:request_status;not null" json:"status"`
    CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
    Services  []RequestService `gorm:"foreignKey:RequestID" json:"services"`
}

type RequestService struct {
    RequestID int64   `gorm:"primaryKey" json:"request_id"`
    ServiceID int     `gorm:"primaryKey" json:"service_id"`
    Quantity  int     `gorm:"not null;default:1" json:"quantity"`
    Service   Service `gorm:"foreignKey:ID;references:ServiceID" json:"service"`
}

// DTO для регистрации
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email" example:"user@example.com"`
    Password string `json:"password" binding:"required,min=6" example:"password123"`
    Name     string `json:"name" example:"John Doe"`
}

// DTO для логина
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email" example:"user@example.com"`
    Password string `json:"password" binding:"required" example:"password123"`
}

// DTO ответа с токеном
type AuthResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
    User  User   `json:"user"`
}
