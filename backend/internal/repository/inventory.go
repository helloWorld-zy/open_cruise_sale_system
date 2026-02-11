package repository

import (
	"backend/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

// InventoryRepository defines the interface for cabin inventory operations
type InventoryRepository interface {
	// GetInventory retrieves inventory for a voyage and cabin type
	GetInventory(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinInventory, error)

	// CreateInventory creates initial inventory record
	CreateInventory(ctx context.Context, inventory *domain.CabinInventory) error

	// LockCabin attempts to lock a cabin for booking with optimistic locking
	LockCabin(ctx context.Context, voyageID, cabinTypeID string, quantity int) error

	// UnlockCabin releases a locked cabin
	UnlockCabin(ctx context.Context, voyageID, cabinTypeID string, quantity int) error

	// ConfirmBooking confirms a locked cabin as booked
	ConfirmBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error

	// CancelBooking releases a booked cabin back to available
	CancelBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error

	// UpdateInventory updates inventory counts directly (for admin)
	UpdateInventory(ctx context.Context, inventory *domain.CabinInventory) error

	// ListInventoryByVoyage lists all inventory for a voyage
	ListInventoryByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinInventory, error)
}

// inventoryRepository implements InventoryRepository
type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetInventory(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinInventory, error) {
	var inventory domain.CabinInventory
	err := r.db.WithContext(ctx).
		Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).
		First(&inventory).Error
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (r *inventoryRepository) CreateInventory(ctx context.Context, inventory *domain.CabinInventory) error {
	return r.db.WithContext(ctx).Create(inventory).Error
}

func (r *inventoryRepository) LockCabin(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var inventory domain.CabinInventory
		err := tx.Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).
			First(&inventory).Error
		if err != nil {
			return err
		}

		// Check if enough cabins available
		if inventory.AvailableCabins < quantity {
			return errors.New("insufficient cabin inventory")
		}

		// Update with optimistic locking
		result := tx.Model(&domain.CabinInventory{}).
			Where("voyage_id = ? AND cabin_type_id = ? AND lock_version = ?", voyageID, cabinTypeID, inventory.LockVersion).
			Updates(map[string]interface{}{
				"available_cabins": inventory.AvailableCabins - quantity,
				"locked_cabins":    inventory.LockedCabins + quantity,
				"lock_version":     inventory.LockVersion + 1,
				"last_updated_at":  gorm.Expr("CURRENT_TIMESTAMP"),
			})

		if result.RowsAffected == 0 {
			return errors.New("concurrent modification detected, please retry")
		}

		return nil
	})
}

func (r *inventoryRepository) UnlockCabin(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.CabinInventory{}).
			Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).
			Updates(map[string]interface{}{
				"available_cabins": gorm.Expr("available_cabins + ?", quantity),
				"locked_cabins":    gorm.Expr("locked_cabins - ?", quantity),
				"lock_version":     gorm.Expr("lock_version + 1"),
				"last_updated_at":  gorm.Expr("CURRENT_TIMESTAMP"),
			})

		if result.RowsAffected == 0 {
			return errors.New("inventory record not found")
		}

		return nil
	})
}

func (r *inventoryRepository) ConfirmBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.CabinInventory{}).
			Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).
			Updates(map[string]interface{}{
				"locked_cabins":   gorm.Expr("locked_cabins - ?", quantity),
				"booked_cabins":   gorm.Expr("booked_cabins + ?", quantity),
				"lock_version":    gorm.Expr("lock_version + 1"),
				"last_updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
			})

		if result.RowsAffected == 0 {
			return errors.New("inventory record not found")
		}

		return nil
	})
}

func (r *inventoryRepository) CancelBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.CabinInventory{}).
			Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).
			Updates(map[string]interface{}{
				"available_cabins": gorm.Expr("available_cabins + ?", quantity),
				"booked_cabins":    gorm.Expr("booked_cabins - ?", quantity),
				"lock_version":     gorm.Expr("lock_version + 1"),
				"last_updated_at":  gorm.Expr("CURRENT_TIMESTAMP"),
			})

		if result.RowsAffected == 0 {
			return errors.New("inventory record not found")
		}

		return nil
	})
}

func (r *inventoryRepository) UpdateInventory(ctx context.Context, inventory *domain.CabinInventory) error {
	return r.db.WithContext(ctx).Save(inventory).Error
}

func (r *inventoryRepository) ListInventoryByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinInventory, error) {
	var inventories []*domain.CabinInventory
	err := r.db.WithContext(ctx).
		Preload("CabinType").
		Where("voyage_id = ?", voyageID).
		Find(&inventories).Error
	return inventories, err
}
