package handler

import (
	"backend/internal/pagination"
	"backend/internal/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminFacilityHandler handles admin HTTP requests for facilities
type AdminFacilityHandler struct {
	service service.FacilityService
}

// AdminFacilityCategoryHandler handles admin HTTP requests for facility categories
type AdminFacilityCategoryHandler struct {
	service service.FacilityCategoryService
}

// NewAdminFacilityHandler creates a new admin facility handler
func NewAdminFacilityHandler(service service.FacilityService) *AdminFacilityHandler {
	return &AdminFacilityHandler{service: service}
}

// NewAdminFacilityCategoryHandler creates a new admin facility category handler
func NewAdminFacilityCategoryHandler(service service.FacilityCategoryService) *AdminFacilityCategoryHandler {
	return &AdminFacilityCategoryHandler{service: service}
}

// Facility handlers

// List godoc
// @Summary Admin list facilities
// @Description List all facilities for admin
// @Tags admin-facilities
// @Accept json
// @Produce json
// @Param cruise_id query string false "Cruise ID"
// @Param category_id query string false "Category ID"
// @Param status query string false "Status"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.Response{data=[]domain.Facility,pagination=pagination.Paginator}
// @Router /admin/facilities [get]
func (h *AdminFacilityHandler) List(c *gin.Context) {
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

// GetByID godoc
// @Summary Admin get facility by ID
// @Description Get facility details
// @Tags admin-facilities
// @Accept json
// @Produce json
// @Param id path string true "Facility ID"
// @Success 200 {object} response.Response{data=domain.Facility}
// @Failure 404 {object} response.Response
// @Router /admin/facilities/{id} [get]
func (h *AdminFacilityHandler) GetByID(c *gin.Context) {
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

// Create godoc
// @Summary Admin create facility
// @Description Create a new facility
// @Tags admin-facilities
// @Accept json
// @Produce json
// @Param request body service.CreateFacilityRequest true "Create facility request"
// @Success 201 {object} response.Response{data=domain.Facility}
// @Failure 400 {object} response.Response
// @Router /admin/facilities [post]
func (h *AdminFacilityHandler) Create(c *gin.Context) {
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

// Update godoc
// @Summary Admin update facility
// @Description Update facility details
// @Tags admin-facilities
// @Accept json
// @Produce json
// @Param id path string true "Facility ID"
// @Param request body service.UpdateFacilityRequest true "Update facility request"
// @Success 200 {object} response.Response{data=domain.Facility}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/facilities/{id} [put]
func (h *AdminFacilityHandler) Update(c *gin.Context) {
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
// @Summary Admin delete facility
// @Description Soft delete a facility
// @Tags admin-facilities
// @Accept json
// @Produce json
// @Param id path string true "Facility ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/facilities/{id} [delete]
func (h *AdminFacilityHandler) Delete(c *gin.Context) {
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

// Facility Category handlers

// ListCategories godoc
// @Summary Admin list facility categories
// @Description List all facility categories
// @Tags admin-facility-categories
// @Accept json
// @Produce json
// @Param cruise_id query string false "Cruise ID"
// @Success 200 {object} response.Response{data=[]domain.FacilityCategory}
// @Router /admin/facility-categories [get]
func (h *AdminFacilityCategoryHandler) ListCategories(c *gin.Context) {
	cruiseID := c.Query("cruise_id")

	categories, err := h.service.ListByCruise(c.Request.Context(), cruiseID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, categories)
}

// CreateCategory godoc
// @Summary Admin create facility category
// @Description Create a new facility category
// @Tags admin-facility-categories
// @Accept json
// @Produce json
// @Param request body service.CreateFacilityCategoryRequest true "Create category request"
// @Success 201 {object} response.Response{data=domain.FacilityCategory}
// @Failure 400 {object} response.Response
// @Router /admin/facility-categories [post]
func (h *AdminFacilityCategoryHandler) CreateCategory(c *gin.Context) {
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

// UpdateCategory godoc
// @Summary Admin update facility category
// @Description Update facility category
// @Tags admin-facility-categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body service.UpdateFacilityCategoryRequest true "Update category request"
// @Success 200 {object} response.Response{data=domain.FacilityCategory}
// @Failure 404 {object} response.Response
// @Router /admin/facility-categories/{id} [put]
func (h *AdminFacilityCategoryHandler) UpdateCategory(c *gin.Context) {
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
// @Summary Admin delete facility category
// @Description Delete a facility category
// @Tags admin-facility-categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/facility-categories/{id} [delete]
func (h *AdminFacilityCategoryHandler) DeleteCategory(c *gin.Context) {
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
