package domain

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
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
	BaseModel
	UserID     uint64         `json:"user_id" gorm:"index"`
	SessionID  string         `json:"session_id" gorm:"index"`
	Type       BehaviorType   `json:"type" gorm:"index"`
	EntityType string         `json:"entity_type,omitempty"` // cruise, voyage, cabin, etc.
	EntityID   uint64         `json:"entity_id,omitempty"`   // ID of the entity
	Data       datatypes.JSON `json:"data,omitempty"`        // Additional event data
	IPAddress  string         `json:"ip_address"`
	UserAgent  string         `json:"user_agent"`
	Referrer   string         `json:"referrer,omitempty"`
}

// TableName returns the table name for BehaviorEvent
func (BehaviorEvent) TableName() string {
	return "behavior_events"
}

// EventData holds structured event data (helper struct, stored as JSON)
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

// UserProfile represents a user's behavioral profile
type UserProfile struct {
	BaseModel
	UserID              uint64         `json:"user_id" gorm:"uniqueIndex"`
	PreferredCruises    datatypes.JSON `json:"preferred_cruises"`     // []uint64
	PreferredCabinTypes datatypes.JSON `json:"preferred_cabin_types"` // []uint64
	PreferredRoutes     datatypes.JSON `json:"preferred_routes"`      // []uint64
	PricePreference     datatypes.JSON `json:"price_preference"`      // *PriceRange
	TravelPatterns      datatypes.JSON `json:"travel_patterns"`       // TravelPatterns
	TotalPageViews      int            `json:"total_page_views"`
	TotalSearches       int            `json:"total_searches"`
	TotalBookings       int            `json:"total_bookings"`
	LastActiveAt        time.Time      `json:"last_active_at"`
	Interests           datatypes.JSON `json:"interests"` // map[string]float64
}

// TableName returns the table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profiles"
}

// GetPreferredCruises Helper methods to unmarshal JSON fields
func (p *UserProfile) GetPreferredCruises() []uint64 {
	var ids []uint64
	if len(p.PreferredCruises) > 0 {
		_ = json.Unmarshal(p.PreferredCruises, &ids)
	}
	return ids
}

func (p *UserProfile) GetPreferredRoutes() []uint64 {
	var ids []uint64
	if len(p.PreferredRoutes) > 0 {
		_ = json.Unmarshal(p.PreferredRoutes, &ids)
	}
	return ids
}

func (p *UserProfile) GetPricePreference() *PriceRange {
	var pr PriceRange
	if len(p.PricePreference) > 0 {
		if err := json.Unmarshal(p.PricePreference, &pr); err == nil {
			return &pr
		}
	}
	return nil
}

func (p *UserProfile) GetTravelPatterns() *TravelPatterns {
	var tp TravelPatterns
	if len(p.TravelPatterns) > 0 {
		if err := json.Unmarshal(p.TravelPatterns, &tp); err == nil {
			return &tp
		}
	}
	return nil
}

// TravelPatterns represents user's travel preferences
type TravelPatterns struct {
	PreferredDepartureMonths []int `json:"preferred_departure_months"` // 1-12
	PreferredDurationDays    []int `json:"preferred_duration_days"`    // Common trip lengths
	PreferredGuestCount      int   `json:"preferred_guest_count"`
	AvgBookingLeadTimeDays   int   `json:"avg_booking_lead_time_days"` // Days between booking and departure
}

// PopularCruise represents a popular cruise with stats (Result struct)
type PopularCruise struct {
	CruiseID        uint64  `json:"cruise_id"`
	ViewCount       int     `json:"view_count"`
	BookingCount    int     `json:"booking_count"`
	WishlistCount   int     `json:"wishlist_count"`
	PopularityScore float64 `json:"popularity_score"`
}

// TrendingSearch represents a trending search query (Result struct)
type TrendingSearch struct {
	Query     string  `json:"query"`
	Count     int     `json:"count"`
	Trend     string  `json:"trend"` // up, down, stable
	ChangePct float64 `json:"change_pct"`
}
