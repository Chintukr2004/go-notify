package models

import "time"

type NotificationType string

const (
	Email NotificationType = "EMAIL"
	SMS   NotificationType = "SMS"
	Push  NotificationType = "PUSH"
)

type NotificationRequest struct {
	ID        string           `json:"id"`
	UserID    string           `json:"user_id"` // Unique ID for idempotency
	Type      NotificationType `json:"type"`    //EMAIL, SMS, PUSH
	Target    string           `json:"target"`  // Email adress or Phone number
	Title     string           `json:"title"`
	Body      string           `json:"body"`
	CreatedAt time.Time        `json:"created_at"`
}
