package main

import (
	"go-notify/internal/broker"
	"go-notify/internal/cache"
	"go-notify/internal/handlers"
	"log"
	"net/http"
)

func main() {
	cache.InitRedis()
	broker.InitNATS()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/notify", handlers.SendNotificationHandler)

	log.Println("Notification API Service running on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
