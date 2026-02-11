package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"

	"gorm.io/gorm"
)

// CabinRepository defines the interface for cabin data operations
type CabinRepository interface {
	Create(ctx context.Context, cabin *domain.Cabin) error
	GetByID(ctx context.Context, id string) (*domain.Cabin, error)
	GetByCabinNumber(ctx context.Context, voyageID, cabinNumber string) (*domain.Cabin, error)
	List(ctx context.Context, filters CabinFilters, paginator *pagination.Paginator) ([]*domain.Cabin, error)
	Count(ctx context.Context, filters CabinFilters) (int64, error)
	ListByVoyage(ctx context.Context, voyageID string) ([]*domain.Cabin, error)
	ListByVoyageAndType(ctx context.Context, voyageID, cabinTypeID string) ([]*domain.Cabin, error)
	ListAvailableByVoyage(ctx context.Context, voyageID string) ([]*domain.Cabin, error)
	Update(ctx context.Context, cabin *domain.Cabin) error
	UpdateStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
	BatchCreate(ctx context.Context, cabins []*domain.Cabin) error
}

// CabinFilters represents filters for cabin queries
type CabinFilters struct {
	VoyageID    string
	CabinTypeID string
	Status      string
	DeckNumber  int
	Section     string
}

// cabinRepository implements CabinRepository
type cabinRepository struct {
	db *gorm.DB
}

// NewCabinRepository creates a new cabin repository
func NewCabinRepository(db *gorm.DB) CabinRepository {
	return &cabinRepository{db: db}
}

func (r *cabinRepository) Create(ctx context.Context, cabin *domain.Cabin) error {
	return r.db.WithContext(ctx).Create(cabin).Error
}

func (r *cabinRepository) GetByID(ctx context.Context, id string) (*domain.Cabin, error) {
	var cabin domain.Cabin
	err := r.db.WithContext(ctx).
		Preload("Voyage").
		Preload("Voyage.Route").
		Preload("CabinType").
		First(&cabin, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cabin, nil
}

func (r *cabinRepository) GetByCabinNumber(ctx context.Context, voyageID, cabinNumber string) (*domain.Cabin, error) {
	var cabin domain.Cabin
	err := r.db.WithContext(ctx).
		Preload("Voyage").
		Preload("CabinType").
		First(&cabin, "voyage_id = ? AND cabin_number = ?", voyageID, cabinNumber).Error
	if err != nil {
		return nil, err
	}
	return &cabin, nil
}

func (r *cabinRepository) List(ctx context.Context, filters CabinFilters, paginator *pagination.Paginator) ([]*domain.Cabin, error) {
	query := r.buildQuery(filters)

	var cabins []*domain.Cabin
	err := query.WithContext(ctx).
		Preload("CabinType").
		Order("deck_number ASC, cabin_number ASC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&cabins).Error

	return cabins, err
}

func (r *cabinRepository) Count(ctx context.Context, filters CabinFilters) (int64, error) {
	query := r.buildQuery(filters)

	var count int64
	err := query.WithContext(ctx).Model(&domain.Cabin{}).Count(&count).Error
	return count, err
}

func (r *cabinRepository) buildQuery(filters CabinFilters) *gorm.DB {
	query := r.db.Model(&domain.Cabin{})

	if filters.VoyageID != "" {
		query = query.Where("voyage_id = ?", filters.VoyageID)
	}
	if filters.CabinTypeID != "" {
		query = query.Where("cabin_type_id = ?", filters.CabinTypeID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.DeckNumber > 0 {
		query = query.Where("deck_number = ?", filters.DeckNumber)
	}
	if filters.Section != "" {
		query = query.Where("section = ?", filters.Section)
	}

	return query
}

func (r *cabinRepository) ListByVoyage(ctx context.Context, voyageID string) ([]*domain.Cabin, error) {
	var cabins []*domain.Cabin
	err := r.db.WithContext(ctx).
		Preload("CabinType").
		Where("voyage_id = ?", voyageID).
		Order("deck_number ASC, cabin_number ASC").
		Find(&cabins).Error
	return cabins, err
}

func (r *cabinRepository) ListByVoyageAndType(ctx context.Context, voyageID, cabinTypeID string) ([]*domain.Cabin, error) {
	var cabins []*domain.Cabin
	err := r.db.WithContext(ctx).
		Preload("CabinType").
		Where("voyage_id = ? AND cabin_type_id = ? AND status = ?", voyageID, cabinTypeID, domain.CabinStatusAvailable).
		Order("deck_number ASC, cabin_number ASC").
		Find(&cabins).Error
	return cabins, err
}

func (r *cabinRepository) ListAvailableByVoyage(ctx context.Context, voyageID string) ([]*domain.Cabin, error) {
	var cabins []*domain.Cabin
	err := r.db.WithContext(ctx).
		Preload("CabinType").
		Where("voyage_id = ? AND status = ?", voyageID, domain.CabinStatusAvailable).
		Order("deck_number ASC, cabin_number ASC").
		Find(&cabins).Error
	return cabins, err
}

func (r *cabinRepository) Update(ctx context.Context, cabin *domain.Cabin) error {
	return r.db.WithContext(ctx).Save(cabin).Error
}

func (r *cabinRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Cabin{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *cabinRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Cabin{}, "id = ?", id).Error
}

func (r *cabinRepository) BatchCreate(ctx context.Context, cabins []*domain.Cabin) error {
	return r.db.WithContext(ctx).Create(cabins).Error
}
