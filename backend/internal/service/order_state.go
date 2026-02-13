package service

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"
)

var (
	ErrInvalidOrderTransition = errors.New("invalid order status transition")
	ErrOrderNotCancellable    = errors.New("order cannot be cancelled")
	ErrOrderAlreadyPaid       = errors.New("order is already paid")
	ErrOrderNotPaid           = errors.New("order is not paid")
	ErrOrderExpired           = errors.New("order has expired")
)

// OrderStateMachine defines order state transitions
type OrderStateMachine struct {
	currentState string
}

// NewOrderStateMachine creates a new order state machine
func NewOrderStateMachine(initialState string) *OrderStateMachine {
	if initialState == "" {
		initialState = domain.OrderStatusPending
	}
	return &OrderStateMachine{
		currentState: initialState,
	}
}

// CurrentState returns the current state
func (sm *OrderStateMachine) CurrentState() string {
	return sm.currentState
}

// CanTransition checks if a transition is valid
func (sm *OrderStateMachine) CanTransition(toState string) bool {
	validTransitions := sm.getValidTransitions()
	allowedStates, exists := validTransitions[sm.currentState]
	if !exists {
		return false
	}

	for _, state := range allowedStates {
		if state == toState {
			return true
		}
	}
	return false
}

// Transition performs a state transition
func (sm *OrderStateMachine) Transition(toState string) error {
	if !sm.CanTransition(toState) {
		return fmt.Errorf("cannot transition from %s to %s: %w", sm.currentState, toState, ErrInvalidOrderTransition)
	}
	sm.currentState = toState
	return nil
}

// getValidTransitions returns the state transition map
func (sm *OrderStateMachine) getValidTransitions() map[string][]string {
	return map[string][]string{
		domain.OrderStatusPending: {
			domain.OrderStatusPaid,
			domain.OrderStatusCancelled,
		},
		domain.OrderStatusPaid: {
			domain.OrderStatusConfirmed,
			domain.OrderStatusCancelled,
			domain.OrderStatusRefunded,
		},
		domain.OrderStatusConfirmed: {
			domain.OrderStatusCompleted,
			domain.OrderStatusCancelled,
			domain.OrderStatusRefunded,
		},
		domain.OrderStatusCompleted: {
			domain.OrderStatusRefunded,
		},
		domain.OrderStatusCancelled: {},
		domain.OrderStatusRefunded:  {},
	}
}

// OrderStateService defines the interface for order state management
type OrderStateService interface {
	// CanPay checks if order can be paid
	CanPay(order *domain.Order) error

	// CanConfirm checks if order can be confirmed
	CanConfirm(order *domain.Order) error

	// CanComplete checks if order can be completed
	CanComplete(order *domain.Order) error

	// CanCancel checks if order can be cancelled
	CanCancel(order *domain.Order) error

	// CanRefund checks if order can be refunded
	CanRefund(order *domain.Order) error

	// TransitionToPaid transitions order to paid status
	TransitionToPaid(ctx context.Context, order *domain.Order) error

	// TransitionToConfirmed transitions order to confirmed status
	TransitionToConfirmed(ctx context.Context, order *domain.Order) error

	// TransitionToCompleted transitions order to completed status
	TransitionToCompleted(ctx context.Context, order *domain.Order) error

	// TransitionToCancelled transitions order to cancelled status
	TransitionToCancelled(ctx context.Context, order *domain.Order) error

	// TransitionToRefunded transitions order to refunded status
	TransitionToRefunded(ctx context.Context, order *domain.Order) error
}

// orderStateService implements OrderStateService
type orderStateService struct {
	orderRepo     repository.OrderRepository
	inventoryRepo repository.InventoryRepository
}

// NewOrderStateService creates a new order state service
func NewOrderStateService(
	orderRepo repository.OrderRepository,
	inventoryRepo repository.InventoryRepository,
) OrderStateService {
	return &orderStateService{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *orderStateService) CanPay(order *domain.Order) error {
	sm := NewOrderStateMachine(order.Status)
	if !sm.CanTransition(domain.OrderStatusPaid) {
		return fmt.Errorf("order cannot be paid from status %s: %w", order.Status, ErrInvalidOrderTransition)
	}

	if order.IsExpired() {
		return ErrOrderExpired
	}

	if order.PaymentStatus == domain.PaymentStatusPaid {
		return ErrOrderAlreadyPaid
	}

	return nil
}

func (s *orderStateService) CanConfirm(order *domain.Order) error {
	sm := NewOrderStateMachine(order.Status)
	if !sm.CanTransition(domain.OrderStatusConfirmed) {
		return fmt.Errorf("order cannot be confirmed from status %s: %w", order.Status, ErrInvalidOrderTransition)
	}

	if !order.IsPaid() {
		return ErrOrderNotPaid
	}

	return nil
}

func (s *orderStateService) CanComplete(order *domain.Order) error {
	sm := NewOrderStateMachine(order.Status)
	if !sm.CanTransition(domain.OrderStatusCompleted) {
		return fmt.Errorf("order cannot be completed from status %s: %w", order.Status, ErrInvalidOrderTransition)
	}

	if order.Status != domain.OrderStatusConfirmed {
		return fmt.Errorf("order must be confirmed before completion: %w", ErrInvalidOrderTransition)
	}

	return nil
}

func (s *orderStateService) CanCancel(order *domain.Order) error {
	sm := NewOrderStateMachine(order.Status)
	if !sm.CanTransition(domain.OrderStatusCancelled) {
		return fmt.Errorf("order cannot be cancelled from status %s: %w", order.Status, ErrOrderNotCancellable)
	}

	return nil
}

func (s *orderStateService) CanRefund(order *domain.Order) error {
	sm := NewOrderStateMachine(order.Status)
	if !sm.CanTransition(domain.OrderStatusRefunded) {
		return fmt.Errorf("order cannot be refunded from status %s: %w", order.Status, ErrInvalidOrderTransition)
	}

	if order.Status != domain.OrderStatusPaid && order.Status != domain.OrderStatusConfirmed {
		return fmt.Errorf("only paid or confirmed orders can be refunded: %w", ErrInvalidOrderTransition)
	}

	return nil
}

func (s *orderStateService) TransitionToPaid(ctx context.Context, order *domain.Order) error {
	if err := s.CanPay(order); err != nil {
		return err
	}

	sm := NewOrderStateMachine(order.Status)
	if err := sm.Transition(domain.OrderStatusPaid); err != nil {
		return err
	}

	order.Status = sm.CurrentState()
	return s.orderRepo.UpdateStatus(ctx, order.ID.String(), order.Status)
}

func (s *orderStateService) TransitionToConfirmed(ctx context.Context, order *domain.Order) error {
	if err := s.CanConfirm(order); err != nil {
		return err
	}

	sm := NewOrderStateMachine(order.Status)
	if err := sm.Transition(domain.OrderStatusConfirmed); err != nil {
		return err
	}

	order.Status = sm.CurrentState()

	// Confirm inventory bookings for each order item
	items, err := s.orderRepo.ListOrderItemsByOrder(ctx, order.ID.String())
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := s.inventoryRepo.ConfirmBooking(ctx, item.VoyageID, item.CabinTypeID, 1); err != nil {
			return fmt.Errorf("failed to confirm booking for cabin %s: %w", item.CabinTypeID, err)
		}
	}

	return s.orderRepo.UpdateStatus(ctx, order.ID.String(), order.Status)
}

func (s *orderStateService) TransitionToCompleted(ctx context.Context, order *domain.Order) error {
	if err := s.CanComplete(order); err != nil {
		return err
	}

	sm := NewOrderStateMachine(order.Status)
	if err := sm.Transition(domain.OrderStatusCompleted); err != nil {
		return err
	}

	order.Status = sm.CurrentState()
	return s.orderRepo.UpdateStatus(ctx, order.ID.String(), order.Status)
}

func (s *orderStateService) TransitionToCancelled(ctx context.Context, order *domain.Order) error {
	if err := s.CanCancel(order); err != nil {
		return err
	}

	// Release locked/confirmed inventory
	items, err := s.orderRepo.ListOrderItemsByOrder(ctx, order.ID.String())
	if err != nil {
		return err
	}

	for _, item := range items {
		// Determine if we need to unlock or cancel booking based on current status
		if order.Status == domain.OrderStatusPaid || order.Status == domain.OrderStatusConfirmed {
			if err := s.inventoryRepo.CancelBooking(ctx, item.VoyageID, item.CabinTypeID, 1); err != nil {
				// Log error but continue with cancellation
				log.Printf("[WARN] Failed to cancel booking for cabin %s: %v", item.CabinTypeID, err)
			}
		} else if order.Status == domain.OrderStatusPending {
			if err := s.inventoryRepo.UnlockCabin(ctx, item.VoyageID, item.CabinTypeID, 1); err != nil {
				log.Printf("[WARN] Failed to unlock cabin %s: %v", item.CabinTypeID, err)
			}
		}
	}

	sm := NewOrderStateMachine(order.Status)
	if err := sm.Transition(domain.OrderStatusCancelled); err != nil {
		return err
	}

	order.Status = sm.CurrentState()
	return s.orderRepo.UpdateStatus(ctx, order.ID.String(), order.Status)
}

func (s *orderStateService) TransitionToRefunded(ctx context.Context, order *domain.Order) error {
	if err := s.CanRefund(order); err != nil {
		return err
	}

	sm := NewOrderStateMachine(order.Status)
	if err := sm.Transition(domain.OrderStatusRefunded); err != nil {
		return err
	}

	order.Status = sm.CurrentState()
	return s.orderRepo.UpdateStatus(ctx, order.ID.String(), order.Status)
}
