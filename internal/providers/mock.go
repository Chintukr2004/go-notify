package providers

import (
	"errors"
	"log"
	"strings"

	"go-notify/internal/models"
)

// MockProvider simulates delivery for local testing and DLQ demonstrations.
type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) Send(req models.NotificationRequest) error {
	// Intentionally trigger a failure if the target contains "fail" to demonstrate DLQ routing
	if strings.Contains(req.Target, "fail") {
		log.Printf("[MOCK ERROR] Simulated delivery rejection for: %s", req.Target)
		return errors.New("simulated provider network failure")
	}

	log.Printf("[MOCK SUCCESS] Simulated instant delivery to: %s", req.Target)
	return nil
}