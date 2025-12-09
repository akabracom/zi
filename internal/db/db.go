package db

import (
    "fmt"
    "os"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func DSNFromEnv() string {
    // Отладочный вывод всех переменных окружения
    log.Printf("=== DEBUG: Environment Variables ===")
    log.Printf("DB_HOST from env: '%s'", os.Getenv("DB_HOST"))
    log.Printf("DB_PORT from env: '%s'", os.Getenv("DB_PORT"))
    log.Printf("DB_USER from env: '%s'", os.Getenv("DB_USER"))
    log.Printf("DB_PASSWORD from env: '%s'", os.Getenv("DB_PASSWORD"))
    log.Printf("DB_NAME from env: '%s'", os.Getenv("DB_NAME"))
    log.Printf("===================================")
    
    host := os.Getenv("DB_HOST")
    if host == "" {
        host = "deposits-db"
        log.Printf("WARNING: DB_HOST is empty, using default: deposits-db")
    }
    port := os.Getenv("DB_PORT")
    if port == "" {
        port = "5432"
        log.Printf("WARNING: DB_PORT is empty, using default: 5432")
    }
    user := os.Getenv("DB_USER")
    if user == "" {
        user = "user"
        log.Printf("WARNING: DB_USER is empty, using default: user")
    }
    password := os.Getenv("DB_PASSWORD")
    if password == "" {
        password = "password"
        log.Printf("WARNING: DB_PASSWORD is empty, using default: password")
    }
    dbname := os.Getenv("DB_NAME")
    if dbname == "" {
        dbname = "deposits"
        log.Printf("WARNING: DB_NAME is empty, using default: deposits")
    }

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    log.Printf("DSN: host=%s port=%s user=%s dbname=%s", host, port, user, dbname)
    
    return dsn
}

func New(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get database instance: %v", err)
    }

    if err = sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %v", err)
    }

    log.Println("Successfully connected to database")
    return db, nil
}
