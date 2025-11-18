package notification

import (
	"context"
	"errors"

	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, id uint, showDeleted *bool) (*Notification, error)
	GetAllByUserIDPaginated(showDeleted *bool, userID uint, page, pageSize int, sortBy, sortOrder string) ([]Notification, int64, error)
	MarkAsRead(notificationID uint, userID uint) error
	MarkAllAsRead(userID uint) error
	GetUnreadCount(userID uint) (int64, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

func (r *notificationRepository) Create(ctx context.Context, notification *Notification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		l.Error("Failed to create notification", "method", "Create", "error", err, "userID", notification.UserID)
		return err
	}
	return nil
}

func (r *notificationRepository) GetByID(ctx context.Context, id uint, showDeleted *bool) (*Notification, error) {
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
		l.Error("Failed to get notification by ID", "method", "GetByID", "error", err, "id", id)
		return nil, err
	}

	return &notification, nil
}

func (r *notificationRepository) GetAllByUserIDPaginated(showDeleted *bool, userID uint, page, pageSize int, sortBy, sortOrder string) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64

	query := r.db.Model(&Notification{}).Where(c.NotificationUserID+" = ?", userID)
	if showDeleted != nil && *showDeleted {
		query = query.Unscoped()
	}

	if orderClause := utils.BuildSortingOrder(sortBy, sortOrder, nil); orderClause != "" {
		query = query.Order(orderClause)
	}

	if err := query.Count(&total).Error; err != nil {
		l.Error("Failed to count categories", "method", "GetAllCategoriesPaginated", "error", err, "userID", userID, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
		return nil, 0, err
	}

	err := query.Offset(utils.GetOffset(page, pageSize)).Limit(pageSize).Find(&notifications).Error
	if err != nil {
		l.Error("Failed to fetch categories paginated", "method", "GetAllCategoriesPaginated", "error", err, "userID", userID, "page", page, "pageSize", pageSize, "sortBy", sortBy, "sortOrder", sortOrder)
	}

	return notifications, total, err
}

func (r *notificationRepository) MarkAsRead(notificationID uint, userID uint) error {
	notification := &Notification{
		IsRead: true,
	}
	notification.UpdatedBy = &userID

	result := r.db.Model(&Notification{}).
		Where(c.FieldID+" = ? AND "+c.NotificationUserID+" = ?", notificationID, userID).
		Updates(notification)

	if result.Error != nil {
		l.Error("Failed to mark notification as read", "method", "MarkAsRead", "error", result.Error, "notificationID", notificationID, "userID", userID)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *notificationRepository) MarkAllAsRead(userID uint) error {
	notification := &Notification{
		IsRead: true,
	}
	notification.UpdatedBy = &userID

	result := r.db.Model(&Notification{}).
		Where(c.NotificationUserID+" = ? AND "+c.NotificationIsRead+" = ?", userID, false).
		Updates(notification)

	if result.Error != nil {
		l.Error("Failed to mark all notifications as read", "method", "MarkAllAsRead", "error", result.Error, "userID", userID)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *notificationRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64

	err := r.db.Model(&Notification{}).
		Where(c.NotificationUserID+" = ? AND "+c.NotificationIsRead+" = ?", userID, false).
		Count(&count).Error

	if err != nil {
		l.Error("Failed to get unread notifications count", "method", "GetUnreadCount", "error", err, "userID", userID)
		return 0, err
	}

	return count, nil
}
