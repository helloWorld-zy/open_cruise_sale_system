package handler

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/response"
	"backend/internal/service"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminCruiseHandler handles admin HTTP requests for cruises
type AdminCruiseHandler struct {
	service       service.CruiseService
	uploadService service.StorageService
}

// NewAdminCruiseHandler creates a new admin cruise handler
func NewAdminCruiseHandler(service service.CruiseService, uploadService service.StorageService) *AdminCruiseHandler {
	return &AdminCruiseHandler{
		service:       service,
		uploadService: uploadService,
	}
}

// List godoc
// @Summary Admin list cruises
// @Description List all cruises for admin with full details
// @Tags admin-cruises
// @Accept json
// @Produce json
// @Param company_id query string false "Company ID"
// @Param status query string false "Status"
// @Param keyword query string false "Search keyword"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} response.Response{data=[]domain.Cruise,pagination=pagination.Paginator}
// @Router /admin/cruises [get]
func (h *AdminCruiseHandler) List(c *gin.Context) {
	var req service.ListCruisesRequest
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
// @Summary Admin get cruise by ID
// @Description Get full cruise details including relations
// @Tags admin-cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Success 200 {object} response.Response{data=domain.Cruise}
// @Failure 404 {object} response.Response
// @Router /admin/cruises/{id} [get]
func (h *AdminCruiseHandler) GetByID(c *gin.Context) {
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

// Create godoc
// @Summary Admin create cruise
// @Description Create a new cruise with image uploads
// @Tags admin-cruises
// @Accept multipart/form-data
// @Produce json
// @Param company_id formData string true "Company ID"
// @Param name_cn formData string true "Chinese name"
// @Param name_en formData string false "English name"
// @Param code formData string true "Cruise code"
// @Param gross_tonnage formData int false "Gross tonnage"
// @Param passenger_capacity formData int false "Passenger capacity"
// @Param crew_count formData int false "Crew count"
// @Param built_year formData int false "Built year"
// @Param renovated_year formData int false "Renovated year"
// @Param length_meters formData number false "Length in meters"
// @Param width_meters formData number false "Width in meters"
// @Param deck_count formData int false "Deck count"
// @Param status formData string false "Status"
// @Param sort_weight formData int false "Sort weight"
// @Param images formData file false "Cover images"
// @Success 201 {object} response.Response{data=domain.Cruise}
// @Failure 400 {object} response.Response
// @Router /admin/cruises [post]
func (h *AdminCruiseHandler) Create(c *gin.Context) {
	var req service.CreateCruiseRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Handle image uploads
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["images"]
		for _, file := range files {
			opened, err := file.Open()
			if err != nil {
				response.BadRequest(c, "Failed to process image: "+err.Error())
				return
			}
			defer opened.Close()

			url, err := h.uploadService.UploadImage(c.Request.Context(), opened, file.Filename, file.Size)
			if err != nil {
				response.BadRequest(c, "Failed to upload image: "+err.Error())
				return
			}
			req.CoverImages = append(req.CoverImages, url)
		}
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

// Update godoc
// @Summary Admin update cruise
// @Description Update cruise details with optional image uploads
// @Tags admin-cruises
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Cruise ID"
// @Param name_cn formData string false "Chinese name"
// @Param name_en formData string false "English name"
// @Param gross_tonnage formData int false "Gross tonnage"
// @Param passenger_capacity formData int false "Passenger capacity"
// @Param crew_count formData int false "Crew count"
// @Param built_year formData int false "Built year"
// @Param renovated_year formData int false "Renovated year"
// @Param length_meters formData number false "Length in meters"
// @Param width_meters formData number false "Width in meters"
// @Param deck_count formData int false "Deck count"
// @Param status formData string false "Status"
// @Param sort_weight formData int false "Sort weight"
// @Param images formData file false "New cover images"
// @Param remove_images formData string false "JSON array of image URLs to remove"
// @Success 200 {object} response.Response{data=domain.Cruise}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/cruises/{id} [put]
func (h *AdminCruiseHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateCruiseRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get existing cruise to merge images
	existing, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrCruiseNotFound {
			response.NotFound(c, "cruise not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Start with existing images
	var currentImages []string
	if existing.CoverImages != nil {
		currentImages = existing.CoverImages
	}

	// Handle image removals
	removeImages := c.PostForm("remove_images")
	if removeImages != "" {
		// Parse and remove specified images
		// Implementation would parse JSON array and remove matching URLs
	}

	// Handle new image uploads
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["images"]
		for _, file := range files {
			opened, err := file.Open()
			if err != nil {
				response.BadRequest(c, "Failed to process image: "+err.Error())
				return
			}
			defer opened.Close()

			url, err := h.uploadService.UploadImage(c.Request.Context(), opened, file.Filename, file.Size)
			if err != nil {
				response.BadRequest(c, "Failed to upload image: "+err.Error())
				return
			}
			currentImages = append(currentImages, url)
		}
	}

	req.CoverImages = currentImages

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
// @Summary Admin delete cruise
// @Description Soft delete a cruise
// @Tags admin-cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/cruises/{id} [delete]
func (h *AdminCruiseHandler) Delete(c *gin.Context) {
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
// @Summary Admin restore cruise
// @Description Restore a soft-deleted cruise
// @Tags admin-cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Success 200 {object} response.Response
// @Router /admin/cruises/{id}/restore [post]
func (h *AdminCruiseHandler) Restore(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Restore(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus godoc
// @Summary Admin update cruise status
// @Description Update cruise status
// @Tags admin-cruises
// @Accept json
// @Produce json
// @Param id path string true "Cruise ID"
// @Param status body UpdateStatusRequest true "Status update"
// @Success 200 {object} response.Response
// @Router /admin/cruises/{id}/status [put]
func (h *AdminCruiseHandler) UpdateStatus(c *gin.Context) {
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

// UploadImages godoc
// @Summary Upload cruise images
// @Description Upload multiple images for a cruise
// @Tags admin-cruises
// @Accept multipart/form-data
// @Produce json
// @Param images formData file true "Images to upload"
// @Success 200 {object} response.Response{data=[]string}
// @Router /admin/cruises/upload [post]
func (h *AdminCruiseHandler) UploadImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "Failed to parse form: "+err.Error())
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		response.BadRequest(c, "No images provided")
		return
	}

	var urls []string
	for _, file := range files {
		opened, err := file.Open()
		if err != nil {
			response.BadRequest(c, "Failed to process image: "+err.Error())
			return
		}

		url, err := h.uploadService.UploadImage(c.Request.Context(), opened, file.Filename, file.Size)
		opened.Close()

		if err != nil {
			response.BadRequest(c, "Failed to upload image: "+err.Error())
			return
		}
		urls = append(urls, url)
	}

	response.Success(c, urls)
}
