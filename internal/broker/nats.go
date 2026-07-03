package broker

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

var JS nats.JetStreamContext

func InitNATS() {
	//1. Connect to local NATS server
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	//2. Enable JetStream
	JS, err = nc.JetStream()
	if err != nil {
		log.Fatalf("Failed to create JetStream Context: %v", err)
	}

	// 3. Create a "Stream" (Queue) for our notifications
	streamNamme := "NOTIFICATIONS"
	_, err = JS.StreamInfo(streamNamme)
	if err != nil {
		_, err = JS.AddStream(&nats.StreamConfig{
			Name:     streamNamme,
			Subjects: []string{"NOTIFY.*"},
			Storage:  nats.FileStorage,
		})
		if err != nil {
			log.Fatalf("Failed to create NATS stream: %v", err)
		}
		log.Println("Created NATS JetStream: NOTIFICATIONS")
	} else {
		log.Println("Connected to existing NATS JetStream: NOTIFICATIONS")
	}

}
