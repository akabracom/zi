// internal/db/db.go
package db

import (
	"fmt"
	"os"
)

func DSNFromEnv() string {
	// ИСПРАВЛЯЕМ: используем имя сервиса, а не localhost
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "deposits-db" // ← НЕ localhost!
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "deposits"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}