package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"

	"gorm.io/gorm"
)

// PriceRepository defines the interface for price data operations
type PriceRepository interface {
	Create(ctx context.Context, price *domain.CabinPrice) error
	GetByID(ctx context.Context, id string) (*domain.CabinPrice, error)
	GetByVoyageAndCabinType(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error)
	List(ctx context.Context, filters PriceFilters, paginator *pagination.Paginator) ([]*domain.CabinPrice, error)
	Count(ctx context.Context, filters PriceFilters) (int64, error)
	ListByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinPrice, error)
	ListByVoyageWithTypes(ctx context.Context, voyageID string, cabinTypeIDs []string) ([]*domain.CabinPrice, error)
	GetCurrentPrice(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error)
	Update(ctx context.Context, price *domain.CabinPrice) error
	UpdatePrice(ctx context.Context, id string, adultPrice, childPrice, infantPrice float64) error
	Delete(ctx context.Context, id string) error
	BatchCreate(ctx context.Context, prices []*domain.CabinPrice) error
}

// PriceFilters represents filters for price queries
type PriceFilters struct {
	VoyageID    string
	CabinTypeID string
	PriceType   string
	IsPromotion *bool
}

// priceRepository implements PriceRepository
type priceRepository struct {
	db *gorm.DB
}

// NewPriceRepository creates a new price repository
func NewPriceRepository(db *gorm.DB) PriceRepository {
	return &priceRepository{db: db}
}

func (r *priceRepository) Create(ctx context.Context, price *domain.CabinPrice) error {
	return r.db.WithContext(ctx).Create(price).Error
}

func (r *priceRepository) GetByID(ctx context.Context, id string) (*domain.CabinPrice, error) {
	var price domain.CabinPrice
	err := r.db.WithContext(ctx).
		Preload("Voyage").
		Preload("CabinType").
		First(&price, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *priceRepository) GetByVoyageAndCabinType(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error) {
	var price domain.CabinPrice
	err := r.db.WithContext(ctx).
		Preload("Voyage").
		Preload("CabinType").
		First(&price, "voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *priceRepository) List(ctx context.Context, filters PriceFilters, paginator *pagination.Paginator) ([]*domain.CabinPrice, error) {
	query := r.buildQuery(filters)

	var prices []*domain.CabinPrice
	err := query.WithContext(ctx).
		Preload("CabinType").
		Order("adult_price ASC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&prices).Error

	return prices, err
}

func (r *priceRepository) Count(ctx context.Context, filters PriceFilters) (int64, error) {
	query := r.buildQuery(filters)

	var count int64
	err := query.WithContext(ctx).Model(&domain.CabinPrice{}).Count(&count).Error
	return count, err
}

func (r *priceRepository) buildQuery(filters PriceFilters) *gorm.DB {
	query := r.db.Model(&domain.CabinPrice{})

	if filters.VoyageID != "" {
		query = query.Where("voyage_id = ?", filters.VoyageID)
	}
	if filters.CabinTypeID != "" {
		query = query.Where("cabin_type_id = ?", filters.CabinTypeID)
	}
	if filters.PriceType != "" {
		query = query.Where("price_type = ?", filters.PriceType)
	}
	if filters.IsPromotion != nil {
		query = query.Where("is_promotion = ?", *filters.IsPromotion)
	}

	return query
}

func (r *priceRepository) ListByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinPrice, error) {
	var prices []*domain.CabinPrice
	err := r.db.WithContext(ctx).
		Preload("CabinType").
		Where("voyage_id = ?", voyageID).
		Order("adult_price ASC").
		Find(&prices).Error
	return prices, err
}

func (r *priceRepository) ListByVoyageWithTypes(ctx context.Context, voyageID string, cabinTypeIDs []string) ([]*domain.CabinPrice, error) {
	var prices []*domain.CabinPrice
	query := r.db.WithContext(ctx).
		Preload("CabinType").
		Where("voyage_id = ?", voyageID)

	if len(cabinTypeIDs) > 0 {
		query = query.Where("cabin_type_id IN ?", cabinTypeIDs)
	}

	err := query.Order("adult_price ASC").Find(&prices).Error
	return prices, err
}

func (r *priceRepository) GetCurrentPrice(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error) {
	var price domain.CabinPrice
	err := r.db.WithContext(ctx).
		Preload("CabinType").
		Where("voyage_id = ? AND cabin_type_id = ?", voyageID, cabinTypeID).
		First(&price).Error
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *priceRepository) Update(ctx context.Context, price *domain.CabinPrice) error {
	return r.db.WithContext(ctx).Save(price).Error
}

func (r *priceRepository) UpdatePrice(ctx context.Context, id string, adultPrice, childPrice, infantPrice float64) error {
	updates := map[string]interface{}{
		"adult_price":  adultPrice,
		"child_price":  childPrice,
		"infant_price": infantPrice,
	}
	return r.db.WithContext(ctx).Model(&domain.CabinPrice{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *priceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.CabinPrice{}, "id = ?", id).Error
}

func (r *priceRepository) BatchCreate(ctx context.Context, prices []*domain.CabinPrice) error {
	return r.db.WithContext(ctx).Create(prices).Error
}
