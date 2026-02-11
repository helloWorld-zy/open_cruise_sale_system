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
	ErrCabinTypeNotFound    = errors.New("cabin type not found")
	ErrInvalidCabinTypeData = errors.New("invalid cabin type data")
)

// CabinTypeService defines the interface for cabin type business logic
type CabinTypeService interface {
	// Create creates a new cabin type
	Create(ctx context.Context, req CreateCabinTypeRequest) (*domain.CabinType, error)

	// GetByID retrieves a cabin type by ID
	GetByID(ctx context.Context, id string) (*domain.CabinType, error)

	// List retrieves a paginated list of cabin types
	List(ctx context.Context, req ListCabinTypesRequest) (*pagination.Result, error)

	// ListByCruise retrieves all cabin types for a cruise
	ListByCruise(ctx context.Context, cruiseID string) ([]*domain.CabinType, error)

	// Update updates a cabin type
	Update(ctx context.Context, id string, req UpdateCabinTypeRequest) (*domain.CabinType, error)

	// Delete soft deletes a cabin type
	Delete(ctx context.Context, id string) error

	// Restore restores a soft-deleted cabin type
	Restore(ctx context.Context, id string) error

	// UpdateStatus updates the status of a cabin type
	UpdateStatus(ctx context.Context, id string, status string) error

	// UpdateSortWeight updates the sort weight of a cabin type
	UpdateSortWeight(ctx context.Context, id string, sortWeight int) error
}

// CreateCabinTypeRequest represents a request to create a cabin type
type CreateCabinTypeRequest struct {
	CruiseID       string   `json:"cruise_id" validate:"required"`
	Name           string   `json:"name" validate:"required,max=100"`
	Code           string   `json:"code" validate:"required,max=50"`
	MinAreaSqm     float64  `json:"min_area_sqm"`
	MaxAreaSqm     float64  `json:"max_area_sqm"`
	StandardGuests int      `json:"standard_guests"`
	MaxGuests      int      `json:"max_guests"`
	BedTypes       string   `json:"bed_types"`
	FeatureTags    []string `json:"feature_tags"`
	Description    string   `json:"description"`
	Images         []string `json:"images"`
	FloorPlanURL   string   `json:"floor_plan_url"`
	Amenities      []string `json:"amenities"`
	SortWeight     int      `json:"sort_weight"`
	Status         string   `json:"status"`
}

// UpdateCabinTypeRequest represents a request to update a cabin type
type UpdateCabinTypeRequest struct {
	Name           string   `json:"name" validate:"max=100"`
	MinAreaSqm     float64  `json:"min_area_sqm"`
	MaxAreaSqm     float64  `json:"max_area_sqm"`
	StandardGuests int      `json:"standard_guests"`
	MaxGuests      int      `json:"max_guests"`
	BedTypes       string   `json:"bed_types"`
	FeatureTags    []string `json:"feature_tags"`
	Description    string   `json:"description"`
	Images         []string `json:"images"`
	FloorPlanURL   string   `json:"floor_plan_url"`
	Amenities      []string `json:"amenities"`
	SortWeight     int      `json:"sort_weight"`
	Status         string   `json:"status"`
}

// ListCabinTypesRequest represents a request to list cabin types
type ListCabinTypesRequest struct {
	CruiseID   string  `form:"cruise_id"`
	Status     string  `form:"status"`
	MinGuests  int     `form:"min_guests"`
	MaxGuests  int     `form:"max_guests"`
	MinArea    float64 `form:"min_area"`
	MaxArea    float64 `form:"max_area"`
	BedType    string  `form:"bed_type"`
	FeatureTag string  `form:"feature_tag"`
	pagination.Paginator
}

// cabinTypeService implements CabinTypeService
type cabinTypeService struct {
	repo repository.CabinTypeRepository
}

// NewCabinTypeService creates a new cabin type service
func NewCabinTypeService(repo repository.CabinTypeRepository) CabinTypeService {
	return &cabinTypeService{repo: repo}
}

func (s *cabinTypeService) Create(ctx context.Context, req CreateCabinTypeRequest) (*domain.CabinType, error) {
	if req.Status != "" && !isValidCabinTypeStatus(req.Status) {
		return nil, ErrInvalidCabinTypeData
	}

	if req.Status == "" {
		req.Status = domain.CabinTypeStatusActive
	}

	cabinType := &domain.CabinType{
		BaseModel: domain.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		CruiseID:       req.CruiseID,
		Name:           req.Name,
		Code:           req.Code,
		MinAreaSqm:     req.MinAreaSqm,
		MaxAreaSqm:     req.MaxAreaSqm,
		StandardGuests: req.StandardGuests,
		MaxGuests:      req.MaxGuests,
		BedTypes:       req.BedTypes,
		Description:    req.Description,
		FloorPlanURL:   req.FloorPlanURL,
		SortWeight:     req.SortWeight,
		Status:         req.Status,
	}

	if err := s.repo.Create(ctx, cabinType); err != nil {
		return nil, err
	}

	return cabinType, nil
}

func (s *cabinTypeService) GetByID(ctx context.Context, id string) (*domain.CabinType, error) {
	cabinType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCabinTypeNotFound
		}
		return nil, err
	}
	return cabinType, nil
}

func (s *cabinTypeService) List(ctx context.Context, req ListCabinTypesRequest) (*pagination.Result, error) {
	filters := repository.CabinTypeFilters{
		CruiseID:   req.CruiseID,
		Status:     req.Status,
		MinGuests:  req.MinGuests,
		MaxGuests:  req.MaxGuests,
		MinArea:    req.MinArea,
		MaxArea:    req.MaxArea,
		BedType:    req.BedType,
		FeatureTag: req.FeatureTag,
	}

	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	paginator := &req.Paginator
	paginator.SetTotal(count)

	cabinTypes, err := s.repo.List(ctx, filters, paginator)
	if err != nil {
		return nil, err
	}

	return &pagination.Result{
		Data:       cabinTypes,
		Pagination: *paginator,
	}, nil
}

func (s *cabinTypeService) ListByCruise(ctx context.Context, cruiseID string) ([]*domain.CabinType, error) {
	return s.repo.ListByCruise(ctx, cruiseID)
}

func (s *cabinTypeService) Update(ctx context.Context, id string, req UpdateCabinTypeRequest) (*domain.CabinType, error) {
	cabinType, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCabinTypeNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		cabinType.Name = req.Name
	}
	if req.MinAreaSqm > 0 {
		cabinType.MinAreaSqm = req.MinAreaSqm
	}
	if req.MaxAreaSqm > 0 {
		cabinType.MaxAreaSqm = req.MaxAreaSqm
	}
	if req.StandardGuests > 0 {
		cabinType.StandardGuests = req.StandardGuests
	}
	if req.MaxGuests > 0 {
		cabinType.MaxGuests = req.MaxGuests
	}
	if req.BedTypes != "" {
		cabinType.BedTypes = req.BedTypes
	}
	if req.Description != "" {
		cabinType.Description = req.Description
	}
	if req.FloorPlanURL != "" {
		cabinType.FloorPlanURL = req.FloorPlanURL
	}
	if req.Status != "" {
		if !isValidCabinTypeStatus(req.Status) {
			return nil, ErrInvalidCabinTypeData
		}
		cabinType.Status = req.Status
	}
	if req.SortWeight != 0 {
		cabinType.SortWeight = req.SortWeight
	}

	cabinType.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, cabinType); err != nil {
		return nil, err
	}

	return cabinType, nil
}

func (s *cabinTypeService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCabinTypeNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *cabinTypeService) Restore(ctx context.Context, id string) error {
	return s.repo.Restore(ctx, id)
}

func (s *cabinTypeService) UpdateStatus(ctx context.Context, id string, status string) error {
	if !isValidCabinTypeStatus(status) {
		return ErrInvalidCabinTypeData
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCabinTypeNotFound
		}
		return err
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *cabinTypeService) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCabinTypeNotFound
		}
		return err
	}

	return s.repo.UpdateSortWeight(ctx, id, sortWeight)
}

func isValidCabinTypeStatus(status string) bool {
	switch status {
	case domain.CabinTypeStatusActive, domain.CabinTypeStatusInactive:
		return true
	}
	return false
}
