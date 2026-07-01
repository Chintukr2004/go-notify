package handlers

import (
	"encoding/json"
	"go-notify/internal/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func SendNotificationHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.NotificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// basic validation
	if req.UserID == "" || req.Target == "" || req.Body == "" {
		http.Error(w, "Missing required fields: user_id, target, or body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	req.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":          "queued",
		"notification_id": req.ID,
	})
}
