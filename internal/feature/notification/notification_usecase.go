package notification

import (
	"context"
	"errors"

	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type NotificationUseCase struct {
	nofiticationRepo *NotificationRepository
	fcmService       *FCMService
}

func NewNotificationUseCase() *NotificationUseCase {
	return &NotificationUseCase{
		nofiticationRepo: NewNotificationRepository(),
		fcmService:       NewFCMService(),
	}
}

func (uc *NotificationUseCase) CreateNotification(ctx context.Context, userID uint, title, message string, notifType string, data map[string]interface{}) (*Notification, error) {
	notification := &Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notifType,
	}

	if err := uc.nofiticationRepo.Create(ctx, notification); err != nil {
		logger.Logger.Error("Failed to create notification", "method", "CreateNotification", "error", err, "userID", userID)
		return nil, err
	}

	return notification, nil
}

func (uc *NotificationUseCase) GetUserNotifications(showDeleted *bool, userID uint, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*models.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	notifications, total, err := uc.nofiticationRepo.GetAllByUserIDPaginated(showDeleted, userID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		logger.Logger.Error("Failed to get user notifications", "method", "GetUserNotifications", "error", err, "userID", userID)
		return nil, err
	}

	return utils.BuildPaginatedResponse(notifications, total, page, pageSize), nil
}

func (uc *NotificationUseCase) MarkNotificationAsRead(notificationID uint, userID uint) error {
	if err := uc.nofiticationRepo.MarkAsRead(notificationID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("notification not found with this id")
		}

		logger.Logger.Error("Failed to mark notification as read", "method", "MarkNotificationAsRead", "error", err, "notificationID", notificationID, "userID", userID)
		return err
	}

	return nil
}

func (uc *NotificationUseCase) MarkAllNotificationsAsRead(userID uint) error {
	if err := uc.nofiticationRepo.MarkAllAsRead(userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("notification not found with this id")
		}

		logger.Logger.Error("Failed to mark all notifications as read", "method", "MarkAllNotificationsAsRead", "error", err, "userID", userID)
		return err
	}

	return nil
}

func (uc *NotificationUseCase) GetUnreadNotificationsCount(userID uint) (int64, error) {
	count, err := uc.nofiticationRepo.GetUnreadCount(userID)
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
