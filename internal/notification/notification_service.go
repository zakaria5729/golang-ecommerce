package notification

import (
	"context"
	"errors"

	r "github.com/easy-comerce/backend/pkg/response"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, userID uint, title, message string, notifType string, data map[string]interface{}) (*Notification, error)
	GetUserNotifications(showDeleted *bool, userID uint, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*r.PaginatedResponse, error)
	MarkNotificationAsRead(notificationID uint, userID uint) error
	MarkAllNotificationsAsRead(userID uint) error
	GetUnreadNotificationsCount(userID uint) (int64, error)
	SendPushNotification(ctx context.Context, deviceToken string, notification *Notification) error
}

type notificationService struct {
	repo       NotificationRepository
	fcmService *FCMService
}

func NewNotificationService(
	repo NotificationRepository,
	fcmService *FCMService,
) NotificationService {
	return &notificationService{
		repo:       repo,
		fcmService: fcmService,
	}
}

func (s *notificationService) CreateNotification(ctx context.Context, userID uint, title, message string, notifType string, data map[string]interface{}) (*Notification, error) {
	notification := &Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notifType,
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, errors.New("failed to create notification")
	}

	return notification, nil
}

func (s *notificationService) GetUserNotifications(showDeleted *bool, userID uint, pageStr string, pageSizeStr string, sortBy, sortOrder string) (*r.PaginatedResponse, error) {
	page, pageSize := utils.ParsePagination(pageStr, pageSizeStr)

	notifications, total, err := s.repo.GetAllByUserIDPaginated(showDeleted, userID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, errors.New("failed to get user notifications")
	}

	return utils.BuildPaginatedResponse(notifications, total, page, pageSize), nil
}

func (s *notificationService) MarkNotificationAsRead(notificationID uint, userID uint) error {
	if err := s.repo.MarkAsRead(notificationID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("notification not found with this id")
		}

		return errors.New("failed to mark notification as read")
	}

	return nil
}

func (s *notificationService) MarkAllNotificationsAsRead(userID uint) error {
	if err := s.repo.MarkAllAsRead(userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("notification not found with this id")
		}

		return errors.New("failed to mark all notifications as read")
	}

	return nil
}

func (s *notificationService) GetUnreadNotificationsCount(userID uint) (int64, error) {
	count, err := s.repo.GetUnreadCount(userID)
	if err != nil {
		return 0, errors.New("failed to get unread notifications count")
	}

	return count, nil
}

func (s *notificationService) SendPushNotification(ctx context.Context, deviceToken string, notification *Notification) error {
	if deviceToken == "" {
		return errors.New("device token is required")
	}

	if notification == nil {
		return errors.New("notification is required")
	}

	if err := s.fcmService.SendNotification(deviceToken, notification); err != nil {
		return errors.New("failed to send push notification")
	}

	return nil
}
