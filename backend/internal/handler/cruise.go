package handler

import (
	"backend/internal/pagination"
	"backend/internal/response"
	"backend/internal/service"
	"backend/internal/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CruiseHandler handles HTTP requests for cruises
type CruiseHandler struct {
	service service.CruiseService
}

// NewCruiseHandler creates a new cruise handler
func NewCruiseHandler(service service.CruiseService) *CruiseHandler {
	return &CruiseHandler{service: service}
}

// Create godoc
// @Summary Create a new cruise
// @Description Create a new cruise ship
// @Tags cruises
// @Accept json
// @Produce json
// @Param request body service.CreateCruiseRequest true "Create cruise request"
// @Success 201 {object} response.Response{data=domain.Cruise}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /cruises [post]
func (h *CruiseHandler) Create(c *gin.Context) {
	var req service.CreateCruiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cruise, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidCruiseData {
			response.BadRequest(c, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, cruise)
}

// GetByID godoc
// @Summary Get a cruise by ID
// @Description Get a cruise by its ID
// @Tags cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Success 200 {object} response.Response{data=domain.Cruise}
// @Failure 404 {object} response.Response
// @Router /cruises/{id} [get]
func (h *CruiseHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	cruise, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrCruiseNotFound {
			response.NotFound(c, "cruise not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, cruise)
}

// GetByCode godoc
// @Summary Get a cruise by code
// @Description Get a cruise by its unique code
// @Tags cruises
// @Accept json
// @Produce json
// @Param code path string true "Cruise code"
// @Success 200 {object} response.Response{data=domain.Cruise}
// @Failure 404 {object} response.Response
// @Router /cruises/code/{code} [get]
func (h *CruiseHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	cruise, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		if err == service.ErrCruiseNotFound {
			response.NotFound(c, "cruise not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, cruise)
}

// List godoc
// @Summary List cruises
// @Description List cruises with pagination and filters
// @Tags cruises
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param status query string false "Status (active, inactive, maintenance)"
// @Param keyword query string false "Search keyword"
// @Param has_cabin_type query bool false "Filter by having cabin types"
// @Param min_capacity query int false "Minimum passenger capacity"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} response.Response{data=[]domain.Cruise,pagination=pagination.Paginator}
// @Router /cruises [get]
func (h *CruiseHandler) List(c *gin.Context) {
	var req service.ListCruisesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Set default pagination
	req.Paginator = *pagination.NewPaginator(c)

	result, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// Update godoc
// @Summary Update a cruise
// @Description Update a cruise by ID
// @Tags cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Param request body service.UpdateCruiseRequest true "Update cruise request"
// @Success 200 {object} response.Response{data=domain.Cruise}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /cruises/{id} [put]
func (h *CruiseHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateCruiseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cruise, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err == service.ErrCruiseNotFound {
			response.NotFound(c, "cruise not found")
			return
		}
		if err == service.ErrInvalidCruiseData {
			response.BadRequest(c, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, cruise)
}

// Delete godoc
// @Summary Delete a cruise
// @Description Soft delete a cruise by ID
// @Tags cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /cruises/{id} [delete]
func (h *CruiseHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == service.ErrCruiseNotFound {
			response.NotFound(c, "cruise not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Restore godoc
// @Summary Restore a deleted cruise
// @Description Restore a soft-deleted cruise by ID
// @Tags cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /cruises/{id}/restore [post]
func (h *CruiseHandler) Restore(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Restore(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus godoc
// @Summary Update cruise status
// @Description Update the status of a cruise
// @Tags cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Param status body UpdateStatusRequest true "Status update"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /cruises/{id}/status [put]
func (h *CruiseHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		if err == service.ErrCruiseNotFound {
			response.NotFound(c, "cruise not found")
			return
		}
		if err == service.ErrInvalidCruiseData {
			response.BadRequest(c, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateSortWeight godoc
// @Summary Update cruise sort weight
// @Description Update the sort weight of a cruise
// @Tags cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Param weight body UpdateSortWeightRequest true "Sort weight update"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /cruises/{id}/sort-weight [put]
func (h *CruiseHandler) UpdateSortWeight(c *gin.Context) {
	id := c.Param("id")

	var req UpdateSortWeightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateSortWeight(c.Request.Context(), id, req.SortWeight); err != nil {
		if err == service.ErrCruiseNotFound {
			response.NotFound(c, "cruise not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatusRequest represents a status update request
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateSortWeightRequest represents a sort weight update request
type UpdateSortWeightRequest struct {
	SortWeight int `json:"sort_weight" binding:"required"`
}
