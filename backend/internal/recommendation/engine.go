package recommendation

import (
	"backend/internal/analytics"
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// RecommendationType represents different types of recommendations
type RecommendationType string

const (
	RecommendationTypePersonalized   RecommendationType = "personalized"
	RecommendationTypePopular        RecommendationType = "popular"
	RecommendationTypeSimilar        RecommendationType = "similar"
	RecommendationTypeTrending       RecommendationType = "trending"
	RecommendationTypePriceDrop      RecommendationType = "price_drop"
	RecommendationTypeLastMinute     RecommendationType = "last_minute"
	RecommendationTypeFamilyFriendly RecommendationType = "family_friendly"
	RecommendationTypeLuxury         RecommendationType = "luxury"
	RecommendationTypeBudget         RecommendationType = "budget"
)

// Recommendation represents a single cruise recommendation
type Recommendation struct {
	VoyageID        string        `json:"voyage_id"`
	CruiseID        string        `json:"cruise_id"`
	CruiseName      string        `json:"cruise_name"`
	CruiseImage     string        `json:"cruise_image,omitempty"`
	VoyageNumber    string        `json:"voyage_number"`
	DepartureDate   string        `json:"departure_date"`
	ArrivalDate     string        `json:"arrival_date"`
	RouteName       string        `json:"route_name"`
	MinPrice        float64       `json:"min_price"`
	MaxPrice        float64       `json:"max_price"`
	Currency        string        `json:"currency"`
	DurationDays    int           `json:"duration_days"`
	Reason          string        `json:"reason"`
	ReasonType      string        `json:"reason_type"`
	Score           float64       `json:"score"`
	MatchFactors    []MatchFactor `json:"match_factors,omitempty"`
	Tags            []string      `json:"tags,omitempty"`
	AvailableCabins int           `json:"available_cabins"`
	DiscountPercent *int          `json:"discount_percent,omitempty"`
}

// MatchFactor represents why a recommendation matches
type MatchFactor struct {
	Factor string  `json:"factor"`
	Score  float64 `json:"score"`
}

// RecommendationRequest represents a request for recommendations
type RecommendationRequest struct {
	UserID           uint64             `json:"user_id,omitempty"`
	Type             RecommendationType `json:"type" validate:"required"`
	Limit            int                `json:"limit,omitempty"`
	ExcludeVoyageIDs []string           `json:"exclude_voyage_ids,omitempty"`
	MinDepartureDate string             `json:"min_departure_date,omitempty"`
	MaxDepartureDate string             `json:"max_departure_date,omitempty"`
	PriceRange       *PriceRange        `json:"price_range,omitempty"`
	CabinTypeIDs     []string           `json:"cabin_type_ids,omitempty"`
	RouteID          string             `json:"route_id,omitempty"`
	DurationDays     *int               `json:"duration_days,omitempty"`
}

// PriceRange for filtering recommendations
type PriceRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// RecommendationEngine defines the interface for generating recommendations
type RecommendationEngine interface {
	GetRecommendations(ctx context.Context, req RecommendationRequest) ([]Recommendation, error)
	GetPersonalizedRecommendations(ctx context.Context, userID uint64, limit int) ([]Recommendation, error)
	GetSimilarVoyages(ctx context.Context, voyageID string, limit int) ([]Recommendation, error)
	GetPopularVoyages(ctx context.Context, limit int) ([]Recommendation, error)
	GetLastMinuteDeals(ctx context.Context, limit int) ([]Recommendation, error)
}

// recommendationEngine implements RecommendationEngine
type recommendationEngine struct {
	voyageRepo      repository.VoyageRepository
	cruiseRepo      repository.CruiseRepository
	cabinRepo       repository.CabinRepository
	priceRepo       repository.PriceRepository
	inventoryRepo   repository.InventoryRepository
	behaviorTracker analytics.BehaviorTracker
}

// NewRecommendationEngine creates a new recommendation engine
func NewRecommendationEngine(
	voyageRepo repository.VoyageRepository,
	cruiseRepo repository.CruiseRepository,
	cabinRepo repository.CabinRepository,
	priceRepo repository.PriceRepository,
	inventoryRepo repository.InventoryRepository,
	behaviorTracker analytics.BehaviorTracker,
) RecommendationEngine {
	return &recommendationEngine{
		voyageRepo:      voyageRepo,
		cruiseRepo:      cruiseRepo,
		cabinRepo:       cabinRepo,
		priceRepo:       priceRepo,
		inventoryRepo:   inventoryRepo,
		behaviorTracker: behaviorTracker,
	}
}

// GetRecommendations generates recommendations based on request
func (e *recommendationEngine) GetRecommendations(ctx context.Context, req RecommendationRequest) ([]Recommendation, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 50 {
		req.Limit = 50
	}

	switch req.Type {
	case RecommendationTypePersonalized:
		return e.GetPersonalizedRecommendations(ctx, req.UserID, req.Limit)
	case RecommendationTypePopular:
		return e.GetPopularVoyages(ctx, req.Limit)
	case RecommendationTypeSimilar:
		if len(req.ExcludeVoyageIDs) > 0 {
			return e.GetSimilarVoyages(ctx, req.ExcludeVoyageIDs[0], req.Limit)
		}
		return e.GetPopularVoyages(ctx, req.Limit)
	case RecommendationTypeLastMinute:
		return e.GetLastMinuteDeals(ctx, req.Limit)
	default:
		return e.GetPopularVoyages(ctx, req.Limit)
	}
}

// GetPersonalizedRecommendations gets recommendations tailored to a specific user
func (e *recommendationEngine) GetPersonalizedRecommendations(ctx context.Context, userID uint64, limit int) ([]Recommendation, error) {
	if userID == 0 {
		return e.GetPopularVoyages(ctx, limit)
	}

	profile, err := e.behaviorTracker.GetUserProfile(ctx, userID)
	if err != nil {
		return e.GetPopularVoyages(ctx, limit)
	}

	// Get voyages using filters
	filters := repository.VoyageFilters{
		BookingStatus: domain.BookingStatusOpen,
	}

	voyages, err := e.voyageRepo.List(ctx, filters, nil)
	if err != nil {
		return nil, err
	}

	// Score and rank voyages
	scored := make([]struct {
		voyage *domain.Voyage
		score  float64
	}, 0)

	for _, voyage := range voyages {
		score := e.calculateScore(voyage, profile)
		if score > 0.3 {
			scored = append(scored, struct {
				voyage *domain.Voyage
				score  float64
			}{voyage, score})
		}
	}

	// Sort by score
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Build recommendations
	recommendations := make([]Recommendation, 0, limit)
	for i, s := range scored {
		if i >= limit {
			break
		}
		rec, err := e.buildRecommendation(ctx, s.voyage, s.score, profile)
		if err == nil {
			recommendations = append(recommendations, *rec)
		}
	}

	return recommendations, nil
}

// GetSimilarVoyages gets voyages similar to a reference voyage
func (e *recommendationEngine) GetSimilarVoyages(ctx context.Context, voyageID string, limit int) ([]Recommendation, error) {
	reference, err := e.voyageRepo.GetByID(ctx, voyageID)
	if err != nil {
		return nil, err
	}

	// Find voyages on same route
	voyages, err := e.voyageRepo.ListByRoute(ctx, reference.RouteID)
	if err != nil {
		return nil, err
	}

	recommendations := make([]Recommendation, 0, limit)
	for _, voyage := range voyages {
		if voyage.ID.String() == voyageID {
			continue
		}

		rec, err := e.buildRecommendation(ctx, voyage, 0.7, nil)
		if err != nil {
			continue
		}

		rec.Reason = "与您查看的航次相似"
		rec.ReasonType = "similar"
		recommendations = append(recommendations, *rec)

		if len(recommendations) >= limit {
			break
		}
	}

	return recommendations, nil
}

// GetPopularVoyages gets currently popular voyages
func (e *recommendationEngine) GetPopularVoyages(ctx context.Context, limit int) ([]Recommendation, error) {
	filters := repository.VoyageFilters{
		BookingStatus: domain.BookingStatusOpen,
	}

	voyages, err := e.voyageRepo.List(ctx, filters, nil)
	if err != nil {
		return nil, err
	}

	// Sort by view count/popularity (simplified - just return first ones)
	recommendations := make([]Recommendation, 0, limit)
	for i, voyage := range voyages {
		if i >= limit {
			break
		}

		rec, err := e.buildRecommendation(ctx, voyage, 0.6, nil)
		if err != nil {
			continue
		}

		rec.Reason = "热门推荐"
		rec.ReasonType = "popular"
		recommendations = append(recommendations, *rec)
	}

	return recommendations, nil
}

// GetLastMinuteDeals gets voyages departing soon with available inventory
func (e *recommendationEngine) GetLastMinuteDeals(ctx context.Context, limit int) ([]Recommendation, error) {
	now := time.Now()
	deadline := now.AddDate(0, 0, 30).Format("2006-01-02")

	filters := repository.VoyageFilters{
		BookingStatus: domain.BookingStatusOpen,
		DepartureTo:   deadline,
	}

	voyages, err := e.voyageRepo.List(ctx, filters, nil)
	if err != nil {
		return nil, err
	}

	recommendations := make([]Recommendation, 0, limit)
	for _, voyage := range voyages {
		// Parse departure date
		depDate, err := time.Parse("2006-01-02", voyage.DepartureDate)
		if err != nil {
			continue
		}

		daysUntil := int(depDate.Sub(now).Hours() / 24)
		if daysUntil < 0 || daysUntil > 30 {
			continue
		}

		// Check inventory - use preloaded inventory from voyage
		totalAvailable := 0
		for _, inv := range voyage.Inventory {
			totalAvailable += inv.AvailableCabins
		}

		if totalAvailable == 0 {
			continue
		}

		score := 0.6 + float64(30-daysUntil)*0.01
		rec, err := e.buildRecommendation(ctx, voyage, score, nil)
		if err != nil {
			continue
		}

		rec.Reason = fmt.Sprintf("%d天后出发，还剩%d个舱位", daysUntil, totalAvailable)
		rec.ReasonType = "last_minute"
		rec.AvailableCabins = totalAvailable
		recommendations = append(recommendations, *rec)

		if len(recommendations) >= limit {
			break
		}
	}

	return recommendations, nil
}

// calculateScore calculates a recommendation score for a voyage based on user profile
func (e *recommendationEngine) calculateScore(voyage *domain.Voyage, profile *analytics.UserProfile) float64 {
	score := 0.3

	if profile == nil {
		return score
	}

	// Check cruise preference
	for _, cruiseID := range profile.PreferredCruises {
		cid := strconv.FormatUint(cruiseID, 10)
		if cid == voyage.CruiseID {
			score += 0.25
			break
		}
	}

	// Check route preference
	for _, routeID := range profile.PreferredRoutes {
		rid := strconv.FormatUint(routeID, 10)
		if rid == voyage.RouteID {
			score += 0.2
			break
		}
	}

	// Check price match
	if profile.PricePreference != nil {
		if voyage.MinPrice >= profile.PricePreference.Min && voyage.MinPrice <= profile.PricePreference.Max {
			score += 0.2
		}
	}

	// Check departure month
	if len(profile.TravelPatterns.PreferredDepartureMonths) > 0 {
		if depDate, err := time.Parse("2006-01-02", voyage.DepartureDate); err == nil {
			month := int(depDate.Month())
			for _, m := range profile.TravelPatterns.PreferredDepartureMonths {
				if m == month {
					score += 0.1
					break
				}
			}
		}
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

// buildRecommendation builds a Recommendation from a voyage
func (e *recommendationEngine) buildRecommendation(ctx context.Context, voyage *domain.Voyage, score float64, profile *analytics.UserProfile) (*Recommendation, error) {
	// Use preloaded Cruise data
	cruise := voyage.Cruise
	if cruise.ID == uuid.Nil {
		// If not preloaded, fetch it
		var err error
		cruisePtr, err := e.cruiseRepo.GetByID(ctx, voyage.CruiseID)
		if err != nil {
			return nil, err
		}
		cruise = *cruisePtr
	}

	// Use preloaded Route data
	route := voyage.Route
	if route.ID == uuid.Nil {
		// Route not preloaded, use default
		route = domain.Route{Name: "航线"}
	}

	// Calculate duration
	depDate, _ := time.Parse("2006-01-02", voyage.DepartureDate)
	arrDate, _ := time.Parse("2006-01-02", voyage.ArrivalDate)
	durationDays := int(arrDate.Sub(depDate).Hours() / 24)

	// Get cruise image
	var cruiseImage string
	if len(cruise.CoverImages) > 0 {
		// Parse JSON array
		var images []string
		if err := json.Unmarshal([]byte(cruise.CoverImages), &images); err == nil && len(images) > 0 {
			cruiseImage = images[0]
		}
	}

	rec := &Recommendation{
		VoyageID:        voyage.ID.String(),
		CruiseID:        voyage.CruiseID,
		CruiseName:      cruise.NameCN,
		CruiseImage:     cruiseImage,
		VoyageNumber:    voyage.VoyageNumber,
		DepartureDate:   voyage.DepartureDate,
		ArrivalDate:     voyage.ArrivalDate,
		RouteName:       route.Name,
		MinPrice:        voyage.MinPrice,
		MaxPrice:        voyage.MaxPrice,
		Currency:        "CNY",
		DurationDays:    durationDays,
		Score:           score,
		AvailableCabins: 0,
	}

	return rec, nil
}
