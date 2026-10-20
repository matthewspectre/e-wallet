package router

import (
	"net/http"

	"e-wallet/internal/handler"
)

func New(h *handler.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	return mux
}
