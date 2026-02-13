package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ReviewFilters defines filters for listing reviews
type ReviewFilters struct {
	CruiseID   *string
	VoyageID   *string
	UserID     *string
	Status     string
	IsVerified *bool
	MinRating  *int
	MaxRating  *int
}

// ReviewRepository defines the interface for review data operations
type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetByID(ctx context.Context, id string) (*domain.Review, error)
	GetByOrderID(ctx context.Context, orderID string) (*domain.Review, error)
	List(ctx context.Context, filters ReviewFilters, paginator *pagination.Paginator) (*pagination.Result, error)
	Update(ctx context.Context, review *domain.Review) error
	Delete(ctx context.Context, id string) error
	HasMarkedHelpful(ctx context.Context, reviewID string, userID string) bool
	GetStats(ctx context.Context, cruiseID *string, voyageID *string) (*domain.ReviewStats, error)
}

// reviewRepository implements ReviewRepository
type reviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *domain.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *reviewRepository) GetByID(ctx context.Context, id string) (*domain.Review, error) {
	var review domain.Review
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // Return error to let service handle it (it expects ErrReviewNotFound usually, but gorm error is fine)
		}
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Review, error) {
	var review domain.Review
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil for not found by order ID?
		}
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) List(ctx context.Context, filters ReviewFilters, paginator *pagination.Paginator) (*pagination.Result, error) {
	// Stub implementation
	return nil, errors.New("not implemented")
}

func (r *reviewRepository) Update(ctx context.Context, review *domain.Review) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *reviewRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Review{}).Error
}

func (r *reviewRepository) HasMarkedHelpful(ctx context.Context, reviewID string, userID string) bool {
	// Stub
	return false
}

func (r *reviewRepository) GetStats(ctx context.Context, cruiseID *string, voyageID *string) (*domain.ReviewStats, error) {
	// Stub
	return nil, errors.New("not implemented")
}
