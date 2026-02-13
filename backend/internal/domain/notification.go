package domain

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// Notification represents a user notification
type Notification struct {
	BaseModel
	UserID      string         `gorm:"not null;index" json:"user_id"`
	Type        string         `gorm:"size:50;not null" json:"type"` // 'order', 'payment', 'inventory', 'system'
	Title       string         `gorm:"size:200;not null" json:"title"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	Data        datatypes.JSON `gorm:"type:jsonb" json:"data,omitempty"`                  // Additional data
	Channel     string         `gorm:"size:20;not null;default:'in_app'" json:"channel"`  // 'in_app', 'wechat', 'sms', 'email'
	Priority    string         `gorm:"size:10;not null;default:'normal'" json:"priority"` // 'low', 'normal', 'high', 'urgent'
	IsRead      bool           `gorm:"not null;default:false" json:"is_read"`
	IsArchived  bool           `gorm:"not null;default:false" json:"is_archived"`
	ReadAt      *time.Time     `json:"read_at,omitempty"`
	ActionURL   string         `gorm:"size:500" json:"action_url,omitempty"`
	ActionType  string         `gorm:"size:50" json:"action_type,omitempty"` // 'view_order', 'view_voyage', etc.
	SourceID    *string        `json:"source_id,omitempty"`
	SourceType  string         `gorm:"size:50" json:"source_type,omitempty"`
	RetryCount  int            `gorm:"not null;default:0" json:"retry_count"`
	SentAt      *time.Time     `json:"sent_at,omitempty"`
	DeliveredAt *time.Time     `json:"delivered_at,omitempty"`
}

// TableName returns the table name
func (Notification) TableName() string {
	return "notifications"
}

// Notification type constants
const (
	NotificationTypeOrder     = "order"
	NotificationTypePayment   = "payment"
	NotificationTypeInventory = "inventory"
	NotificationTypeSystem    = "system"
	NotificationTypeRefund    = "refund"
	NotificationTypeVoyage    = "voyage"
	NotificationTypePromotion = "promotion"
)

// Notification channel constants
const (
	NotificationChannelInApp  = "in_app"
	NotificationChannelWechat = "wechat"
	NotificationChannelSMS    = "sms"
	NotificationChannelEmail  = "email"
)

// Notification priority constants
const (
	NotificationPriorityLow    = "low"
	NotificationPriorityNormal = "normal"
	NotificationPriorityHigh   = "high"
	NotificationPriorityUrgent = "urgent"
)

// Action type constants
const (
	NotificationActionViewOrder  = "view_order"
	NotificationActionViewVoyage = "view_voyage"
	NotificationActionViewCabin  = "view_cabin"
	NotificationActionViewRefund = "view_refund"
	NotificationActionExternal   = "external_link"
	NotificationActionNone       = "none"
)

// MarkAsRead marks notification as read
func (n *Notification) MarkAsRead() {
	if !n.IsRead {
		n.IsRead = true
		now := time.Now()
		n.ReadAt = &now
	}
}

// MarkAsArchived marks notification as archived
func (n *Notification) MarkAsArchived() {
	n.IsArchived = true
}

// MarkAsDelivered marks external channel notification as delivered
func (n *Notification) MarkAsDelivered() {
	now := time.Now()
	n.DeliveredAt = &now
}

// IncrementRetry increments retry count
func (n *Notification) IncrementRetry() {
	n.RetryCount++
}

// IsDeliverable checks if notification should be sent via external channel
func (n *Notification) IsDeliverable() bool {
	return n.Channel != NotificationChannelInApp && n.DeliveredAt == nil && n.RetryCount < 3
}

// NotificationData holds structured notification data
type NotificationData struct {
	OrderID      *string `json:"order_id,omitempty"`
	OrderNo      string  `json:"order_no,omitempty"`
	VoyageID     *string `json:"voyage_id,omitempty"`
	VoyageName   string  `json:"voyage_name,omitempty"`
	CabinID      *string `json:"cabin_id,omitempty"`
	CabinNumber  string  `json:"cabin_number,omitempty"`
	Amount       float64 `json:"amount,omitempty"`
	Currency     string  `json:"currency,omitempty"`
	RefundID     *string `json:"refund_id,omitempty"`
	RefundAmount float64 `json:"refund_amount,omitempty"`
	OldStatus    string  `json:"old_status,omitempty"`
	NewStatus    string  `json:"new_status,omitempty"`
	Count        int     `json:"count,omitempty"`
	Days         int     `json:"days,omitempty"`
}

// SetData sets notification data
func (n *Notification) SetData(data *NotificationData) error {
	if data == nil {
		n.Data = nil
		return nil
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	n.Data = jsonData
	return nil
}

// GetData retrieves notification data
func (n *Notification) GetData() (*NotificationData, error) {
	if n.Data == nil {
		return nil, nil
	}
	var data NotificationData
	if err := json.Unmarshal(n.Data, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// NotificationSetting represents user notification preferences
type NotificationSetting struct {
	BaseModel
	UserID            string  `gorm:"uniqueIndex;not null" json:"user_id"`
	OrderEnabled      bool    `gorm:"not null;default:true" json:"order_enabled"`
	PaymentEnabled    bool    `gorm:"not null;default:true" json:"payment_enabled"`
	InventoryEnabled  bool    `gorm:"not null;default:true" json:"inventory_enabled"`
	SystemEnabled     bool    `gorm:"not null;default:true" json:"system_enabled"`
	RefundEnabled     bool    `gorm:"not null;default:true" json:"refund_enabled"`
	VoyageEnabled     bool    `gorm:"not null;default:true" json:"voyage_enabled"`
	PromotionEnabled  bool    `gorm:"not null;default:false" json:"promotion_enabled"`
	WechatEnabled     bool    `gorm:"not null;default:true" json:"wechat_enabled"`
	SMSEnabled        bool    `gorm:"not null;default:false" json:"sms_enabled"`
	EmailEnabled      bool    `gorm:"not null;default:false" json:"email_enabled"`
	QuietHoursStart   *string `json:"quiet_hours_start,omitempty"` // "22:00"
	QuietHoursEnd     *string `json:"quiet_hours_end,omitempty"`   // "08:00"
	QuietHoursEnabled bool    `gorm:"not null;default:true" json:"quiet_hours_enabled"`
}

// TableName returns the table name
func (NotificationSetting) TableName() string {
	return "notification_settings"
}

// IsChannelEnabled checks if a channel is enabled for this user
func (s *NotificationSetting) IsChannelEnabled(channel string) bool {
	switch channel {
	case NotificationChannelWechat:
		return s.WechatEnabled
	case NotificationChannelSMS:
		return s.SMSEnabled
	case NotificationChannelEmail:
		return s.EmailEnabled
	default:
		return true // in_app is always enabled
	}
}

// IsTypeEnabled checks if a notification type is enabled
func (s *NotificationSetting) IsTypeEnabled(nType string) bool {
	switch nType {
	case NotificationTypeOrder:
		return s.OrderEnabled
	case NotificationTypePayment:
		return s.PaymentEnabled
	case NotificationTypeInventory:
		return s.InventoryEnabled
	case NotificationTypeSystem:
		return s.SystemEnabled
	case NotificationTypeRefund:
		return s.RefundEnabled
	case NotificationTypeVoyage:
		return s.VoyageEnabled
	case NotificationTypePromotion:
		return s.PromotionEnabled
	default:
		return true
	}
}

// IsInQuietHours checks if current time is in quiet hours
func (s *NotificationSetting) IsInQuietHours() bool {
	if !s.QuietHoursEnabled || s.QuietHoursStart == nil || s.QuietHoursEnd == nil {
		return false
	}

	now := time.Now()
	start, _ := time.Parse("15:04", *s.QuietHoursStart)
	end, _ := time.Parse("15:04", *s.QuietHoursEnd)

	currentMinutes := now.Hour()*60 + now.Minute()
	startMinutes := start.Hour()*60 + start.Minute()
	endMinutes := end.Hour()*60 + end.Minute()

	if startMinutes < endMinutes {
		return currentMinutes >= startMinutes && currentMinutes < endMinutes
	}
	// Quiet hours span midnight (e.g., 22:00 - 08:00)
	return currentMinutes >= startMinutes || currentMinutes < endMinutes
}
