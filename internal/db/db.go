package db

import (
	"awesomeProject3/internal/config"
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	dbURL := config.DatabaseURL
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	if !strings.Contains(dbURL, "sslmode=") {
		if strings.Contains(dbURL, "?") {
			dbURL += "&sslmode=disable"
		} else {
			dbURL += "?sslmode=disable"
		}
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}

	log.Println("✅ Connected to Postgres!")
}
