package jobs

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"backend/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// OrderTimeoutJob handles order timeout and cleanup
type OrderTimeoutJob struct {
	orderRepo     repository.OrderRepository
	inventoryRepo repository.InventoryRepository
	stateService  service.OrderStateService
	natsConn      *nats.Conn
	ticker        *time.Ticker
	quit          chan bool
}

// NewOrderTimeoutJob creates a new order timeout job
func NewOrderTimeoutJob(
	orderRepo repository.OrderRepository,
	inventoryRepo repository.InventoryRepository,
	stateService service.OrderStateService,
	natsConn *nats.Conn,
) *OrderTimeoutJob {
	return &OrderTimeoutJob{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		stateService:  stateService,
		natsConn:      natsConn,
		quit:          make(chan bool),
	}
}

// Start starts the order timeout job
func (j *OrderTimeoutJob) Start(interval time.Duration) {
	j.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-j.ticker.C:
				j.processExpiredOrders()
			case <-j.quit:
				j.ticker.Stop()
				return
			}
		}
	}()

	log.Println("Order timeout job started")
}

// Stop stops the order timeout job
func (j *OrderTimeoutJob) Stop() {
	close(j.quit)
	log.Println("Order timeout job stopped")
}

// processExpiredOrders processes expired orders
func (j *OrderTimeoutJob) processExpiredOrders() {
	ctx := context.Background()

	// Find expired pending orders

	// Get all pending orders that have expired
	orders, err := j.getExpiredOrders(ctx)
	if err != nil {
		log.Printf("Failed to get expired orders: %v", err)
		return
	}

	log.Printf("Processing %d expired orders", len(orders))

	for _, order := range orders {
		if err := j.cancelExpiredOrder(ctx, order); err != nil {
			log.Printf("Failed to cancel expired order %s: %v", order.ID, err)
		}
	}

	// Publish metrics
	j.publishMetrics(len(orders))
}

// getExpiredOrders gets all expired pending orders
func (j *OrderTimeoutJob) getExpiredOrders(ctx context.Context) ([]*domain.Order, error) {
	paginator := &pagination.Paginator{Page: 1, PageSize: 100}

	filters := repository.OrderFilters{
		Status: domain.OrderStatusPending,
	}

	orders, err := j.orderRepo.List(ctx, filters, paginator)
	if err != nil {
		return nil, err
	}

	var expiredOrders []*domain.Order
	now := time.Now()

	for _, order := range orders {
		if order.IsExpired() || j.isOrderExpired(order, now) {
			expiredOrders = append(expiredOrders, order)
		}
	}

	return expiredOrders, nil
}

// isOrderExpired checks if an order is expired
func (j *OrderTimeoutJob) isOrderExpired(order *domain.Order, now time.Time) bool {
	if order.ExpiresAt == "" {
		// Default expiration: 15 minutes after creation
		createdAt := order.CreatedAt
		return now.Sub(createdAt) > 15*time.Minute
	}

	expiresAt, err := time.Parse(time.RFC3339, order.ExpiresAt)
	if err != nil {
		// If parsing fails, use default expiration
		return now.Sub(order.CreatedAt) > 15*time.Minute
	}

	return now.After(expiresAt)
}

// cancelExpiredOrder cancels an expired order and releases inventory
func (j *OrderTimeoutJob) cancelExpiredOrder(ctx context.Context, order *domain.Order) error {
	log.Printf("Cancelling expired order: %s (OrderNumber: %s)", order.ID, order.OrderNumber)

	// Get order items
	items, err := j.orderRepo.ListOrderItemsByOrder(ctx, order.ID.String())
	if err != nil {
		return fmt.Errorf("failed to get order items: %w", err)
	}

	// Release inventory for each item
	for _, item := range items {
		if err := j.inventoryRepo.UnlockCabin(ctx, item.VoyageID, item.CabinTypeID, 1); err != nil {
			log.Printf("Failed to unlock cabin for order %s, item %s: %v", order.ID, item.ID, err)
			// Continue with other items
		}
	}

	// Update order status to cancelled
	if err := j.orderRepo.UpdateStatus(ctx, order.ID.String(), domain.OrderStatusCancelled); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Publish event
	j.publishOrderCancelled(order)

	log.Printf("Successfully cancelled expired order: %s", order.ID)
	return nil
}

// publishOrderCancelled publishes order cancelled event
func (j *OrderTimeoutJob) publishOrderCancelled(order *domain.Order) {
	if j.natsConn == nil {
		return
	}

	event := map[string]interface{}{
		"type":         "order.cancelled",
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
		"reason":       "expired",
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	data, _ := json.Marshal(event)
	if err := j.natsConn.Publish("order.cancelled", data); err != nil {
		log.Printf("Failed to publish order cancelled event: %v", err)
	}
}

// publishMetrics publishes job metrics
func (j *OrderTimeoutJob) publishMetrics(count int) {
	if j.natsConn == nil {
		return
	}

	metrics := map[string]interface{}{
		"type":           "job.metrics",
		"job_name":       "order_timeout",
		"expired_orders": count,
		"timestamp":      time.Now().Format(time.RFC3339),
	}

	data, _ := json.Marshal(metrics)
	if err := j.natsConn.Publish("job.metrics", data); err != nil {
		log.Printf("Failed to publish job metrics: %v", err)
	}
}

// SubscribeToEvents subscribes to NATS events
func (j *OrderTimeoutJob) SubscribeToEvents() {
	if j.natsConn == nil {
		return
	}

	// Subscribe to payment timeout events
	j.natsConn.Subscribe("payment.timeout", func(msg *nats.Msg) {
		var event struct {
			OrderID string `json:"order_id"`
		}

		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal payment timeout event: %v", err)
			return
		}

		// Cancel the order
		ctx := context.Background()
		order, err := j.orderRepo.GetByID(ctx, event.OrderID)
		if err != nil {
			log.Printf("Failed to get order for timeout: %v", err)
			return
		}

		if err := j.cancelExpiredOrder(ctx, order); err != nil {
			log.Printf("Failed to cancel order on payment timeout: %v", err)
		}
	})

	log.Println("Order timeout job subscribed to NATS events")
}

// RunOnce runs the job once for testing
func (j *OrderTimeoutJob) RunOnce() {
	j.processExpiredOrders()
}
