package service

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrFacilityNotFound    = errors.New("facility not found")
	ErrInvalidFacilityData = errors.New("invalid facility data")
)

// FacilityService defines the interface for facility business logic
type FacilityService interface {
	// Create creates a new facility
	Create(ctx context.Context, req CreateFacilityRequest) (*domain.Facility, error)

	// GetByID retrieves a facility by ID
	GetByID(ctx context.Context, id string) (*domain.Facility, error)

	// List retrieves a paginated list of facilities
	List(ctx context.Context, req ListFacilitiesRequest) (*pagination.Result, error)

	// ListByCruise retrieves all facilities for a cruise
	ListByCruise(ctx context.Context, cruiseID string) ([]*domain.Facility, error)

	// ListByCruiseGrouped retrieves facilities grouped by category
	ListByCruiseGrouped(ctx context.Context, cruiseID string) ([]*domain.FacilityCategory, error)

	// Update updates a facility
	Update(ctx context.Context, id string, req UpdateFacilityRequest) (*domain.Facility, error)

	// Delete soft deletes a facility
	Delete(ctx context.Context, id string) error

	// Restore restores a soft-deleted facility
	Restore(ctx context.Context, id string) error

	// UpdateStatus updates the status of a facility
	UpdateStatus(ctx context.Context, id string, status string) error

	// UpdateSortWeight updates the sort weight of a facility
	UpdateSortWeight(ctx context.Context, id string, sortWeight int) error
}

// FacilityCategoryService defines the interface for facility category business logic
type FacilityCategoryService interface {
	// Create creates a new facility category
	Create(ctx context.Context, req CreateFacilityCategoryRequest) (*domain.FacilityCategory, error)

	// GetByID retrieves a category by ID
	GetByID(ctx context.Context, id string) (*domain.FacilityCategory, error)

	// ListByCruise retrieves all categories for a cruise
	ListByCruise(ctx context.Context, cruiseID string) ([]*domain.FacilityCategory, error)

	// Update updates a category
	Update(ctx context.Context, id string, req UpdateFacilityCategoryRequest) (*domain.FacilityCategory, error)

	// Delete deletes a category
	Delete(ctx context.Context, id string) error

	// UpdateSortWeight updates the sort weight of a category
	UpdateSortWeight(ctx context.Context, id string, sortWeight int) error
}

// CreateFacilityRequest represents a request to create a facility
type CreateFacilityRequest struct {
	CruiseID     string   `json:"cruise_id" validate:"required"`
	CategoryID   string   `json:"category_id"`
	Name         string   `json:"name" validate:"required,max=100"`
	DeckNumber   int      `json:"deck_number"`
	OpenTime     string   `json:"open_time"`
	IsFree       bool     `json:"is_free"`
	Price        float64  `json:"price"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
	SuitableTags []string `json:"suitable_tags"`
	SortWeight   int      `json:"sort_weight"`
	Status       string   `json:"status"`
}

// UpdateFacilityRequest represents a request to update a facility
type UpdateFacilityRequest struct {
	CategoryID   string   `json:"category_id"`
	Name         string   `json:"name" validate:"max=100"`
	DeckNumber   int      `json:"deck_number"`
	OpenTime     string   `json:"open_time"`
	IsFree       *bool    `json:"is_free"`
	Price        float64  `json:"price"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
	SuitableTags []string `json:"suitable_tags"`
	SortWeight   int      `json:"sort_weight"`
	Status       string   `json:"status"`
}

// ListFacilitiesRequest represents a request to list facilities
type ListFacilitiesRequest struct {
	CruiseID   string `form:"cruise_id"`
	CategoryID string `form:"category_id"`
	Status     string `form:"status"`
	IsFree     *bool  `form:"is_free"`
	DeckNumber int    `form:"deck_number"`
	pagination.Paginator
}

// CreateFacilityCategoryRequest represents a request to create a facility category
type CreateFacilityCategoryRequest struct {
	CruiseID   string `json:"cruise_id" validate:"required"`
	Name       string `json:"name" validate:"required,max=100"`
	Icon       string `json:"icon"`
	SortWeight int    `json:"sort_weight"`
}

// UpdateFacilityCategoryRequest represents a request to update a facility category
type UpdateFacilityCategoryRequest struct {
	Name       string `json:"name" validate:"max=100"`
	Icon       string `json:"icon"`
	SortWeight int    `json:"sort_weight"`
}

// facilityService implements FacilityService
type facilityService struct {
	facilityRepo repository.FacilityRepository
	categoryRepo repository.FacilityCategoryRepository
}

// facilityCategoryService implements FacilityCategoryService
type facilityCategoryService struct {
	repo repository.FacilityCategoryRepository
}

// NewFacilityService creates a new facility service
func NewFacilityService(facilityRepo repository.FacilityRepository, categoryRepo repository.FacilityCategoryRepository) FacilityService {
	return &facilityService{
		facilityRepo: facilityRepo,
		categoryRepo: categoryRepo,
	}
}

// NewFacilityCategoryService creates a new facility category service
func NewFacilityCategoryService(repo repository.FacilityCategoryRepository) FacilityCategoryService {
	return &facilityCategoryService{repo: repo}
}

// Facility Service Implementation

func (s *facilityService) Create(ctx context.Context, req CreateFacilityRequest) (*domain.Facility, error) {
	if req.Status != "" && !isValidFacilityStatus(req.Status) {
		return nil, ErrInvalidFacilityData
	}

	if req.Status == "" {
		req.Status = domain.FacilityStatusVisible
	}

	facility := &domain.Facility{
		BaseModel: domain.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		CruiseID:    req.CruiseID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		DeckNumber:  req.DeckNumber,
		OpenTime:    req.OpenTime,
		IsFree:      req.IsFree,
		Price:       req.Price,
		Description: req.Description,
		SortWeight:  req.SortWeight,
		Status:      req.Status,
	}

	if err := s.facilityRepo.Create(ctx, facility); err != nil {
		return nil, err
	}

	return facility, nil
}

func (s *facilityService) GetByID(ctx context.Context, id string) (*domain.Facility, error) {
	facility, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFacilityNotFound
		}
		return nil, err
	}
	return facility, nil
}

func (s *facilityService) List(ctx context.Context, req ListFacilitiesRequest) (*pagination.Result, error) {
	filters := repository.FacilityFilters{
		CruiseID:   req.CruiseID,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		IsFree:     req.IsFree,
		DeckNumber: req.DeckNumber,
	}

	count, err := s.facilityRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	paginator := &req.Paginator
	paginator.SetTotal(count)

	facilities, err := s.facilityRepo.List(ctx, filters, paginator)
	if err != nil {
		return nil, err
	}

	return &pagination.Result{
		Data:       facilities,
		Pagination: *paginator,
	}, nil
}

func (s *facilityService) ListByCruise(ctx context.Context, cruiseID string) ([]*domain.Facility, error) {
	return s.facilityRepo.ListByCruise(ctx, cruiseID)
}

func (s *facilityService) ListByCruiseGrouped(ctx context.Context, cruiseID string) ([]*domain.FacilityCategory, error) {
	return s.facilityRepo.ListByCruiseGrouped(ctx, cruiseID)
}

func (s *facilityService) Update(ctx context.Context, id string, req UpdateFacilityRequest) (*domain.Facility, error) {
	facility, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFacilityNotFound
		}
		return nil, err
	}

	if req.CategoryID != "" {
		facility.CategoryID = req.CategoryID
	}
	if req.Name != "" {
		facility.Name = req.Name
	}
	if req.DeckNumber > 0 {
		facility.DeckNumber = req.DeckNumber
	}
	if req.OpenTime != "" {
		facility.OpenTime = req.OpenTime
	}
	if req.IsFree != nil {
		facility.IsFree = *req.IsFree
	}
	if req.Price > 0 {
		facility.Price = req.Price
	}
	if req.Description != "" {
		facility.Description = req.Description
	}
	if req.Status != "" {
		if !isValidFacilityStatus(req.Status) {
			return nil, ErrInvalidFacilityData
		}
		facility.Status = req.Status
	}
	if req.SortWeight != 0 {
		facility.SortWeight = req.SortWeight
	}

	facility.UpdatedAt = time.Now().UTC()

	if err := s.facilityRepo.Update(ctx, facility); err != nil {
		return nil, err
	}

	return facility, nil
}

func (s *facilityService) Delete(ctx context.Context, id string) error {
	_, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFacilityNotFound
		}
		return err
	}

	return s.facilityRepo.Delete(ctx, id)
}

func (s *facilityService) Restore(ctx context.Context, id string) error {
	return s.facilityRepo.Restore(ctx, id)
}

func (s *facilityService) UpdateStatus(ctx context.Context, id string, status string) error {
	if !isValidFacilityStatus(status) {
		return ErrInvalidFacilityData
	}

	_, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFacilityNotFound
		}
		return err
	}

	return s.facilityRepo.UpdateStatus(ctx, id, status)
}

func (s *facilityService) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	_, err := s.facilityRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFacilityNotFound
		}
		return err
	}

	return s.facilityRepo.UpdateSortWeight(ctx, id, sortWeight)
}

func isValidFacilityStatus(status string) bool {
	switch status {
	case domain.FacilityStatusVisible, domain.FacilityStatusHidden:
		return true
	}
	return false
}

// Facility Category Service Implementation

func (s *facilityCategoryService) Create(ctx context.Context, req CreateFacilityCategoryRequest) (*domain.FacilityCategory, error) {
	category := &domain.FacilityCategory{
		BaseModel: domain.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		CruiseID:   req.CruiseID,
		Name:       req.Name,
		Icon:       req.Icon,
		SortWeight: req.SortWeight,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *facilityCategoryService) GetByID(ctx context.Context, id string) (*domain.FacilityCategory, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFacilityNotFound
		}
		return nil, err
	}
	return category, nil
}

func (s *facilityCategoryService) ListByCruise(ctx context.Context, cruiseID string) ([]*domain.FacilityCategory, error) {
	return s.repo.ListByCruise(ctx, cruiseID)
}

func (s *facilityCategoryService) Update(ctx context.Context, id string, req UpdateFacilityCategoryRequest) (*domain.FacilityCategory, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFacilityNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.SortWeight != 0 {
		category.SortWeight = req.SortWeight
	}

	category.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *facilityCategoryService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFacilityNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *facilityCategoryService) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFacilityNotFound
		}
		return err
	}

	return s.repo.UpdateSortWeight(ctx, id, sortWeight)
}
