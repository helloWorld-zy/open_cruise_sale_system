package service

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrInvalidOrderData      = errors.New("invalid order data")
	ErrOrderCreationFailed   = errors.New("failed to create order")
	ErrInvalidPassengerCount = errors.New("invalid passenger count")
	ErrCabinNotAvailable     = errors.New("cabin is not available")
)

// OrderService defines the interface for order business logic
type OrderService interface {
	// Create creates a new order with inventory locking
	Create(ctx context.Context, req CreateOrderRequest) (*domain.Order, error)

	// GetByID retrieves an order by ID
	GetByID(ctx context.Context, id string) (*domain.Order, error)

	// GetByOrderNumber retrieves an order by order number
	GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error)

	// GetWithDetails retrieves an order with all related details
	GetWithDetails(ctx context.Context, id string) (*domain.Order, error)

	// List retrieves a paginated list of orders
	List(ctx context.Context, req ListOrdersRequest) (*pagination.Result, error)

	// ListByUser retrieves orders for a specific user
	ListByUser(ctx context.Context, userID string, paginator *pagination.Paginator) (*pagination.Result, error)

	// Update updates an order
	Update(ctx context.Context, id string, req UpdateOrderRequest) (*domain.Order, error)

	// Cancel cancels an order and releases inventory
	Cancel(ctx context.Context, id string) error

	// Confirm confirms a paid order
	Confirm(ctx context.Context, id string) error

	// Complete completes a confirmed order
	Complete(ctx context.Context, id string) error

	// Delete deletes an order (admin only)
	Delete(ctx context.Context, id string) error

	// CalculateTotal calculates order total from items
	CalculateTotal(ctx context.Context, items []OrderItemRequest) (OrderCalculation, error)
}

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	UserID       string             `json:"user_id,omitempty"`
	VoyageID     string             `json:"voyage_id" validate:"required"`
	CruiseID     string             `json:"cruise_id" validate:"required"`
	Items        []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
	Passengers   []PassengerRequest `json:"passengers" validate:"required,min=1,dive"`
	ContactName  string             `json:"contact_name" validate:"required"`
	ContactPhone string             `json:"contact_phone" validate:"required"`
	ContactEmail string             `json:"contact_email" validate:"required,email"`
	Remark       string             `json:"remark,omitempty"`
}

// OrderItemRequest represents an item in an order
type OrderItemRequest struct {
	CabinID     string `json:"cabin_id" validate:"required"`
	CabinTypeID string `json:"cabin_type_id" validate:"required"`
	AdultCount  int    `json:"adult_count" validate:"required,min=1"`
	ChildCount  int    `json:"child_count" validate:"gte=0"`
	InfantCount int    `json:"infant_count" validate:"gte=0"`
}

// PassengerRequest represents a passenger in an order
type PassengerRequest struct {
	OrderItemID           string `json:"order_item_id"`
	Name                  string `json:"name" validate:"required"`
	Surname               string `json:"surname" validate:"required"`
	GivenName             string `json:"given_name,omitempty"`
	Gender                string `json:"gender" validate:"required,oneof=male female"`
	BirthDate             string `json:"birth_date" validate:"required"`
	Nationality           string `json:"nationality,omitempty"`
	PassportNumber        string `json:"passport_number,omitempty"`
	PassportExpiry        string `json:"passport_expiry,omitempty"`
	IDNumber              string `json:"id_number,omitempty"`
	Phone                 string `json:"phone,omitempty"`
	Email                 string `json:"email,omitempty"`
	PassengerType         string `json:"passenger_type" validate:"required,oneof=adult child infant"`
	EmergencyContactName  string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone string `json:"emergency_contact_phone,omitempty"`
	DietaryRequirements   string `json:"dietary_requirements,omitempty"`
	MedicalNotes          string `json:"medical_notes,omitempty"`
}

// UpdateOrderRequest represents a request to update an order
type UpdateOrderRequest struct {
	ContactName  string `json:"contact_name,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
	Remark       string `json:"remark,omitempty"`
}

// ListOrdersRequest represents a request to list orders
type ListOrdersRequest struct {
	UserID        string `form:"user_id"`
	VoyageID      string `form:"voyage_id"`
	Status        string `form:"status"`
	PaymentStatus string `form:"payment_status"`
	OrderNumber   string `form:"order_number"`
	DateFrom      string `form:"date_from"`
	DateTo        string `form:"date_to"`
	pagination.Paginator
}

// OrderCalculation represents the result of order calculation
type OrderCalculation struct {
	Subtotal       float64           `json:"subtotal"`
	PortFee        float64           `json:"port_fee"`
	ServiceFee     float64           `json:"service_fee"`
	DiscountAmount float64           `json:"discount_amount"`
	TotalAmount    float64           `json:"total_amount"`
	Items          []ItemCalculation `json:"items"`
}

// ItemCalculation represents calculation for a single item
type ItemCalculation struct {
	CabinTypeID string  `json:"cabin_type_id"`
	CabinNumber string  `json:"cabin_number"`
	AdultPrice  float64 `json:"adult_price"`
	ChildPrice  float64 `json:"child_price"`
	InfantPrice float64 `json:"infant_price"`
	PortFee     float64 `json:"port_fee"`
	ServiceFee  float64 `json:"service_fee"`
	Subtotal    float64 `json:"subtotal"`
}

// orderService implements OrderService
type orderService struct {
	orderRepo     repository.OrderRepository
	voyageRepo    repository.VoyageRepository
	cabinRepo     repository.CabinRepository
	priceRepo     repository.PriceRepository
	inventoryRepo repository.InventoryRepository
	stateService  OrderStateService
	redis         *redis.Client
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo repository.OrderRepository,
	voyageRepo repository.VoyageRepository,
	cabinRepo repository.CabinRepository,
	priceRepo repository.PriceRepository,
	inventoryRepo repository.InventoryRepository,
	redisClients ...*redis.Client,
) OrderService {
	stateService := NewOrderStateService(orderRepo, inventoryRepo)
	var redisClient *redis.Client
	if len(redisClients) > 0 {
		redisClient = redisClients[0]
	}
	return &orderService{
		orderRepo:     orderRepo,
		voyageRepo:    voyageRepo,
		cabinRepo:     cabinRepo,
		priceRepo:     priceRepo,
		inventoryRepo: inventoryRepo,
		stateService:  stateService,
		redis:         redisClient,
	}
}

// itemSubtotal holds the breakdown of a calculated item subtotal
// CS-003: Extracted from duplicated logic in Create and CalculateTotal
type itemSubtotal struct {
	AdultTotal  float64
	ChildTotal  float64
	InfantTotal float64
	PortFee     float64
	ServiceFee  float64
	Subtotal    float64
}

// calculateItemSubtotal computes the price breakdown for a single order item
func calculateItemSubtotal(price *domain.CabinPrice, adultCount, childCount, infantCount int) itemSubtotal {
	adultTotal := price.AdultPrice * float64(adultCount)
	childTotal := price.ChildPrice * float64(childCount)
	infantTotal := price.InfantPrice * float64(infantCount)
	portFee := price.PortFee * float64(adultCount+childCount)
	serviceFee := price.ServiceFee * float64(adultCount+childCount)
	return itemSubtotal{
		AdultTotal:  adultTotal,
		ChildTotal:  childTotal,
		InfantTotal: infantTotal,
		PortFee:     portFee,
		ServiceFee:  serviceFee,
		Subtotal:    adultTotal + childTotal + infantTotal + portFee + serviceFee,
	}
}

func (s *orderService) Create(ctx context.Context, req CreateOrderRequest) (*domain.Order, error) {
	// Validate voyage exists
	voyage, err := s.voyageRepo.GetByID(ctx, req.VoyageID)
	if err != nil {
		return nil, fmt.Errorf("voyage not found: %w", err)
	}

	// Check voyage is open for booking
	if voyage.BookingStatus != domain.BookingStatusOpen {
		return nil, errors.New("voyage is not open for booking")
	}

	// Validate total passenger count
	totalPassengers := 0
	for _, item := range req.Items {
		totalPassengers += item.AdultCount + item.ChildCount
	}
	if totalPassengers != len(req.Passengers) {
		return nil, ErrInvalidPassengerCount
	}

	now := time.Now().UTC()
	expiresAt := now.Add(15 * time.Minute) // 15 minutes to pay (aligned with lock timeout)

	order := &domain.Order{
		OrderNumber:    generateOrderNumber(),
		VoyageID:       req.VoyageID,
		CruiseID:       req.CruiseID,
		TotalAmount:    0,
		Status:         domain.OrderStatusPending,
		PaymentStatus:  domain.PaymentStatusUnpaid,
		PassengerCount: totalPassengers,
		ContactName:    req.ContactName,
		ContactPhone:   req.ContactPhone,
		ContactEmail:   req.ContactEmail,
		Remark:         req.Remark,
		ExpiresAt:      expiresAt.Format(time.RFC3339),
	}

	if req.UserID != "" {
		order.UserID = &req.UserID
	}

	// DD-001: Wrap entire order creation in a database transaction
	err = s.orderRepo.WithTransaction(ctx, func(txRepo repository.OrderRepository, tx *gorm.DB) error {
		txInventoryRepo := repository.NewInventoryRepository(tx)

		// Create order
		if err := txRepo.Create(ctx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		var totalAmount float64
		var cabinCount int

		// Process each item and lock inventory
		for _, itemReq := range req.Items {
			unlock, err := s.acquireInventoryLock(ctx, req.VoyageID, itemReq.CabinTypeID, order.ID.String())
			if err != nil {
				return fmt.Errorf("failed to acquire inventory lock: %w", err)
			}
			defer unlock()

			// Get cabin
			cabin, err := s.cabinRepo.GetByID(ctx, itemReq.CabinID)
			if err != nil {
				return fmt.Errorf("cabin not found: %w", err)
			}

			if cabin.Status != domain.CabinStatusAvailable {
				return ErrCabinNotAvailable
			}

			// Get price
			price, err := s.priceRepo.GetCurrentPrice(ctx, req.VoyageID, itemReq.CabinTypeID)
			if err != nil {
				return fmt.Errorf("price not found: %w", err)
			}

			// Lock inventory
			if err := txInventoryRepo.LockCabin(ctx, req.VoyageID, itemReq.CabinTypeID, 1); err != nil {
				return fmt.Errorf("failed to lock cabin: %w", err)
			}

			// CS-003: Use extracted helper for price calculation
			calc := calculateItemSubtotal(price, itemReq.AdultCount, itemReq.ChildCount, itemReq.InfantCount)

			orderItem := &domain.OrderItem{
				OrderID:       order.ID.String(),
				CabinID:       itemReq.CabinID,
				CabinTypeID:   itemReq.CabinTypeID,
				VoyageID:      req.VoyageID,
				CabinNumber:   cabin.CabinNumber,
				PriceSnapshot: price.AdultPrice,
				AdultCount:    itemReq.AdultCount,
				ChildCount:    itemReq.ChildCount,
				InfantCount:   itemReq.InfantCount,
				AdultPrice:    price.AdultPrice,
				ChildPrice:    price.ChildPrice,
				InfantPrice:   price.InfantPrice,
				PortFee:       calc.PortFee,
				ServiceFee:    calc.ServiceFee,
				Subtotal:      calc.Subtotal,
				Status:        domain.OrderItemStatusConfirmed,
			}

			if err := txRepo.CreateOrderItem(ctx, orderItem); err != nil {
				// Transaction will auto-rollback; attempt unlock for extra safety
				_ = txInventoryRepo.UnlockCabin(ctx, req.VoyageID, itemReq.CabinTypeID, 1)
				return fmt.Errorf("failed to create order item: %w", err)
			}

			totalAmount += calc.Subtotal
			cabinCount++
		}

		// Update order total
		order.TotalAmount = totalAmount
		order.CabinCount = cabinCount
		if err := txRepo.Update(ctx, order); err != nil {
			return fmt.Errorf("failed to update order total: %w", err)
		}

		// Create passengers
		var passengers []*domain.Passenger
		for i, p := range req.Passengers {
			passenger := &domain.Passenger{
				OrderID:               order.ID.String(),
				Name:                  p.Name,
				Surname:               p.Surname,
				GivenName:             p.GivenName,
				Gender:                p.Gender,
				BirthDate:             p.BirthDate,
				Nationality:           p.Nationality,
				PassportNumber:        p.PassportNumber,
				PassportExpiry:        p.PassportExpiry,
				IDNumber:              p.IDNumber,
				Phone:                 p.Phone,
				Email:                 p.Email,
				PassengerType:         p.PassengerType,
				EmergencyContactName:  p.EmergencyContactName,
				EmergencyContactPhone: p.EmergencyContactPhone,
				DietaryRequirements:   p.DietaryRequirements,
				MedicalNotes:          p.MedicalNotes,
			}

			// Associate passenger with corresponding order item
			if i < len(order.Items) {
				passenger.OrderItemID = order.Items[i].ID.String()
			}

			passengers = append(passengers, passenger)
		}

		if err := txRepo.BatchCreatePassengers(ctx, passengers); err != nil {
			return fmt.Errorf("failed to create passengers: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *orderService) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *orderService) GetWithDetails(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.orderRepo.GetOrderWithDetails(ctx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *orderService) List(ctx context.Context, req ListOrdersRequest) (*pagination.Result, error) {
	filters := repository.OrderFilters{
		UserID:        req.UserID,
		VoyageID:      req.VoyageID,
		Status:        req.Status,
		PaymentStatus: req.PaymentStatus,
		OrderNumber:   req.OrderNumber,
		DateFrom:      req.DateFrom,
		DateTo:        req.DateTo,
	}

	count, err := s.orderRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	paginator := &req.Paginator
	paginator.SetTotal(count)

	orders, err := s.orderRepo.List(ctx, filters, paginator)
	if err != nil {
		return nil, err
	}

	return &pagination.Result{
		Data:       orders,
		Pagination: *paginator,
	}, nil
}

func (s *orderService) ListByUser(ctx context.Context, userID string, paginator *pagination.Paginator) (*pagination.Result, error) {
	orders, err := s.orderRepo.ListByUser(ctx, userID, paginator)
	if err != nil {
		return nil, err
	}

	return &pagination.Result{
		Data:       orders,
		Pagination: *paginator,
	}, nil
}

func (s *orderService) Update(ctx context.Context, id string, req UpdateOrderRequest) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	// Only allow updates for pending orders
	if order.Status != domain.OrderStatusPending {
		return nil, errors.New("only pending orders can be updated")
	}

	// Update fields if provided
	if req.ContactName != "" {
		order.ContactName = req.ContactName
	}
	if req.ContactPhone != "" {
		order.ContactPhone = req.ContactPhone
	}
	if req.ContactEmail != "" {
		order.ContactEmail = req.ContactEmail
	}
	if req.Remark != "" {
		order.Remark = req.Remark
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) Cancel(ctx context.Context, id string) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return ErrOrderNotFound
	}

	return s.stateService.TransitionToCancelled(ctx, order)
}

func (s *orderService) Confirm(ctx context.Context, id string) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return ErrOrderNotFound
	}

	return s.stateService.TransitionToConfirmed(ctx, order)
}

func (s *orderService) Complete(ctx context.Context, id string) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return ErrOrderNotFound
	}

	return s.stateService.TransitionToCompleted(ctx, order)
}

func (s *orderService) Delete(ctx context.Context, id string) error {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return ErrOrderNotFound
	}

	// Only allow deletion of pending or cancelled orders
	if order.Status != domain.OrderStatusPending && order.Status != domain.OrderStatusCancelled {
		return errors.New("only pending or cancelled orders can be deleted")
	}

	return s.orderRepo.Delete(ctx, id)
}

func (s *orderService) CalculateTotal(ctx context.Context, items []OrderItemRequest) (OrderCalculation, error) {
	var calculation OrderCalculation

	for _, item := range items {
		price, err := s.priceRepo.GetCurrentPrice(ctx, "", item.CabinTypeID)
		if err != nil {
			return calculation, fmt.Errorf("price not found for cabin type %s: %w", item.CabinTypeID, err)
		}

		// CS-003: Use extracted helper for price calculation
		calc := calculateItemSubtotal(price, item.AdultCount, item.ChildCount, item.InfantCount)

		calculation.Items = append(calculation.Items, ItemCalculation{
			CabinTypeID: item.CabinTypeID,
			AdultPrice:  calc.AdultTotal,
			ChildPrice:  calc.ChildTotal,
			InfantPrice: calc.InfantTotal,
			PortFee:     calc.PortFee,
			ServiceFee:  calc.ServiceFee,
			Subtotal:    calc.Subtotal,
		})

		calculation.Subtotal += calc.AdultTotal + calc.ChildTotal + calc.InfantTotal
		calculation.PortFee += calc.PortFee
		calculation.ServiceFee += calc.ServiceFee
	}

	calculation.TotalAmount = calculation.Subtotal + calculation.PortFee + calculation.ServiceFee - calculation.DiscountAmount

	return calculation, nil
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD%s%s", time.Now().Format("20060102"), uuid.New().String()[:8])
}

func (s *orderService) acquireInventoryLock(ctx context.Context, voyageID, cabinTypeID, orderID string) (func(), error) {
	if s.redis == nil {
		return func() {}, nil
	}

	lockKey := fmt.Sprintf("lock:inventory:%s:%s:%s", voyageID, cabinTypeID, orderID)
	lockValue := uuid.New().String()

	var acquired bool
	for i := 0; i < 3; i++ {
		ok, err := s.redis.SetNX(ctx, lockKey, lockValue, 30*time.Second).Result()
		if err != nil {
			return nil, err
		}
		if ok {
			acquired = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !acquired {
		return nil, errors.New("inventory lock timeout")
	}

	unlock := func() {
		script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		end
		return 0`
		_, _ = s.redis.Eval(ctx, script, []string{lockKey}, lockValue).Result()
	}

	return unlock, nil
}
