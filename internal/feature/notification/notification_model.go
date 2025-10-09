package notification

import "time"

type NotificationType string

const (
	TypeInfo    NotificationType = "INFO"
	TypeSuccess NotificationType = "SUCCESS"
	TypeWarning NotificationType = "WARNING"
	TypeError   NotificationType = "ERROR"
)

type Notification struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	UserID    uint            `json:"userId" gorm:"index"`
	Title     string          `json:"title"`
	Message   string          `json:"message"`
	Type      NotificationType `json:"type"`
	IsRead    bool            `json:"isRead" gorm:"default:false"`
	Data      JSONMap         `json:"data,omitempty" gorm:"type:jsonb"`
	CreatedAt time.Time       `json:"createdAt"`
}

type JSONMap map[string]interface{}

type FCMMessage struct {
	To           string      `json:"to,omitempty"`
	Notification *FCMPayload `json:"notification,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type FCMPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
