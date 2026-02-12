package service

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/datatypes"
)

var (
	ErrReviewNotFound     = errors.New("review not found")
	ErrInvalidReviewData  = errors.New("invalid review data")
	ErrDuplicateReview    = errors.New("duplicate review for this order")
	ErrOrderNotCompleted  = errors.New("order must be completed before reviewing")
	ErrUnauthorizedReview = errors.New("unauthorized to create this review")
)

// ReviewService defines the interface for review business logic
type ReviewService interface {
	// Create creates a new review
	Create(ctx context.Context, req CreateReviewRequest) (*domain.Review, error)

	// GetByID retrieves a review by ID
	GetByID(ctx context.Context, id uint64) (*domain.Review, error)

	// List retrieves reviews with filters
	List(ctx context.Context, req ListReviewsRequest) (*pagination.Result, error)

	// ListByUser retrieves reviews by a specific user
	ListByUser(ctx context.Context, userID uint64, paginator *pagination.Paginator) (*pagination.Result, error)

	// ListByCruise retrieves reviews for a cruise
	ListByCruise(ctx context.Context, cruiseID string, paginator *pagination.Paginator) (*pagination.Result, error)

	// ListByVoyage retrieves reviews for a voyage
	ListByVoyage(ctx context.Context, voyageID string, paginator *pagination.Paginator) (*pagination.Result, error)

	// Update updates a review
	Update(ctx context.Context, id uint64, req UpdateReviewRequest) (*domain.Review, error)

	// Delete deletes a review
	Delete(ctx context.Context, id uint64) error

	// Helpful marks a review as helpful
	Helpful(ctx context.Context, reviewID uint64, userID uint64) error

	// Unhelpful removes helpful mark
	Unhelpful(ctx context.Context, reviewID uint64, userID uint64) error

	// Reply adds admin reply to review
	Reply(ctx context.Context, reviewID uint64, reply string) error

	// Approve approves a pending review
	Approve(ctx context.Context, reviewID uint64) error

	// Reject rejects a pending review
	Reject(ctx context.Context, reviewID uint64, reason string) error

	// GetReviewStats gets review statistics for a cruise or voyage
	GetReviewStats(ctx context.Context, cruiseID *string, voyageID *string) (*ReviewStats, error)
}

// CreateReviewRequest represents a request to create a review
type CreateReviewRequest struct {
	UserID      uint64   `json:"user_id" validate:"required"`
	VoyageID    *uint64  `json:"voyage_id,omitempty"`
	CruiseID    *uint64  `json:"cruise_id,omitempty"`
	OrderID     *uint64  `json:"order_id,omitempty"`
	Rating      int      `json:"rating" validate:"required,min=1,max=5"`
	Title       string   `json:"title" validate:"required,max=200"`
	Content     string   `json:"content" validate:"required,min=10,max=5000"`
	Pros        []string `json:"pros,omitempty"`
	Cons        []string `json:"cons,omitempty"`
	Images      []string `json:"images,omitempty"`
	IsAnonymous bool     `json:"is_anonymous"`
}

// UpdateReviewRequest represents a request to update a review
type UpdateReviewRequest struct {
	Rating  *int     `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Title   *string  `json:"title,omitempty" validate:"omitempty,max=200"`
	Content *string  `json:"content,omitempty" validate:"omitempty,min=10,max=5000"`
	Pros    []string `json:"pros,omitempty"`
	Cons    []string `json:"cons,omitempty"`
	Images  []string `json:"images,omitempty"`
}

// ListReviewsRequest represents filters for listing reviews
type ListReviewsRequest struct {
	CruiseID   *string
	VoyageID   *string
	UserID     *uint64
	Status     string
	IsVerified *bool
	MinRating  *int
	MaxRating  *int
	SortBy     string // rating, date, helpful
	SortOrder  string // asc, desc
	*pagination.Paginator
}

// ReviewStats represents review statistics
type ReviewStats struct {
	TotalReviews       int         `json:"total_reviews"`
	AverageRating      float64     `json:"average_rating"`
	FiveStarCount      int         `json:"five_star_count"`
	FourStarCount      int         `json:"four_star_count"`
	ThreeStarCount     int         `json:"three_star_count"`
	TwoStarCount       int         `json:"two_star_count"`
	OneStarCount       int         `json:"one_star_count"`
	RatingDistribution map[int]int `json:"rating_distribution"`
	VerifiedCount      int         `json:"verified_count"`
}

// reviewService implements ReviewService
type reviewService struct {
	reviewRepo repository.ReviewRepository
	orderRepo  repository.OrderRepository
	userRepo   repository.UserRepository
}

// NewReviewService creates a new review service
func NewReviewService(
	reviewRepo repository.ReviewRepository,
	orderRepo repository.OrderRepository,
	userRepo repository.UserRepository,
) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		orderRepo:  orderRepo,
		userRepo:   userRepo,
	}
}

// Create creates a new review
func (s *reviewService) Create(ctx context.Context, req CreateReviewRequest) (*domain.Review, error) {
	// Verify order if provided
	if req.OrderID != nil {
		order, err := s.orderRepo.GetByID(ctx, fmt.Sprintf("%d", *req.OrderID))
		if err != nil {
			return nil, ErrOrderNotFound
		}

		// Check if user owns this order
		if order.UserID != fmt.Sprintf("%d", req.UserID) {
			return nil, ErrUnauthorizedReview
		}

		// Check if order is completed
		if order.Status != domain.OrderStatusCompleted {
			return nil, ErrOrderNotCompleted
		}

		// Check for duplicate review
		existing, _ := s.reviewRepo.GetByOrderID(ctx, *req.OrderID)
		if existing != nil {
			return nil, ErrDuplicateReview
		}

		// Set verified purchase
		req.CruiseID = &order.CruiseID
	}

	// Marshal arrays to JSON
	prosJSON, _ := json.Marshal(req.Pros)
	consJSON, _ := json.Marshal(req.Cons)
	imagesJSON, _ := json.Marshal(req.Images)

	review := &domain.Review{
		UserID:      req.UserID,
		VoyageID:    req.VoyageID,
		CruiseID:    req.CruiseID,
		OrderID:     req.OrderID,
		Rating:      req.Rating,
		Title:       req.Title,
		Content:     req.Content,
		Pros:        datatypes.JSON(prosJSON),
		Cons:        datatypes.JSON(consJSON),
		Images:      datatypes.JSON(imagesJSON),
		IsVerified:  req.OrderID != nil,
		IsAnonymous: req.IsAnonymous,
		Status:      domain.ReviewStatusPending,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

// GetByID retrieves a review by ID
func (s *reviewService) GetByID(ctx context.Context, id uint64) (*domain.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrReviewNotFound
	}
	return review, nil
}

// List retrieves reviews with filters
func (s *reviewService) List(ctx context.Context, req ListReviewsRequest) (*pagination.Result, error) {
	filters := repository.ReviewFilters{
		CruiseID:   req.CruiseID,
		VoyageID:   req.VoyageID,
		UserID:     req.UserID,
		Status:     req.Status,
		IsVerified: req.IsVerified,
		MinRating:  req.MinRating,
		MaxRating:  req.MaxRating,
	}

	return s.reviewRepo.List(ctx, filters, req.Paginator)
}

// ListByUser retrieves reviews by a specific user
func (s *reviewService) ListByUser(ctx context.Context, userID uint64, paginator *pagination.Paginator) (*pagination.Result, error) {
	filters := repository.ReviewFilters{
		UserID: &userID,
		Status: domain.ReviewStatusApproved,
	}
	return s.reviewRepo.List(ctx, filters, paginator)
}

// ListByCruise retrieves reviews for a cruise
func (s *reviewService) ListByCruise(ctx context.Context, cruiseID string, paginator *pagination.Paginator) (*pagination.Result, error) {
	filters := repository.ReviewFilters{
		CruiseID: &cruiseID,
		Status:   domain.ReviewStatusApproved,
	}
	return s.reviewRepo.List(ctx, filters, paginator)
}

// ListByVoyage retrieves reviews for a voyage
func (s *reviewService) ListByVoyage(ctx context.Context, voyageID string, paginator *pagination.Paginator) (*pagination.Result, error) {
	filters := repository.ReviewFilters{
		VoyageID: &voyageID,
		Status:   domain.ReviewStatusApproved,
	}
	return s.reviewRepo.List(ctx, filters, paginator)
}

// Update updates a review
func (s *reviewService) Update(ctx context.Context, id uint64, req UpdateReviewRequest) (*domain.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrReviewNotFound
	}

	// Update fields
	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Title != nil {
		review.Title = *req.Title
	}
	if req.Content != nil {
		review.Content = *req.Content
	}
	if req.Pros != nil {
		prosJSON, _ := json.Marshal(req.Pros)
		review.Pros = datatypes.JSON(prosJSON)
	}
	if req.Cons != nil {
		consJSON, _ := json.Marshal(req.Cons)
		review.Cons = datatypes.JSON(consJSON)
	}
	if req.Images != nil {
		imagesJSON, _ := json.Marshal(req.Images)
		review.Images = datatypes.JSON(imagesJSON)
	}

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	return review, nil
}

// Delete deletes a review
func (s *reviewService) Delete(ctx context.Context, id uint64) error {
	_, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return ErrReviewNotFound
	}

	return s.reviewRepo.Delete(ctx, id)
}

// Helpful marks a review as helpful
func (s *reviewService) Helpful(ctx context.Context, reviewID uint64, userID uint64) error {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return ErrReviewNotFound
	}

	// Check if user already marked helpful
	if s.reviewRepo.HasMarkedHelpful(ctx, reviewID, userID) {
		return nil // Already marked, ignore
	}

	review.IncrementHelpful()
	return s.reviewRepo.Update(ctx, review)
}

// Unhelpful removes helpful mark
func (s *reviewService) Unhelpful(ctx context.Context, reviewID uint64, userID uint64) error {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return ErrReviewNotFound
	}

	// Check if user marked helpful
	if !s.reviewRepo.HasMarkedHelpful(ctx, reviewID, userID) {
		return nil // Not marked, ignore
	}

	// Decrement
	if review.HelpfulCount > 0 {
		review.HelpfulCount--
	}
	return s.reviewRepo.Update(ctx, review)
}

// Reply adds admin reply to review
func (s *reviewService) Reply(ctx context.Context, reviewID uint64, reply string) error {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return ErrReviewNotFound
	}

	review.Reply(reply)
	return s.reviewRepo.Update(ctx, review)
}

// Approve approves a pending review
func (s *reviewService) Approve(ctx context.Context, reviewID uint64) error {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return ErrReviewNotFound
	}

	review.Approve()
	return s.reviewRepo.Update(ctx, review)
}

// Reject rejects a pending review
func (s *reviewService) Reject(ctx context.Context, reviewID uint64, reason string) error {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return ErrReviewNotFound
	}

	review.Reject()
	// Store rejection reason in admin_reply temporarily
	review.AdminReply = fmt.Sprintf("Rejected: %s", reason)
	return s.reviewRepo.Update(ctx, review)
}

// GetReviewStats gets review statistics
func (s *reviewService) GetReviewStats(ctx context.Context, cruiseID *string, voyageID *string) (*ReviewStats, error) {
	return s.reviewRepo.GetStats(ctx, cruiseID, voyageID)
}
