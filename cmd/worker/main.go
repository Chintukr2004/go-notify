package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-notify/internal/broker"
	"go-notify/internal/models"

	"github.com/nats-io/nats.go"
)

func main() {
	// 1. Connect to NATS
	broker.InitNATS()

	// 2. Subscribe to the stream
	sub, err := broker.JS.QueueSubscribe("NOTIFY.*", "worker-group", processMessage, nats.ManualAck())
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	log.Println("Worker service is running and waiting for jobs...")

	// 3. Keep the application running until we press Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Worker shutting down cleanly...")
}

// processMessage handles the actual notification logic
func processMessage(msg *nats.Msg) {
	var req models.NotificationRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		log.Printf("Error decoding message: %v\n", err)
		// Nak (Negative Acknowledgement) tells NATS to retry this later
		msg.Nak()
		return
	}

	log.Printf("[START] Processing %s for %s (ID: %s)\n", req.Type, req.Target, req.ID)

	// Simulate calling an external API like SendGrid or Twilio (Network Delay)
	time.Sleep(2 * time.Second)

	log.Printf("[DONE] Successfully delivered %s to %s\n", req.Type, req.Target)

	// Ack (Acknowledgement) tells NATS the job is complete and can be deleted from the queue
	if err := msg.Ack(); err != nil {
		log.Printf("Failed to acknowledge message: %v\n", err)
	}
}