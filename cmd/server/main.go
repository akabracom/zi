package main

import (
    "log"

    "backend/internal/api"
    "backend/internal/cache"
    "backend/internal/db"
    "backend/internal/handler"
    "backend/internal/repository"
    "backend/internal/storage/minio"
    
    _ "backend/docs"
)

// @title Deposits API
// @version 1.0
// @description API для управления заявками и услугами с JWT авторизацией и Redis
// @host localhost:3001
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
    // --- PostgreSQL ---
    dsn := db.DSNFromEnv()
    gormDB, err := db.New(dsn)
    if err != nil {
        log.Fatalf("db connect: %v", err)
    }

    // --- Redis ---
    if err := cache.InitRedis(); err != nil {
        log.Fatalf("redis connect: %v", err)
    }

    // --- MinIO ---
    minioClient, err := minio.NewMinioClient()
    if err != nil {
        log.Fatalf("minio connect: %v", err)
    }

    // --- Repository ---
    repo := repository.New(gormDB)

    // --- Handler ---
    h := handler.New(repo, minioClient)

    // --- API ---
    api.Setup(h)
    api.Run()
}
