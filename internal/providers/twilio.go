package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"go-notify/internal/models"
)

type TwilioProvider struct {
	AccountSID  string
	AuthToken   string
	FromNumber  string
}

func NewTwilioProvider(sid, token, fromNumber string) *TwilioProvider {
	return &TwilioProvider{
		AccountSID: sid,
		AuthToken:  token,
		FromNumber: fromNumber,
	}
}

func (t *TwilioProvider) Send(req models.NotificationRequest) error {
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.AccountSID)

	data := url.Values{}
	data.Set("To", req.Target)
	data.Set("From", t.FromNumber)
	data.Set("Body", fmt.Sprintf("%s: %s", req.Title, req.Body))

	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	httpReq.SetBasicAuth(t.AccountSID, t.AuthToken)
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("twilio network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("twilio failed (%d): %v", resp.StatusCode, errResp)
	}

	log.Printf("[TWILIO SUCCESS] Dispatched SMS to %s", req.Target)
	return nil
}