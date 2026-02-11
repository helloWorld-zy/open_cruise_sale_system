package repository

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"context"
	"errors"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// User operations
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByWechatOpenID(ctx context.Context, openID string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters UserFilters, paginator *pagination.Paginator) ([]*domain.User, error)
	Count(ctx context.Context, filters UserFilters) (int64, error)

	// FrequentPassenger operations
	CreateFrequentPassenger(ctx context.Context, passenger *domain.FrequentPassenger) error
	GetFrequentPassengerByID(ctx context.Context, id string) (*domain.FrequentPassenger, error)
	ListFrequentPassengersByUser(ctx context.Context, userID string) ([]*domain.FrequentPassenger, error)
	UpdateFrequentPassenger(ctx context.Context, passenger *domain.FrequentPassenger) error
	DeleteFrequentPassenger(ctx context.Context, id string) error
}

// UserFilters represents filters for user queries
type UserFilters struct {
	Phone  string
	Email  string
	Status string
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("FrequentPassengers").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByWechatOpenID(ctx context.Context, openID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Where("wechat_openid = ?", openID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepository) List(ctx context.Context, filters UserFilters, paginator *pagination.Paginator) ([]*domain.User, error) {
	query := r.buildQuery(filters)

	var users []*domain.User
	err := query.WithContext(ctx).
		Order("created_at DESC").
		Offset(paginator.Offset()).
		Limit(paginator.Limit()).
		Find(&users).Error

	return users, err
}

func (r *userRepository) Count(ctx context.Context, filters UserFilters) (int64, error) {
	query := r.buildQuery(filters)

	var count int64
	err := query.WithContext(ctx).Model(&domain.User{}).Count(&count).Error
	return count, err
}

func (r *userRepository) buildQuery(filters UserFilters) *gorm.DB {
	query := r.db.Model(&domain.User{})

	if filters.Phone != "" {
		query = query.Where("phone = ?", filters.Phone)
	}
	if filters.Email != "" {
		query = query.Where("email = ?", filters.Email)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	return query
}

// ==================== FrequentPassenger Operations ====================

func (r *userRepository) CreateFrequentPassenger(ctx context.Context, passenger *domain.FrequentPassenger) error {
	return r.db.WithContext(ctx).Create(passenger).Error
}

func (r *userRepository) GetFrequentPassengerByID(ctx context.Context, id string) (*domain.FrequentPassenger, error) {
	var passenger domain.FrequentPassenger
	err := r.db.WithContext(ctx).
		First(&passenger, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &passenger, nil
}

func (r *userRepository) ListFrequentPassengersByUser(ctx context.Context, userID string) ([]*domain.FrequentPassenger, error) {
	var passengers []*domain.FrequentPassenger
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&passengers).Error
	return passengers, err
}

func (r *userRepository) UpdateFrequentPassenger(ctx context.Context, passenger *domain.FrequentPassenger) error {
	return r.db.WithContext(ctx).Save(passenger).Error
}

func (r *userRepository) DeleteFrequentPassenger(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.FrequentPassenger{}, "id = ?", id).Error
}

// ErrUserNotFound is returned when a user is not found
var ErrUserNotFound = errors.New("user not found")
