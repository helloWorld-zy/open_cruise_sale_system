package analytics

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// BehaviorType represents different types of user behaviors
type BehaviorType string

const (
	BehaviorTypePageView        BehaviorType = "page_view"
	BehaviorTypeCruiseView      BehaviorType = "cruise_view"
	BehaviorTypeSearch          BehaviorType = "search"
	BehaviorTypeCabinView       BehaviorType = "cabin_view"
	BehaviorTypePriceCheck      BehaviorType = "price_check"
	BehaviorTypeAddToCart       BehaviorType = "add_to_cart"
	BehaviorTypeBookingStart    BehaviorType = "booking_start"
	BehaviorTypeBookingComplete BehaviorType = "booking_complete"
	BehaviorTypeWishlistAdd     BehaviorType = "wishlist_add"
	BehaviorTypeShare           BehaviorType = "share"
	BehaviorTypeFilterUse       BehaviorType = "filter_use"
)

// BehaviorEvent represents a single user behavior event
type BehaviorEvent struct {
	ID         uint64       `json:"id"`
	UserID     uint64       `json:"user_id"`
	SessionID  string       `json:"session_id"`
	Type       BehaviorType `json:"type"`
	EntityType string       `json:"entity_type,omitempty"` // cruise, voyage, cabin, etc.
	EntityID   uint64       `json:"entity_id,omitempty"`   // ID of the entity
	Data       EventData    `json:"data,omitempty"`        // Additional event data
	IPAddress  string       `json:"ip_address"`
	UserAgent  string       `json:"user_agent"`
	Referrer   string       `json:"referrer,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
}

// EventData holds structured event data
type EventData struct {
	SearchQuery   string                 `json:"search_query,omitempty"`
	Filters       map[string]interface{} `json:"filters,omitempty"`
	Duration      int                    `json:"duration,omitempty"` // Time spent in seconds
	PriceRange    *PriceRange            `json:"price_range,omitempty"`
	CabinTypeIDs  []uint64               `json:"cabin_type_ids,omitempty"`
	RouteID       uint64                 `json:"route_id,omitempty"`
	VoyageID      uint64                 `json:"voyage_id,omitempty"`
	DepartureDate string                 `json:"departure_date,omitempty"`
	ReturnDate    string                 `json:"return_date,omitempty"`
	GuestCount    int                    `json:"guest_count,omitempty"`
	Source        string                 `json:"source,omitempty"`  // web, mini-program, etc.
	Channel       string                 `json:"channel,omitempty"` // organic, ad, referral, etc.
}

// PriceRange represents a price range filter
type PriceRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// TrackRequest represents a request to track a behavior event
type TrackRequest struct {
	UserID     uint64       `json:"user_id,omitempty"`
	SessionID  string       `json:"session_id" validate:"required"`
	Type       BehaviorType `json:"type" validate:"required"`
	EntityType string       `json:"entity_type,omitempty"`
	EntityID   uint64       `json:"entity_id,omitempty"`
	Data       EventData    `json:"data,omitempty"`
}

// UserProfile represents a user's behavioral profile
type UserProfile struct {
	UserID              uint64             `json:"user_id"`
	PreferredCruises    []uint64           `json:"preferred_cruises"`     // Most viewed cruises
	PreferredCabinTypes []uint64           `json:"preferred_cabin_types"` // Most viewed cabin types
	PreferredRoutes     []uint64           `json:"preferred_routes"`      // Most viewed routes
	PricePreference     *PriceRange        `json:"price_preference"`      // Typical price range
	TravelPatterns      TravelPatterns     `json:"travel_patterns"`       // Date preferences
	TotalPageViews      int                `json:"total_page_views"`
	TotalSearches       int                `json:"total_searches"`
	TotalBookings       int                `json:"total_bookings"`
	LastActiveAt        time.Time          `json:"last_active_at"`
	Interests           map[string]float64 `json:"interests"` // Interest scores
}

// TravelPatterns represents user's travel preferences
type TravelPatterns struct {
	PreferredDepartureMonths []int `json:"preferred_departure_months"` // 1-12
	PreferredDurationDays    []int `json:"preferred_duration_days"`    // Common trip lengths
	PreferredGuestCount      int   `json:"preferred_guest_count"`
	AvgBookingLeadTimeDays   int   `json:"avg_booking_lead_time_days"` // Days between booking and departure
}

// BehaviorTracker defines the interface for tracking user behavior
type BehaviorTracker interface {
	// Track records a behavior event
	Track(ctx context.Context, req TrackRequest, ip, userAgent, referrer string) error

	// TrackBatch records multiple events at once
	TrackBatch(ctx context.Context, events []TrackRequest, ip, userAgent string) error

	// GetUserProfile retrieves a user's behavioral profile
	GetUserProfile(ctx context.Context, userID uint64) (*UserProfile, error)

	// GetPopularCruises gets the most popular cruises based on views/bookings
	GetPopularCruises(ctx context.Context, limit int) ([]PopularCruise, error)

	// GetTrendingSearches gets current trending search queries
	GetTrendingSearches(ctx context.Context, limit int) ([]TrendingSearch, error)

	// GetUserSimilarCruises gets cruises similar to what the user has viewed
	GetUserSimilarCruises(ctx context.Context, userID uint64, limit int) ([]uint64, error)

	// AggregateUserBehavior aggregates raw events into user profiles
	AggregateUserBehavior(ctx context.Context, userID uint64) error
}

// PopularCruise represents a popular cruise with stats
type PopularCruise struct {
	CruiseID        uint64  `json:"cruise_id"`
	ViewCount       int     `json:"view_count"`
	BookingCount    int     `json:"booking_count"`
	WishlistCount   int     `json:"wishlist_count"`
	PopularityScore float64 `json:"popularity_score"`
}

// TrendingSearch represents a trending search query
type TrendingSearch struct {
	Query     string  `json:"query"`
	Count     int     `json:"count"`
	Trend     string  `json:"trend"` // up, down, stable
	ChangePct float64 `json:"change_pct"`
}

// behaviorTracker implements BehaviorTracker
type behaviorTracker struct {
	eventRepo  repository.BehaviorEventRepository
	userRepo   repository.UserRepository
	cruiseRepo repository.CruiseRepository
}

// NewBehaviorTracker creates a new behavior tracker
func NewBehaviorTracker(
	eventRepo repository.BehaviorEventRepository,
	userRepo repository.UserRepository,
	cruiseRepo repository.CruiseRepository,
) BehaviorTracker {
	return &behaviorTracker{
		eventRepo:  eventRepo,
		userRepo:   userRepo,
		cruiseRepo: cruiseRepo,
	}
}

// Track records a behavior event
func (t *behaviorTracker) Track(ctx context.Context, req TrackRequest, ip, userAgent, referrer string) error {
	event := &BehaviorEvent{
		UserID:     req.UserID,
		SessionID:  req.SessionID,
		Type:       req.Type,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		Data:       req.Data,
		IPAddress:  ip,
		UserAgent:  userAgent,
		Referrer:   referrer,
		CreatedAt:  time.Now(),
	}

	if err := t.eventRepo.Create(ctx, event); err != nil {
		return fmt.Errorf("failed to create behavior event: %w", err)
	}

	// For high-value events, trigger real-time profile update
	if req.Type == BehaviorTypeBookingComplete || req.Type == BehaviorTypeWishlistAdd {
		if req.UserID > 0 {
			go t.AggregateUserBehavior(context.Background(), req.UserID)
		}
	}

	return nil
}

// TrackBatch records multiple events at once
func (t *behaviorTracker) TrackBatch(ctx context.Context, events []TrackRequest, ip, userAgent string) error {
	behaviorEvents := make([]*BehaviorEvent, len(events))
	now := time.Now()

	for i, req := range events {
		behaviorEvents[i] = &BehaviorEvent{
			UserID:     req.UserID,
			SessionID:  req.SessionID,
			Type:       req.Type,
			EntityType: req.EntityType,
			EntityID:   req.EntityID,
			Data:       req.Data,
			IPAddress:  ip,
			UserAgent:  userAgent,
			CreatedAt:  now,
		}
	}

	if err := t.eventRepo.CreateBatch(ctx, behaviorEvents); err != nil {
		return fmt.Errorf("failed to batch create behavior events: %w", err)
	}

	return nil
}

// GetUserProfile retrieves a user's behavioral profile
func (t *behaviorTracker) GetUserProfile(ctx context.Context, userID uint64) (*UserProfile, error) {
	// Try to get cached profile first
	profile, err := t.eventRepo.GetUserProfile(ctx, userID)
	if err != nil {
		// Generate profile from raw events
		if err := t.AggregateUserBehavior(ctx, userID); err != nil {
			return nil, err
		}
		profile, err = t.eventRepo.GetUserProfile(ctx, userID)
		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}

// GetPopularCruises gets the most popular cruises
func (t *behaviorTracker) GetPopularCruises(ctx context.Context, limit int) ([]PopularCruise, error) {
	return t.eventRepo.GetPopularCruises(ctx, limit)
}

// GetTrendingSearches gets current trending search queries
func (t *behaviorTracker) GetTrendingSearches(ctx context.Context, limit int) ([]TrendingSearch, error) {
	return t.eventRepo.GetTrendingSearches(ctx, limit)
}

// GetUserSimilarCruises gets cruises similar to what the user has viewed
func (t *behaviorTracker) GetUserSimilarCruises(ctx context.Context, userID uint64, limit int) ([]uint64, error) {
	profile, err := t.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get cruises similar to user's preferences
	return t.eventRepo.FindSimilarCruises(ctx, profile, limit)
}

// AggregateUserBehavior aggregates raw events into user profiles
func (t *behaviorTracker) AggregateUserBehavior(ctx context.Context, userID uint64) error {
	// Get user's recent events (last 90 days)
	since := time.Now().AddDate(0, 0, -90)
	events, err := t.eventRepo.GetUserEvents(ctx, userID, since)
	if err != nil {
		return err
	}

	profile := &UserProfile{
		UserID:    userID,
		Interests: make(map[string]float64),
	}

	// Count events
	cruiseViews := make(map[uint64]int)
	cabinTypeViews := make(map[uint64]int)
	routeViews := make(map[uint64]int)
	var totalPrice float64
	var priceCount int
	departureMonths := make(map[int]int)
	durationDays := make(map[int]int)

	for _, event := range events {
		switch event.Type {
		case BehaviorTypePageView:
			profile.TotalPageViews++
		case BehaviorTypeSearch:
			profile.TotalSearches++
		case BehaviorTypeBookingComplete:
			profile.TotalBookings++
		case BehaviorTypeCruiseView:
			if event.EntityID > 0 {
				cruiseViews[event.EntityID]++
			}
		case BehaviorTypeCabinView:
			if len(event.Data.CabinTypeIDs) > 0 {
				for _, id := range event.Data.CabinTypeIDs {
					cabinTypeViews[id]++
				}
			}
		}

		// Track price preferences
		if event.Data.PriceRange != nil {
			totalPrice += (event.Data.PriceRange.Min + event.Data.PriceRange.Max) / 2
			priceCount++
		}

		// Track travel patterns
		if event.Data.DepartureDate != "" {
			if date, err := time.Parse("2006-01-02", event.Data.DepartureDate); err == nil {
				departureMonths[int(date.Month())]++
			}
		}

		// Interests scoring
		if event.EntityType == "cruise" && event.EntityID > 0 {
			profile.Interests[fmt.Sprintf("cruise_%d", event.EntityID)] += 1.0
		}
	}

	// Build top preferences
	profile.PreferredCruises = getTopNFromMap(cruiseViews, 5)
	profile.PreferredCabinTypes = getTopNFromMap(cabinTypeViews, 3)
	profile.PreferredRoutes = getTopNFromMap(routeViews, 3)

	// Calculate price preference
	if priceCount > 0 {
		avgPrice := totalPrice / float64(priceCount)
		profile.PricePreference = &PriceRange{
			Min: avgPrice * 0.7,
			Max: avgPrice * 1.3,
		}
	}

	// Build travel patterns
	profile.TravelPatterns.PreferredDepartureMonths = getTopNFromMap(departureMonths, 3)
	profile.TravelPatterns.PreferredDurationDays = getTopNFromMap(durationDays, 3)

	profile.LastActiveAt = time.Now()

	// Save profile
	return t.eventRepo.SaveUserProfile(ctx, profile)
}

// getTopNFromMap returns top N keys from a map by value
func getTopNFromMap(m map[uint64]int, n int) []uint64 {
	type pair struct {
		key   uint64
		value int
	}

	pairs := make([]pair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, pair{k, v})
	}

	// Sort by value descending
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].value > pairs[i].value {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	result := make([]uint64, 0, n)
	for i := 0; i < len(pairs) && i < n; i++ {
		result = append(result, pairs[i].key)
	}

	return result
}

// String helpers for map operations
func getTopNFromMapInt(m map[int]int, n int) []int {
	type pair struct {
		key   int
		value int
	}

	pairs := make([]pair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, pair{k, v})
	}

	// Sort by value descending
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].value > pairs[i].value {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	result := make([]int, 0, n)
	for i := 0; i < len(pairs) && i < n; i++ {
		result = append(result, pairs[i].key)
	}

	return result
}

// Serialize EventData to JSON for storage
func (d EventData) MarshalJSON() ([]byte, error) {
	type Alias EventData
	return json.Marshal((*Alias)(&d))
}

// Deserialize EventData from JSON
func (d *EventData) UnmarshalJSON(data []byte) error {
	type Alias EventData
	return json.Unmarshal(data, (*Alias)(d))
}
