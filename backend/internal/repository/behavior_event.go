package repository

import (
	"backend/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// BehaviorEventRepository defines the interface for behavior event operations
type BehaviorEventRepository interface {
	Create(ctx context.Context, event *domain.BehaviorEvent) error
	CreateBatch(ctx context.Context, events []*domain.BehaviorEvent) error
	GetUserProfile(ctx context.Context, userID uint64) (*domain.UserProfile, error)
	GetPopularCruises(ctx context.Context, limit int) ([]domain.PopularCruise, error)
	GetTrendingSearches(ctx context.Context, limit int) ([]domain.TrendingSearch, error)
	FindSimilarCruises(ctx context.Context, profile *domain.UserProfile, limit int) ([]uint64, error)
	GetUserEvents(ctx context.Context, userID uint64, since time.Time) ([]*domain.BehaviorEvent, error)
	SaveUserProfile(ctx context.Context, profile *domain.UserProfile) error
}

type behaviorEventRepository struct {
	db *gorm.DB
}

func NewBehaviorEventRepository(db *gorm.DB) BehaviorEventRepository {
	return &behaviorEventRepository{db: db}
}

func (r *behaviorEventRepository) Create(ctx context.Context, event *domain.BehaviorEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *behaviorEventRepository) CreateBatch(ctx context.Context, events []*domain.BehaviorEvent) error {
	return r.db.WithContext(ctx).CreateInBatches(events, 100).Error
}

func (r *behaviorEventRepository) GetUserProfile(ctx context.Context, userID uint64) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *behaviorEventRepository) GetPopularCruises(ctx context.Context, limit int) ([]domain.PopularCruise, error) {
	// For MVP, just return empty list or mock query
	// Complex aggregation query needed here
	return []domain.PopularCruise{}, nil
}

func (r *behaviorEventRepository) GetTrendingSearches(ctx context.Context, limit int) ([]domain.TrendingSearch, error) {
	return []domain.TrendingSearch{}, nil
}

func (r *behaviorEventRepository) FindSimilarCruises(ctx context.Context, profile *domain.UserProfile, limit int) ([]uint64, error) {
	return []uint64{}, nil
}

func (r *behaviorEventRepository) GetUserEvents(ctx context.Context, userID uint64, since time.Time) ([]*domain.BehaviorEvent, error) {
	var events []*domain.BehaviorEvent
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Find(&events).Error
	return events, err
}

func (r *behaviorEventRepository) SaveUserProfile(ctx context.Context, profile *domain.UserProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}
