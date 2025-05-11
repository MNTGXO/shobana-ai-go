package api

import (
    "io"
    "log"
    "net/http"
    "strings"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func init() {
    var err error
    bot, err = tgbotapi.NewBotAPI(Cfg.Token)
    if err != nil {
        log.Fatalf("bot init failed: %v", err)
    }
    // set webhook
    hookURL := fmt.Sprintf("https://%s/api/%s", Cfg.VercelURL, Cfg.WebhookSecret)
    _, err = bot.Request(tgbotapi.DeleteWebhookConfig{})
    if err != nil {
        log.Println("warning: could not delete old webhook:", err)
    }
    _, err = bot.Request(tgbotapi.NewWebhook(hookURL))
    if err != nil {
        log.Fatalf("webhook setup failed: %v", err)
    }
}

// Handler is the Vercel function entrypoint.
func Handler(w http.ResponseWriter, r *http.Request) {
    // verify secret
    if !strings.HasSuffix(r.URL.Path, Cfg.WebhookSecret) {
        http.NotFound(w, r)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "bad request", 400)
        return
    }

    update, err := tgbotapi.ParseUpdate(body)
    if err != nil {
        http.Error(w, "invalid update", 400)
        return
    }

    go processUpdate(update)
    w.WriteHeader(http.StatusOK)
}

func processUpdate(update tgbotapi.Update) {
    if update.Message == nil {
        return
    }

    msg := update.Message
    if msg.From != nil && msg.From.IsBot {
        return
    }

    // Handle /start
    if msg.IsCommand() && msg.Command() == "start" {
        text := "👋 Hey! I’m Shobana, your ChatGPT-powered assistant."
        bot.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
        return
    }

    // Handle replies to the bot
    if msg.ReplyToMessage != nil && msg.ReplyToMessage.From != nil && msg.ReplyToMessage.From.ID == bot.Self.ID {
        reply, err := FetchAIResponse(r.Context(), msg.Text)
        if err != nil {
            bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "❌ Error: "+err.Error()))
        } else {
            bot.Send(tgbotapi.NewMessage(msg.Chat.ID, reply))
        }
        return
    }

    // default echo
    bot.Send(tgbotapi.NewMessage(msg.Chat.ID, msg.Text))
}
