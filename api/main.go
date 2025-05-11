package api

import (
    _ "embed"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })
    // route webhook: /<secret>
    http.HandleFunc("/"+Cfg.WebhookSecret, Handler)
    http.ListenAndServe(":3000", nil)
}
