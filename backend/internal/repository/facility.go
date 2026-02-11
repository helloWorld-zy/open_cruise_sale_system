package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"

	"gorm.io/gorm"
)

// FacilityRepository defines the interface for facility data operations
type FacilityRepository interface {
	// Create creates a new facility
	Create(ctx context.Context, facility *domain.Facility) error

	// GetByID retrieves a facility by ID
	GetByID(ctx context.Context, id string) (*domain.Facility, error)

	// List retrieves a paginated list of facilities
	List(ctx context.Context, filters FacilityFilters, paginator *pagination.Paginator) ([]*domain.Facility, error)

	// Count returns the total count of facilities matching the filters
	Count(ctx context.Context, filters FacilityFilters) (int64, error)

	// ListByCruise retrieves all facilities for a cruise
	ListByCruise(ctx context.Context, cruiseID string) ([]*domain.Facility, error)

	// ListByCruiseGrouped retrieves facilities grouped by category for a cruise
	ListByCruiseGrouped(ctx context.Context, cruiseID string) ([]*domain.FacilityCategory, error)

	// Update updates a facility
	Update(ctx context.Context, facility *domain.Facility) error

	// Delete soft deletes a facility
	Delete(ctx context.Context, id string) error

	// Restore restores a soft-deleted facility
	Restore(ctx context.Context, id string) error

	// UpdateStatus updates the status of a facility
	UpdateStatus(ctx context.Context, id string, status string) error

	// UpdateSortWeight updates the sort weight of a facility
	UpdateSortWeight(ctx context.Context, id string, sortWeight int) error

	// DeleteByCruise soft deletes all facilities for a cruise
	DeleteByCruise(ctx context.Context, cruiseID string) error
}

// FacilityCategoryRepository defines the interface for facility category operations
type FacilityCategoryRepository interface {
	// Create creates a new facility category
	Create(ctx context.Context, category *domain.FacilityCategory) error

	// GetByID retrieves a category by ID
	GetByID(ctx context.Context, id string) (*domain.FacilityCategory, error)

	// ListByCruise retrieves all categories for a cruise
	ListByCruise(ctx context.Context, cruiseID string) ([]*domain.FacilityCategory, error)

	// Update updates a category
	Update(ctx context.Context, category *domain.FacilityCategory) error

	// Delete deletes a category
	Delete(ctx context.Context, id string) error

	// UpdateSortWeight updates the sort weight of a category
	UpdateSortWeight(ctx context.Context, id string, sortWeight int) error

	// DeleteByCruise deletes all categories for a cruise
	DeleteByCruise(ctx context.Context, cruiseID string) error
}

// FacilityFilters represents filters for facility list queries
type FacilityFilters struct {
	CruiseID   string
	CategoryID string
	Status     string
	IsFree     *bool
	DeckNumber int
}

// facilityRepository implements FacilityRepository
type facilityRepository struct {
	db *gorm.DB
}

// facilityCategoryRepository implements FacilityCategoryRepository
type facilityCategoryRepository struct {
	db *gorm.DB
}

// NewFacilityRepository creates a new facility repository
func NewFacilityRepository(db *gorm.DB) FacilityRepository {
	return &facilityRepository{db: db}
}

// NewFacilityCategoryRepository creates a new facility category repository
func NewFacilityCategoryRepository(db *gorm.DB) FacilityCategoryRepository {
	return &facilityCategoryRepository{db: db}
}

// Facility Repository Implementation

func (r *facilityRepository) Create(ctx context.Context, facility *domain.Facility) error {
	return r.db.WithContext(ctx).Create(facility).Error
}

func (r *facilityRepository) GetByID(ctx context.Context, id string) (*domain.Facility, error) {
	var facility domain.Facility
	err := r.db.WithContext(ctx).
		Preload("Cruise").
		Preload("Category").
		First(&facility, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &facility, nil
}

func (r *facilityRepository) List(ctx context.Context, filters FacilityFilters, paginator *pagination.Paginator) ([]*domain.Facility, error) {
	query := r.buildListQuery(filters)

	var facilities []*domain.Facility
	err := query.WithContext(ctx).
		Preload("Cruise").
		Preload("Category").
		Order("sort_weight DESC, created_at DESC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&facilities).Error

	return facilities, err
}

func (r *facilityRepository) Count(ctx context.Context, filters FacilityFilters) (int64, error) {
	query := r.buildListQuery(filters)

	var count int64
	err := query.WithContext(ctx).Model(&domain.Facility{}).Count(&count).Error
	return count, err
}

func (r *facilityRepository) buildListQuery(filters FacilityFilters) *gorm.DB {
	query := r.db.Model(&domain.Facility{})

	if filters.CruiseID != "" {
		query = query.Where("cruise_id = ?", filters.CruiseID)
	}

	if filters.CategoryID != "" {
		query = query.Where("category_id = ?", filters.CategoryID)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.IsFree != nil {
		query = query.Where("is_free = ?", *filters.IsFree)
	}

	if filters.DeckNumber > 0 {
		query = query.Where("deck_number = ?", filters.DeckNumber)
	}

	return query
}

func (r *facilityRepository) ListByCruise(ctx context.Context, cruiseID string) ([]*domain.Facility, error) {
	var facilities []*domain.Facility
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("cruise_id = ? AND status = ?", cruiseID, domain.FacilityStatusVisible).
		Order("sort_weight DESC, created_at DESC").
		Find(&facilities).Error
	return facilities, err
}

func (r *facilityRepository) ListByCruiseGrouped(ctx context.Context, cruiseID string) ([]*domain.FacilityCategory, error) {
	var categories []*domain.FacilityCategory
	err := r.db.WithContext(ctx).
		Preload("Facilities", "status = ?", domain.FacilityStatusVisible).
		Where("cruise_id = ?", cruiseID).
		Order("sort_weight ASC, created_at ASC").
		Find(&categories).Error
	return categories, err
}

func (r *facilityRepository) Update(ctx context.Context, facility *domain.Facility) error {
	return r.db.WithContext(ctx).Save(facility).Error
}

func (r *facilityRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Facility{}, "id = ?", id).Error
}

func (r *facilityRepository) Restore(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.Facility{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

func (r *facilityRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Facility{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *facilityRepository) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	return r.db.WithContext(ctx).Model(&domain.Facility{}).
		Where("id = ?", id).
		Update("sort_weight", sortWeight).Error
}

func (r *facilityRepository) DeleteByCruise(ctx context.Context, cruiseID string) error {
	return r.db.WithContext(ctx).
		Where("cruise_id = ?", cruiseID).
		Delete(&domain.Facility{}).Error
}

// Facility Category Repository Implementation

func (r *facilityCategoryRepository) Create(ctx context.Context, category *domain.FacilityCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *facilityCategoryRepository) GetByID(ctx context.Context, id string) (*domain.FacilityCategory, error) {
	var category domain.FacilityCategory
	err := r.db.WithContext(ctx).
		Preload("Cruise").
		First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *facilityCategoryRepository) ListByCruise(ctx context.Context, cruiseID string) ([]*domain.FacilityCategory, error) {
	var categories []*domain.FacilityCategory
	err := r.db.WithContext(ctx).
		Where("cruise_id = ?", cruiseID).
		Order("sort_weight ASC, created_at ASC").
		Find(&categories).Error
	return categories, err
}

func (r *facilityCategoryRepository) Update(ctx context.Context, category *domain.FacilityCategory) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *facilityCategoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.FacilityCategory{}, "id = ?", id).Error
}

func (r *facilityCategoryRepository) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	return r.db.WithContext(ctx).Model(&domain.FacilityCategory{}).
		Where("id = ?", id).
		Update("sort_weight", sortWeight).Error
}

func (r *facilityCategoryRepository) DeleteByCruise(ctx context.Context, cruiseID string) error {
	return r.db.WithContext(ctx).
		Where("cruise_id = ?", cruiseID).
		Delete(&domain.FacilityCategory{}).Error
}
