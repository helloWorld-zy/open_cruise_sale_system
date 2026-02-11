package domain

// FrequentPassenger represents a frequently used passenger for a user
type FrequentPassenger struct {
	BaseModel
	UserID              string `gorm:"not null;index" json:"user_id"`
	Name                string `gorm:"not null;size:100" json:"name"`
	Surname             string `gorm:"not null;size:100" json:"surname"`
	GivenName           string `gorm:"size:100" json:"given_name,omitempty"`
	Gender              string `gorm:"not null;size:10" json:"gender"`
	BirthDate           string `json:"birth_date"`
	Nationality         string `gorm:"size:50;default:中国" json:"nationality,omitempty"`
	PassportNumber      string `gorm:"size:50;column:passport_number" json:"passport_number,omitempty"`
	PassportExpiry      string `json:"passport_expiry,omitempty"`
	IDNumber            string `gorm:"size:18;column:id_number" json:"id_number,omitempty"`
	Phone               string `gorm:"size:20" json:"phone,omitempty"`
	Email               string `gorm:"size:100" json:"email,omitempty"`
	DietaryRequirements string `json:"dietary_requirements,omitempty"`
	MedicalNotes        string `json:"medical_notes,omitempty"`
	IsDefault           bool   `gorm:"default:false;column:is_default" json:"is_default"`
}

// TableName returns the table name for FrequentPassenger
func (FrequentPassenger) TableName() string {
	return "frequent_passengers"
}

// GetDisplayName returns the display name for passenger
func (fp *FrequentPassenger) GetDisplayName() string {
	if fp.GivenName != "" {
		return fp.Surname + " " + fp.GivenName + " / " + fp.Name
	}
	return fp.Name
}
