package domain

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// Travelogue represents a user-generated travel story/blog post
type Travelogue struct {
	BaseModel
	UserID          string         `gorm:"not null;index" json:"user_id"`
	User            User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	VoyageID        *string        `json:"voyage_id,omitempty"`
	Voyage          *Voyage        `gorm:"foreignKey:VoyageID" json:"voyage,omitempty"`
	CruiseID        *string        `json:"cruise_id,omitempty"`
	Cruise          *Cruise        `gorm:"foreignKey:CruiseID" json:"cruise,omitempty"`
	Title           string         `gorm:"not null;size:200" json:"title"`
	Content         string         `gorm:"not null;type:text" json:"content"`
	Summary         string         `gorm:"type:text" json:"summary,omitempty"`
	CoverImage      string         `json:"cover_image,omitempty"`
	Images          datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"images,omitempty"`
	Tags            datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"tags,omitempty"`
	DestinationTags datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"destination_tags,omitempty"`
	TravelDate      *time.Time     `json:"travel_date,omitempty"`
	DurationDays    int            `json:"duration_days,omitempty"`
	Companions      string         `gorm:"size:50" json:"companions,omitempty"` // solo, couple, family, friends, group
	Rating          float64        `gorm:"type:decimal(2,1)" json:"rating,omitempty"`
	ViewCount       int            `gorm:"not null;default:0" json:"view_count"`
	LikeCount       int            `gorm:"not null;default:0" json:"like_count"`
	CommentCount    int            `gorm:"not null;default:0" json:"comment_count"`
	ShareCount      int            `gorm:"not null;default:0" json:"share_count"`
	IsFeatured      bool           `gorm:"not null;default:false" json:"is_featured"`
	IsPublished     bool           `gorm:"not null;default:false" json:"is_published"`
	PublishedAt     *time.Time     `json:"published_at,omitempty"`
	Status          string         `gorm:"not null;size:20;default:'draft'" json:"status"` // draft, pending_review, published, rejected, archived
	FeaturedOrder   int            `json:"featured_order,omitempty"`
	MetaTitle       string         `gorm:"size:200" json:"meta_title,omitempty"`
	MetaDescription string         `gorm:"type:text" json:"meta_description,omitempty"`

	// Relations
	Comments []TravelogueComment `gorm:"foreignKey:TravelogueID" json:"comments,omitempty"`
}

// TableName returns the table name
func (Travelogue) TableName() string {
	return "travelogues"
}

// Travelogue status constants
const (
	TravelogueStatusDraft         = "draft"
	TravelogueStatusPendingReview = "pending_review"
	TravelogueStatusPublished     = "published"
	TravelogueStatusRejected      = "rejected"
	TravelogueStatusArchived      = "archived"
)

// Companion type constants
const (
	CompanionSolo    = "solo"
	CompanionCouple  = "couple"
	CompanionFamily  = "family"
	CompanionFriends = "friends"
	CompanionGroup   = "group"
)

// IncrementView increments view count
func (t *Travelogue) IncrementView() {
	t.ViewCount++
}

// IncrementLike increments like count
func (t *Travelogue) IncrementLike() {
	t.LikeCount++
}

// DecrementLike decrements like count
func (t *Travelogue) DecrementLike() {
	if t.LikeCount > 0 {
		t.LikeCount--
	}
}

// IncrementComment increments comment count
func (t *Travelogue) IncrementComment() {
	t.CommentCount++
}

// DecrementComment decrements comment count
func (t *Travelogue) DecrementComment() {
	if t.CommentCount > 0 {
		t.CommentCount--
	}
}

// Publish publishes the travelogue
func (t *Travelogue) Publish() {
	t.Status = TravelogueStatusPublished
	t.IsPublished = true
	now := time.Now()
	t.PublishedAt = &now
}

// Unpublish unpublishes the travelogue
func (t *Travelogue) Unpublish() {
	t.Status = TravelogueStatusDraft
	t.IsPublished = false
	t.PublishedAt = nil
}

// Feature marks the travelogue as featured
func (t *Travelogue) Feature(order int) {
	t.IsFeatured = true
	t.FeaturedOrder = order
}

// Unfeature removes featured status
func (t *Travelogue) Unfeature() {
	t.IsFeatured = false
	t.FeaturedOrder = 0
}

// GetImages returns images as slice
func (t *Travelogue) GetImages() ([]string, error) {
	var images []string
	if t.Images == nil {
		return images, nil
	}
	err := json.Unmarshal(t.Images, &images)
	return images, err
}

// SetImages sets images from slice
func (t *Travelogue) SetImages(images []string) error {
	if images == nil {
		t.Images = datatypes.JSON("[]")
		return nil
	}
	data, err := json.Marshal(images)
	if err != nil {
		return err
	}
	t.Images = datatypes.JSON(data)
	return nil
}

// GetTags returns tags as slice
func (t *Travelogue) GetTags() ([]string, error) {
	var tags []string
	if t.Tags == nil {
		return tags, nil
	}
	err := json.Unmarshal(t.Tags, &tags)
	return tags, err
}

// SetTags sets tags from slice
func (t *Travelogue) SetTags(tags []string) error {
	if tags == nil {
		t.Tags = datatypes.JSON("[]")
		return nil
	}
	data, err := json.Marshal(tags)
	if err != nil {
		return err
	}
	t.Tags = datatypes.JSON(data)
	return nil
}

// CanEdit checks if the travelogue can be edited
func (t *Travelogue) CanEdit(userID string, isAdmin bool) bool {
	if isAdmin {
		return true
	}
	return t.UserID == userID && (t.Status == TravelogueStatusDraft || t.Status == TravelogueStatusRejected)
}

// TravelogueComment represents a comment on a travelogue
type TravelogueComment struct {
	BaseModel
	TravelogueID string              `gorm:"not null;index" json:"travelogue_id"`
	Travelogue   Travelogue          `gorm:"foreignKey:TravelogueID" json:"travelogue,omitempty"`
	UserID       string              `gorm:"not null;index" json:"user_id"`
	User         User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentID     *string             `json:"parent_id,omitempty"`
	Parent       *TravelogueComment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Content      string              `gorm:"not null;type:text" json:"content"`
	LikeCount    int                 `gorm:"not null;default:0" json:"like_count"`
	IsApproved   bool                `gorm:"not null;default:true" json:"is_approved"`
	Replies      []TravelogueComment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// TableName returns the table name
func (TravelogueComment) TableName() string {
	return "travelogue_comments"
}

// IsReply checks if this comment is a reply to another comment
func (c *TravelogueComment) IsReply() bool {
	return c.ParentID != nil
}

// TravelogueLike represents a like on a travelogue
type TravelogueLike struct {
	BaseModel
	TravelogueID string `gorm:"not null;index" json:"travelogue_id"`
	UserID       string `gorm:"not null;index" json:"user_id"`
	User         User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name
func (TravelogueLike) TableName() string {
	return "travelogue_likes"
}

// Review represents a user review for a cruise or voyage
type Review struct {
	BaseModel
	UserID       string         `gorm:"not null;index" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	VoyageID     *string        `json:"voyage_id,omitempty"`
	Voyage       *Voyage        `gorm:"foreignKey:VoyageID" json:"voyage,omitempty"`
	CruiseID     *string        `json:"cruise_id,omitempty"`
	Cruise       *Cruise        `gorm:"foreignKey:CruiseID" json:"cruise,omitempty"`
	OrderID      *string        `json:"order_id,omitempty"`
	Order        *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Rating       int            `gorm:"not null" json:"rating"` // 1-5
	Title        string         `gorm:"size:200" json:"title"`
	Content      string         `gorm:"not null;type:text" json:"content"`
	Pros         datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"pros,omitempty"`
	Cons         datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"cons,omitempty"`
	Images       datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"images,omitempty"`
	IsVerified   bool           `gorm:"not null;default:false" json:"is_verified"` // Verified purchase
	IsAnonymous  bool           `gorm:"not null;default:false" json:"is_anonymous"`
	HelpfulCount int            `gorm:"not null;default:0" json:"helpful_count"`
	Status       string         `gorm:"not null;size:20;default:'pending'" json:"status"` // pending, approved, rejected
	AdminReply   string         `gorm:"type:text" json:"admin_reply,omitempty"`
	RepliedAt    *time.Time     `json:"replied_at,omitempty"`
}

// ReviewStats represents review statistics
type ReviewStats struct {
	TotalReviews       int         `json:"total_reviews"`
	AverageRating      float64     `json:"average_rating"`
	FiveStarCount      int         `json:"five_star_count"`
	FourStarCount      int         `json:"four_star_count"`
	ThreeStarCount     int         `json:"three_star_count"`
	TwoStarCount       int         `json:"two_star_count"`
	OneStarCount       int         `json:"one_star_count"`
	RatingDistribution map[int]int `json:"rating_distribution"`
	VerifiedCount      int         `json:"verified_count"`
}

// TableName returns the table name
func (Review) TableName() string {
	return "reviews"
}

// Review status constants
const (
	ReviewStatusPending  = "pending"
	ReviewStatusApproved = "approved"
	ReviewStatusRejected = "rejected"
)

// IsPositive checks if the review is positive (4-5 stars)
func (r *Review) IsPositive() bool {
	return r.Rating >= 4
}

// IsNegative checks if the review is negative (1-2 stars)
func (r *Review) IsNegative() bool {
	return r.Rating <= 2
}

// IncrementHelpful increments helpful count
func (r *Review) IncrementHelpful() {
	r.HelpfulCount++
}

// Approve approves the review
func (r *Review) Approve() {
	r.Status = ReviewStatusApproved
}

// Reject rejects the review
func (r *Review) Reject() {
	r.Status = ReviewStatusRejected
}

// Reply adds admin reply
func (r *Review) Reply(reply string) {
	r.AdminReply = reply
	now := time.Now()
	r.RepliedAt = &now
}
