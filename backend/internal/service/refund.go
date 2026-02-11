package service

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/payment"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
)

var (
	ErrRefundNotFound       = errors.New("refund request not found")
	ErrInvalidRefundData    = errors.New("invalid refund data")
	ErrOrderNotRefundable   = errors.New("order is not refundable")
	ErrRefundNotReviewable  = errors.New("refund request cannot be reviewed")
	ErrRefundNotProcessable = errors.New("refund request cannot be processed")
	ErrRefundAmountInvalid  = errors.New("refund amount is invalid")
)

// RefundService defines the interface for refund workflow
type RefundService interface {
	// CreateRefundRequest creates a new refund request
	CreateRefundRequest(ctx context.Context, req CreateRefundRequest) (*domain.RefundRequest, error)

	// GetRefundByID retrieves a refund request by ID
	GetRefundByID(ctx context.Context, id string) (*domain.RefundRequest, error)

	// ListRefunds retrieves a paginated list of refunds
	ListRefunds(ctx context.Context, filters repository.RefundFilters, paginator *pagination.Paginator) ([]*domain.RefundRequest, error)

	// ApproveRefund approves a refund request
	ApproveRefund(ctx context.Context, id string, reviewerID, note string) error

	// RejectRefund rejects a refund request
	RejectRefund(ctx context.Context, id string, reviewerID, note string) error

	// ProcessRefund processes the actual refund payment
	ProcessRefund(ctx context.Context, id string) error

	// GetRefundableAmount calculates the maximum refundable amount for an order
	GetRefundableAmount(ctx context.Context, orderID string) (float64, error)
}

// CreateRefundRequest represents a request to create a refund
type CreateRefundRequest struct {
	OrderID            string  `json:"order_id" validate:"required"`
	OrderItemID        *string `json:"order_item_id,omitempty"`
	UserID             *string `json:"user_id,omitempty"`
	RefundAmount       float64 `json:"refund_amount" validate:"required,gt=0"`
	RefundReason       string  `json:"refund_reason" validate:"required"`
	RefundType         string  `json:"refund_type" validate:"omitempty,oneof=full partial"`
	RefundMethod       string  `json:"refund_method" validate:"omitempty,oneof=original bank alipay wechat"`
	BankName           string  `json:"bank_name,omitempty"`
	BankAccount        string  `json:"bank_account,omitempty"`
	AccountHolder      string  `json:"account_holder,omitempty"`
	CancellationReason string  `json:"cancellation_reason" validate:"omitempty,oneof=customer_request voyage_cancelled cabin_upgrade other"`
}

// refundService implements RefundService
type refundService struct {
	repo           repository.OrderRepository
	paymentService payment.PaymentService
	orderService   OrderService
}

// NewRefundService creates a new refund service
func NewRefundService(
	repo repository.OrderRepository,
	paymentService payment.PaymentService,
	orderService OrderService,
) RefundService {
	return &refundService{
		repo:           repo,
		paymentService: paymentService,
		orderService:   orderService,
	}
}

func (s *refundService) CreateRefundRequest(ctx context.Context, req CreateRefundRequest) (*domain.RefundRequest, error) {
	// Get order
	order, err := s.orderService.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Check if order is refundable
	if !order.CanCancel() {
		return nil, ErrOrderNotRefundable
	}

	// Check refund amount
	maxRefundable, err := s.GetRefundableAmount(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	if req.RefundAmount > maxRefundable {
		return nil, ErrRefundAmountInvalid
	}

	// Determine refund type
	refundType := req.RefundType
	if refundType == "" {
		if req.RefundAmount >= order.TotalAmount {
			refundType = domain.RefundTypeFull
		} else {
			refundType = domain.RefundTypePartial
		}
	}

	refund := &domain.RefundRequest{
		OrderID:            req.OrderID,
		OrderItemID:        req.OrderItemID,
		UserID:             req.UserID,
		RefundAmount:       req.RefundAmount,
		RefundReason:       req.RefundReason,
		RefundType:         refundType,
		RefundMethod:       req.RefundMethod,
		BankName:           req.BankName,
		BankAccount:        req.BankAccount,
		AccountHolder:      req.AccountHolder,
		CancellationReason: req.CancellationReason,
		Status:             domain.RefundStatusPending,
	}

	if err := s.repo.CreateRefundRequest(ctx, refund); err != nil {
		return nil, fmt.Errorf("failed to create refund request: %w", err)
	}

	return refund, nil
}

func (s *refundService) GetRefundByID(ctx context.Context, id string) (*domain.RefundRequest, error) {
	refund, err := s.repo.GetRefundRequestByID(ctx, id)
	if err != nil {
		return nil, ErrRefundNotFound
	}
	return refund, nil
}

func (s *refundService) ListRefunds(ctx context.Context, filters repository.RefundFilters, paginator *pagination.Paginator) ([]*domain.RefundRequest, error) {
	return s.repo.ListRefundRequests(ctx, filters, paginator)
}

func (s *refundService) ApproveRefund(ctx context.Context, id string, reviewerID, note string) error {
	refund, err := s.repo.GetRefundRequestByID(ctx, id)
	if err != nil {
		return ErrRefundNotFound
	}

	if !refund.CanReview() {
		return ErrRefundNotReviewable
	}

	refund.Approve(reviewerID, note)

	if err := s.repo.UpdateRefundRequest(ctx, refund); err != nil {
		return fmt.Errorf("failed to approve refund: %w", err)
	}

	return nil
}

func (s *refundService) RejectRefund(ctx context.Context, id string, reviewerID, note string) error {
	refund, err := s.repo.GetRefundRequestByID(ctx, id)
	if err != nil {
		return ErrRefundNotFound
	}

	if !refund.CanReview() {
		return ErrRefundNotReviewable
	}

	refund.Reject(reviewerID, note)

	if err := s.repo.UpdateRefundRequest(ctx, refund); err != nil {
		return fmt.Errorf("failed to reject refund: %w", err)
	}

	return nil
}

func (s *refundService) ProcessRefund(ctx context.Context, id string) error {
	refund, err := s.repo.GetRefundRequestByID(ctx, id)
	if err != nil {
		return ErrRefundNotFound
	}

	if !refund.CanProcess() {
		return ErrRefundNotProcessable
	}

	// Mark as processing
	refund.MarkProcessing()
	if err := s.repo.UpdateRefundRequest(ctx, refund); err != nil {
		return err
	}

	// Get payment for the order
	payment, err := s.paymentService.GetPaymentByOrder(ctx, refund.OrderID)
	if err != nil {
		refund.MarkFailed()
		s.repo.UpdateRefundRequest(ctx, refund)
		return fmt.Errorf("failed to get payment: %w", err)
	}

	// Process refund through payment provider
	if err := s.paymentService.Refund(ctx, payment.ID, refund.RefundAmount, refund.RefundReason); err != nil {
		refund.MarkFailed()
		s.repo.UpdateRefundRequest(ctx, refund)
		return fmt.Errorf("failed to process refund payment: %w", err)
	}

	// Mark as completed
	refund.MarkCompleted("", "") // IDs will be updated by payment callback
	if err := s.repo.UpdateRefundRequest(ctx, refund); err != nil {
		return err
	}

	// Update order status to refunded
	order, _ := s.orderService.GetByID(ctx, refund.OrderID)
	if order != nil {
		s.repo.UpdateStatus(ctx, order.ID, domain.OrderStatusRefunded)
	}

	return nil
}

func (s *refundService) GetRefundableAmount(ctx context.Context, orderID string) (float64, error) {
	order, err := s.orderService.GetByID(ctx, orderID)
	if err != nil {
		return 0, err
	}

	// Check if order is paid
	if !order.IsPaid() {
		return 0, nil
	}

	// Get existing refund requests
	filters := repository.RefundFilters{
		OrderID: orderID,
	}
	paginator := &pagination.Paginator{Page: 1, PageSize: 1000}

	refunds, err := s.repo.ListRefundRequests(ctx, filters, paginator)
	if err != nil {
		return 0, err
	}

	// Calculate already refunded amount
	var alreadyRefunded float64
	for _, refund := range refunds {
		if refund.Status == domain.RefundStatusCompleted || refund.Status == domain.RefundStatusProcessing {
			alreadyRefunded += refund.RefundAmount
		}
	}

	// Maximum refundable is total amount minus already refunded
	maxRefundable := order.TotalAmount - alreadyRefunded
	if maxRefundable < 0 {
		maxRefundable = 0
	}

	return maxRefundable, nil
}
