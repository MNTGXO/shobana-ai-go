package api

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Token       string
    OwnerID     int64
    ChatAPIURL  string
    WebhookSecret string
    VercelURL   string
}

var Cfg *Config

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, reading from environment")
    }

    owner := int64(0)
    fmt.Sscan(os.Getenv("OWNER_ID"), &owner)

    Cfg = &Config{
        Token:        os.Getenv("TELEGRAM_TOKEN"),
        OwnerID:      owner,
        ChatAPIURL:   os.Getenv("CHAT_API_URL"),
        WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
        VercelURL:    os.Getenv("VERCEL_URL"),
    }
    if Cfg.Token == "" || Cfg.WebhookSecret == "" || Cfg.VercelURL == "" {
        log.Fatal("Required env vars: TELEGRAM_TOKEN, WEBHOOK_SECRET, VERCEL_URL")
    }
}
