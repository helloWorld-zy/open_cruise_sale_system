package domain

import "gorm.io/datatypes"

// CabinType represents a type of cabin on a cruise ship
type CabinType struct {
	BaseModel
	CruiseID       string         `gorm:"not null;index" json:"cruise_id"`
	Cruise         Cruise         `gorm:"foreignKey:CruiseID" json:"cruise,omitempty"`
	Name           string         `gorm:"not null" json:"name"`
	Code           string         `gorm:"not null" json:"code"`
	MinAreaSqm     float64        `json:"min_area_sqm,omitempty"`
	MaxAreaSqm     float64        `json:"max_area_sqm,omitempty"`
	StandardGuests int            `json:"standard_guests,omitempty"`
	MaxGuests      int            `json:"max_guests,omitempty"`
	BedTypes       string         `json:"bed_types,omitempty"`
	FeatureTags    datatypes.JSON `gorm:"default:'[]'" json:"feature_tags"`
	Description    string         `json:"description,omitempty"`
	Images         datatypes.JSON `gorm:"default:'[]'" json:"images"`
	FloorPlanURL   string         `json:"floor_plan_url,omitempty"`
	Amenities      datatypes.JSON `gorm:"default:'[]'" json:"amenities"`
	SortWeight     int            `gorm:"default:0" json:"sort_weight"`
	Status         string         `gorm:"default:active" json:"status"`
}

// TableName returns the table name for CabinType
func (CabinType) TableName() string {
	return "cabin_types"
}

// CabinTypeStatus constants
const (
	CabinTypeStatusActive   = "active"
	CabinTypeStatusInactive = "inactive"
)
