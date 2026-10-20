package router

import (
	"encoding/json"
	"net/http"
)

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(struct {
			Status string `json:"status"`
		}{Status: "ok"}); err != nil {
			http.Error(w, "failed to encode health response", http.StatusInternalServerError)
		}
	})
	return mux
}
