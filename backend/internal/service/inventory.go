package service

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrInventoryNotFound      = errors.New("inventory not found")
	ErrInsufficientInventory  = errors.New("insufficient cabin inventory")
	ErrConcurrentModification = errors.New("concurrent modification detected, please retry")
	ErrInvalidInventoryData   = errors.New("invalid inventory data")
	ErrCabinAlreadyLocked     = errors.New("cabin is already locked")
)

// InventoryService defines the interface for inventory business logic
type InventoryService interface {
	// InitializeInventory creates initial inventory for a voyage
	InitializeInventory(ctx context.Context, voyageID string, cabinTypeCounts map[string]int) error

	// LockCabins attempts to lock cabins for booking
	LockCabins(ctx context.Context, voyageID, cabinTypeID string, quantity int) error

	// UnlockCabins releases locked cabins
	UnlockCabins(ctx context.Context, voyageID, cabinTypeID string, quantity int) error

	// ConfirmBooking confirms a locked booking
	ConfirmBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error

	// CancelBooking cancels a confirmed booking and releases cabins
	CancelBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error

	// GetInventory retrieves inventory for a voyage and cabin type
	GetInventory(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinInventory, error)

	// ListInventoryByVoyage lists all inventory for a voyage
	ListInventoryByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinInventory, error)

	// CheckAvailability checks if cabins are available
	CheckAvailability(ctx context.Context, voyageID, cabinTypeID string, quantity int) (bool, int, error)

	// UpdateInventory updates inventory (for admin use)
	UpdateInventory(ctx context.Context, inventory *domain.CabinInventory) error
}

// LockRequest represents a request to lock cabins
type LockRequest struct {
	VoyageID    string `json:"voyage_id" validate:"required"`
	CabinTypeID string `json:"cabin_type_id" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required,min=1"`
}

// inventoryService implements InventoryService
type inventoryService struct {
	inventoryRepo repository.InventoryRepository
	voyageRepo    repository.VoyageRepository
	cabinRepo     repository.CabinRepository
}

// NewInventoryService creates a new inventory service
func NewInventoryService(
	inventoryRepo repository.InventoryRepository,
	voyageRepo repository.VoyageRepository,
	cabinRepo repository.CabinRepository,
) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		voyageRepo:    voyageRepo,
		cabinRepo:     cabinRepo,
	}
}

func (s *inventoryService) InitializeInventory(ctx context.Context, voyageID string, cabinTypeCounts map[string]int) error {
	// Verify voyage exists
	_, err := s.voyageRepo.GetByID(ctx, voyageID)
	if err != nil {
		return fmt.Errorf("voyage not found: %w", err)
	}

	for cabinTypeID, count := range cabinTypeCounts {
		if count <= 0 {
			continue
		}

		inventory := &domain.CabinInventory{
			VoyageID:        voyageID,
			CabinTypeID:     cabinTypeID,
			TotalCabins:     count,
			AvailableCabins: count,
			LastUpdatedAt:   time.Now().Format(time.RFC3339),
		}

		if err := s.inventoryRepo.CreateInventory(ctx, inventory); err != nil {
			return fmt.Errorf("failed to create inventory for cabin type %s: %w", cabinTypeID, err)
		}
	}

	return nil
}

func (s *inventoryService) LockCabins(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidInventoryData
	}

	// Check availability first
	available, _, err := s.CheckAvailability(ctx, voyageID, cabinTypeID, quantity)
	if err != nil {
		return err
	}
	if !available {
		return ErrInsufficientInventory
	}

	// Attempt to lock with optimistic locking
	err = s.inventoryRepo.LockCabin(ctx, voyageID, cabinTypeID, quantity)
	if err != nil {
		if err.Error() == "concurrent modification detected, please retry" {
			return ErrConcurrentModification
		}
		return err
	}

	return nil
}

func (s *inventoryService) UnlockCabins(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidInventoryData
	}

	err := s.inventoryRepo.UnlockCabin(ctx, voyageID, cabinTypeID, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (s *inventoryService) ConfirmBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidInventoryData
	}

	err := s.inventoryRepo.ConfirmBooking(ctx, voyageID, cabinTypeID, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (s *inventoryService) CancelBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidInventoryData
	}

	err := s.inventoryRepo.CancelBooking(ctx, voyageID, cabinTypeID, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (s *inventoryService) GetInventory(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinInventory, error) {
	inventory, err := s.inventoryRepo.GetInventory(ctx, voyageID, cabinTypeID)
	if err != nil {
		if errors.Is(err, errors.New("record not found")) {
			return nil, ErrInventoryNotFound
		}
		return nil, err
	}
	return inventory, nil
}

func (s *inventoryService) ListInventoryByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinInventory, error) {
	return s.inventoryRepo.ListInventoryByVoyage(ctx, voyageID)
}

func (s *inventoryService) CheckAvailability(ctx context.Context, voyageID, cabinTypeID string, quantity int) (bool, int, error) {
	inventory, err := s.inventoryRepo.GetInventory(ctx, voyageID, cabinTypeID)
	if err != nil {
		return false, 0, err
	}

	available := inventory.AvailableCabins >= quantity
	return available, inventory.AvailableCabins, nil
}

func (s *inventoryService) UpdateInventory(ctx context.Context, inventory *domain.CabinInventory) error {
	// Validate inventory counts
	if inventory.TotalCabins < 0 || inventory.AvailableCabins < 0 || inventory.LockedCabins < 0 || inventory.BookedCabins < 0 {
		return ErrInvalidInventoryData
	}

	inventory.LastUpdatedAt = time.Now().Format(time.RFC3339)
	return s.inventoryRepo.UpdateInventory(ctx, inventory)
}
