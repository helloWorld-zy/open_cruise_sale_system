package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Paginator represents pagination parameters
type Paginator struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
	Pages    int   `json:"pages"`
}

// NewPaginator creates a new paginator from gin context
func NewPaginator(c *gin.Context) *Paginator {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return &Paginator{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset returns the offset for database query
func (p *Paginator) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the limit for database query
func (p *Paginator) Limit() int {
	return p.PageSize
}

// SetTotal sets the total count and calculates pages
func (p *Paginator) SetTotal(total int64) {
	p.Total = total
	p.Pages = int(math.Ceil(float64(total) / float64(p.PageSize)))
}

// Paginate applies pagination to a GORM query
func Paginate(db *gorm.DB, p *Paginator) *gorm.DB {
	return db.Offset(p.Offset()).Limit(p.Limit())
}

// Result represents a paginated result
type Result struct {
	Data       interface{} `json:"data"`
	Pagination Paginator   `json:"pagination"`
}

// NewResult creates a new paginated result
func NewResult(data interface{}, paginator Paginator) Result {
	return Result{
		Data:       data,
		Pagination: paginator,
	}
}
