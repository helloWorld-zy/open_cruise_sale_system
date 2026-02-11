package handler

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/response"
	"backend/internal/service"
	"backend/internal/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	service service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// Create godoc
// @Summary Create a new order
// @Description Create a new booking order with inventory locking
// @Tags orders
// @Accept json
// @Produce json
// @Param request body service.CreateOrderRequest true "Create order request"
// @Success 201 {object} response.Response{data=domain.Order}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Extract user ID from context if authenticated
	if userID, exists := c.Get("userID"); exists {
		req.UserID = userID.(string)
	}

	order, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrInvalidOrderData || err == service.ErrInvalidPassengerCount {
			response.BadRequest(c, err.Error())
			return
		}
		if err == service.ErrCabinNotAvailable {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, order)
}

// Calculate godoc
// @Summary Calculate order total
// @Description Calculate total price for order items before creating order
// @Tags orders
// @Accept json
// @Produce json
// @Param request body CalculateRequest true "Calculate request"
// @Success 200 {object} response.Response{data=service.OrderCalculation}
// @Failure 400 {object} response.Response
// @Router /orders/calculate [post]
func (h *OrderHandler) Calculate(c *gin.Context) {
	var req CalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	calculation, err := h.service.CalculateTotal(c.Request.Context(), req.Items)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, calculation)
}

// GetByID godoc
// @Summary Get an order by ID
// @Description Get an order by its ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response{data=domain.Order}
// @Failure 404 {object} response.Response
// @Router /orders/{id} [get]
func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	order, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "order not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, order)
}

// GetByOrderNumber godoc
// @Summary Get an order by order number
// @Description Get an order by its unique order number
// @Tags orders
// @Accept json
// @Produce json
// @Param orderNumber path string true "Order Number"
// @Success 200 {object} response.Response{data=domain.Order}
// @Failure 404 {object} response.Response
// @Router /orders/number/{orderNumber} [get]
func (h *OrderHandler) GetByOrderNumber(c *gin.Context) {
	orderNumber := c.Param("orderNumber")

	order, err := h.service.GetByOrderNumber(c.Request.Context(), orderNumber)
	if err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "order not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, order)
}

// GetWithDetails godoc
// @Summary Get order with all details
// @Description Get order with items, passengers, and payments
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response{data=domain.Order}
// @Failure 404 {object} response.Response
// @Router /orders/{id}/details [get]
func (h *OrderHandler) GetWithDetails(c *gin.Context) {
	id := c.Param("id")

	order, err := h.service.GetWithDetails(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "order not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, order)
}

// List godoc
// @Summary List orders
// @Description List orders with pagination and filters
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Param voyage_id query string false "Voyage ID"
// @Param status query string false "Status (pending, paid, confirmed, completed, cancelled, refunded)"
// @Param payment_status query string false "Payment status (unpaid, partial, paid, refunded)"
// @Param order_number query string false "Order number"
// @Param date_from query string false "Date from (RFC3339)"
// @Param date_to query string false "Date to (RFC3339)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} response.Response{data=[]domain.Order,pagination=pagination.Paginator}
// @Router /orders [get]
func (h *OrderHandler) List(c *gin.Context) {
	var req service.ListOrdersRequest
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

// ListMyOrders godoc
// @Summary List my orders
// @Description List orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} response.Response{data=[]domain.Order,pagination=pagination.Paginator}
// @Router /orders/my [get]
func (h *OrderHandler) ListMyOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	paginator := pagination.NewPaginator(c)
	result, err := h.service.ListByUser(c.Request.Context(), userID.(string), paginator)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// Update godoc
// @Summary Update an order
// @Description Update order contact information (only pending orders)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body service.UpdateOrderRequest true "Update order request"
// @Success 200 {object} response.Response{data=domain.Order}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /orders/{id} [put]
func (h *OrderHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	order, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "order not found")
			return
		}
		if err == service.ErrInvalidOrderData {
			response.BadRequest(c, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, order)
}

// Cancel godoc
// @Summary Cancel an order
// @Description Cancel an order and release inventory
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Cancel(c.Request.Context(), id); err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "order not found")
			return
		}
		if err == service.ErrOrderNotCancellable {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Confirm godoc
// @Summary Confirm an order
// @Description Confirm a paid order and confirm inventory booking
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /orders/{id}/confirm [post]
func (h *OrderHandler) Confirm(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Confirm(c.Request.Context(), id); err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "order not found")
			return
		}
		if err == service.ErrInvalidOrderTransition {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Complete godoc
// @Summary Complete an order
// @Description Complete a confirmed order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /orders/{id}/complete [post]
func (h *OrderHandler) Complete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Complete(c.Request.Context(), id); err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "order not found")
			return
		}
		if err == service.ErrInvalidOrderTransition {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete godoc
// @Summary Delete an order
// @Description Delete an order (admin only, only pending or cancelled)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /orders/{id} [delete]
func (h *OrderHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "order not found")
			return
		}
		if err == service.ErrInvalidOrderData {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// CalculateRequest represents a calculation request
type CalculateRequest struct {
	Items []service.OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}
