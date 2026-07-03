package providers

import (
	"fmt"
	"log"

	"go-notify/internal/models"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridProvider struct {
	APIKey      string
	SenderEmail string
	SenderName  string
}

func NewSendGridProvider(apiKey, senderEmail, senderName string) *SendGridProvider {
	return &SendGridProvider{
		APIKey:      apiKey,
		SenderEmail: senderEmail,
		SenderName:  senderName,
	}
}

func (s *SendGridProvider) Send(req models.NotificationRequest) error {
	from := mail.NewEmail(s.SenderName, s.SenderEmail)
	to := mail.NewEmail("User", req.Target)
	message := mail.NewSingleEmail(from, req.Title, to, req.Body, req.Body)

	client := sendgrid.NewSendClient(s.APIKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("sendgrid API request failed: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid rejected message with status %d: %s", response.StatusCode, response.Body)
	}

	log.Printf("[SENDGRID SUCCESS] Dispatched email to %s (Status: %d)", req.Target, response.StatusCode)
	return nil
}
