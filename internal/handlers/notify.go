package handlers

import (
	"encoding/json"
	"go-notify/internal/broker"
	"go-notify/internal/cache"
	"go-notify/internal/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const RateLimitMax = 3
const RateLimitWindow = time.Minute

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

	//idempotency check
	isNew, err := cache.Client.SetNX(cache.Ctx, "idem:"+req.ID, "processing", 24*time.Hour).Result()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !isNew {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":          "already_processed",
			"notification_id": req.ID,
		})
	}

	//2. Rate limiting check

	rateKey := "rate:" + req.ID
	count, err := cache.Client.Incr(cache.Ctx, rateKey).Result()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if count == 1 {
		cache.Client.Expire(cache.Ctx, rateKey, RateLimitWindow)
	}

	if count > RateLimitMax {
		http.Error(w, "Too many requests. Please try again later", http.StatusTooManyRequests)
		return
	}

	// 3. Publish to Message Broker 
	// Convert our struct back to a JSON byte array for the queue
	payload, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to process payload", http.StatusInternalServerError)
		return
	}

	// Route based on type (e.g., "NOTIFY.EMAIL" or "NOTIFY.SMS")
	subject := "NOTIFY." + string(req.Type)

	// Publish the message to NATS
	_, err = broker.JS.Publish(subject, payload)
	if err != nil {
		http.Error(w, "Failed to queue notification", http.StatusInternalServerError)
		return
	}

	// 4. Success Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":          "queued",
		"notification_id": req.ID,
	})
}
