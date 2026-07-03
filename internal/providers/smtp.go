package providers

import (
	"fmt"
	"log"
	"net/smtp"

	"go-notify/internal/models"
)

type SMTPProvider struct {
	EmailAddress string
	AppPassword  string
	Host         string
	Port         string
}

func NewSMTPProvider(email, password string) *SMTPProvider {
	return &SMTPProvider{
		EmailAddress: "chintukr1904@gmail.com",
		AppPassword:  "lkvt sffr grip wait",
		Host:         "smtp.gmail.com",
		Port:         "587",
	}
}

func (s *SMTPProvider) Send(req models.NotificationRequest) error {
	subject := fmt.Sprintf("Subject: %s\n", req.Title)
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	message := []byte(subject + mime + req.Body)

	auth := smtp.PlainAuth("", s.EmailAddress, s.AppPassword, s.Host)
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	err := smtp.SendMail(addr, auth, s.EmailAddress, []string{req.Target}, message)
	if err != nil {
		return fmt.Errorf("SMTP delivery failed: %w", err)
	}

	log.Printf("[SMTP SUCCESS] Delivered email to %s", req.Target)
	return nil
}
