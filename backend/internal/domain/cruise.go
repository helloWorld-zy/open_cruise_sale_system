package domain

import "gorm.io/datatypes"

// Cruise represents a cruise ship
type Cruise struct {
	BaseModel
	CompanyID         string         `gorm:"not null;index" json:"company_id"`
	Company           CruiseCompany  `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	NameCN            string         `gorm:"not null" json:"name_cn"`
	NameEN            string         `json:"name_en,omitempty"`
	Code              string         `gorm:"not null;uniqueIndex" json:"code"`
	GrossTonnage      int            `json:"gross_tonnage,omitempty"`
	PassengerCapacity int            `json:"passenger_capacity,omitempty"`
	CrewCount         int            `json:"crew_count,omitempty"`
	BuiltYear         int            `json:"built_year,omitempty"`
	RenovatedYear     int            `json:"renovated_year,omitempty"`
	LengthMeters      float64        `json:"length_meters,omitempty"`
	WidthMeters       float64        `json:"width_meters,omitempty"`
	DeckCount         int            `json:"deck_count,omitempty"`
	CoverImages       datatypes.JSON `gorm:"default:'[]'" json:"cover_images"`
	Status            string         `gorm:"default:active" json:"status"`
	SortWeight        int            `gorm:"default:0" json:"sort_weight"`

	// Relations
	CabinTypes []CabinType `gorm:"foreignKey:CruiseID" json:"cabin_types,omitempty"`
	Facilities []Facility  `gorm:"foreignKey:CruiseID" json:"facilities,omitempty"`
}

// TableName returns the table name for Cruise
func (Cruise) TableName() string {
	return "cruises"
}

// CruiseStatus constants
const (
	CruiseStatusActive      = "active"
	CruiseStatusInactive    = "inactive"
	CruiseStatusMaintenance = "maintenance"
)
