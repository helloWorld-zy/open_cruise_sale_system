package service

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	ErrCruiseNotFound    = errors.New("cruise not found")
	ErrInvalidCruiseData = errors.New("invalid cruise data")
)

// CruiseService defines the interface for cruise business logic
type CruiseService interface {
	// Create creates a new cruise
	Create(ctx context.Context, req CreateCruiseRequest) (*domain.Cruise, error)

	// GetByID retrieves a cruise by ID
	GetByID(ctx context.Context, id string) (*domain.Cruise, error)

	// GetByCode retrieves a cruise by code
	GetByCode(ctx context.Context, code string) (*domain.Cruise, error)

	// List retrieves a paginated list of cruises
	List(ctx context.Context, req ListCruisesRequest) (*pagination.Result, error)

	// Update updates a cruise
	Update(ctx context.Context, id string, req UpdateCruiseRequest) (*domain.Cruise, error)

	// Delete soft deletes a cruise
	Delete(ctx context.Context, id string) error

	// Restore restores a soft-deleted cruise
	Restore(ctx context.Context, id string) error

	// UpdateStatus updates the status of a cruise
	UpdateStatus(ctx context.Context, id string, status string) error

	// UpdateSortWeight updates the sort weight of a cruise
	UpdateSortWeight(ctx context.Context, id string, sortWeight int) error
}

// CreateCruiseRequest represents a request to create a cruise
type CreateCruiseRequest struct {
	CompanyID         string   `json:"company_id" validate:"required"`
	NameCN            string   `json:"name_cn" validate:"required,max=100"`
	NameEN            string   `json:"name_en" validate:"max=100"`
	Code              string   `json:"code" validate:"required,max=50"`
	GrossTonnage      int      `json:"gross_tonnage"`
	PassengerCapacity int      `json:"passenger_capacity"`
	CrewCount         int      `json:"crew_count"`
	BuiltYear         int      `json:"built_year"`
	RenovatedYear     int      `json:"renovated_year"`
	LengthMeters      float64  `json:"length_meters"`
	WidthMeters       float64  `json:"width_meters"`
	DeckCount         int      `json:"deck_count"`
	CoverImages       []string `json:"cover_images"`
	Status            string   `json:"status"`
	SortWeight        int      `json:"sort_weight"`
}

// UpdateCruiseRequest represents a request to update a cruise
type UpdateCruiseRequest struct {
	NameCN            string   `json:"name_cn" validate:"max=100"`
	NameEN            string   `json:"name_en" validate:"max=100"`
	GrossTonnage      int      `json:"gross_tonnage"`
	PassengerCapacity int      `json:"passenger_capacity"`
	CrewCount         int      `json:"crew_count"`
	BuiltYear         int      `json:"built_year"`
	RenovatedYear     int      `json:"renovated_year"`
	LengthMeters      float64  `json:"length_meters"`
	WidthMeters       float64  `json:"width_meters"`
	DeckCount         int      `json:"deck_count"`
	CoverImages       []string `json:"cover_images"`
	Status            string   `json:"status"`
	SortWeight        int      `json:"sort_weight"`
}

// ListCruisesRequest represents a request to list cruises
type ListCruisesRequest struct {
	CompanyID    string `form:"company_id"`
	Status       string `form:"status"`
	Keyword      string `form:"keyword"`
	HasCabinType bool   `form:"has_cabin_type"`
	MinCapacity  int    `form:"min_capacity"`
	pagination.Paginator
}

// cruiseService implements CruiseService
type cruiseService struct {
	repo repository.CruiseRepository
}

// NewCruiseService creates a new cruise service
func NewCruiseService(repo repository.CruiseRepository) CruiseService {
	return &cruiseService{repo: repo}
}

func (s *cruiseService) Create(ctx context.Context, req CreateCruiseRequest) (*domain.Cruise, error) {
	// Validate status
	if req.Status != "" && !isValidCruiseStatus(req.Status) {
		return nil, ErrInvalidCruiseData
	}

	if req.Status == "" {
		req.Status = domain.CruiseStatusActive
	}

	cruise := &domain.Cruise{
		BaseModel: domain.BaseModel{
			ID:        uuid.New().String(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		CompanyID:         req.CompanyID,
		NameCN:            req.NameCN,
		NameEN:            req.NameEN,
		Code:              req.Code,
		GrossTonnage:      req.GrossTonnage,
		PassengerCapacity: req.PassengerCapacity,
		CrewCount:         req.CrewCount,
		BuiltYear:         req.BuiltYear,
		RenovatedYear:     req.RenovatedYear,
		LengthMeters:      req.LengthMeters,
		WidthMeters:       req.WidthMeters,
		DeckCount:         req.DeckCount,
		Status:            req.Status,
		SortWeight:        req.SortWeight,
	}

	if err := s.repo.Create(ctx, cruise); err != nil {
		return nil, err
	}

	return cruise, nil
}

func (s *cruiseService) GetByID(ctx context.Context, id string) (*domain.Cruise, error) {
	cruise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCruiseNotFound
		}
		return nil, err
	}
	return cruise, nil
}

func (s *cruiseService) GetByCode(ctx context.Context, code string) (*domain.Cruise, error) {
	cruise, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCruiseNotFound
		}
		return nil, err
	}
	return cruise, nil
}

func (s *cruiseService) List(ctx context.Context, req ListCruisesRequest) (*pagination.Result, error) {
	filters := repository.CruiseFilters{
		CompanyID:    req.CompanyID,
		Status:       req.Status,
		Keyword:      req.Keyword,
		HasCabinType: req.HasCabinType,
		MinCapacity:  req.MinCapacity,
	}

	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	paginator := &req.Paginator
	paginator.SetTotal(count)

	cruises, err := s.repo.List(ctx, filters, paginator)
	if err != nil {
		return nil, err
	}

	return &pagination.Result{
		Data:       cruises,
		Pagination: *paginator,
	}, nil
}

func (s *cruiseService) Update(ctx context.Context, id string, req UpdateCruiseRequest) (*domain.Cruise, error) {
	cruise, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCruiseNotFound
		}
		return nil, err
	}

	// Update fields if provided
	if req.NameCN != "" {
		cruise.NameCN = req.NameCN
	}
	if req.NameEN != "" {
		cruise.NameEN = req.NameEN
	}
	if req.GrossTonnage > 0 {
		cruise.GrossTonnage = req.GrossTonnage
	}
	if req.PassengerCapacity > 0 {
		cruise.PassengerCapacity = req.PassengerCapacity
	}
	if req.CrewCount > 0 {
		cruise.CrewCount = req.CrewCount
	}
	if req.BuiltYear > 0 {
		cruise.BuiltYear = req.BuiltYear
	}
	if req.RenovatedYear > 0 {
		cruise.RenovatedYear = req.RenovatedYear
	}
	if req.LengthMeters > 0 {
		cruise.LengthMeters = req.LengthMeters
	}
	if req.WidthMeters > 0 {
		cruise.WidthMeters = req.WidthMeters
	}
	if req.DeckCount > 0 {
		cruise.DeckCount = req.DeckCount
	}
	if len(req.CoverImages) > 0 {
		cruise.CoverImages = datatypes.JSON(req.CoverImages)
	}
	if req.Status != "" {
		if !isValidCruiseStatus(req.Status) {
			return nil, ErrInvalidCruiseData
		}
		cruise.Status = req.Status
	}
	if req.SortWeight != 0 {
		cruise.SortWeight = req.SortWeight
	}

	cruise.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, cruise); err != nil {
		return nil, err
	}

	return cruise, nil
}

func (s *cruiseService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCruiseNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *cruiseService) Restore(ctx context.Context, id string) error {
	return s.repo.Restore(ctx, id)
}

func (s *cruiseService) UpdateStatus(ctx context.Context, id string, status string) error {
	if !isValidCruiseStatus(status) {
		return ErrInvalidCruiseData
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCruiseNotFound
		}
		return err
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *cruiseService) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCruiseNotFound
		}
		return err
	}

	return s.repo.UpdateSortWeight(ctx, id, sortWeight)
}

func isValidCruiseStatus(status string) bool {
	switch status {
	case domain.CruiseStatusActive, domain.CruiseStatusInactive, domain.CruiseStatusMaintenance:
		return true
	}
	return false
}
