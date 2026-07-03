package providers

import "go-notify/internal/models"

type NotificationProvider interface {
	Send(req models.NotificationRequest) error 
}