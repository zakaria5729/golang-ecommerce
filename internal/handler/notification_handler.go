package handler

import (
	"encoding/json"
	"net/http"

	"github.com/easy-comerce/backend/internal/feature/notification"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
)

type NotificationHandler struct {
	notificationUC *notification.NotificationUseCase
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notificationUC: notification.NewNotificationUseCase(),
	}
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication Required", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()
	pageStr := q.Get(c.Page)
	pageSizeStr := q.Get(c.PageSize)
	sortBy := q.Get(c.SortBy)
	sortOrder := q.Get(c.SortOrder)
	showDeleted := utils.ParseBoolPtr(q.Get(c.ShowDeleted))

	notifications, err := h.notificationUC.GetUserNotifications(showDeleted, *userID, pageStr, pageSizeStr, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to get notifications", "error", err, "userID", userID)
		response.SendErrorJSON(w, "Failed to get notifications", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, notifications)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication Required", http.StatusUnauthorized)
		return
	}

	id, err := utils.ParseUint(r.PathValue(c.FieldID))
	if err != nil {
		response.SendErrorJSON(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	if err := h.notificationUC.MarkNotificationAsRead(*id, *userID); err != nil {
		response.SendErrorJSON(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication Required", http.StatusUnauthorized)
		return
	}

	if err := h.notificationUC.MarkAllNotificationsAsRead(*userID); err != nil {
		logger.Logger.Error("Failed to mark all notifications as read", "error", err, "userID", userID)
		response.SendErrorJSON(w, "Failed to mark all notifications as read", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "All notifications marked as read"})
}

func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil || userID == nil || *userID == 0 {
		response.SendErrorJSON(w, "Authentication Required", http.StatusUnauthorized)
		return
	}

	count, err := h.notificationUC.GetUnreadNotificationsCount(*userID)
	if err != nil {
		logger.Logger.Error("Failed to get unread notifications count", "error", err, "userID", userID)
		response.SendErrorJSON(w, "Failed to get unread notifications count", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]int64{"count": count})
}

func (h *NotificationHandler) SendPush(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.GetUserFromContext(ctx)
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

	notification, err := h.notificationUC.CreateNotification(ctx, user.ID, req.Title, req.Message, req.Type, req.Data)
	if err != nil {
		logger.Logger.Error("Failed to create notification", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	if err := h.notificationUC.SendPushNotification(ctx, req.DeviceToken, notification); err != nil {
		logger.Logger.Error("Failed to send push notification", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, "Failed to send push notification", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Push notification sent successfully"})
}
