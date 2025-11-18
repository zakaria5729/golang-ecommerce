package notification

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/easy-comerce/backend/pkg/logger"
)

type FCMService struct {
	serverKey string
	fcmURL    string
}

func NewFCMService() *FCMService {
	return &FCMService{
		serverKey: os.Getenv("FCM_SERVER_KEY"),
		fcmURL:    os.Getenv("FCM_URL"),
	}
}

func (s *FCMService) SendNotification(deviceToken string, notification *Notification) error {
	if s.serverKey == "" || s.fcmURL == "" {
		logger.Warn("FCM configuration is missing")
		return nil
	}

	msg := FCMMessage{
		To: deviceToken,
		Notification: &FCMPayload{
			Title: notification.Title,
			Body:  notification.Message,
		},
		// Data: notification.Data,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.fcmURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "key="+s.serverKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}
		logger.Error("Failed to send FCM notification", "error", result)
		return err
	}

	return nil
}
