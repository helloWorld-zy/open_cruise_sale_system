package domain

// CruiseCompany represents a cruise line company
type CruiseCompany struct {
	BaseModel
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	NameEN      string `json:"name_en,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
	Website     string `json:"website,omitempty"`
	Description string `json:"description,omitempty"`
}

// TableName returns the table name for CruiseCompany
func (CruiseCompany) TableName() string {
	return "cruise_companies"
}
