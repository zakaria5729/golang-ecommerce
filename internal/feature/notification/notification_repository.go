package notification

import (
	"context"
	"errors"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/pkg/logger"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		db: db.GetDB(),
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *Notification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		logger.Logger.Error("Failed to create notification", "method", "Create", "error", err, "userID", notification.UserID)
		return err
	}
	return nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uint, showDeleted *bool) (*Notification, error) {
	var notification Notification
	query := r.db.WithContext(ctx).Model(&Notification{})

	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	err := query.First(&notification, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Logger.Error("Failed to get notification by ID", "method", "GetByID", "error", err, "id", id)
		return nil, err
	}

	return &notification, nil
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uint, page, limit int, showDeleted *bool) ([]*Notification, error) {
	var notifications []*Notification
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if showDeleted != nil && !*showDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error

	if err != nil {
		logger.Logger.Error("Failed to get user notifications", "method", "GetByUserID", "error", err, "userID", userID)
		return nil, err
	}

	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID, userID uint) error {
	result := r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)

	if result.Error != nil {
		logger.Logger.Error("Failed to mark notification as read", "method", "MarkAsRead", "error", result.Error, "notificationID", notificationID, "userID", userID)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	result := r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	if result.Error != nil {
		logger.Logger.Error("Failed to mark all notifications as read", "method", "MarkAllAsRead", "error", result.Error, "userID", userID)
		return result.Error
	}

	return nil
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error

	if err != nil {
		logger.Logger.Error("Failed to get unread notifications count", "method", "GetUnreadCount", "error", err, "userID", userID)
		return 0, err
	}

	return count, nil
}
