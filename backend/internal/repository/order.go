package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"
	"errors"

	"gorm.io/gorm"
)

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
	List(ctx context.Context, filters OrderFilters, paginator *pagination.Paginator) ([]*domain.Order, error)
	Count(ctx context.Context, filters OrderFilters) (int64, error)
	ListByUser(ctx context.Context, userID string, paginator *pagination.Paginator) ([]*domain.Order, error)
	ListByStatus(ctx context.Context, status string, paginator *pagination.Paginator) ([]*domain.Order, error)
	ListByVoyage(ctx context.Context, voyageID string) ([]*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdatePaymentStatus(ctx context.Context, id string, paymentStatus string, paidAmount float64) error
	Delete(ctx context.Context, id string) error

	// OrderItem operations
	CreateOrderItem(ctx context.Context, item *domain.OrderItem) error
	GetOrderItemByID(ctx context.Context, id string) (*domain.OrderItem, error)
	ListOrderItemsByOrder(ctx context.Context, orderID string) ([]*domain.OrderItem, error)
	UpdateOrderItem(ctx context.Context, item *domain.OrderItem) error
	UpdateOrderItemStatus(ctx context.Context, id string, status string) error
	DeleteOrderItem(ctx context.Context, id string) error

	// Passenger operations
	CreatePassenger(ctx context.Context, passenger *domain.Passenger) error
	GetPassengerByID(ctx context.Context, id string) (*domain.Passenger, error)
	ListPassengersByOrder(ctx context.Context, orderID string) ([]*domain.Passenger, error)
	ListPassengersByOrderItem(ctx context.Context, orderItemID string) ([]*domain.Passenger, error)
	UpdatePassenger(ctx context.Context, passenger *domain.Passenger) error
	DeletePassenger(ctx context.Context, id string) error
	BatchCreatePassengers(ctx context.Context, passengers []*domain.Passenger) error

	// Complex queries
	GetOrderWithDetails(ctx context.Context, id string) (*domain.Order, error)

	// Payment operations
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error)
	GetPaymentByNo(ctx context.Context, paymentNo string) (*domain.Payment, error)
	UpdatePayment(ctx context.Context, payment *domain.Payment) error

	// Refund operations
	CreateRefundRequest(ctx context.Context, refund *domain.RefundRequest) error
	GetRefundRequestByID(ctx context.Context, id string) (*domain.RefundRequest, error)
	ListRefundRequests(ctx context.Context, filters RefundFilters, paginator *pagination.Paginator) ([]*domain.RefundRequest, error)
	UpdateRefundRequest(ctx context.Context, refund *domain.RefundRequest) error
}

// RefundFilters represents filters for refund queries
type RefundFilters struct {
	OrderID    string
	UserID     string
	Status     string
	ReviewedBy string
}

// OrderFilters represents filters for order queries
type OrderFilters struct {
	UserID        string
	VoyageID      string
	Status        string
	PaymentStatus string
	OrderNumber   string
	ContactPhone  string
	ContactEmail  string
	DateFrom      string
	DateTo        string
}

// orderRepository implements OrderRepository
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// ==================== Order Operations ====================

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		First(&order, "order_number = ?", orderNumber).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) List(ctx context.Context, filters OrderFilters, paginator *pagination.Paginator) ([]*domain.Order, error) {
	query := r.buildOrderQuery(filters)

	var orders []*domain.Order
	err := query.WithContext(ctx).
		Order("created_at DESC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&orders).Error

	return orders, err
}

func (r *orderRepository) Count(ctx context.Context, filters OrderFilters) (int64, error) {
	query := r.buildOrderQuery(filters)

	var count int64
	err := query.WithContext(ctx).Model(&domain.Order{}).Count(&count).Error
	return count, err
}

func (r *orderRepository) buildOrderQuery(filters OrderFilters) *gorm.DB {
	query := r.db.Model(&domain.Order{})

	if filters.UserID != "" {
		query = query.Where("user_id = ?", filters.UserID)
	}
	if filters.VoyageID != "" {
		query = query.Where("voyage_id = ?", filters.VoyageID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filters.PaymentStatus)
	}
	if filters.OrderNumber != "" {
		query = query.Where("order_number LIKE ?", "%"+filters.OrderNumber+"%")
	}
	if filters.ContactPhone != "" {
		query = query.Where("contact_phone = ?", filters.ContactPhone)
	}
	if filters.ContactEmail != "" {
		query = query.Where("contact_email = ?", filters.ContactEmail)
	}
	if filters.DateFrom != "" {
		query = query.Where("created_at >= ?", filters.DateFrom)
	}
	if filters.DateTo != "" {
		query = query.Where("created_at <= ?", filters.DateTo)
	}

	return query
}

func (r *orderRepository) ListByUser(ctx context.Context, userID string, paginator *pagination.Paginator) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) ListByStatus(ctx context.Context, status string, paginator *pagination.Paginator) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) ListByVoyage(ctx context.Context, voyageID string) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.WithContext(ctx).
		Where("voyage_id = ?", voyageID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) Update(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	now := getCurrentTimestamp()
	updates := map[string]interface{}{
		"status": status,
	}

	if status == domain.OrderStatusConfirmed {
		updates["confirmed_at"] = now
	} else if status == domain.OrderStatusCancelled {
		updates["cancelled_at"] = now
	}

	return r.db.WithContext(ctx).Model(&domain.Order{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *orderRepository) UpdatePaymentStatus(ctx context.Context, id string, paymentStatus string, paidAmount float64) error {
	now := getCurrentTimestamp()
	updates := map[string]interface{}{
		"payment_status": paymentStatus,
		"paid_amount":    paidAmount,
	}

	if paymentStatus == domain.PaymentStatusPaid {
		updates["paid_at"] = now
		updates["status"] = domain.OrderStatusPaid
	}

	return r.db.WithContext(ctx).Model(&domain.Order{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *orderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Order{}, "id = ?", id).Error
}

// ==================== OrderItem Operations ====================

func (r *orderRepository) CreateOrderItem(ctx context.Context, item *domain.OrderItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *orderRepository) GetOrderItemByID(ctx context.Context, id string) (*domain.OrderItem, error) {
	var item domain.OrderItem
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("Cabin").
		Preload("CabinType").
		First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *orderRepository) ListOrderItemsByOrder(ctx context.Context, orderID string) ([]*domain.OrderItem, error) {
	var items []*domain.OrderItem
	err := r.db.WithContext(ctx).
		Preload("Cabin").
		Preload("CabinType").
		Where("order_id = ?", orderID).
		Find(&items).Error
	return items, err
}

func (r *orderRepository) UpdateOrderItem(ctx context.Context, item *domain.OrderItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *orderRepository) UpdateOrderItemStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.OrderItem{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *orderRepository) DeleteOrderItem(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.OrderItem{}, "id = ?", id).Error
}

// ==================== Passenger Operations ====================

func (r *orderRepository) CreatePassenger(ctx context.Context, passenger *domain.Passenger) error {
	return r.db.WithContext(ctx).Create(passenger).Error
}

func (r *orderRepository) GetPassengerByID(ctx context.Context, id string) (*domain.Passenger, error) {
	var passenger domain.Passenger
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("OrderItem").
		First(&passenger, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &passenger, nil
}

func (r *orderRepository) ListPassengersByOrder(ctx context.Context, orderID string) ([]*domain.Passenger, error) {
	var passengers []*domain.Passenger
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Find(&passengers).Error
	return passengers, err
}

func (r *orderRepository) ListPassengersByOrderItem(ctx context.Context, orderItemID string) ([]*domain.Passenger, error) {
	var passengers []*domain.Passenger
	err := r.db.WithContext(ctx).
		Where("order_item_id = ?", orderItemID).
		Find(&passengers).Error
	return passengers, err
}

func (r *orderRepository) UpdatePassenger(ctx context.Context, passenger *domain.Passenger) error {
	return r.db.WithContext(ctx).Save(passenger).Error
}

func (r *orderRepository) DeletePassenger(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Passenger{}, "id = ?", id).Error
}

func (r *orderRepository) BatchCreatePassengers(ctx context.Context, passengers []*domain.Passenger) error {
	return r.db.WithContext(ctx).Create(passengers).Error
}

// ==================== Complex Queries ====================

func (r *orderRepository) GetOrderWithDetails(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Voyage").
		Preload("Voyage.Route").
		Preload("Cruise").
		Preload("Items").
		Preload("Items.Cabin").
		Preload("Items.CabinType").
		Preload("Passengers").
		Preload("Payments").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ==================== Payment Operations ====================

func (r *orderRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *orderRepository) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *orderRepository) GetPaymentByNo(ctx context.Context, paymentNo string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		First(&payment, "payment_no = ?", paymentNo).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *orderRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

// ==================== Refund Operations ====================

func (r *orderRepository) CreateRefundRequest(ctx context.Context, refund *domain.RefundRequest) error {
	return r.db.WithContext(ctx).Create(refund).Error
}

func (r *orderRepository) GetRefundRequestByID(ctx context.Context, id string) (*domain.RefundRequest, error) {
	var refund domain.RefundRequest
	err := r.db.WithContext(ctx).
		Preload("Order").
		Preload("OrderItem").
		First(&refund, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

func (r *orderRepository) ListRefundRequests(ctx context.Context, filters RefundFilters, paginator *pagination.Paginator) ([]*domain.RefundRequest, error) {
	query := r.db.Model(&domain.RefundRequest{})

	if filters.OrderID != "" {
		query = query.Where("order_id = ?", filters.OrderID)
	}
	if filters.UserID != "" {
		query = query.Where("user_id = ?", filters.UserID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.ReviewedBy != "" {
		query = query.Where("reviewed_by = ?", filters.ReviewedBy)
	}

	var refunds []*domain.RefundRequest
	err := query.WithContext(ctx).
		Preload("Order").
		Order("requested_at DESC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&refunds).Error

	return refunds, err
}

func (r *orderRepository) UpdateRefundRequest(ctx context.Context, refund *domain.RefundRequest) error {
	return r.db.WithContext(ctx).Save(refund).Error
}

// Helper function
func getCurrentTimestamp() string {
	return "" // This should be implemented properly in real code
}

// ErrOrderNotFound is returned when an order is not found
var ErrOrderNotFound = errors.New("order not found")
