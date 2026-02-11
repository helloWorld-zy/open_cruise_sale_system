package domain

import "time"

// User represents a user account
type User struct {
	BaseModel
	Phone              string              `gorm:"uniqueIndex;size:20" json:"phone,omitempty"`
	Email              string              `gorm:"uniqueIndex;size:100" json:"email,omitempty"`
	PasswordHash       string              `gorm:"size:255" json:"-"` // Never expose password
	WechatOpenID       string              `gorm:"uniqueIndex;size:100;column:wechat_openid" json:"wechat_openid,omitempty"`
	WechatUnionID      string              `gorm:"size:100;column:wechat_unionid" json:"wechat_unionid,omitempty"`
	WechatNickname     string              `gorm:"size:100;column:wechat_nickname" json:"wechat_nickname,omitempty"`
	WechatAvatarURL    string              `gorm:"column:wechat_avatar_url" json:"wechat_avatar_url,omitempty"`
	Nickname           string              `gorm:"size:100" json:"nickname,omitempty"`
	AvatarURL          string              `gorm:"column:avatar_url" json:"avatar_url,omitempty"`
	RealName           string              `gorm:"size:100;column:real_name" json:"real_name,omitempty"`
	Gender             string              `gorm:"size:10" json:"gender,omitempty"`
	BirthDate          string              `json:"birth_date,omitempty"`
	Status             string              `gorm:"size:20;default:active" json:"status"`
	PhoneVerified      bool                `gorm:"default:false;column:phone_verified" json:"phone_verified"`
	EmailVerified      bool                `gorm:"default:false;column:email_verified" json:"email_verified"`
	IdentityVerified   bool                `gorm:"default:false;column:identity_verified" json:"identity_verified"`
	IDNumber           string              `gorm:"size:18;column:id_number" json:"id_number,omitempty"`
	IDFrontImage       string              `gorm:"column:id_front_image" json:"id_front_image,omitempty"`
	IDBackImage        string              `gorm:"column:id_back_image" json:"id_back_image,omitempty"`
	LastLoginAt        *time.Time          `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	LastLoginIP        string              `gorm:"column:last_login_ip" json:"last_login_ip,omitempty"`
	FrequentPassengers []FrequentPassenger `gorm:"foreignKey:UserID" json:"frequent_passengers,omitempty"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// User status constants
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBanned   = "banned"
)

// Gender constants
const (
	GenderMale    = "male"
	GenderFemale  = "female"
	GenderUnknown = "unknown"
)

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsVerified checks if user has verified phone or email
func (u *User) IsVerified() bool {
	return u.PhoneVerified || u.EmailVerified
}

// UpdateLastLogin updates last login info
func (u *User) UpdateLastLogin(ip string) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastLoginIP = ip
}

// CanLogin checks if user can login
func (u *User) CanLogin() bool {
	return u.Status == UserStatusActive
}
