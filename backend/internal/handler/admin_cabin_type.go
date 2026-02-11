package handler

import (
	"backend/internal/pagination"
	"backend/internal/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminCabinTypeHandler handles admin HTTP requests for cabin types
type AdminCabinTypeHandler struct {
	service service.CabinTypeService
}

// NewAdminCabinTypeHandler creates a new admin cabin type handler
func NewAdminCabinTypeHandler(service service.CabinTypeService) *AdminCabinTypeHandler {
	return &AdminCabinTypeHandler{service: service}
}

// List godoc
// @Summary Admin list cabin types
// @Description List all cabin types for admin
// @Tags admin-cabin-types
// @Accept json
// @Produce json
// @Param cruise_id query string false "Cruise ID"
// @Param status query string false "Status"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.Response{data=[]domain.CabinType,pagination=pagination.Paginator}
// @Router /admin/cabin-types [get]
func (h *AdminCabinTypeHandler) List(c *gin.Context) {
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

// GetByID godoc
// @Summary Admin get cabin type by ID
// @Description Get cabin type details
// @Tags admin-cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin Type ID"
// @Success 200 {object} response.Response{data=domain.CabinType}
// @Failure 404 {object} response.Response
// @Router /admin/cabin-types/{id} [get]
func (h *AdminCabinTypeHandler) GetByID(c *gin.Context) {
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

// Create godoc
// @Summary Admin create cabin type
// @Description Create a new cabin type
// @Tags admin-cabin-types
// @Accept json
// @Produce json
// @Param request body service.CreateCabinTypeRequest true "Create cabin type request"
// @Success 201 {object} response.Response{data=domain.CabinType}
// @Failure 400 {object} response.Response
// @Router /admin/cabin-types [post]
func (h *AdminCabinTypeHandler) Create(c *gin.Context) {
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

// Update godoc
// @Summary Admin update cabin type
// @Description Update cabin type details
// @Tags admin-cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin Type ID"
// @Param request body service.UpdateCabinTypeRequest true "Update cabin type request"
// @Success 200 {object} response.Response{data=domain.CabinType}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/cabin-types/{id} [put]
func (h *AdminCabinTypeHandler) Update(c *gin.Context) {
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
// @Summary Admin delete cabin type
// @Description Soft delete a cabin type
// @Tags admin-cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin Type ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/cabin-types/{id} [delete]
func (h *AdminCabinTypeHandler) Delete(c *gin.Context) {
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
// @Summary Admin restore cabin type
// @Description Restore a soft-deleted cabin type
// @Tags admin-cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin Type ID"
// @Success 200 {object} response.Response
// @Router /admin/cabin-types/{id}/restore [post]
func (h *AdminCabinTypeHandler) Restore(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Restore(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus godoc
// @Summary Admin update cabin type status
// @Description Update cabin type status
// @Tags admin-cabin-types
// @Accept json
// @Produce json
// @Param id path string true "Cabin Type ID"
// @Param status body UpdateStatusRequest true "Status update"
// @Success 200 {object} response.Response
// @Router /admin/cabin-types/{id}/status [put]
func (h *AdminCabinTypeHandler) UpdateStatus(c *gin.Context) {
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
