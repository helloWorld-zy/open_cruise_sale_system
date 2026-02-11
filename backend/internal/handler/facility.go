package handler

import (
	"backend/internal/pagination"
	"backend/internal/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FacilityHandler handles HTTP requests for facilities
type FacilityHandler struct {
	service service.FacilityService
}

// FacilityCategoryHandler handles HTTP requests for facility categories
type FacilityCategoryHandler struct {
	service service.FacilityCategoryService
}

// NewFacilityHandler creates a new facility handler
func NewFacilityHandler(service service.FacilityService) *FacilityHandler {
	return &FacilityHandler{service: service}
}

// NewFacilityCategoryHandler creates a new facility category handler
func NewFacilityCategoryHandler(service service.FacilityCategoryService) *FacilityCategoryHandler {
	return &FacilityCategoryHandler{service: service}
}

// Create godoc
// @Summary Create a new facility
// @Description Create a new facility for a cruise
// @Tags facilities
// @Accept json
// @Produce json
// @Param request body service.CreateFacilityRequest true "Create facility request"
// @Success 201 {object} response.Response{data=domain.Facility}
// @Failure 400 {object} response.Response
// @Router /facilities [post]
func (h *FacilityHandler) Create(c *gin.Context) {
	var req service.CreateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	facility, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidFacilityData {
			response.BadRequest(c, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, facility)
}

// GetByID godoc
// @Summary Get a facility by ID
// @Description Get a facility by its ID
// @Tags facilities
// @Accept json
// @Produce json
// @Param id path string true "Facility ID"
// @Success 200 {object} response.Response{data=domain.Facility}
// @Failure 404 {object} response.Response
// @Router /facilities/{id} [get]
func (h *FacilityHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	facility, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrFacilityNotFound {
			response.NotFound(c, "facility not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, facility)
}

// List godoc
// @Summary List facilities
// @Description List facilities with pagination and filters
// @Tags facilities
// @Accept json
// @Produce json
// @Param cruise_id query string false "Cruise ID"
// @Param category_id query string false "Category ID"
// @Param status query string false "Status"
// @Param is_free query bool false "Is free"
// @Param deck_number query int false "Deck number"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.Response{data=[]domain.Facility,pagination=pagination.Paginator}
// @Router /facilities [get]
func (h *FacilityHandler) List(c *gin.Context) {
	var req service.ListFacilitiesRequest
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
// @Summary List facilities by cruise
// @Description Get all facilities for a specific cruise
// @Tags facilities
// @Accept json
// @Produce json
// @Param cruise_id path string true "Cruise ID"
// @Success 200 {object} response.Response{data=[]domain.Facility}
// @Router /cruises/{cruise_id}/facilities [get]
func (h *FacilityHandler) ListByCruise(c *gin.Context) {
	cruiseID := c.Param("cruise_id")

	facilities, err := h.service.ListByCruise(c.Request.Context(), cruiseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, facilities)
}

// ListByCruiseGrouped godoc
// @Summary List facilities grouped by category
// @Description Get all facilities for a cruise grouped by category
// @Tags facilities
// @Accept json
// @Produce json
// @Param cruise_id path string true "Cruise ID"
// @Success 200 {object} response.Response{data=[]domain.FacilityCategory}
// @Router /cruises/{cruise_id}/facilities/grouped [get]
func (h *FacilityHandler) ListByCruiseGrouped(c *gin.Context) {
	cruiseID := c.Param("cruise_id")

	categories, err := h.service.ListByCruiseGrouped(c.Request.Context(), cruiseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, categories)
}

// Update godoc
// @Summary Update a facility
// @Description Update a facility by ID
// @Tags facilities
// @Accept json
// @Produce json
// @Param id path string true "Facility ID"
// @Param request body service.UpdateFacilityRequest true "Update facility request"
// @Success 200 {object} response.Response{data=domain.Facility}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /facilities/{id} [put]
func (h *FacilityHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	facility, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err == service.ErrFacilityNotFound {
			response.NotFound(c, "facility not found")
			return
		}
		if err == service.ErrInvalidFacilityData {
			response.BadRequest(c, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, facility)
}

// Delete godoc
// @Summary Delete a facility
// @Description Soft delete a facility by ID
// @Tags facilities
// @Accept json
// @Produce json
// @Param id path string true "Facility ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /facilities/{id} [delete]
func (h *FacilityHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == service.ErrFacilityNotFound {
			response.NotFound(c, "facility not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Facility Category Handlers

// CreateCategory godoc
// @Summary Create a facility category
// @Description Create a new facility category
// @Tags facility-categories
// @Accept json
// @Produce json
// @Param request body service.CreateFacilityCategoryRequest true "Create category request"
// @Success 201 {object} response.Response{data=domain.FacilityCategory}
// @Failure 400 {object} response.Response
// @Router /facility-categories [post]
func (h *FacilityCategoryHandler) CreateCategory(c *gin.Context) {
	var req service.CreateFacilityCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, category)
}

// ListCategoriesByCruise godoc
// @Summary List facility categories by cruise
// @Description Get all facility categories for a cruise
// @Tags facility-categories
// @Accept json
// @Produce json
// @Param cruise_id path string true "Cruise ID"
// @Success 200 {object} response.Response{data=[]domain.FacilityCategory}
// @Router /cruises/{cruise_id}/facility-categories [get]
func (h *FacilityCategoryHandler) ListCategoriesByCruise(c *gin.Context) {
	cruiseID := c.Param("cruise_id")

	categories, err := h.service.ListByCruise(c.Request.Context(), cruiseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, categories)
}

// UpdateCategory godoc
// @Summary Update a facility category
// @Description Update a facility category by ID
// @Tags facility-categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body service.UpdateFacilityCategoryRequest true "Update category request"
// @Success 200 {object} response.Response{data=domain.FacilityCategory}
// @Failure 404 {object} response.Response
// @Router /facility-categories/{id} [put]
func (h *FacilityCategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateFacilityCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	category, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err == service.ErrFacilityNotFound {
			response.NotFound(c, "facility category not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, category)
}

// DeleteCategory godoc
// @Summary Delete a facility category
// @Description Delete a facility category by ID
// @Tags facility-categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /facility-categories/{id} [delete]
func (h *FacilityCategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == service.ErrFacilityNotFound {
			response.NotFound(c, "facility category not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
