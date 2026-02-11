package service

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidNotification  = errors.New("invalid notification data")
)

// NotificationSender defines interface for external notification channels
type NotificationSender interface {
	// Send sends a notification via the channel
	Send(ctx context.Context, notification *domain.Notification, user *domain.User) error
	// IsAvailable checks if the sender is properly configured
	IsAvailable() bool
}

// NotificationService defines the interface for notification business logic
type NotificationService interface {
	// Create creates a new notification
	Create(ctx context.Context, req CreateNotificationRequest) (*domain.Notification, error)

	// CreateAndSend creates a notification and sends via appropriate channels
	CreateAndSend(ctx context.Context, req CreateNotificationRequest) (*domain.Notification, error)

	// GetByID retrieves a notification by ID
	GetByID(ctx context.Context, id uint64) (*domain.Notification, error)

	// List retrieves notifications for a user
	List(ctx context.Context, userID uint64, paginator *pagination.Paginator, unreadOnly bool) (*pagination.Result, error)

	// GetUnreadCount gets the count of unread notifications
	GetUnreadCount(ctx context.Context, userID uint64) (int64, error)

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id uint64) error

	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(ctx context.Context, userID uint64) error

	// Archive archives a notification
	Archive(ctx context.Context, id uint64) error

	// Delete deletes a notification
	Delete(ctx context.Context, id uint64) error

	// GetSettings retrieves notification settings for a user
	GetSettings(ctx context.Context, userID uint64) (*domain.NotificationSetting, error)

	// UpdateSettings updates notification settings
	UpdateSettings(ctx context.Context, userID uint64, req UpdateSettingsRequest) (*domain.NotificationSetting, error)

	// SendOrderCreatedNotification sends notification when order is created
	SendOrderCreatedNotification(ctx context.Context, userID uint64, order *domain.Order) error

	// SendOrderPaidNotification sends notification when order is paid
	SendOrderPaidNotification(ctx context.Context, userID uint64, order *domain.Order) error

	// SendOrderConfirmedNotification sends notification when order is confirmed
	SendOrderConfirmedNotification(ctx context.Context, userID uint64, order *domain.Order) error

	// SendRefundApprovedNotification sends notification when refund is approved
	SendRefundApprovedNotification(ctx context.Context, userID uint64, refund *domain.RefundRequest) error

	// SendInventoryAlertNotification sends inventory alert to admins
	SendInventoryAlertNotification(ctx context.Context, voyageID uint64, cabinTypeID uint64, remaining int) error
}

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
	UserID       uint64                      `json:"user_id" validate:"required"
	Type         string                      `json:"type" validate:"required,oneof=order payment inventory system refund voyage promotion"
	Title        string                      `json:"title" validate:"required,max=200"
	Content      string                      `json:"content" validate:"required"
	Data         *domain.NotificationData    `json:"data,omitempty"`
	Priority     string                      `json:"priority,omitempty" validate:"omitempty,oneof=low normal high urgent"`
	Channel      string                      `json:"channel,omitempty" validate:"omitempty,oneof=in_app wechat sms email"`
	ActionURL    string                      `json:"action_url,omitempty"`
	ActionType   string                      `json:"action_type,omitempty"`
	SourceID     *uint64                     `json:"source_id,omitempty"`
	SourceType   string                      `json:"source_type,omitempty"`
}

// UpdateSettingsRequest represents a request to update notification settings
type UpdateSettingsRequest struct {
	OrderEnabled      *bool   `json:"order_enabled,omitempty"`
	PaymentEnabled    *bool   `json:"payment_enabled,omitempty"`
	InventoryEnabled  *bool   `json:"inventory_enabled,omitempty"`
	SystemEnabled     *bool   `json:"system_enabled,omitempty"`
	RefundEnabled     *bool   `json:"refund_enabled,omitempty"`
	VoyageEnabled     *bool   `json:"voyage_enabled,omitempty"`
	PromotionEnabled  *bool   `json:"promotion_enabled,omitempty"`
	WechatEnabled     *bool   `json:"wechat_enabled,omitempty"`
	SMSEnabled        *bool   `json:"sms_enabled,omitempty"`
	EmailEnabled      *bool   `json:"email_enabled,omitempty"`
	QuietHoursStart   *string `json:"quiet_hours_start,omitempty" validate:"omitempty,datetime=15:04"`
	QuietHoursEnd     *string `json:"quiet_hours_end,omitempty" validate:"omitempty,datetime=15:04"`
	QuietHoursEnabled *bool   `json:"quiet_hours_enabled,omitempty"`
}

// notificationService implements NotificationService
type notificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	wechatSender     NotificationSender
	smsSender        NotificationSender
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	wechatSender NotificationSender,
	smsSender NotificationSender,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		wechatSender:     wechatSender,
		smsSender:        smsSender,
	}
}

// Create creates a new notification
func (s *notificationService) Create(ctx context.Context, req CreateNotificationRequest) (*domain.Notification, error) {
	// Validate user exists
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get user settings to determine channels
	settings, err := s.notificationRepo.GetOrCreateSettings(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification settings: %w", err)
	}

	// Check if notification type is enabled
	if !settings.IsTypeEnabled(req.Type) {
		return nil, nil // Silently skip disabled notification types
	}

	// Set defaults
	priority := domain.NotificationPriorityNormal
	if req.Priority != "" {
		priority = req.Priority
	}

	channel := domain.NotificationChannelInApp
	if req.Channel != "" {
		channel = req.Channel
	}

	notification := &domain.Notification{
		UserID:     req.UserID,
		Type:       req.Type,
		Title:      req.Title,
		Content:    req.Content,
		Priority:   priority,
		Channel:    channel,
		ActionURL:  req.ActionURL,
		ActionType: req.ActionType,
		SourceID:   req.SourceID,
		SourceType: req.SourceType,
	}

	if err := notification.SetData(req.Data); err != nil {
		return nil, fmt.Errorf("failed to set notification data: %w", err)
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Send to external channels if configured and available
	if notification.IsDeliverable() {
		if err := s.sendExternalNotification(ctx, notification, user, settings); err != nil {
			// Log error but don't fail - notification is already saved in-app
			// Retry will be handled by background job
			notification.IncrementRetry()
			s.notificationRepo.Update(ctx, notification)
		}
	}

	return notification, nil
}

// CreateAndSend creates a notification and sends via appropriate channels
func (s *notificationService) CreateAndSend(ctx context.Context, req CreateNotificationRequest) (*domain.Notification, error) {
	return s.Create(ctx, req)
}

// GetByID retrieves a notification by ID
func (s *notificationService) GetByID(ctx context.Context, id uint64) (*domain.Notification, error) {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrNotificationNotFound
		}
		return nil, err
	}
	return notification, nil
}

// List retrieves notifications for a user
func (s *notificationService) List(ctx context.Context, userID uint64, paginator *pagination.Paginator, unreadOnly bool) (*pagination.Result, error) {
	return s.notificationRepo.List(ctx, userID, paginator, unreadOnly)
}

// GetUnreadCount gets the count of unread notifications
func (s *notificationService) GetUnreadCount(ctx context.Context, userID uint64) (int64, error) {
	return s.notificationRepo.GetUnreadCount(ctx, userID)
}

// MarkAsRead marks a notification as read
func (s *notificationService) MarkAsRead(ctx context.Context, id uint64) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return ErrNotificationNotFound
		}
		return err
	}

	notification.MarkAsRead()
	return s.notificationRepo.Update(ctx, notification)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *notificationService) MarkAllAsRead(ctx context.Context, userID uint64) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

// Archive archives a notification
func (s *notificationService) Archive(ctx context.Context, id uint64) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return ErrNotificationNotFound
		}
		return err
	}

	notification.MarkAsArchived()
	return s.notificationRepo.Update(ctx, notification)
}

// Delete deletes a notification
func (s *notificationService) Delete(ctx context.Context, id uint64) error {
	return s.notificationRepo.Delete(ctx, id)
}

// GetSettings retrieves notification settings for a user
func (s *notificationService) GetSettings(ctx context.Context, userID uint64) (*domain.NotificationSetting, error) {
	return s.notificationRepo.GetOrCreateSettings(ctx, userID)
}

// UpdateSettings updates notification settings
func (s *notificationService) UpdateSettings(ctx context.Context, userID uint64, req UpdateSettingsRequest) (*domain.NotificationSetting, error) {
	settings, err := s.notificationRepo.GetOrCreateSettings(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.OrderEnabled != nil {
		settings.OrderEnabled = *req.OrderEnabled
	}
	if req.PaymentEnabled != nil {
		settings.PaymentEnabled = *req.PaymentEnabled
	}
	if req.InventoryEnabled != nil {
		settings.InventoryEnabled = *req.InventoryEnabled
	}
	if req.SystemEnabled != nil {
		settings.SystemEnabled = *req.SystemEnabled
	}
	if req.RefundEnabled != nil {
		settings.RefundEnabled = *req.RefundEnabled
	}
	if req.VoyageEnabled != nil {
		settings.VoyageEnabled = *req.VoyageEnabled
	}
	if req.PromotionEnabled != nil {
		settings.PromotionEnabled = *req.PromotionEnabled
	}
	if req.WechatEnabled != nil {
		settings.WechatEnabled = *req.WechatEnabled
	}
	if req.SMSEnabled != nil {
		settings.SMSEnabled = *req.SMSEnabled
	}
	if req.EmailEnabled != nil {
		settings.EmailEnabled = *req.EmailEnabled
	}
	if req.QuietHoursStart != nil {
		settings.QuietHoursStart = req.QuietHoursStart
	}
	if req.QuietHoursEnd != nil {
		settings.QuietHoursEnd = req.QuietHoursEnd
	}
	if req.QuietHoursEnabled != nil {
		settings.QuietHoursEnabled = *req.QuietHoursEnabled
	}

	if err := s.notificationRepo.UpdateSettings(ctx, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// SendOrderCreatedNotification sends notification when order is created
func (s *notificationService) SendOrderCreatedNotification(ctx context.Context, userID uint64, order *domain.Order) error {
	req := CreateNotificationRequest{
		UserID:  userID,
		Type:    domain.NotificationTypeOrder,
		Title:   "订单已创建",
		Content: fmt.Sprintf("您的订单 %s 已创建，请在30分钟内完成支付。", order.OrderNumber),
		Data: &domain.NotificationData{
			OrderID:  &order.ID,
			OrderNo:  order.OrderNumber,
			Amount:   order.TotalAmount,
			Currency: order.Currency,
		},
		Priority:   domain.NotificationPriorityHigh,
		ActionType: domain.NotificationActionViewOrder,
		SourceID:   &order.ID,
		SourceType: domain.NotificationTypeOrder,
	}

	_, err := s.Create(ctx, req)
	return err
}

// SendOrderPaidNotification sends notification when order is paid
func (s *notificationService) SendOrderPaidNotification(ctx context.Context, userID uint64, order *domain.Order) error {
	req := CreateNotificationRequest{
		UserID:  userID,
		Type:    domain.NotificationTypePayment,
		Title:   "订单支付成功",
		Content: fmt.Sprintf("您的订单 %s 支付成功，金额: %.2f %s", order.OrderNumber, order.TotalAmount, order.Currency),
		Data: &domain.NotificationData{
			OrderID:  &order.ID,
			OrderNo:  order.OrderNumber,
			Amount:   order.TotalAmount,
			Currency: order.Currency,
		},
		Priority:   domain.NotificationPriorityNormal,
		ActionType: domain.NotificationActionViewOrder,
		SourceID:   &order.ID,
		SourceType: domain.NotificationTypeOrder,
	}

	_, err := s.Create(ctx, req)
	return err
}

// SendOrderConfirmedNotification sends notification when order is confirmed
func (s *notificationService) SendOrderConfirmedNotification(ctx context.Context, userID uint64, order *domain.Order) error {
	req := CreateNotificationRequest{
		UserID:  userID,
		Type:    domain.NotificationTypeOrder,
		Title:   "订单已确认",
		Content: fmt.Sprintf("您的订单 %s 已确认，祝您旅途愉快！", order.OrderNumber),
		Data: &domain.NotificationData{
			OrderID:  &order.ID,
			OrderNo:  order.OrderNumber,
		},
		Priority:   domain.NotificationPriorityNormal,
		ActionType: domain.NotificationActionViewOrder,
		SourceID:   &order.ID,
		SourceType: domain.NotificationTypeOrder,
	}

	_, err := s.Create(ctx, req)
	return err
}

// SendRefundApprovedNotification sends notification when refund is approved
func (s *notificationService) SendRefundApprovedNotification(ctx context.Context, userID uint64, refund *domain.RefundRequest) error {
	req := CreateNotificationRequest{
		UserID:  userID,
		Type:    domain.NotificationTypeRefund,
		Title:   "退款已批准",
		Content: fmt.Sprintf("您的退款申请已批准，退款金额: %.2f，将在3-7个工作日内到账。", refund.Amount),
		Data: &domain.NotificationData{
			RefundID:     &refund.ID,
			RefundAmount: refund.Amount,
		},
		Priority:   domain.NotificationPriorityHigh,
		ActionType: domain.NotificationActionViewRefund,
		SourceID:   &refund.ID,
		SourceType: domain.NotificationTypeRefund,
	}

	_, err := s.Create(ctx, req)
	return err
}

// SendInventoryAlertNotification sends inventory alert to admins
func (s *notificationService) SendInventoryAlertNotification(ctx context.Context, voyageID uint64, cabinTypeID uint64, remaining int) error {
	// This would typically send to admin users
	// For now, create a system notification
	req := CreateNotificationRequest{
		UserID:  1, // Admin user ID
		Type:    domain.NotificationTypeInventory,
		Title:   "库存预警",
		Content: fmt.Sprintf("航次 %d 的房型 %d 仅剩 %d 间，请及时补充库存。", voyageID, cabinTypeID, remaining),
		Data: &domain.NotificationData{
			Count: remaining,
		},
		Priority:   domain.NotificationPriorityUrgent,
		SourceID:   &voyageID,
		SourceType: domain.NotificationTypeVoyage,
	}

	_, err := s.Create(ctx, req)
	return err
}

// sendExternalNotification sends notification via external channels
func (s *notificationService) sendExternalNotification(ctx context.Context, notification *domain.Notification, user *domain.User, settings *domain.NotificationSetting) error {
	// Check quiet hours
	if settings.IsInQuietHours() && notification.Priority != domain.NotificationPriorityUrgent {
		return nil // Skip non-urgent notifications during quiet hours
	}

	// Send via appropriate channel
	switch notification.Channel {
	case domain.NotificationChannelWechat:
		if s.wechatSender != nil && s.wechatSender.IsAvailable() && settings.WechatEnabled {
			if err := s.wechatSender.Send(ctx, notification, user); err != nil {
				return fmt.Errorf("wechat send failed: %w", err)
			}
			notification.MarkAsDelivered()
		}
	case domain.NotificationChannelSMS:
		if s.smsSender != nil && s.smsSender.IsAvailable() && settings.SMSEnabled {
			if err := s.smsSender.Send(ctx, notification, user); err != nil {
				return fmt.Errorf("sms send failed: %w", err)
			}
			notification.MarkAsDelivered()
		}
	}

	return nil
}
