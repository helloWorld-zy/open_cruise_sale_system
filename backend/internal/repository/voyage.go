package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"

	"gorm.io/gorm"
)

// VoyageRepository defines the interface for voyage data operations
type VoyageRepository interface {
	Create(ctx context.Context, voyage *domain.Voyage) error
	GetByID(ctx context.Context, id string) (*domain.Voyage, error)
	GetByVoyageNumber(ctx context.Context, voyageNumber string) (*domain.Voyage, error)
	List(ctx context.Context, filters VoyageFilters, paginator *pagination.Paginator) ([]*domain.Voyage, error)
	Count(ctx context.Context, filters VoyageFilters) (int64, error)
	ListByCruise(ctx context.Context, cruiseID string) ([]*domain.Voyage, error)
	ListByRoute(ctx context.Context, routeID string) ([]*domain.Voyage, error)
	Update(ctx context.Context, voyage *domain.Voyage) error
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateBookingStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
}

// VoyageFilters represents filters for voyage queries
type VoyageFilters struct {
	CruiseID      string
	RouteID       string
	Status        string
	BookingStatus string
	DepartureFrom string
	DepartureTo   string
}

// voyageRepository implements VoyageRepository
type voyageRepository struct {
	db *gorm.DB
}

// NewVoyageRepository creates a new voyage repository
func NewVoyageRepository(db *gorm.DB) VoyageRepository {
	return &voyageRepository{db: db}
}

func (r *voyageRepository) Create(ctx context.Context, voyage *domain.Voyage) error {
	return r.db.WithContext(ctx).Create(voyage).Error
}

func (r *voyageRepository) GetByID(ctx context.Context, id string) (*domain.Voyage, error) {
	var voyage domain.Voyage
	err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Cruise").
		Preload("Cabins").
		Preload("Prices").
		First(&voyage, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &voyage, nil
}

func (r *voyageRepository) GetByVoyageNumber(ctx context.Context, voyageNumber string) (*domain.Voyage, error) {
	var voyage domain.Voyage
	err := r.db.WithContext(ctx).
		Preload("Route").
		Preload("Cruise").
		Preload("Prices").
		First(&voyage, "voyage_number = ?", voyageNumber).Error
	if err != nil {
		return nil, err
	}
	return &voyage, nil
}

func (r *voyageRepository) List(ctx context.Context, filters VoyageFilters, paginator *pagination.Paginator) ([]*domain.Voyage, error) {
	query := r.buildQuery(filters)

	var voyages []*domain.Voyage
	err := query.WithContext(ctx).
		Preload("Route").
		Preload("Cruise").
		Order("departure_date ASC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&voyages).Error

	return voyages, err
}

func (r *voyageRepository) Count(ctx context.Context, filters VoyageFilters) (int64, error) {
	query := r.buildQuery(filters)

	var count int64
	err := query.WithContext(ctx).Model(&domain.Voyage{}).Count(&count).Error
	return count, err
}

func (r *voyageRepository) buildQuery(filters VoyageFilters) *gorm.DB {
	query := r.db.Model(&domain.Voyage{})

	if filters.CruiseID != "" {
		query = query.Where("cruise_id = ?", filters.CruiseID)
	}
	if filters.RouteID != "" {
		query = query.Where("route_id = ?", filters.RouteID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.BookingStatus != "" {
		query = query.Where("booking_status = ?", filters.BookingStatus)
	}
	if filters.DepartureFrom != "" {
		query = query.Where("departure_date >= ?", filters.DepartureFrom)
	}
	if filters.DepartureTo != "" {
		query = query.Where("departure_date <= ?", filters.DepartureTo)
	}

	return query
}

func (r *voyageRepository) ListByCruise(ctx context.Context, cruiseID string) ([]*domain.Voyage, error) {
	var voyages []*domain.Voyage
	err := r.db.WithContext(ctx).
		Where("cruise_id = ? AND status = ? AND booking_status = ?", cruiseID, domain.VoyageStatusScheduled, domain.BookingStatusOpen).
		Order("departure_date ASC").
		Find(&voyages).Error
	return voyages, err
}

func (r *voyageRepository) ListByRoute(ctx context.Context, routeID string) ([]*domain.Voyage, error) {
	var voyages []*domain.Voyage
	err := r.db.WithContext(ctx).
		Where("route_id = ?", routeID).
		Order("departure_date ASC").
		Find(&voyages).Error
	return voyages, err
}

func (r *voyageRepository) Update(ctx context.Context, voyage *domain.Voyage) error {
	return r.db.WithContext(ctx).Save(voyage).Error
}

func (r *voyageRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Voyage{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *voyageRepository) UpdateBookingStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Voyage{}).
		Where("id = ?", id).
		Update("booking_status", status).Error
}

func (r *voyageRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Voyage{}, "id = ?", id).Error
}
