package http

import (
	"asai/internal/core"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func Run(ctx context.Context, a *core.Agent) {
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID int64  `json:"user_id"`
			Input  string `json:"input"`
		}

		json.NewDecoder(r.Body).Decode(&req)
		resp, err := a.HandleInput(ctx, req.UserID, req.Input)
		if err != nil {
			http.Error(w, "Error", 500) //сделать ошибку
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"reply": resp})
	})

	log.Println("[http] сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
