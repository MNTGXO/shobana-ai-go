package utils

import (
    "fmt"
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    Token         string
    OwnerID       int64
    ChatAPIURL    string
    WebhookSecret string
    VercelURL     string
}

var Cfg Config

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, reading from environment")
    }

    ownerID, err := strconv.ParseInt(os.Getenv("OWNER_ID"), 10, 64)
    if err != nil {
        log.Fatal("OWNER_ID must be a valid integer")
    }

    Cfg = Config{
        Token:         os.Getenv("TELEGRAM_TOKEN"),
        OwnerID:       ownerID,
        ChatAPIURL:    os.Getenv("CHAT_API_URL"),
        WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
        VercelURL:     os.Getenv("VERCEL_URL"),
    }

    if Cfg.Token == "" || Cfg.WebhookSecret == "" || Cfg.VercelURL == "" {
        log.Fatal("Required env vars: TELEGRAM_TOKEN, WEBHOOK_SECRET, VERCEL_URL")
    }
    fmt.Println("⚙️  Config loaded successfully")
}
