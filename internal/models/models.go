package models

import "time"

type User struct {
    ID        int64     `gorm:"primaryKey"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Name      string
    CreatedAt time.Time `gorm:"not null"`
}

type DepositOffer struct {
    ID              int       `gorm:"primaryKey"`
    Name            string    `gorm:"uniqueIndex;not null"`
    Days            int       `gorm:"not null"`
    InterestRatePct float64   `gorm:"type:numeric(5,2);not null"`
    StartDate       string
    EndDate         string
    Season          string
    Category        string
    Holidays        string
    ImageKey        string
    Description     string
    IsActive        bool      `gorm:"not null;default:true"`
    CreatedAt       time.Time `gorm:"not null"`
}

type DepositApplication struct {
    ID           int64     `gorm:"primaryKey"`
    UserID       int64     `gorm:"not null;index"`
    Status       string    `gorm:"type:deposit_application_status;not null"`
    CreatedAt    time.Time `gorm:"not null"`
    UpdatedAt    time.Time `gorm:"not null"`
    CommonAmount float64   `gorm:"type:numeric(12,2);not null"`
    TotalAmount  float64   `gorm:"type:numeric(12,2);not null"`
    TotalProfit  float64   `gorm:"type:numeric(12,2);not null"`
    Items        []DepositApplicationItem `gorm:"foreignKey:ApplicationID"`
}

type DepositApplicationItem struct {
    ApplicationID int64   `gorm:"primaryKey"`
    OfferID       int     `gorm:"primaryKey"`
    Quantity      int     `gorm:"not null;default:1"`
    Amount        float64 `gorm:"type:numeric(12,2);not null"`
    Profit        float64 `gorm:"type:numeric(12,2);not null"`

    // Даты вклада для конкретного предложения в рамках заявки
    StartDate string
    EndDate   string

    Offer DepositOffer `gorm:"foreignKey:ID;references:OfferID"`
}
