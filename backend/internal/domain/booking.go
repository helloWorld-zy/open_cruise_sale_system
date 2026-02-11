package domain

import "gorm.io/datatypes"

// Route represents a cruise route/itinerary
type Route struct {
	BaseModel
	CruiseID      string         `gorm:"not null;index" json:"cruise_id"`
	Cruise        Cruise         `gorm:"foreignKey:CruiseID" json:"cruise,omitempty"`
	Name          string         `gorm:"not null" json:"name"`
	Code          string         `gorm:"not null;uniqueIndex" json:"code"`
	DeparturePort string         `gorm:"not null" json:"departure_port"`
	ArrivalPort   string         `gorm:"not null" json:"arrival_port"`
	DurationDays  int            `gorm:"not null;default:1" json:"duration_days"`
	Description   string         `json:"description,omitempty"`
	Itinerary     datatypes.JSON `gorm:"default:'[]'" json:"itinerary"`
	Status        string         `gorm:"default:active" json:"status"`
	SortWeight    int            `gorm:"default:0" json:"sort_weight"`

	// Relations
	Voyages []Voyage `gorm:"foreignKey:RouteID" json:"voyages,omitempty"`
}

// TableName returns the table name for Route
func (Route) TableName() string {
	return "routes"
}

// RouteStatus constants
const (
	RouteStatusActive   = "active"
	RouteStatusInactive = "inactive"
)

// Voyage represents a specific sailing/cruise voyage
type Voyage struct {
	BaseModel
	RouteID       string  `gorm:"not null;index" json:"route_id"`
	Route         Route   `gorm:"foreignKey:RouteID" json:"route,omitempty"`
	CruiseID      string  `gorm:"not null;index" json:"cruise_id"`
	Cruise        Cruise  `gorm:"foreignKey:CruiseID" json:"cruise,omitempty"`
	VoyageNumber  string  `gorm:"not null;uniqueIndex" json:"voyage_number"`
	DepartureDate string  `gorm:"not null" json:"departure_date"`
	ArrivalDate   string  `gorm:"not null" json:"arrival_date"`
	DepartureTime string  `json:"departure_time,omitempty"`
	ArrivalTime   string  `json:"arrival_time,omitempty"`
	Status        string  `gorm:"default:scheduled" json:"status"`
	BookingStatus string  `gorm:"default:open" json:"booking_status"`
	MinPrice      float64 `json:"min_price,omitempty"`
	MaxPrice      float64 `json:"max_price,omitempty"`

	// Relations
	Cabins    []Cabin          `gorm:"foreignKey:VoyageID" json:"cabins,omitempty"`
	Inventory []CabinInventory `gorm:"foreignKey:VoyageID" json:"inventory,omitempty"`
	Prices    []CabinPrice     `gorm:"foreignKey:VoyageID" json:"prices,omitempty"`
}

// TableName returns the table name for Voyage
func (Voyage) TableName() string {
	return "voyages"
}

// VoyageStatus constants
const (
	VoyageStatusScheduled = "scheduled"
	VoyageStatusActive    = "active"
	VoyageStatusCompleted = "completed"
	VoyageStatusCancelled = "cancelled"
)

// BookingStatus constants
const (
	BookingStatusOpen   = "open"
	BookingStatusFull   = "full"
	BookingStatusClosed = "closed"
)

// Cabin represents a specific cabin instance for a voyage
type Cabin struct {
	BaseModel
	VoyageID     string    `gorm:"not null;index" json:"voyage_id"`
	Voyage       Voyage    `gorm:"foreignKey:VoyageID" json:"voyage,omitempty"`
	CabinTypeID  string    `gorm:"not null;index" json:"cabin_type_id"`
	CabinType    CabinType `gorm:"foreignKey:CabinTypeID" json:"cabin_type,omitempty"`
	CabinNumber  string    `gorm:"not null" json:"cabin_number"`
	DeckNumber   int       `json:"deck_number,omitempty"`
	Section      string    `json:"section,omitempty"`
	Status       string    `gorm:"default:available" json:"status"`
	IsAccessible bool      `gorm:"default:false" json:"is_accessible"`
	IsConnecting bool      `gorm:"default:false" json:"is_connecting"`

	// Relations
	OrderItems []OrderItem `gorm:"foreignKey:CabinID" json:"order_items,omitempty"`
}

// TableName returns the table name for Cabin
func (Cabin) TableName() string {
	return "cabins"
}

// CabinStatus constants
const (
	CabinStatusAvailable   = "available"
	CabinStatusOccupied    = "occupied"
	CabinStatusMaintenance = "maintenance"
	CabinStatusLocked      = "locked"
)

// CabinInventory represents the inventory tracking for cabin types per voyage
type CabinInventory struct {
	BaseModel
	VoyageID          string    `gorm:"not null;index" json:"voyage_id"`
	Voyage            Voyage    `gorm:"foreignKey:VoyageID" json:"voyage,omitempty"`
	CabinTypeID       string    `gorm:"not null;index" json:"cabin_type_id"`
	CabinType         CabinType `gorm:"foreignKey:CabinTypeID" json:"cabin_type,omitempty"`
	TotalCabins       int       `gorm:"not null;default:0" json:"total_cabins"`
	AvailableCabins   int       `gorm:"not null;default:0" json:"available_cabins"`
	ReservedCabins    int       `gorm:"not null;default:0" json:"reserved_cabins"`
	BookedCabins      int       `gorm:"not null;default:0" json:"booked_cabins"`
	LockedCabins      int       `gorm:"not null;default:0" json:"locked_cabins"`
	MaintenanceCabins int       `gorm:"not null;default:0" json:"maintenance_cabins"`
	LockVersion       int       `gorm:"default:0" json:"lock_version"` // Optimistic locking
	LastUpdatedAt     string    `json:"last_updated_at,omitempty"`
}

// TableName returns the table name for CabinInventory
func (CabinInventory) TableName() string {
	return "cabin_inventory"
}

// CabinPrice represents pricing for a cabin type on a specific voyage
type CabinPrice struct {
	BaseModel
	VoyageID           string    `gorm:"not null;index" json:"voyage_id"`
	Voyage             Voyage    `gorm:"foreignKey:VoyageID" json:"voyage,omitempty"`
	CabinTypeID        string    `gorm:"not null;index" json:"cabin_type_id"`
	CabinType          CabinType `gorm:"foreignKey:CabinTypeID" json:"cabin_type,omitempty"`
	PriceType          string    `gorm:"not null;default:standard" json:"price_type"`
	AdultPrice         float64   `gorm:"not null" json:"adult_price"`
	ChildPrice         float64   `json:"child_price,omitempty"`
	InfantPrice        float64   `json:"infant_price,omitempty"`
	SingleSupplement   float64   `json:"single_supplement,omitempty"`
	PortFee            float64   `gorm:"default:0" json:"port_fee"`
	ServiceFee         float64   `gorm:"default:0" json:"service_fee"`
	IsPromotion        bool      `gorm:"default:false" json:"is_promotion"`
	PromotionStartDate string    `json:"promotion_start_date,omitempty"`
	PromotionEndDate   string    `json:"promotion_end_date,omitempty"`
	MinPassengers      int       `gorm:"default:1" json:"min_passengers"`
	MaxPassengers      int       `gorm:"default:4" json:"max_passengers"`
}

// TableName returns the table name for CabinPrice
func (CabinPrice) TableName() string {
	return "cabin_prices"
}

// PriceType constants
const (
	PriceTypeStandard   = "standard"
	PriceTypeEarlyBird  = "early_bird"
	PriceTypeLastMinute = "last_minute"
	PriceTypeGroup      = "group"
)

// CalculateTotal calculates total price for given passengers
func (cp *CabinPrice) CalculateTotal(adults, children, infants int) float64 {
	total := cp.AdultPrice * float64(adults)
	if cp.ChildPrice > 0 {
		total += cp.ChildPrice * float64(children)
	}
	if cp.InfantPrice > 0 {
		total += cp.InfantPrice * float64(infants)
	}
	total += cp.PortFee * float64(adults+children)
	total += cp.ServiceFee * float64(adults+children)
	return total
}
