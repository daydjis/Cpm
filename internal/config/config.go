package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	TelegramToken  string
	DatabaseURL    string
	ScrapeInterval time.Duration
	ScrapeBaseURL  string
)

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	TelegramToken = getEnv("TG_BOT_TOKEN", "")
	if TelegramToken == "" {
		TelegramToken = getEnv("TELEGRAM_TOKEN", "")
	}
	DatabaseURL = getEnv("DATABASE_URL", "")
	ScrapeBaseURL = getEnv("SCRAPE_BASE_URL", "https://minifreemarket.com/catalog/games-with-miniatures/warhammer-40000")
	if ScrapeBaseURL != "" && ScrapeBaseURL[len(ScrapeBaseURL)-1] == '/' {
		ScrapeBaseURL = ScrapeBaseURL[:len(ScrapeBaseURL)-1]
	}

	intervalStr := getEnv("SCRAPE_INTERVAL", "10m")
	d, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Printf("Invalid SCRAPE_INTERVAL %q, using 10m: %v", intervalStr, err)
		d = 10 * time.Minute
	}
	ScrapeInterval = d
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
