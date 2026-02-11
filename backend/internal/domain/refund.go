package domain

import "time"

// RefundRequest represents a refund application from customer
type RefundRequest struct {
	BaseModel
	OrderID            string     `gorm:"not null;index" json:"order_id"`
	Order              Order      `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	OrderItemID        *string    `gorm:"index" json:"order_item_id,omitempty"`
	OrderItem          OrderItem  `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
	UserID             *string    `gorm:"index" json:"user_id,omitempty"`
	RefundAmount       float64    `gorm:"not null" json:"refund_amount"`
	RefundReason       string     `gorm:"not null" json:"refund_reason"`
	RefundType         string     `gorm:"default:partial" json:"refund_type"`
	Status             string     `gorm:"default:pending" json:"status"`
	RequestedAt        time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"requested_at"`
	ReviewedAt         *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy         *string    `json:"reviewed_by,omitempty"`
	ReviewNote         string     `json:"review_note,omitempty"`
	ProcessedAt        *time.Time `json:"processed_at,omitempty"`
	PaymentRefundID    string     `json:"payment_refund_id,omitempty"`
	ThirdPartyRefundID string     `json:"third_party_refund_id,omitempty"`
	RefundMethod       string     `gorm:"default:original" json:"refund_method"`
	BankName           string     `json:"bank_name,omitempty"`
	BankAccount        string     `json:"bank_account,omitempty"`
	AccountHolder      string     `json:"account_holder,omitempty"`
	CancellationReason string     `json:"cancellation_reason,omitempty"`
}

// TableName returns the table name for RefundRequest
func (RefundRequest) TableName() string {
	return "refund_requests"
}

// RefundStatus constants
const (
	RefundStatusPending    = "pending"
	RefundStatusApproved   = "approved"
	RefundStatusRejected   = "rejected"
	RefundStatusProcessing = "processing"
	RefundStatusCompleted  = "completed"
	RefundStatusFailed     = "failed"
)

// RefundType constants
const (
	RefundTypeFull    = "full"
	RefundTypePartial = "partial"
)

// RefundMethod constants
const (
	RefundMethodOriginal = "original"
	RefundMethodBank     = "bank"
	RefundMethodAlipay   = "alipay"
	RefundMethodWechat   = "wechat"
)

// CancellationReason constants
const (
	CancellationReasonCustomerRequest = "customer_request"
	CancellationReasonVoyageCancelled = "voyage_cancelled"
	CancellationReasonCabinUpgrade    = "cabin_upgrade"
	CancellationReasonOther           = "other"
)

// CanReview checks if refund request can be reviewed
func (r *RefundRequest) CanReview() bool {
	return r.Status == RefundStatusPending
}

// CanProcess checks if refund request can be processed
func (r *RefundRequest) CanProcess() bool {
	return r.Status == RefundStatusApproved
}

// IsCompleted checks if refund is completed
func (r *RefundRequest) IsCompleted() bool {
	return r.Status == RefundStatusCompleted
}

// IsRejected checks if refund is rejected
func (r *RefundRequest) IsRejected() bool {
	return r.Status == RefundStatusRejected
}

// Approve approves the refund request
func (r *RefundRequest) Approve(reviewerID, note string) {
	r.Status = RefundStatusApproved
	r.ReviewedBy = &reviewerID
	r.ReviewNote = note
	now := time.Now()
	r.ReviewedAt = &now
	r.UpdatedAt = now
}

// Reject rejects the refund request
func (r *RefundRequest) Reject(reviewerID, note string) {
	r.Status = RefundStatusRejected
	r.ReviewedBy = &reviewerID
	r.ReviewNote = note
	now := time.Now()
	r.ReviewedAt = &now
	r.UpdatedAt = now
}

// MarkProcessing marks refund as processing
func (r *RefundRequest) MarkProcessing() {
	r.Status = RefundStatusProcessing
	now := time.Now()
	r.ProcessedAt = &now
	r.UpdatedAt = now
}

// MarkCompleted marks refund as completed
func (r *RefundRequest) MarkCompleted(paymentRefundID, thirdPartyRefundID string) {
	r.Status = RefundStatusCompleted
	r.PaymentRefundID = paymentRefundID
	r.ThirdPartyRefundID = thirdPartyRefundID
	r.UpdatedAt = time.Now()
}

// MarkFailed marks refund as failed
func (r *RefundRequest) MarkFailed() {
	r.Status = RefundStatusFailed
	r.UpdatedAt = time.Now()
}
