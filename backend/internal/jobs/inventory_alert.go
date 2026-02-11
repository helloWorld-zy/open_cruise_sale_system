package jobs

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"backend/internal/service"
	"context"
	"fmt"
	"log"
	"time"
)

// InventoryAlertConfig holds configuration for inventory alerts
type InventoryAlertConfig struct {
	LowInventoryThreshold int           // Alert when inventory falls below this number
	CheckInterval         time.Duration // How often to check inventory
	AlertCooldown         time.Duration // Minimum time between alerts for same voyage/cabin
}

// InventoryAlert represents an inventory alert record
type InventoryAlert struct {
	VoyageID    uint64
	CabinTypeID uint64
	Remaining   int
	AlertedAt   time.Time
}

// InventoryAlertJob monitors cabin inventory and sends alerts when low
type InventoryAlertJob struct {
	inventoryRepo       repository.InventoryRepository
	cabinRepo           repository.CabinRepository
	voyageRepo          repository.VoyageRepository
	notificationService service.NotificationService
	config              InventoryAlertConfig
	ticker              *time.Ticker
	quit                chan bool
	alertCache          map[string]time.Time // Cache to prevent duplicate alerts
}

// NewInventoryAlertJob creates a new inventory alert job
func NewInventoryAlertJob(
	inventoryRepo repository.InventoryRepository,
	cabinRepo repository.CabinRepository,
	voyageRepo repository.VoyageRepository,
	notificationService service.NotificationService,
	config InventoryAlertConfig,
) *InventoryAlertJob {
	return &InventoryAlertJob{
		inventoryRepo:       inventoryRepo,
		cabinRepo:           cabinRepo,
		voyageRepo:          voyageRepo,
		notificationService: notificationService,
		config:              config,
		quit:                make(chan bool),
		alertCache:          make(map[string]time.Time),
	}
}

// DefaultInventoryAlertConfig returns default configuration
func DefaultInventoryAlertConfig() InventoryAlertConfig {
	return InventoryAlertConfig{
		LowInventoryThreshold: 5,
		CheckInterval:         15 * time.Minute,
		AlertCooldown:         4 * time.Hour,
	}
}

// Start starts the inventory alert job
func (j *InventoryAlertJob) Start() {
	j.ticker = time.NewTicker(j.config.CheckInterval)

	go func() {
		for {
			select {
			case <-j.ticker.C:
				j.checkInventory()
			case <-j.quit:
				j.ticker.Stop()
				return
			}
		}
	}()

	log.Println("Inventory alert job started")
}

// Stop stops the inventory alert job
func (j *InventoryAlertJob) Stop() {
	close(j.quit)
	log.Println("Inventory alert job stopped")
}

// checkInventory checks inventory levels and sends alerts
func (j *InventoryAlertJob) checkInventory() {
	ctx := context.Background()

	// Get active voyages (departing in the future)
	activeVoyages, err := j.getActiveVoyages(ctx)
	if err != nil {
		log.Printf("Failed to get active voyages: %v", err)
		return
	}

	log.Printf("Checking inventory for %d active voyages", len(activeVoyages))

	alertCount := 0
	for _, voyage := range activeVoyages {
		// Get cabin inventory for this voyage
		inventories, err := j.inventoryRepo.GetByVoyageID(ctx, voyage.ID)
		if err != nil {
			log.Printf("Failed to get inventory for voyage %d: %v", voyage.ID, err)
			continue
		}

		// Check each cabin type's inventory
		for _, inv := range inventories {
			// Skip if inventory is not low
			if inv.Remaining > j.config.LowInventoryThreshold {
				continue
			}

			// Check if we should alert (cooldown period)
			alertKey := fmt.Sprintf("%d:%d", voyage.ID, inv.CabinTypeID)
			if !j.shouldAlert(alertKey) {
				continue
			}

			// Send alert
			if err := j.sendInventoryAlert(ctx, voyage, inv); err != nil {
				log.Printf("Failed to send inventory alert: %v", err)
				continue
			}

			// Update cache
			j.alertCache[alertKey] = time.Now()
			alertCount++
		}
	}

	if alertCount > 0 {
		log.Printf("Sent %d inventory alerts", alertCount)
	}

	// Clean up old cache entries
	j.cleanAlertCache()
}

// getActiveVoyages gets voyages that are active and upcoming
func (j *InventoryAlertJob) getActiveVoyages(ctx context.Context) ([]*domain.Voyage, error) {
	now := time.Now()
	// Get voyages that depart in the future and are not completed
	// This is a simplified implementation - in production, you'd use repository methods
	// For now, return empty list - the actual implementation should query the database
	return []*domain.Voyage{}, nil
}

// shouldAlert checks if we should send an alert based on cooldown
func (j *InventoryAlertJob) shouldAlert(alertKey string) bool {
	lastAlert, exists := j.alertCache[alertKey]
	if !exists {
		return true
	}

	return time.Since(lastAlert) >= j.config.AlertCooldown
}

// sendInventoryAlert sends an inventory alert
func (j *InventoryAlertJob) sendInventoryAlert(ctx context.Context, voyage *domain.Voyage, inv *domain.CabinInventory) error {
	// Get cabin type details
	cabinType, err := j.cabinRepo.GetCabinTypeByID(ctx, inv.CabinTypeID)
	if err != nil {
		return fmt.Errorf("failed to get cabin type: %w", err)
	}

	// Send alert via notification service
	if err := j.notificationService.SendInventoryAlertNotification(
		ctx,
		voyage.ID,
		inv.CabinTypeID,
		inv.Remaining,
	); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	log.Printf("Inventory alert sent: voyage=%d, cabin=%s, remaining=%d",
		voyage.ID, cabinType.Name, inv.Remaining)

	return nil
}

// cleanAlertCache removes old cache entries
func (j *InventoryAlertJob) cleanAlertCache() {
	cutoff := time.Now().Add(-24 * time.Hour) // Keep 24 hours of cache
	for key, alertTime := range j.alertCache {
		if alertTime.Before(cutoff) {
			delete(j.alertCache, key)
		}
	}
}

// GetAlertStats returns statistics about current alerts
func (j *InventoryAlertJob) GetAlertStats() map[string]interface{} {
	j.cleanAlertCache()

	return map[string]interface{}{
		"cached_alerts":  len(j.alertCache),
		"threshold":      j.config.LowInventoryThreshold,
		"check_interval": j.config.CheckInterval.String(),
		"alert_cooldown": j.config.AlertCooldown.String(),
	}
}
