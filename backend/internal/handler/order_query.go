package handler

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OrderQueryHandler handles order query operations
type OrderQueryHandler struct {
	service service.OrderService
	repo    repository.OrderRepository
}

// NewOrderQueryHandler creates a new order query handler
func NewOrderQueryHandler(service service.OrderService, repo repository.OrderRepository) *OrderQueryHandler {
	return &OrderQueryHandler{
		service: service,
		repo:    repo,
	}
}

// GetMyOrders godoc
// @Summary Get my orders
// @Description Get orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} response.Response{data=[]domain.Order,pagination=pagination.Paginator}
// @Failure 401 {object} response.Response
// @Router /orders/my [get]
func (h *OrderQueryHandler) GetMyOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req service.ListOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Force filter by user
	req.UserID = userID.(string)
	req.Paginator = *pagination.NewPaginator(c)

	result, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetOrderDetail godoc
// @Summary Get order detail
// @Description Get order detail with all related information
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response{data=domain.Order}
// @Failure 404 {object} response.Response
// @Router /orders/{id}/detail [get]
func (h *OrderQueryHandler) GetOrderDetail(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	order, err := h.service.GetWithDetails(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "订单不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Check permission - only order owner or admin can view
	if order.UserID != nil && *order.UserID != userID.(string) {
		role, hasRole := c.Get("role")
		if !hasRole {
			response.Forbidden(c, "无权查看此订单")
			return
		}

		roleStr, ok := role.(string)
		if !ok || (roleStr != "super_admin" && roleStr != "operations" && roleStr != "finance" && roleStr != "customer_service") {
			response.Forbidden(c, "无权查看此订单")
			return
		}
	}

	response.Success(c, order)
}

// GetOrderStatistics godoc
// @Summary Get order statistics
// @Description Get order statistics for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=OrderStatisticsResponse}
// @Failure 401 {object} response.Response
// @Router /orders/statistics [get]
func (h *OrderQueryHandler) GetOrderStatistics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	// Get all orders for user
	// Get all orders for user
	paginator := &pagination.Paginator{Page: 1, PageSize: 1000}

	orders, err := h.repo.ListByUser(c.Request.Context(), userID.(string), paginator)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate statistics
	stats := calculateOrderStatistics(orders)

	response.Success(c, stats)
}

// OrderStatisticsResponse represents order statistics response
type OrderStatisticsResponse struct {
	TotalOrders     int64   `json:"total_orders"`
	PendingOrders   int64   `json:"pending_orders"`
	PaidOrders      int64   `json:"paid_orders"`
	ConfirmedOrders int64   `json:"confirmed_orders"`
	CompletedOrders int64   `json:"completed_orders"`
	CancelledOrders int64   `json:"cancelled_orders"`
	TotalSpent      float64 `json:"total_spent"`
	TotalRefunded   float64 `json:"total_refunded"`
}

func calculateOrderStatistics(orders []*domain.Order) OrderStatisticsResponse {
	var stats OrderStatisticsResponse
	stats.TotalOrders = int64(len(orders))

	for _, order := range orders {
		switch order.Status {
		case domain.OrderStatusPending:
			stats.PendingOrders++
		case domain.OrderStatusPaid:
			stats.PaidOrders++
		case domain.OrderStatusConfirmed:
			stats.ConfirmedOrders++
		case domain.OrderStatusCompleted:
			stats.CompletedOrders++
		case domain.OrderStatusCancelled, domain.OrderStatusRefunded:
			stats.CancelledOrders++
		}

		if order.Status != domain.OrderStatusCancelled && order.Status != domain.OrderStatusRefunded {
			stats.TotalSpent += order.TotalAmount
		}
	}

	return stats
}
