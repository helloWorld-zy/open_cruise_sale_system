package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"

	"gorm.io/gorm"
)

// CruiseRepository defines the interface for cruise data operations
type CruiseRepository interface {
	// Create creates a new cruise
	Create(ctx context.Context, cruise *domain.Cruise) error

	// GetByID retrieves a cruise by ID with its relations
	GetByID(ctx context.Context, id string) (*domain.Cruise, error)

	// GetByCode retrieves a cruise by code
	GetByCode(ctx context.Context, code string) (*domain.Cruise, error)

	// List retrieves a paginated list of cruises
	List(ctx context.Context, filters CruiseFilters, paginator *pagination.Paginator) ([]*domain.Cruise, error)

	// Count returns the total count of cruises matching the filters
	Count(ctx context.Context, filters CruiseFilters) (int64, error)

	// Update updates a cruise
	Update(ctx context.Context, cruise *domain.Cruise) error

	// Delete soft deletes a cruise
	Delete(ctx context.Context, id string) error

	// Restore restores a soft-deleted cruise
	Restore(ctx context.Context, id string) error

	// ListByCompany retrieves cruises by company ID
	ListByCompany(ctx context.Context, companyID string) ([]*domain.Cruise, error)

	// UpdateStatus updates the status of a cruise
	UpdateStatus(ctx context.Context, id string, status string) error

	// UpdateSortWeight updates the sort weight of a cruise
	UpdateSortWeight(ctx context.Context, id string, sortWeight int) error
}

// CruiseFilters represents filters for cruise list queries
type CruiseFilters struct {
	CompanyID    string
	Status       string
	Keyword      string
	HasCabinType bool
	MinCapacity  int
}

// cruiseRepository implements CruiseRepository
type cruiseRepository struct {
	db *gorm.DB
}

// NewCruiseRepository creates a new cruise repository
func NewCruiseRepository(db *gorm.DB) CruiseRepository {
	return &cruiseRepository{db: db}
}

func (r *cruiseRepository) Create(ctx context.Context, cruise *domain.Cruise) error {
	return r.db.WithContext(ctx).Create(cruise).Error
}

func (r *cruiseRepository) GetByID(ctx context.Context, id string) (*domain.Cruise, error) {
	var cruise domain.Cruise
	err := r.db.WithContext(ctx).
		Preload("Company").
		Preload("CabinTypes").
		Preload("Facilities.Category").
		First(&cruise, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cruise, nil
}

func (r *cruiseRepository) GetByCode(ctx context.Context, code string) (*domain.Cruise, error) {
	var cruise domain.Cruise
	err := r.db.WithContext(ctx).
		Preload("Company").
		Preload("CabinTypes").
		Preload("Facilities.Category").
		First(&cruise, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &cruise, nil
}

func (r *cruiseRepository) List(ctx context.Context, filters CruiseFilters, paginator *pagination.Paginator) ([]*domain.Cruise, error) {
	query := r.buildListQuery(filters)

	var cruises []*domain.Cruise
	err := query.WithContext(ctx).
		Preload("Company").
		Order("sort_weight DESC, created_at DESC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&cruises).Error

	return cruises, err
}

func (r *cruiseRepository) Count(ctx context.Context, filters CruiseFilters) (int64, error) {
	query := r.buildListQuery(filters)

	var count int64
	err := query.WithContext(ctx).Model(&domain.Cruise{}).Count(&count).Error
	return count, err
}

func (r *cruiseRepository) buildListQuery(filters CruiseFilters) *gorm.DB {
	query := r.db.Model(&domain.Cruise{})

	if filters.CompanyID != "" {
		query = query.Where("company_id = ?", filters.CompanyID)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.Keyword != "" {
		query = query.Where(
			"name_cn ILIKE ? OR name_en ILIKE ? OR code ILIKE ?",
			"%"+filters.Keyword+"%",
			"%"+filters.Keyword+"%",
			"%"+filters.Keyword+"%",
		)
	}

	if filters.HasCabinType {
		query = query.Where("EXISTS (SELECT 1 FROM cabin_types WHERE cabin_types.cruise_id = cruises.id AND cabin_types.deleted_at IS NULL)")
	}

	if filters.MinCapacity > 0 {
		query = query.Where("passenger_capacity >= ?", filters.MinCapacity)
	}

	return query
}

func (r *cruiseRepository) Update(ctx context.Context, cruise *domain.Cruise) error {
	return r.db.WithContext(ctx).Save(cruise).Error
}

func (r *cruiseRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Cruise{}, "id = ?", id).Error
}

func (r *cruiseRepository) Restore(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.Cruise{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

func (r *cruiseRepository) ListByCompany(ctx context.Context, companyID string) ([]*domain.Cruise, error) {
	var cruises []*domain.Cruise
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("sort_weight DESC, created_at DESC").
		Find(&cruises).Error
	return cruises, err
}

func (r *cruiseRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Cruise{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *cruiseRepository) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	return r.db.WithContext(ctx).Model(&domain.Cruise{}).
		Where("id = ?", id).
		Update("sort_weight", sortWeight).Error
}
