package domain

import "time"

// Order represents a booking order
type Order struct {
	BaseModel
	OrderNumber    string  `gorm:"not null;uniqueIndex" json:"order_number"`
	UserID         *string `gorm:"index" json:"user_id,omitempty"`
	VoyageID       string  `gorm:"not null;index" json:"voyage_id"`
	Voyage         Voyage  `gorm:"foreignKey:VoyageID" json:"voyage,omitempty"`
	CruiseID       string  `gorm:"not null;index" json:"cruise_id"`
	Cruise         Cruise  `gorm:"foreignKey:CruiseID" json:"cruise,omitempty"`
	TotalAmount    float64 `gorm:"not null;default:0" json:"total_amount"`
	DiscountAmount float64 `gorm:"default:0" json:"discount_amount"`
	PaidAmount     float64 `gorm:"default:0" json:"paid_amount"`
	Currency       string  `gorm:"default:CNY" json:"currency"`
	Status         string  `gorm:"default:pending" json:"status"`
	PaymentStatus  string  `gorm:"default:unpaid" json:"payment_status"`
	PassengerCount int     `gorm:"not null;default:1" json:"passenger_count"`
	CabinCount     int     `gorm:"not null;default:1" json:"cabin_count"`
	ContactName    string  `json:"contact_name,omitempty"`
	ContactPhone   string  `json:"contact_phone,omitempty"`
	ContactEmail   string  `json:"contact_email,omitempty"`
	Remark         string  `json:"remark,omitempty"`
	BookedAt       *string `json:"booked_at,omitempty"`
	PaidAt         *string `json:"paid_at,omitempty"`
	ConfirmedAt    *string `json:"confirmed_at,omitempty"`
	ExpiresAt      string  `gorm:"not null" json:"expires_at"`

	// Relations
	Items      []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Passengers []Passenger `gorm:"foreignKey:OrderID" json:"passengers,omitempty"`
	Payments   []Payment   `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}

// TableName returns the table name for Order
func (Order) TableName() string {
	return "orders"
}

// OrderStatus constants
const (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusConfirmed = "confirmed"
	OrderStatusCancelled = "cancelled"
	OrderStatusCompleted = "completed"
	OrderStatusRefunded  = "refunded"
)

// PaymentStatus constants
const (
	PaymentStatusUnpaid   = "unpaid"
	PaymentStatusPartial  = "partial"
	PaymentStatusPaid     = "paid"
	PaymentStatusRefunded = "refunded"
)

// IsPaid checks if order is fully paid
func (o *Order) IsPaid() bool {
	return o.PaymentStatus == PaymentStatusPaid
}

// CanCancel checks if order can be cancelled
func (o *Order) CanCancel() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusPaid
}

// IsExpired checks if order has expired
func (o *Order) IsExpired() bool {
	if o.ExpiresAt == "" {
		return false
	}
	// Simple string comparison works for RFC3339 format
	return o.ExpiresAt < time.Now().Format(time.RFC3339)
}

// OrderItem represents a line item in an order
type OrderItem struct {
	BaseModel
	OrderID       string    `gorm:"not null;index" json:"order_id"`
	Order         Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	CabinID       string    `gorm:"not null;index" json:"cabin_id"`
	Cabin         Cabin     `gorm:"foreignKey:CabinID" json:"cabin,omitempty"`
	CabinTypeID   string    `gorm:"not null;index" json:"cabin_type_id"`
	CabinType     CabinType `gorm:"foreignKey:CabinTypeID" json:"cabin_type,omitempty"`
	VoyageID      string    `gorm:"not null;index" json:"voyage_id"`
	Voyage        Voyage    `gorm:"foreignKey:VoyageID" json:"voyage,omitempty"`
	CabinNumber   string    `json:"cabin_number,omitempty"`
	PriceSnapshot float64   `gorm:"not null" json:"price_snapshot"`
	AdultCount    int       `gorm:"not null;default:2" json:"adult_count"`
	ChildCount    int       `gorm:"default:0" json:"child_count"`
	InfantCount   int       `gorm:"default:0" json:"infant_count"`
	AdultPrice    float64   `gorm:"not null" json:"adult_price"`
	ChildPrice    float64   `json:"child_price,omitempty"`
	InfantPrice   float64   `json:"infant_price,omitempty"`
	PortFee       float64   `gorm:"default:0" json:"port_fee"`
	ServiceFee    float64   `gorm:"default:0" json:"service_fee"`
	Subtotal      float64   `gorm:"not null" json:"subtotal"`
	Status        string    `gorm:"default:confirmed" json:"status"`

	// Relations
	Passengers []Passenger `gorm:"foreignKey:OrderItemID" json:"passengers,omitempty"`
}

// TableName returns the table name for OrderItem
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderItemStatus constants
const (
	OrderItemStatusConfirmed = "confirmed"
	OrderItemStatusCancelled = "cancelled"
	OrderItemStatusChanged   = "changed"
)

// CalculateSubtotal recalculates subtotal based on current values
func (oi *OrderItem) CalculateSubtotal() float64 {
	subtotal := oi.AdultPrice * float64(oi.AdultCount)
	if oi.ChildPrice > 0 {
		subtotal += oi.ChildPrice * float64(oi.ChildCount)
	}
	if oi.InfantPrice > 0 {
		subtotal += oi.InfantPrice * float64(oi.InfantCount)
	}
	subtotal += oi.PortFee * float64(oi.AdultCount+oi.ChildCount)
	subtotal += oi.ServiceFee * float64(oi.AdultCount+oi.ChildCount)
	return subtotal
}

// Passenger represents a passenger in an order
type Passenger struct {
	BaseModel
	OrderID               string    `gorm:"not null;index" json:"order_id"`
	Order                 Order     `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	OrderItemID           string    `gorm:"not null;index" json:"order_item_id"`
	OrderItem             OrderItem `gorm:"foreignKey:OrderItemID" json:"order_item,omitempty"`
	Name                  string    `gorm:"not null" json:"name"`
	Surname               string    `gorm:"not null" json:"surname"`
	GivenName             string    `json:"given_name,omitempty"`
	Gender                string    `gorm:"not null" json:"gender"`
	BirthDate             string    `gorm:"not null" json:"birth_date"`
	Nationality           string    `json:"nationality,omitempty"`
	PassportNumber        string    `json:"passport_number,omitempty"`
	PassportExpiry        string    `json:"passport_expiry,omitempty"`
	IDNumber              string    `json:"id_number,omitempty"`
	Phone                 string    `json:"phone,omitempty"`
	Email                 string    `json:"email,omitempty"`
	PassengerType         string    `gorm:"not null;default:adult" json:"passenger_type"`
	EmergencyContactName  string    `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone string    `json:"emergency_contact_phone,omitempty"`
	DietaryRequirements   string    `json:"dietary_requirements,omitempty"`
	MedicalNotes          string    `json:"medical_notes,omitempty"`
}

// TableName returns the table name for Passenger
func (Passenger) TableName() string {
	return "passengers"
}

// PassengerType constants
const (
	PassengerTypeAdult  = "adult"
	PassengerTypeChild  = "child"
	PassengerTypeInfant = "infant"
)

// Payment represents a payment record
type Payment struct {
	BaseModel
	OrderID                 string  `gorm:"not null;index" json:"order_id"`
	Order                   Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	PaymentNo               string  `gorm:"not null;uniqueIndex" json:"payment_no"`
	PaymentMethod           string  `gorm:"not null" json:"payment_method"`
	PaymentChannel          string  `json:"payment_channel,omitempty"`
	Amount                  float64 `gorm:"not null" json:"amount"`
	Currency                string  `gorm:"default:CNY" json:"currency"`
	Status                  string  `gorm:"default:pending" json:"status"`
	ThirdPartyTransactionID string  `gorm:"index" json:"third_party_transaction_id,omitempty"`
	ThirdPartyResponse      string  `json:"third_party_response,omitempty"`
	PaidAt                  *string `json:"paid_at,omitempty"`
	NotifyAt                *string `json:"notify_at,omitempty"`
	NotifyData              string  `json:"notify_data,omitempty"`
	RetryCount              int     `gorm:"default:0" json:"retry_count"`
	ErrorMessage            string  `json:"error_message,omitempty"`
}

// TableName returns the table name for Payment
func (Payment) TableName() string {
	return "payments"
}

// PaymentMethod constants
const (
	PaymentMethodWechat = "wechat"
	PaymentMethodAlipay = "alipay"
	PaymentMethodCard   = "card"
)

// PaymentStatus constants
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusSuccess    = "success"
	PaymentStatusFailed     = "failed"
	PaymentStatusCancelled  = "cancelled"
)

// IsSuccessful checks if payment was successful
func (p *Payment) IsSuccessful() bool {
	return p.Status == PaymentStatusSuccess
}

// CanRetry checks if payment can be retried
func (p *Payment) CanRetry() bool {
	return p.Status == PaymentStatusFailed && p.RetryCount < 3
}
