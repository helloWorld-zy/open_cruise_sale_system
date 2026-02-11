package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"

	"gorm.io/gorm"
)

// CabinTypeRepository defines the interface for cabin type data operations
type CabinTypeRepository interface {
	// Create creates a new cabin type
	Create(ctx context.Context, cabinType *domain.CabinType) error

	// GetByID retrieves a cabin type by ID
	GetByID(ctx context.Context, id string) (*domain.CabinType, error)

	// GetByCode retrieves a cabin type by code and cruise ID
	GetByCode(ctx context.Context, cruiseID string, code string) (*domain.CabinType, error)

	// List retrieves a paginated list of cabin types
	List(ctx context.Context, filters CabinTypeFilters, paginator *pagination.Paginator) ([]*domain.CabinType, error)

	// Count returns the total count of cabin types matching the filters
	Count(ctx context.Context, filters CabinTypeFilters) (int64, error)

	// ListByCruise retrieves all cabin types for a cruise
	ListByCruise(ctx context.Context, cruiseID string) ([]*domain.CabinType, error)

	// Update updates a cabin type
	Update(ctx context.Context, cabinType *domain.CabinType) error

	// Delete soft deletes a cabin type
	Delete(ctx context.Context, id string) error

	// Restore restores a soft-deleted cabin type
	Restore(ctx context.Context, id string) error

	// UpdateStatus updates the status of a cabin type
	UpdateStatus(ctx context.Context, id string, status string) error

	// UpdateSortWeight updates the sort weight of a cabin type
	UpdateSortWeight(ctx context.Context, id string, sortWeight int) error

	// DeleteByCruise soft deletes all cabin types for a cruise
	DeleteByCruise(ctx context.Context, cruiseID string) error
}

// CabinTypeFilters represents filters for cabin type list queries
type CabinTypeFilters struct {
	CruiseID   string
	Status     string
	MinGuests  int
	MaxGuests  int
	MinArea    float64
	MaxArea    float64
	BedType    string
	FeatureTag string
}

// cabinTypeRepository implements CabinTypeRepository
type cabinTypeRepository struct {
	db *gorm.DB
}

// NewCabinTypeRepository creates a new cabin type repository
func NewCabinTypeRepository(db *gorm.DB) CabinTypeRepository {
	return &cabinTypeRepository{db: db}
}

func (r *cabinTypeRepository) Create(ctx context.Context, cabinType *domain.CabinType) error {
	return r.db.WithContext(ctx).Create(cabinType).Error
}

func (r *cabinTypeRepository) GetByID(ctx context.Context, id string) (*domain.CabinType, error) {
	var cabinType domain.CabinType
	err := r.db.WithContext(ctx).
		Preload("Cruise").
		First(&cabinType, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cabinType, nil
}

func (r *cabinTypeRepository) GetByCode(ctx context.Context, cruiseID string, code string) (*domain.CabinType, error) {
	var cabinType domain.CabinType
	err := r.db.WithContext(ctx).
		Preload("Cruise").
		Where("cruise_id = ? AND code = ?", cruiseID, code).
		First(&cabinType).Error
	if err != nil {
		return nil, err
	}
	return &cabinType, nil
}

func (r *cabinTypeRepository) List(ctx context.Context, filters CabinTypeFilters, paginator *pagination.Paginator) ([]*domain.CabinType, error) {
	query := r.buildListQuery(filters)

	var cabinTypes []*domain.CabinType
	err := query.WithContext(ctx).
		Preload("Cruise").
		Order("sort_weight DESC, created_at DESC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&cabinTypes).Error

	return cabinTypes, err
}

func (r *cabinTypeRepository) Count(ctx context.Context, filters CabinTypeFilters) (int64, error) {
	query := r.buildListQuery(filters)

	var count int64
	err := query.WithContext(ctx).Model(&domain.CabinType{}).Count(&count).Error
	return count, err
}

func (r *cabinTypeRepository) buildListQuery(filters CabinTypeFilters) *gorm.DB {
	query := r.db.Model(&domain.CabinType{})

	if filters.CruiseID != "" {
		query = query.Where("cruise_id = ?", filters.CruiseID)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.MinGuests > 0 {
		query = query.Where("max_guests >= ?", filters.MinGuests)
	}

	if filters.MaxGuests > 0 {
		query = query.Where("standard_guests <= ?", filters.MaxGuests)
	}

	if filters.MinArea > 0 {
		query = query.Where("max_area_sqm >= ? OR (max_area_sqm IS NULL AND min_area_sqm >= ?)", filters.MinArea, filters.MinArea)
	}

	if filters.MaxArea > 0 {
		query = query.Where("min_area_sqm <= ?", filters.MaxArea)
	}

	if filters.BedType != "" {
		query = query.Where("bed_types ILIKE ?", "%"+filters.BedType+"%")
	}

	if filters.FeatureTag != "" {
		query = query.Where("feature_tags::text ILIKE ?", "%"+filters.FeatureTag+"%")
	}

	return query
}

func (r *cabinTypeRepository) ListByCruise(ctx context.Context, cruiseID string) ([]*domain.CabinType, error) {
	var cabinTypes []*domain.CabinType
	err := r.db.WithContext(ctx).
		Where("cruise_id = ? AND status = ?", cruiseID, domain.CabinTypeStatusActive).
		Order("sort_weight DESC, created_at DESC").
		Find(&cabinTypes).Error
	return cabinTypes, err
}

func (r *cabinTypeRepository) Update(ctx context.Context, cabinType *domain.CabinType) error {
	return r.db.WithContext(ctx).Save(cabinType).Error
}

func (r *cabinTypeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.CabinType{}, "id = ?", id).Error
}

func (r *cabinTypeRepository) Restore(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.CabinType{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

func (r *cabinTypeRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.CabinType{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *cabinTypeRepository) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	return r.db.WithContext(ctx).Model(&domain.CabinType{}).
		Where("id = ?", id).
		Update("sort_weight", sortWeight).Error
}

func (r *cabinTypeRepository) DeleteByCruise(ctx context.Context, cruiseID string) error {
	return r.db.WithContext(ctx).
		Where("cruise_id = ?", cruiseID).
		Delete(&domain.CabinType{}).Error
}
