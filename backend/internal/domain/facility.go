package domain

import "gorm.io/datatypes"

// FacilityCategory represents a category of facilities on a cruise ship
type FacilityCategory struct {
	BaseModel
	CruiseID   string `gorm:"not null;index" json:"cruise_id"`
	Cruise     Cruise `gorm:"foreignKey:CruiseID" json:"cruise,omitempty"`
	Name       string `gorm:"not null" json:"name"`
	Icon       string `json:"icon,omitempty"`
	SortWeight int    `gorm:"default:0" json:"sort_weight"`
}

// TableName returns the table name for FacilityCategory
func (FacilityCategory) TableName() string {
	return "facility_categories"
}

// Facility represents a facility or amenity on a cruise ship
type Facility struct {
	BaseModel
	CruiseID     string           `gorm:"not null;index" json:"cruise_id"`
	Cruise       Cruise           `gorm:"foreignKey:CruiseID" json:"cruise,omitempty"`
	CategoryID   string           `gorm:"index" json:"category_id,omitempty"`
	Category     FacilityCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name         string           `gorm:"not null" json:"name"`
	DeckNumber   int              `json:"deck_number,omitempty"`
	OpenTime     string           `json:"open_time,omitempty"`
	IsFree       bool             `gorm:"default:true" json:"is_free"`
	Price        float64          `json:"price,omitempty"`
	Description  string           `json:"description,omitempty"`
	Images       datatypes.JSON   `gorm:"default:'[]'" json:"images"`
	SuitableTags datatypes.JSON   `gorm:"default:'[]'" json:"suitable_tags"`
	SortWeight   int              `gorm:"default:0" json:"sort_weight"`
	Status       string           `gorm:"default:visible" json:"status"`
}

// TableName returns the table name for Facility
func (Facility) TableName() string {
	return "facilities"
}

// FacilityStatus constants
const (
	FacilityStatusVisible = "visible"
	FacilityStatusHidden  = "hidden"
)
