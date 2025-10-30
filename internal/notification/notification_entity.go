package notification

import (
	"github.com/easy-comerce/backend/pkg/base"
	// "gorm.io/datatypes"
)

type Notification struct {
	base.AuditEntity
	UserID  uint   `json:"user_id" gorm:"index; column:user_id"`
	Title   string `json:"title" gorm:"column:message"`
	Message string `json:"message" gorm:"column:message"`
	Type    string `json:"type" gorm:"column:type"`
	IsRead  bool   `json:"is_read" gorm:"default:false; column:is_read"`
	// Data    *datatypes.JSON `json:"data" gorm:"type:json; column:data"`
}

type FCMMessage struct {
	To           string      `json:"to,omitempty"`
	Notification *FCMPayload `json:"notification,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type FCMPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
