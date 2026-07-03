package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-notify/internal/broker"
	"go-notify/internal/models"
	"go-notify/internal/providers"

	"github.com/nats-io/nats.go"
)

var emailProvider providers.NotificationProvider
var smsProvider providers.NotificationProvider

func main() {
	broker.InitNATS()

	providerType := os.Getenv("EMAIL_PROVIDER")
	switch providerType {
	case "sendgrid":
		emailProvider = providers.NewSendGridProvider(
			os.Getenv("SENDGRID_API_KEY"),
			os.Getenv("SENDGRID_SENDER_EMAIL"),
			"Distributed Notify",
		)
	case "smtp":
		emailProvider = providers.NewSMTPProvider(
			os.Getenv("GMAIL_ADDRESS"),
			os.Getenv("GMAIL_APP_PASSWORD"),
		)
		log.Println("Using Gmail SMTP for Email Delivery")
	default:
		emailProvider = providers.NewMockProvider()
		log.Println("Using MockProvider for local testing/DLQ simulation")
	}

	if os.Getenv("TWILIO_ACCOUNT_SID") == "" {
		smsProvider = providers.NewMockProvider()
		log.Println("Using MockProvider for SMS Delivery (No Twilio keys found)")
	} else {
		smsProvider = providers.NewTwilioProvider(
			os.Getenv("TWILIO_ACCOUNT_SID"),
			os.Getenv("TWILIO_AUTH_TOKEN"),
			os.Getenv("TWILIO_FROM_NUMBER"),
		)
		log.Println("Using Twilio for SMS Delivery")
	}

	sub, err := broker.JS.QueueSubscribe("NOTIFY.*", "worker-group", processMessage, nats.ManualAck())
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	log.Println("Worker service running...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}

func processMessage(msg *nats.Msg) {
	var req models.NotificationRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		msg.Ack()
		return
	}

	meta, err := msg.Metadata()
	if err == nil && meta.NumDelivered > 3 {
		log.Printf("[DLQ] Message %s failed 3 times. Routing to DLQ.", req.ID)
		broker.JS.Publish("DLQ.NOTIFY", msg.Data)
		msg.Ack()
		return
	}

	var deliveryErr error
	switch req.Type {
	case models.Email:
		deliveryErr = emailProvider.Send(req)
	case models.SMS:
		deliveryErr = smsProvider.Send(req)
	}

	if deliveryErr != nil {
		log.Printf("[FAILED] Delivery failed: %v. Re-queuing...", deliveryErr)
		msg.Nak()
		return
	}

	log.Printf("[DONE] Job completed for %s", req.Target)
	msg.Ack()
}
