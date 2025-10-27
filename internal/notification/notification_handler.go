package notification

import (
	"encoding/json"
	"net/http"

	c "github.com/easy-comerce/backend/pkg/constants"
	cu "github.com/easy-comerce/backend/pkg/contextutil"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type NotificationHandler interface {
	GetNotifications(w http.ResponseWriter, r *http.Request)
	MarkAsRead(w http.ResponseWriter, r *http.Request)
	MarkAllAsRead(w http.ResponseWriter, r *http.Request)
	GetUnreadCount(w http.ResponseWriter, r *http.Request)
	SendPush(w http.ResponseWriter, r *http.Request)
}

type notificationHandler struct {
	service NotificationService
}

func NewNotificationHandler(service NotificationService) NotificationHandler {
	return &notificationHandler{
		service: service,
	}
}

func (h *notificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication Required", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))

	notifications, err := h.service.GetUserNotifications(showDeleted, *userID, pageStr, pageSizeStr, sortBy, sortOrder)
	response.SendResponse(w, notifications, err, http.StatusInternalServerError)
}

func (h *notificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication Required", http.StatusUnauthorized)
		return
	}

	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil {
		response.SendErrorJSON(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	err = h.service.MarkNotificationAsRead(*id, *userID)
	response.SendResponse(w, "Notification marked as read", err, http.StatusInternalServerError)
}

func (h *notificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication Required", http.StatusUnauthorized)
		return
	}

	err = h.service.MarkAllNotificationsAsRead(*userID)
	response.SendResponse(w, "All notifications marked as read", err, http.StatusInternalServerError)
}

func (h *notificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, err := cu.GetUserIDFromContext(r.Context())
	if err != nil {
		response.SendErrorJSON(w, "Authentication Required", http.StatusUnauthorized)
		return
	}

	count, err := h.service.GetUnreadNotificationsCount(*userID)
	response.SendResponse(w, map[string]int64{"count": count}, err, http.StatusInternalServerError)
}

func (h *notificationHandler) SendPush(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := cu.GetUserIDFromContext(ctx)
	if err != nil {
		response.SendErrorJSON(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		Title       string                 `json:"title"`
		Message     string                 `json:"message"`
		Type        string                 `json:"type"`
		Data        map[string]interface{} `json:"data,omitempty"`
		DeviceToken string                 `json:"deviceToken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.DeviceToken == "" {
		response.SendErrorJSON(w, "Device token is required", http.StatusBadRequest)
		return
	}

	notification, err := h.service.CreateNotification(ctx, *userID, req.Title, req.Message, req.Type, req.Data)
	if err != nil {
		response.SendErrorJSON(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	err = h.service.SendPushNotification(ctx, req.DeviceToken, notification)
	response.SendResponse(w, "Push notification sent successfully", err, http.StatusInternalServerError)
}
