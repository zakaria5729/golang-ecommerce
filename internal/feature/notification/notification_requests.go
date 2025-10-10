package notification

var CreateNotification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	// Type        NotificationType       `json:"type"`
	Data        map[string]interface{} `json:"data,omitempty"`
	DeviceToken string                 `json:"deviceToken"`
}
