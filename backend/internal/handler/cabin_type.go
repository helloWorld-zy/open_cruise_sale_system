package handler

import (
	"backend/internal/pagination"
	"backend/internal/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CabinTypeHandler handles HTTP requests for cabin types
type CabinTypeHandler struct {
	service service.CabinTypeService
}

// NewCabinTypeHandler creates a new cabin type handler
func NewCabinTypeHandler(service service.CabinTypeService) *CabinTypeHandler {
	return &CabinTypeHandler{service: service}
}

// Create godoc
// @Summary Create a new cabin type
// @Description Create a new cabin type for a cruise
// @Tags cabin-types
// @Accept json
// @Produce json
// @Param request body service.CreateCabinTypeRequest true "Create cabin type request"
// @Success 201 {object} response.Response{data=domain.CabinType}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /cabin-types [post]
func (h *CabinTypeHandler) Create(c *gin.Context) {
	var req service.CreateCabinTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cabinType, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidCabinTypeData {
			response.BadRequest(c, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, cabinType)
}

// GetByID godoc
// @Summary Get a cabin type by ID
// @Description Get a cabin type by its ID
// @Tags cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin type ID"
// @Success 200 {object} response.Response{data=domain.CabinType}
// @Failure 404 {object} response.Response
// @Router /cabin-types/{id} [get]
func (h *CabinTypeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	cabinType, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrCabinTypeNotFound {
			response.NotFound(c, "cabin type not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, cabinType)
}

// List godoc
// @Summary List cabin types
// @Description List cabin types with pagination and filters
// @Tags cabin-types
// @Accept json
// @Produce json
// @Param cruise_id query string false "Cruise ID"
// @Param status query string false "Status (active, inactive)"
// @Param min_guests query int false "Minimum guests"
// @Param max_guests query int false "Maximum guests"
// @Param min_area query number false "Minimum area in sqm"
// @Param max_area query number false "Maximum area in sqm"
// @Param bed_type query string false "Bed type"
// @Param feature_tag query string false "Feature tag"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.Response{data=[]domain.CabinType,pagination=pagination.Paginator}
// @Router /cabin-types [get]
func (h *CabinTypeHandler) List(c *gin.Context) {
	var req service.ListCabinTypesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	req.Paginator = *pagination.NewPaginator(c)

	result, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ListByCruise godoc
// @Summary List cabin types by cruise
// @Description Get all cabin types for a specific cruise
// @Tags cabin-types
// @Accept json
// @Produce json
// @Param cruise_id path string true "Cruise ID"
// @Success 200 {object} response.Response{data=[]domain.CabinType}
// @Failure 404 {object} response.Response
// @Router /cruises/{cruise_id}/cabin-types [get]
func (h *CabinTypeHandler) ListByCruise(c *gin.Context) {
	cruiseID := c.Param("cruise_id")

	cabinTypes, err := h.service.ListByCruise(c.Request.Context(), cruiseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, cabinTypes)
}

// Update godoc
// @Summary Update a cabin type
// @Description Update a cabin type by ID
// @Tags cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin type ID"
// @Param request body service.UpdateCabinTypeRequest true "Update cabin type request"
// @Success 200 {object} response.Response{data=domain.CabinType}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /cabin-types/{id} [put]
func (h *CabinTypeHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateCabinTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cabinType, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err == service.ErrCabinTypeNotFound {
			response.NotFound(c, "cabin type not found")
			return
		}
		if err == service.ErrInvalidCabinTypeData {
			response.BadRequest(c, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, cabinType)
}

// Delete godoc
// @Summary Delete a cabin type
// @Description Soft delete a cabin type by ID
// @Tags cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin type ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /cabin-types/{id} [delete]
func (h *CabinTypeHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == service.ErrCabinTypeNotFound {
			response.NotFound(c, "cabin type not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Restore godoc
// @Summary Restore a deleted cabin type
// @Description Restore a soft-deleted cabin type
// @Tags cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin type ID"
// @Success 200 {object} response.Response
// @Router /cabin-types/{id}/restore [post]
func (h *CabinTypeHandler) Restore(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Restore(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus godoc
// @Summary Update cabin type status
// @Description Update the status of a cabin type
// @Tags cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin type ID"
// @Param status body UpdateStatusRequest true "Status update"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /cabin-types/{id}/status [put]
func (h *CabinTypeHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		if err == service.ErrCabinTypeNotFound {
			response.NotFound(c, "cabin type not found")
			return
		}
		if err == service.ErrInvalidCabinTypeData {
			response.BadRequest(c, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
