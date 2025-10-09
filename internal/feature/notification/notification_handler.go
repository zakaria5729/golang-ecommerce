package notification

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/response"
)

type NotificationHandler struct {
	notificationUC *NotificationUseCase
}

func NewNotificationHandler(notificationUC *NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		notificationUC: notificationUC,
	}
}

type NotificationRequest struct {
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Type    NotificationType       `json:"type"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type SendPushRequest struct {
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Type        NotificationType       `json:"type"`
	Data        map[string]interface{} `json:"data,omitempty"`
	DeviceToken string                 `json:"deviceToken"`
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.SendErrorJSON(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}

	showDeleted := false
	if r.URL.Query().Get("showDeleted") == "true" {
		showDeleted = true
	}

	notifications, err := h.notificationUC.GetUserNotifications(ctx, user.ID, page, limit, &showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to get notifications", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, "Failed to get notifications", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, notifications)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.SendErrorJSON(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	notificationID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.SendErrorJSON(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	if err := h.notificationUC.MarkNotificationAsRead(ctx, uint(notificationID), user.ID); err != nil {
		logger.Logger.Error("Failed to mark notification as read", "error", err, "notificationID", notificationID, "userID", user.ID)
		response.SendErrorJSON(w, "Failed to mark notification as read", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.SendErrorJSON(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if err := h.notificationUC.MarkAllNotificationsAsRead(ctx, user.ID); err != nil {
		logger.Logger.Error("Failed to mark all notifications as read", "error", err, "userID", user.ID)
		response.SendErrorJSON(w, "Failed to mark all notifications as read", http.StatusInternalServerError)
		return
	}

	response.SendSuccessJSON(w, map[string]string{"message": "All notifications marked as read"})
}

func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.GetUserFromContext(ctx)
	if err != nil {
		response.SendErrorJSON(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	count, err := h.notificationUC.GetUnreadNotificationsCount(ctx, user.ID)
	if err != nil {
		logger.Logger.Error("Failed to get unread notifications count", "error", err, "userID", user.ID)
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

	var req SendPushRequest
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
