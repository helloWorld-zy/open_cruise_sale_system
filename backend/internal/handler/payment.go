package handler

import (
	"backend/internal/payment"
	"backend/internal/response"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PaymentHandler handles HTTP requests for payments
type PaymentHandler struct {
	service payment.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(service payment.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// Create godoc
// @Summary Create a payment
// @Description Create a payment for an order
// @Tags payments
// @Accept json
// @Produce json
// @Param request body CreatePaymentRequest true "Create payment request"
// @Success 201 {object} response.Response{data=domain.Payment}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /payments [post]
func (h *PaymentHandler) Create(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if exists {
		_ = userID
		// Could add user verification here
	}

	payment, err := h.service.CreatePayment(c.Request.Context(), req.OrderID, req.Method, req.Description)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, payment)
}

// WechatCallback godoc
// @Summary WeChat Pay callback
// @Description Handle WeChat Pay notification callback
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {string} string "success"
// @Failure 400 {string} string "fail"
// @Router /payments/callback/wechat [post]
func (h *PaymentHandler) WechatCallback(c *gin.Context) {
	// Read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// Get signature from headers
	signature := c.GetHeader("Wechatpay-Signature")
	if signature == "" {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// Process callback
	if err := h.service.ProcessCallback(c.Request.Context(), "wechat", body, signature); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// Return success to WeChat
	c.String(http.StatusOK, "success")
}

// Query godoc
// @Summary Query payment status
// @Description Query payment status from provider
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Success 200 {object} response.Response{data=domain.Payment}
// @Failure 404 {object} response.Response
// @Router /payments/{id} [get]
func (h *PaymentHandler) Query(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.QueryPayment(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// GetByOrder godoc
// @Summary Get payment by order
// @Description Get payment information for an order
// @Tags payments
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} response.Response{data=domain.Payment}
// @Failure 404 {object} response.Response
// @Router /payments/order/{orderId} [get]
func (h *PaymentHandler) GetByOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	payment, err := h.service.GetPaymentByOrder(c.Request.Context(), orderID)
	if err != nil {
		response.NotFound(c, "payment not found")
		return
	}

	response.Success(c, payment)
}

// Refund godoc
// @Summary Refund a payment
// @Description Process a refund for a payment
// @Tags payments
// @Accept json
// @Produce json
// @Param id path string true "Payment ID"
// @Param request body RefundRequest true "Refund request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /payments/{id}/refund [post]
func (h *PaymentHandler) Refund(c *gin.Context) {
	id := c.Param("id")

	var req RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.Refund(c.Request.Context(), id, req.Amount, req.Reason); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreatePaymentRequest represents a create payment request
type CreatePaymentRequest struct {
	OrderID     string `json:"order_id" binding:"required"`
	Method      string `json:"method" binding:"required,oneof=wechat alipay card"`
	Description string `json:"description"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Reason string  `json:"reason"`
}
