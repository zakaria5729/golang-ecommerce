package notification

import (
	"context"
	"errors"

	"github.com/easy-comerce/backend/pkg/logger"
	"gorm.io/gorm"
)

type NotificationUseCase struct {
	nofiticationRepo *NotificationRepository
	fcmService      *FCMService
}

func NewNotificationUseCase(notificationRepo *NotificationRepository, fcmService *FCMService) *NotificationUseCase {
	return &NotificationUseCase{
		nofiticationRepo: notificationRepo,
		fcmService:      fcmService,
	}
}

func (uc *NotificationUseCase) CreateNotification(ctx context.Context, userID uint, title, message string, notifType NotificationType, data map[string]interface{}) (*Notification, error) {
	notification := &Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notifType,
		Data:    data,
	}

	if err := uc.nofiticationRepo.Create(ctx, notification); err != nil {
		logger.Logger.Error("Failed to create notification", "method", "CreateNotification", "error", err, "userID", userID)
		return nil, err
	}

	return notification, nil
}

func (uc *NotificationUseCase) GetUserNotifications(ctx context.Context, userID uint, page, limit int, showDeleted *bool) ([]*Notification, error) {
	notifications, err := uc.nofiticationRepo.GetByUserID(ctx, userID, page, limit, showDeleted)
	if err != nil {
		logger.Logger.Error("Failed to get user notifications", "method", "GetUserNotifications", "error", err, "userID", userID)
		return nil, err
	}

	if notifications == nil {
		return []*Notification{}, nil
	}

	return notifications, nil
}

func (uc *NotificationUseCase) MarkNotificationAsRead(ctx context.Context, notificationID, userID uint) error {
	if err := uc.nofiticationRepo.MarkAsRead(ctx, notificationID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		logger.Logger.Error("Failed to mark notification as read", "method", "MarkNotificationAsRead", "error", err, "notificationID", notificationID, "userID", userID)
		return err
	}

	return nil
}

func (uc *NotificationUseCase) MarkAllNotificationsAsRead(ctx context.Context, userID uint) error {
	if err := uc.nofiticationRepo.MarkAllAsRead(ctx, userID); err != nil {
		logger.Logger.Error("Failed to mark all notifications as read", "method", "MarkAllNotificationsAsRead", "error", err, "userID", userID)
		return err
	}

	return nil
}

func (uc *NotificationUseCase) GetUnreadNotificationsCount(ctx context.Context, userID uint) (int64, error) {
	count, err := uc.nofiticationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		logger.Logger.Error("Failed to get unread notifications count", "method", "GetUnreadNotificationsCount", "error", err, "userID", userID)
		return 0, err
	}

	return count, nil
}

func (uc *NotificationUseCase) SendPushNotification(ctx context.Context, deviceToken string, notification *Notification) error {
	if deviceToken == "" {
		return errors.New("device token is required")
	}

	if notification == nil {
		return errors.New("notification is required")
	}

	if err := uc.fcmService.SendNotification(deviceToken, notification); err != nil {
		logger.Logger.Error("Failed to send push notification", "method", "SendPushNotification", "error", err, "deviceToken", deviceToken)
		return err
	}

	return nil
}
