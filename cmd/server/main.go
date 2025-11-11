// cmd/server/main.go
package main

import (
	"log"

	"backend/internal/api"
	"backend/internal/db"
	"backend/internal/handler" // ДОБАВЛЕНО
	"backend/internal/repository"
	"backend/internal/storage/minio"
)

func main() {
	// --- PostgreSQL ---
	gormDB, err := db.New(db.DSNFromEnv())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// --- MinIO ---
	minioClient := minio.NewMinioClient()

	// --- Repository ---
	repo := repository.New(gormDB)

	// --- Handler ---
	h := handler.New(repo, minioClient)

	// --- API ---
	api.Setup(h)
	api.Run()
}