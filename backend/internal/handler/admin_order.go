package handler

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"backend/internal/response"
	"backend/internal/service"
	"backend/internal/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOrderHandler handles admin order operations
type AdminOrderHandler struct {
	orderService  service.OrderService
	refundService service.RefundService
	repo          repository.OrderRepository
}

// NewAdminOrderHandler creates a new admin order handler
func NewAdminOrderHandler(
	orderService service.OrderService,
	refundService service.RefundService,
	repo repository.OrderRepository,
) *AdminOrderHandler {
	return &AdminOrderHandler{
		orderService:  orderService,
		refundService: refundService,
		repo:          repo,
	}
}

// ListOrders godoc
// @Summary List all orders (Admin)
// @Description List all orders with filters for admin
// @Tags admin-orders
// @Accept json
// @Produce json
// @Param status query string false "Order status"
// @Param payment_status query string false "Payment status"
// @Param order_number query string false "Order number"
// @Param date_from query string false "Date from (RFC3339)"
// @Param date_to query string false "Date to (RFC3339)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} response.Response{data=[]domain.Order,pagination=pagination.Paginator}
// @Failure 403 {object} response.Response
// @Router /admin/orders [get]
func (h *AdminOrderHandler) ListOrders(c *gin.Context) {
	var req service.ListOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	req.Paginator = *pagination.NewPaginator(c)

	result, err := h.orderService.List(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetOrderDetail godoc
// @Summary Get order detail (Admin)
// @Description Get full order details for admin
// @Tags admin-orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} response.Response{data=domain.Order}
// @Failure 404 {object} response.Response
// @Router /admin/orders/{id} [get]
func (h *AdminOrderHandler) GetOrderDetail(c *gin.Context) {
	id := c.Param("id")

	order, err := h.orderService.GetWithDetails(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrOrderNotFound {
			response.NotFound(c, "订单不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, order)
}

// UpdateOrderStatus godoc
// @Summary Update order status (Admin)
// @Description Update order status manually (for admin override)
// @Tags admin-orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body UpdateOrderStatusRequest true "Status update request"
// @Success 200 {object} response.Response{data=domain.Order}
// @Failure 400 {object} response.Response
// @Router /admin/orders/{id}/status [put]
func (h *AdminOrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.Validate(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Validate status
	if !isValidOrderStatus(req.Status) {
		response.BadRequest(c, "无效的订单状态")
		return
	}

	if err := h.repo.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Get updated order
	order, err := h.orderService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, order)
}

// ListRefunds godoc
// @Summary List all refund requests (Admin)
// @Description List all refund requests with filters
// @Tags admin-refunds
// @Accept json
// @Produce json
// @Param status query string false "Refund status"
// @Param order_id query string false "Order ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} response.Response{data=[]domain.RefundRequest,pagination=pagination.Paginator}
// @Router /admin/refunds [get]
func (h *AdminOrderHandler) ListRefunds(c *gin.Context) {
	var filters repository.RefundFilters
	filters.Status = c.Query("status")
	filters.OrderID = c.Query("order_id")

	paginator := pagination.NewPaginator(c)

	refunds, err := h.refundService.ListRefunds(c.Request.Context(), filters, paginator)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Count total
	paginator.SetTotal(int64(len(refunds))) // TODO: Implement proper count

	result := pagination.Result{
		Data:       refunds,
		Pagination: *paginator,
	}

	response.Success(c, result)
}

// GetRefundDetail godoc
// @Summary Get refund detail (Admin)
// @Description Get refund request details
// @Tags admin-refunds
// @Accept json
// @Produce json
// @Param id path string true "Refund ID"
// @Success 200 {object} response.Response{data=domain.RefundRequest}
// @Failure 404 {object} response.Response
// @Router /admin/refunds/{id} [get]
func (h *AdminOrderHandler) GetRefundDetail(c *gin.Context) {
	id := c.Param("id")

	refund, err := h.refundService.GetRefundByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrRefundNotFound {
			response.NotFound(c, "退款申请不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, refund)
}

// ApproveRefund godoc
// @Summary Approve refund request (Admin)
// @Description Approve a pending refund request
// @Tags admin-refunds
// @Accept json
// @Produce json
// @Param id path string true "Refund ID"
// @Param request body ReviewRefundRequest true "Review request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/refunds/{id}/approve [post]
func (h *AdminOrderHandler) ApproveRefund(c *gin.Context) {
	id := c.Param("id")
	adminID, _ := c.Get("userID")

	var req ReviewRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.refundService.ApproveRefund(c.Request.Context(), id, adminID.(string), req.Note); err != nil {
		if err == service.ErrRefundNotFound {
			response.NotFound(c, "退款申请不存在")
			return
		}
		if err == service.ErrRefundNotReviewable {
			response.BadRequest(c, "该退款申请无法审核")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// RejectRefund godoc
// @Summary Reject refund request (Admin)
// @Description Reject a pending refund request
// @Tags admin-refunds
// @Accept json
// @Produce json
// @Param id path string true "Refund ID"
// @Param request body ReviewRefundRequest true "Review request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /admin/refunds/{id}/reject [post]
func (h *AdminOrderHandler) RejectRefund(c *gin.Context) {
	id := c.Param("id")
	adminID, _ := c.Get("userID")

	var req ReviewRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Note == "" {
		response.BadRequest(c, "请填写拒绝原因")
		return
	}

	if err := h.refundService.RejectRefund(c.Request.Context(), id, adminID.(string), req.Note); err != nil {
		if err == service.ErrRefundNotFound {
			response.NotFound(c, "退款申请不存在")
			return
		}
		if err == service.ErrRefundNotReviewable {
			response.BadRequest(c, "该退款申请无法审核")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// ProcessRefund godoc
// @Summary Process refund payment (Admin)
// @Description Execute the actual refund payment after approval
// @Tags admin-refunds
// @Accept json
// @Produce json
// @Param id path string true "Refund ID"
// @Success 200 {object} response.Response{data=domain.RefundRequest}
// @Failure 400 {object} response.Response
// @Router /admin/refunds/{id}/process [post]
func (h *AdminOrderHandler) ProcessRefund(c *gin.Context) {
	id := c.Param("id")

	if err := h.refundService.ProcessRefund(c.Request.Context(), id); err != nil {
		if err == service.ErrRefundNotFound {
			response.NotFound(c, "退款申请不存在")
			return
		}
		if err == service.ErrRefundNotProcessable {
			response.BadRequest(c, "该退款申请无法处理")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Get updated refund
	refund, _ := h.refundService.GetRefundByID(c.Request.Context(), id)
	response.Success(c, refund)
}

// GetOrderStatistics godoc
// @Summary Get order statistics (Admin)
// @Description Get order statistics for dashboard
// @Tags admin-orders
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=OrderStatisticsResponse}
// @Router /admin/orders/statistics [get]
func (h *AdminOrderHandler) GetOrderStatistics(c *gin.Context) {
	// Get all orders for statistics
	filters := repository.OrderFilters{}
	paginator := &pagination.Paginator{Page: 1, PageSize: 10000}

	orders, err := h.repo.List(c.Request.Context(), filters, paginator)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	stats := calculateAdminOrderStatistics(orders)
	response.Success(c, stats)
}

// UpdateOrderStatusRequest represents a status update request
type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required"`
	Note   string `json:"note,omitempty"`
}

// ReviewRefundRequest represents a refund review request
type ReviewRefundRequest struct {
	Note string `json:"note"`
}

// AdminOrderStatistics represents order statistics for admin
type AdminOrderStatistics struct {
	TotalOrders     int64   `json:"total_orders"`
	PendingOrders   int64   `json:"pending_orders"`
	PaidOrders      int64   `json:"paid_orders"`
	ConfirmedOrders int64   `json:"confirmed_orders"`
	CompletedOrders int64   `json:"completed_orders"`
	CancelledOrders int64   `json:"cancelled_orders"`
	RefundedOrders  int64   `json:"refunded_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	TotalRefunded   float64 `json:"total_refunded"`
	PendingRefunds  int64   `json:"pending_refunds"`
	TodayOrders     int64   `json:"today_orders"`
	WeekOrders      int64   `json:"week_orders"`
	MonthOrders     int64   `json:"month_orders"`
}

func calculateAdminOrderStatistics(orders []*domain.Order) AdminOrderStatistics {
	var stats AdminOrderStatistics
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
			stats.TotalRevenue += order.TotalAmount
		case domain.OrderStatusCancelled:
			stats.CancelledOrders++
		case domain.OrderStatusRefunded:
			stats.RefundedOrders++
			stats.TotalRefunded += order.TotalAmount
		}

		if order.Status == domain.OrderStatusPaid || order.Status == domain.OrderStatusConfirmed || order.Status == domain.OrderStatusCompleted {
			stats.TotalRevenue += order.TotalAmount
		}
	}

	return stats
}

func isValidOrderStatus(status string) bool {
	switch status {
	case domain.OrderStatusPending, domain.OrderStatusPaid, domain.OrderStatusConfirmed,
		domain.OrderStatusCompleted, domain.OrderStatusCancelled, domain.OrderStatusRefunded:
		return true
	}
	return false
}
