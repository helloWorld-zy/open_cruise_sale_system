package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"
	"time"

	"gorm.io/gorm"
)

// NotificationRepository defines the interface for notification data operations
type NotificationRepository interface {
	// Notification CRUD
	Create(ctx context.Context, notification *domain.Notification) error
	GetByID(ctx context.Context, id uint64) (*domain.Notification, error)
	List(ctx context.Context, userID uint64, paginator *pagination.Paginator, unreadOnly bool) (*pagination.Result, error)
	GetUnreadCount(ctx context.Context, userID uint64) (int64, error)
	Update(ctx context.Context, notification *domain.Notification) error
	Delete(ctx context.Context, id uint64) error

	// Batch operations
	MarkAllAsRead(ctx context.Context, userID uint64) error
	DeleteAllRead(ctx context.Context, userID uint64) error
	DeleteOldNotifications(ctx context.Context, days int) error

	// Settings
	GetOrCreateSettings(ctx context.Context, userID uint64) (*domain.NotificationSetting, error)
	UpdateSettings(ctx context.Context, settings *domain.NotificationSetting) error

	// For retry jobs
	GetUndeliveredNotifications(ctx context.Context, limit int) ([]*domain.Notification, error)
}

// notificationRepository implements NotificationRepository
type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

// Create creates a new notification
func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// GetByID retrieves a notification by ID
func (r *notificationRepository) GetByID(ctx context.Context, id uint64) (*domain.Notification, error) {
	var notification domain.Notification
	err := r.db.WithContext(ctx).First(&notification, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &notification, nil
}

// List retrieves notifications for a user
func (r *notificationRepository) List(ctx context.Context, userID uint64, paginator *pagination.Paginator, unreadOnly bool) (*pagination.Result, error) {
	query := r.db.WithContext(ctx).Model(&domain.Notification{}).Where("user_id = ? AND is_archived = ?", userID, false)

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	// Order by priority and created_at
	query = query.Order("CASE priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'normal' THEN 3 ELSE 4 END, created_at DESC")

	var notifications []*domain.Notification
	var total int64

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset(paginator.Offset()).Limit(paginator.Limit()).Find(&notifications).Error; err != nil {
		return nil, err
	}

	return &pagination.Result{
		Data:       notifications,
		Total:      total,
		Page:       paginator.Page,
		PageSize:   paginator.PageSize,
		TotalPages: paginator.TotalPages(total),
	}, nil
}

// GetUnreadCount gets the count of unread notifications
func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ? AND is_archived = ?", userID, false, false).
		Count(&count).Error
	return count, err
}

// Update updates a notification
func (r *notificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

// Delete deletes a notification
func (r *notificationRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Notification{}, id).Error
}

// MarkAllAsRead marks all notifications as read for a user
func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// DeleteAllRead deletes all read notifications for a user
func (r *notificationRepository) DeleteAllRead(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND is_read = ?", userID, true).
		Delete(&domain.Notification{}).Error
}

// DeleteOldNotifications deletes notifications older than specified days
func (r *notificationRepository) DeleteOldNotifications(ctx context.Context, days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.WithContext(ctx).
		Where("created_at < ? AND is_read = ?", cutoff, true).
		Delete(&domain.Notification{}).Error
}

// GetOrCreateSettings gets or creates notification settings for a user
func (r *notificationRepository) GetOrCreateSettings(ctx context.Context, userID uint64) (*domain.NotificationSetting, error) {
	var settings domain.NotificationSetting
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create default settings
			settings = domain.NotificationSetting{
				UserID:            userID,
				OrderEnabled:      true,
				PaymentEnabled:    true,
				InventoryEnabled:  true,
				SystemEnabled:     true,
				RefundEnabled:     true,
				VoyageEnabled:     true,
				PromotionEnabled:  false,
				WechatEnabled:     true,
				SMSEnabled:        false,
				EmailEnabled:      false,
				QuietHoursEnabled: true,
			}
			quietStart := "22:00"
			quietEnd := "08:00"
			settings.QuietHoursStart = &quietStart
			settings.QuietHoursEnd = &quietEnd

			if err := r.db.WithContext(ctx).Create(&settings).Error; err != nil {
				return nil, err
			}
			return &settings, nil
		}
		return nil, err
	}
	return &settings, nil
}

// UpdateSettings updates notification settings
func (r *notificationRepository) UpdateSettings(ctx context.Context, settings *domain.NotificationSetting) error {
	return r.db.WithContext(ctx).Save(settings).Error
}

// GetUndeliveredNotifications gets notifications that need to be retried
func (r *notificationRepository) GetUndeliveredNotifications(ctx context.Context, limit int) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	err := r.db.WithContext(ctx).
		Where("channel != ? AND delivered_at IS NULL AND retry_count < ?", domain.NotificationChannelInApp, 3).
		Order("created_at ASC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}
