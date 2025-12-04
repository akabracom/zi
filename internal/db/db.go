package db

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func DSNFromEnv() string {
    return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        getEnv("DB_HOST", "localhost"),
        getEnv("DB_USER", "user"),
        getEnv("DB_PASSWORD", "password"),
        getEnv("DB_NAME", "deposits"),
        getEnv("DB_PORT", "5432"),
    )
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}

func New(dsn string) (*gorm.DB, error) {
    var db *gorm.DB
    var err error
    for i := 0; i < 20; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            if sqlDB, _ := db.DB(); sqlDB.Ping() == nil {
                log.Println("DB connected")
                return db, nil
            }
        }
        log.Printf("DB not ready (%d/20)", i+1)
        time.Sleep(2 * time.Second)
    }
    return nil, fmt.Errorf("failed after 20 retries: %w", err)
}