package payment

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// PaymentService provides high-level payment operations
type PaymentService interface {
	// CreatePayment creates a payment for an order
	CreatePayment(ctx context.Context, orderID string, method string, description string) (*domain.Payment, error)

	// ProcessCallback processes payment provider callback
	ProcessCallback(ctx context.Context, provider string, body []byte, signature string) error

	// QueryPayment queries payment status
	QueryPayment(ctx context.Context, paymentID string) (*domain.Payment, error)

	// Refund processes a refund
	Refund(ctx context.Context, paymentID string, amount float64, reason string) error

	// GetPaymentByOrder gets payment by order ID
	GetPaymentByOrder(ctx context.Context, orderID string) (*domain.Payment, error)
}

// paymentService implements PaymentService
type paymentService struct {
	orderRepo repository.OrderRepository
	providers map[string]PaymentProvider
	natsConn  *nats.Conn
}

// NewPaymentService creates a new payment service
func NewPaymentService(orderRepo repository.OrderRepository, natsConn *nats.Conn) PaymentService {
	return &paymentService{
		orderRepo: orderRepo,
		providers: make(map[string]PaymentProvider),
		natsConn:  natsConn,
	}
}

// RegisterProvider registers a payment provider
func (s *paymentService) RegisterProvider(name string, provider PaymentProvider) {
	s.providers[name] = provider
}

// CreatePayment creates a payment for an order
func (s *paymentService) CreatePayment(ctx context.Context, orderID string, method string, description string) (*domain.Payment, error) {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Check if order can be paid
	if order.Status != domain.OrderStatusPending {
		return nil, fmt.Errorf("order cannot be paid in status %s", order.Status)
	}

	// Get provider
	provider, exists := s.providers[method]
	if !exists {
		return nil, fmt.Errorf("unsupported payment method: %s", method)
	}

	// Create payment with provider
	result, err := provider.CreatePayment(ctx, order, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Get payment from database
	payment, err := s.orderRepo.GetPaymentByNo(ctx, result.PaymentNo)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// ProcessCallback processes payment callback
func (s *paymentService) ProcessCallback(ctx context.Context, provider string, body []byte, signature string) error {
	// Get provider
	prov, exists := s.providers[provider]
	if !exists {
		return fmt.Errorf("unsupported payment provider: %s", provider)
	}

	// Process callback
	result, err := prov.ProcessCallback(ctx, body, signature)
	if err != nil {
		return fmt.Errorf("failed to process callback: %w", err)
	}

	// Get payment
	payment, err := s.orderRepo.GetPaymentByNo(ctx, result.PaymentNo)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// SEC-004: Idempotency check - skip if payment is already in terminal state
	if payment.Status == domain.PaymentStatusSuccess ||
		payment.Status == domain.PaymentStatusRefunded {
		// Payment already processed, skip duplicate callback
		return nil
	}

	// Update payment status
	payment.Status = result.Status
	payment.ThirdPartyTransactionID = result.ThirdPartyID
	payment.PaidAt = &result.PaidAt

	if err := s.orderRepo.UpdatePayment(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	// If payment successful, update order
	if result.Status == domain.PaymentStatusSuccess {
		if err := s.orderRepo.UpdatePaymentStatus(ctx, payment.OrderID, domain.PaymentStatusPaid, result.Amount); err != nil {
			return fmt.Errorf("failed to update order payment status: %w", err)
		}

		// Publish payment success event
		s.publishPaymentEvent("payment.success", payment)
	}

	return nil
}

// QueryPayment queries payment status from provider and updates database
func (s *paymentService) QueryPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
	// Get payment from database
	payment, err := s.orderRepo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// Get provider
	provider, exists := s.providers[payment.PaymentMethod]
	if !exists {
		return nil, fmt.Errorf("payment provider not found: %s", payment.PaymentMethod)
	}

	// Query from provider
	result, err := provider.QueryPayment(ctx, payment.PaymentNo)
	if err != nil {
		return nil, err
	}

	// Update payment if status changed
	if result.Status != payment.Status {
		payment.Status = result.Status
		payment.ThirdPartyTransactionID = result.ThirdPartyID

		if result.Status == domain.PaymentStatusSuccess {
			now := time.Now().Format(time.RFC3339)
			payment.PaidAt = &now

			// Update order
			if err := s.orderRepo.UpdatePaymentStatus(ctx, payment.OrderID, domain.PaymentStatusPaid, result.Amount); err != nil {
				return nil, err
			}
		}

		if err := s.orderRepo.UpdatePayment(ctx, payment); err != nil {
			return nil, err
		}
	}

	return payment, nil
}

// Refund processes a refund
func (s *paymentService) Refund(ctx context.Context, paymentID string, amount float64, reason string) error {
	// Get payment
	payment, err := s.orderRepo.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// Verify payment is successful
	if payment.Status != domain.PaymentStatusSuccess {
		return fmt.Errorf("cannot refund payment in status %s", payment.Status)
	}

	// Get provider
	provider, exists := s.providers[payment.PaymentMethod]
	if !exists {
		return fmt.Errorf("payment provider not found: %s", payment.PaymentMethod)
	}

	// Process refund
	result, err := provider.Refund(ctx, payment, amount, reason)
	if err != nil {
		return fmt.Errorf("refund failed: %w", err)
	}

	// Update payment status
	if result.Status == "SUCCESS" {
		payment.Status = domain.PaymentStatusRefunded
		if err := s.orderRepo.UpdatePayment(ctx, payment); err != nil {
			return err
		}

		// Update order status
		if err := s.orderRepo.UpdateStatus(ctx, payment.OrderID, domain.OrderStatusRefunded); err != nil {
			return err
		}

		// Release inventory
		s.publishPaymentEvent("payment.refunded", payment)
	}

	return nil
}

// GetPaymentByOrder gets payment by order ID
func (s *paymentService) GetPaymentByOrder(ctx context.Context, orderID string) (*domain.Payment, error) {
	order, err := s.orderRepo.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if len(order.Payments) == 0 {
		return nil, fmt.Errorf("no payments found for order")
	}

	// Return most recent payment
	return &order.Payments[len(order.Payments)-1], nil
}

func (s *paymentService) publishPaymentEvent(eventType string, payment *domain.Payment) {
	if s.natsConn == nil {
		return
	}

	event := map[string]interface{}{
		"type":       eventType,
		"payment_id": payment.ID,
		"order_id":   payment.OrderID,
		"amount":     payment.Amount,
		"status":     payment.Status,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	data, _ := json.Marshal(event)
	if err := s.natsConn.Publish(eventType, data); err != nil {
		// Log but don't fail - the payment is already processed
		fmt.Printf("[WARN] Failed to publish payment event %s: %v\n", eventType, err)
	}
}
