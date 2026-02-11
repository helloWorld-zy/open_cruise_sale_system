package service

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrPriceNotFound    = errors.New("price not found")
	ErrInvalidPriceData = errors.New("invalid price data")
	ErrPriceExists      = errors.New("price already exists for this voyage and cabin type")
)

// PriceService defines the interface for price business logic
type PriceService interface {
	// Create creates a new price
	Create(ctx context.Context, req CreatePriceRequest) (*domain.CabinPrice, error)

	// GetByID retrieves a price by ID
	GetByID(ctx context.Context, id string) (*domain.CabinPrice, error)

	// GetByVoyageAndCabinType retrieves price for a voyage and cabin type
	GetByVoyageAndCabinType(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error)

	// List retrieves a paginated list of prices
	List(ctx context.Context, req ListPricesRequest) (*pagination.Result, error)

	// ListByVoyage retrieves all prices for a voyage
	ListByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinPrice, error)

	// GetCurrentPrice retrieves the effective price (considering promotions)
	GetCurrentPrice(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error)

	// Update updates a price
	Update(ctx context.Context, id string, req UpdatePriceRequest) (*domain.CabinPrice, error)

	// Delete deletes a price
	Delete(ctx context.Context, id string) error

	// CalculatePrice calculates the total price for passengers
	CalculatePrice(ctx context.Context, priceID string, adults, children, infants int) (PriceCalculation, error)

	// BatchCreate creates multiple prices at once
	BatchCreate(ctx context.Context, voyageID string, prices []CreatePriceRequest) error
}

// PriceCalculation represents the result of a price calculation
type PriceCalculation struct {
	AdultPrice       float64 `json:"adult_price"`
	ChildPrice       float64 `json:"child_price"`
	InfantPrice      float64 `json:"infant_price"`
	PortFee          float64 `json:"port_fee"`
	ServiceFee       float64 `json:"service_fee"`
	SingleSupplement float64 `json:"single_supplement,omitempty"`
	Subtotal         float64 `json:"subtotal"`
	Total            float64 `json:"total"`
}

// CreatePriceRequest represents a request to create a price
type CreatePriceRequest struct {
	CabinTypeID        string  `json:"cabin_type_id" validate:"required"`
	PriceType          string  `json:"price_type" validate:"required"`
	AdultPrice         float64 `json:"adult_price" validate:"required,gt=0"`
	ChildPrice         float64 `json:"child_price" validate:"gte=0"`
	InfantPrice        float64 `json:"infant_price" validate:"gte=0"`
	SingleSupplement   float64 `json:"single_supplement" validate:"gte=0"`
	PortFee            float64 `json:"port_fee" validate:"gte=0"`
	ServiceFee         float64 `json:"service_fee" validate:"gte=0"`
	IsPromotion        bool    `json:"is_promotion"`
	PromotionStartDate string  `json:"promotion_start_date,omitempty"`
	PromotionEndDate   string  `json:"promotion_end_date,omitempty"`
	MinPassengers      int     `json:"min_passengers" validate:"gte=1"`
	MaxPassengers      int     `json:"max_passengers" validate:"gte=1"`
}

// UpdatePriceRequest represents a request to update a price
type UpdatePriceRequest struct {
	PriceType          string  `json:"price_type,omitempty"`
	AdultPrice         float64 `json:"adult_price,omitempty" validate:"omitempty,gt=0"`
	ChildPrice         float64 `json:"child_price,omitempty" validate:"omitempty,gte=0"`
	InfantPrice        float64 `json:"infant_price,omitempty" validate:"omitempty,gte=0"`
	SingleSupplement   float64 `json:"single_supplement,omitempty" validate:"omitempty,gte=0"`
	PortFee            float64 `json:"port_fee,omitempty" validate:"omitempty,gte=0"`
	ServiceFee         float64 `json:"service_fee,omitempty" validate:"omitempty,gte=0"`
	IsPromotion        *bool   `json:"is_promotion,omitempty"`
	PromotionStartDate string  `json:"promotion_start_date,omitempty"`
	PromotionEndDate   string  `json:"promotion_end_date,omitempty"`
	MinPassengers      int     `json:"min_passengers,omitempty" validate:"omitempty,gte=1"`
	MaxPassengers      int     `json:"max_passengers,omitempty" validate:"omitempty,gte=1"`
}

// ListPricesRequest represents a request to list prices
type ListPricesRequest struct {
	VoyageID    string `form:"voyage_id"`
	CabinTypeID string `form:"cabin_type_id"`
	PriceType   string `form:"price_type"`
	IsPromotion *bool  `form:"is_promotion"`
	pagination.Paginator
}

// priceService implements PriceService
type priceService struct {
	priceRepo     repository.PriceRepository
	voyageRepo    repository.VoyageRepository
	cabinTypeRepo repository.CabinTypeRepository
}

// NewPriceService creates a new price service
func NewPriceService(
	priceRepo repository.PriceRepository,
	voyageRepo repository.VoyageRepository,
	cabinTypeRepo repository.CabinTypeRepository,
) PriceService {
	return &priceService{
		priceRepo:     priceRepo,
		voyageRepo:    voyageRepo,
		cabinTypeRepo: cabinTypeRepo,
	}
}

func (s *priceService) Create(ctx context.Context, req CreatePriceRequest) (*domain.CabinPrice, error) {
	// Check if price already exists for this voyage and cabin type
	// This would need voyageID from context or request
	return nil, fmt.Errorf("voyage ID is required")
}

func (s *priceService) CreateWithVoyage(ctx context.Context, voyageID string, req CreatePriceRequest) (*domain.CabinPrice, error) {
	// Verify voyage exists
	_, err := s.voyageRepo.GetByID(ctx, voyageID)
	if err != nil {
		return nil, fmt.Errorf("voyage not found: %w", err)
	}

	// Verify cabin type exists
	_, err = s.cabinTypeRepo.GetByID(ctx, req.CabinTypeID)
	if err != nil {
		return nil, fmt.Errorf("cabin type not found: %w", err)
	}

	// Check if price already exists
	existing, _ := s.priceRepo.GetByVoyageAndCabinType(ctx, voyageID, req.CabinTypeID)
	if existing != nil {
		return nil, ErrPriceExists
	}

	// Validate price type
	if req.PriceType == "" {
		req.PriceType = domain.PriceTypeStandard
	}
	if !isValidPriceType(req.PriceType) {
		return nil, ErrInvalidPriceData
	}

	price := &domain.CabinPrice{
		VoyageID:           voyageID,
		CabinTypeID:        req.CabinTypeID,
		PriceType:          req.PriceType,
		AdultPrice:         req.AdultPrice,
		ChildPrice:         req.ChildPrice,
		InfantPrice:        req.InfantPrice,
		SingleSupplement:   req.SingleSupplement,
		PortFee:            req.PortFee,
		ServiceFee:         req.ServiceFee,
		IsPromotion:        req.IsPromotion,
		PromotionStartDate: req.PromotionStartDate,
		PromotionEndDate:   req.PromotionEndDate,
		MinPassengers:      req.MinPassengers,
		MaxPassengers:      req.MaxPassengers,
	}

	if price.MinPassengers == 0 {
		price.MinPassengers = 1
	}
	if price.MaxPassengers == 0 {
		price.MaxPassengers = 4
	}

	if err := s.priceRepo.Create(ctx, price); err != nil {
		return nil, err
	}

	return price, nil
}

func (s *priceService) GetByID(ctx context.Context, id string) (*domain.CabinPrice, error) {
	price, err := s.priceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrPriceNotFound
	}
	return price, nil
}

func (s *priceService) GetByVoyageAndCabinType(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error) {
	price, err := s.priceRepo.GetByVoyageAndCabinType(ctx, voyageID, cabinTypeID)
	if err != nil {
		return nil, ErrPriceNotFound
	}
	return price, nil
}

func (s *priceService) List(ctx context.Context, req ListPricesRequest) (*pagination.Result, error) {
	filters := repository.PriceFilters{
		VoyageID:    req.VoyageID,
		CabinTypeID: req.CabinTypeID,
		PriceType:   req.PriceType,
		IsPromotion: req.IsPromotion,
	}

	count, err := s.priceRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	paginator := &req.Paginator
	paginator.SetTotal(count)

	prices, err := s.priceRepo.List(ctx, filters, paginator)
	if err != nil {
		return nil, err
	}

	return &pagination.Result{
		Data:       prices,
		Pagination: *paginator,
	}, nil
}

func (s *priceService) ListByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinPrice, error) {
	return s.priceRepo.ListByVoyage(ctx, voyageID)
}

func (s *priceService) GetCurrentPrice(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error) {
	price, err := s.priceRepo.GetCurrentPrice(ctx, voyageID, cabinTypeID)
	if err != nil {
		return nil, ErrPriceNotFound
	}

	// Check if promotion is active
	if price.IsPromotion && price.PromotionStartDate != "" && price.PromotionEndDate != "" {
		now := time.Now()
		startDate, _ := time.Parse(time.RFC3339, price.PromotionStartDate)
		endDate, _ := time.Parse(time.RFC3339, price.PromotionEndDate)

		if now.Before(startDate) || now.After(endDate) {
			// Promotion is not active, find standard price
			return nil, ErrPriceNotFound
		}
	}

	return price, nil
}

func (s *priceService) Update(ctx context.Context, id string, req UpdatePriceRequest) (*domain.CabinPrice, error) {
	price, err := s.priceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrPriceNotFound
	}

	// Update fields if provided
	if req.PriceType != "" {
		if !isValidPriceType(req.PriceType) {
			return nil, ErrInvalidPriceData
		}
		price.PriceType = req.PriceType
	}
	if req.AdultPrice > 0 {
		price.AdultPrice = req.AdultPrice
	}
	if req.ChildPrice >= 0 {
		price.ChildPrice = req.ChildPrice
	}
	if req.InfantPrice >= 0 {
		price.InfantPrice = req.InfantPrice
	}
	if req.SingleSupplement >= 0 {
		price.SingleSupplement = req.SingleSupplement
	}
	if req.PortFee >= 0 {
		price.PortFee = req.PortFee
	}
	if req.ServiceFee >= 0 {
		price.ServiceFee = req.ServiceFee
	}
	if req.IsPromotion != nil {
		price.IsPromotion = *req.IsPromotion
	}
	if req.PromotionStartDate != "" {
		price.PromotionStartDate = req.PromotionStartDate
	}
	if req.PromotionEndDate != "" {
		price.PromotionEndDate = req.PromotionEndDate
	}
	if req.MinPassengers > 0 {
		price.MinPassengers = req.MinPassengers
	}
	if req.MaxPassengers > 0 {
		price.MaxPassengers = req.MaxPassengers
	}

	if err := s.priceRepo.Update(ctx, price); err != nil {
		return nil, err
	}

	return price, nil
}

func (s *priceService) Delete(ctx context.Context, id string) error {
	_, err := s.priceRepo.GetByID(ctx, id)
	if err != nil {
		return ErrPriceNotFound
	}

	return s.priceRepo.Delete(ctx, id)
}

func (s *priceService) CalculatePrice(ctx context.Context, priceID string, adults, children, infants int) (PriceCalculation, error) {
	price, err := s.priceRepo.GetByID(ctx, priceID)
	if err != nil {
		return PriceCalculation{}, ErrPriceNotFound
	}

	totalPassengers := adults + children
	adultTotal := price.AdultPrice * float64(adults)
	childTotal := price.ChildPrice * float64(children)
	infantTotal := price.InfantPrice * float64(infants)
	portFeeTotal := price.PortFee * float64(totalPassengers)
	serviceFeeTotal := price.ServiceFee * float64(totalPassengers)

	subtotal := adultTotal + childTotal + infantTotal
	total := subtotal + portFeeTotal + serviceFeeTotal

	return PriceCalculation{
		AdultPrice:  adultTotal,
		ChildPrice:  childTotal,
		InfantPrice: infantTotal,
		PortFee:     portFeeTotal,
		ServiceFee:  serviceFeeTotal,
		Subtotal:    subtotal,
		Total:       total,
	}, nil
}

func (s *priceService) BatchCreate(ctx context.Context, voyageID string, prices []CreatePriceRequest) error {
	for _, req := range prices {
		_, err := s.CreateWithVoyage(ctx, voyageID, req)
		if err != nil && err != ErrPriceExists {
			return err
		}
	}
	return nil
}

func isValidPriceType(priceType string) bool {
	switch priceType {
	case domain.PriceTypeStandard, domain.PriceTypeEarlyBird, domain.PriceTypeLastMinute, domain.PriceTypeGroup:
		return true
	}
	return false
}
