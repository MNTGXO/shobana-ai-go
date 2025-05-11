package api

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

func FetchAIResponse(ctx context.Context, prompt string) (string, error) {
    payload := map[string]string{"message": prompt}
    body, _ := json.Marshal(payload)

    req, err := http.NewRequestWithContext(ctx, "POST", Cfg.ChatAPIURL, bytes.NewBuffer(body))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result struct {
        Reply string `json:"reply"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }
    return result.Reply, nil
}
