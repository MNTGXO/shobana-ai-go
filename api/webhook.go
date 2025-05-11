package handler

import (
    "io"
    "log"
    "net/http"
    "strings"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/MNTGXO/shobana-ai-go/utils"
)

// Handler is the Vercel function entrypoint.
func Handler(w http.ResponseWriter, r *http.Request) {
    // 1) Secret validation
    if !strings.HasSuffix(r.URL.Path, utils.Cfg.WebhookSecret) {
        http.NotFound(w, r)
        return
    }

    // 2) Read incoming update
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    update, err := tgbotapi.ParseUpdate(body)
    if err != nil {
        http.Error(w, "invalid update", http.StatusBadRequest)
        return
    }

    // 3) Process in background
    go processUpdate(update)
    w.WriteHeader(http.StatusOK)
}

func processUpdate(update tgbotapi.Update) {
    if update.Message == nil || update.Message.From.IsBot {
        return
    }

    bot, err := tgbotapi.NewBotAPI(utils.Cfg.Token)
    if err != nil {
        log.Printf("bot init error: %v", err)
        return
    }

    msg := update.Message

    // /start command
    if msg.IsCommand() && msg.Command() == "start" {
        bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "👋 Hey! I’m Shobana, your assistant."))
        return
    }

    // If replying to the bot → forward to AI
    if msg.ReplyToMessage != nil && msg.ReplyToMessage.From.ID == bot.Self.ID {
        reply, err := utils.FetchAIResponse(r.Context(), msg.Text)
        if err != nil {
            bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Error: "+err.Error()))
        } else {
            bot.Send(tgbotapi.NewMessage(msg.Chat.ID, reply))
        }
        return
    }

    // Default echo back
    bot.Send(tgbotapi.NewMessage(msg.Chat.ID, msg.Text))
}
