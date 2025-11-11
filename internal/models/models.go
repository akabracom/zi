package models

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Name      string
	CreatedAt time.Time `gorm:"not null"`
}

type Service struct {
	ID          int       `gorm:"primaryKey"`
	Name        string    `gorm:"uniqueIndex;not null"`
	Description string
	Price       float64   `gorm:"type:numeric(12,2);not null"`
	IsActive    bool      `gorm:"not null;default:true"`
	ImageURL    string
	CreatedAt   time.Time `gorm:"not null"`
}

type Request struct {
	ID        int64            `gorm:"primaryKey"`
	UserID    int64            `gorm:"not null;index"`
	Status    string           `gorm:"type:request_status;not null"`
	CreatedAt time.Time        `gorm:"not null"`
	UpdatedAt time.Time        `gorm:"not null"`
	Services  []RequestService `gorm:"foreignKey:RequestID"`
}

type RequestService struct {
	RequestID int64   `gorm:"primaryKey"`
	ServiceID int     `gorm:"primaryKey"`
	Quantity  int     `gorm:"not null;default:1"`
	Service   Service `gorm:"foreignKey:ID;references:ServiceID"`
}
