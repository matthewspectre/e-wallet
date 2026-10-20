package main

import (
	"log"
	"net/http"
	"os"

	"e-wallet/internal/handler"
	"e-wallet/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiRouter := router.New(handler.New())
	server := &http.Server{
		Addr:    ":" + port,
		Handler: apiRouter,
	}

	log.Printf("e-wallet API listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
