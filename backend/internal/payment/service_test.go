package payment

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock Order Repository for Payment
type MockPaymentOrderRepository struct {
	mock.Mock
}

func (m *MockPaymentOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockPaymentOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	args := m.Called(ctx, orderNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockPaymentOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) UpdatePaymentStatus(ctx context.Context, id string, paymentStatus string, paidAmount float64) error {
	args := m.Called(ctx, id, paymentStatus, paidAmount)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) GetOrderWithDetails(ctx context.Context, id string) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockPaymentOrderRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentOrderRepository) GetPaymentByNo(ctx context.Context, paymentNo string) (*domain.Payment, error) {
	args := m.Called(ctx, paymentNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentOrderRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) List(ctx context.Context, filters repository.OrderFilters, paginator *pagination.Paginator) ([]*domain.Order, error) {
	args := m.Called(ctx, filters, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockPaymentOrderRepository) Count(ctx context.Context, filters repository.OrderFilters) (int64, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPaymentOrderRepository) ListByUser(ctx context.Context, userID string, paginator *pagination.Paginator) ([]*domain.Order, error) {
	args := m.Called(ctx, userID, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockPaymentOrderRepository) ListByStatus(ctx context.Context, status string, paginator *pagination.Paginator) ([]*domain.Order, error) {
	args := m.Called(ctx, status, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockPaymentOrderRepository) ListByVoyage(ctx context.Context, voyageID string) ([]*domain.Order, error) {
	args := m.Called(ctx, voyageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockPaymentOrderRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) CreateOrderItem(ctx context.Context, item *domain.OrderItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) GetOrderItemByID(ctx context.Context, id string) (*domain.OrderItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OrderItem), args.Error(1)
}

func (m *MockPaymentOrderRepository) UpdateOrderItem(ctx context.Context, item *domain.OrderItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) UpdateOrderItemStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) DeleteOrderItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) CreatePassenger(ctx context.Context, passenger *domain.Passenger) error {
	args := m.Called(ctx, passenger)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) GetPassengerByID(ctx context.Context, id string) (*domain.Passenger, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Passenger), args.Error(1)
}

func (m *MockPaymentOrderRepository) ListPassengersByOrder(ctx context.Context, orderID string) ([]*domain.Passenger, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Passenger), args.Error(1)
}

func (m *MockPaymentOrderRepository) ListPassengersByOrderItem(ctx context.Context, orderItemID string) ([]*domain.Passenger, error) {
	args := m.Called(ctx, orderItemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Passenger), args.Error(1)
}

func (m *MockPaymentOrderRepository) UpdatePassenger(ctx context.Context, passenger *domain.Passenger) error {
	args := m.Called(ctx, passenger)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) DeletePassenger(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) BatchCreatePassengers(ctx context.Context, passengers []*domain.Passenger) error {
	args := m.Called(ctx, passengers)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) CreateRefundRequest(ctx context.Context, refund *domain.RefundRequest) error {
	args := m.Called(ctx, refund)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) GetRefundRequestByID(ctx context.Context, id string) (*domain.RefundRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefundRequest), args.Error(1)
}

func (m *MockPaymentOrderRepository) ListRefundRequests(ctx context.Context, filters repository.RefundFilters, paginator *pagination.Paginator) ([]*domain.RefundRequest, error) {
	args := m.Called(ctx, filters, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.RefundRequest), args.Error(1)
}

func (m *MockPaymentOrderRepository) UpdateRefundRequest(ctx context.Context, refund *domain.RefundRequest) error {
	args := m.Called(ctx, refund)
	return args.Error(0)
}

func (m *MockPaymentOrderRepository) WithTransaction(ctx context.Context, fn func(repo repository.OrderRepository, tx *gorm.DB) error) error {
	return fn(m, nil)
}

func (m *MockPaymentOrderRepository) ListOrderItemsByOrder(ctx context.Context, orderID string) ([]*domain.OrderItem, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.OrderItem), args.Error(1)
}

// Mock Payment Provider
type MockPaymentProvider struct {
	mock.Mock
}

func (m *MockPaymentProvider) CreatePayment(ctx context.Context, order *domain.Order, description string) (*PaymentResult, error) {
	args := m.Called(ctx, order, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PaymentResult), args.Error(1)
}

func (m *MockPaymentProvider) QueryPayment(ctx context.Context, paymentNo string) (*PaymentQueryResult, error) {
	args := m.Called(ctx, paymentNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PaymentQueryResult), args.Error(1)
}

func (m *MockPaymentProvider) ProcessCallback(ctx context.Context, body []byte, signature string) (*CallbackResult, error) {
	args := m.Called(ctx, body, signature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*CallbackResult), args.Error(1)
}

func (m *MockPaymentProvider) VerifySignature(body []byte, signature string) bool {
	args := m.Called(body, signature)
	return args.Bool(0)
}

func (m *MockPaymentProvider) Refund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*RefundResult, error) {
	args := m.Called(ctx, payment, amount, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RefundResult), args.Error(1)
}

// Mock NATS Connection
type MockNatsConn struct {
	mock.Mock
}

func (m *MockNatsConn) Publish(subj string, data []byte) error {
	args := m.Called(subj, data)
	return args.Error(0)
}

func (m *MockNatsConn) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	args := m.Called(subj, cb)
	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func TestPaymentService_CreatePayment(t *testing.T) {
	mockOrderRepo := new(MockPaymentOrderRepository)
	mockProvider := new(MockPaymentProvider)
	service := NewPaymentService(mockOrderRepo, nil)
	service.RegisterProvider("wechat", mockProvider)

	ctx := context.Background()

	t.Run("should create payment for pending order", func(t *testing.T) {
		order := &domain.Order{
			BaseModel:      domain.BaseModel{ID: uuid.New()},
			Status:         domain.OrderStatusPending,
			TotalAmount:    1000,
			DiscountAmount: 0,
			Currency:       "CNY",
		}

		paymentResult := &PaymentResult{
			PaymentNo: "PAY20240101123456",
			PayURL:    "https://pay.example.com/qr",
		}

		payment := &domain.Payment{
			PaymentNo: "PAY20240101123456",
		}

		mockOrderRepo.On("GetByID", ctx, "order-1").Return(order, nil).Once()
		mockProvider.On("CreatePayment", ctx, order, mock.AnythingOfType("string")).Return(paymentResult, nil).Once()
		mockOrderRepo.On("GetPaymentByNo", ctx, "PAY20240101123456").Return(payment, nil).Once()

		result, err := service.CreatePayment(ctx, "order-1", "wechat", "订单支付")

		assert.NoError(t, err)
		assert.Equal(t, "PAY20240101123456", result.PaymentNo)
		mockOrderRepo.AssertExpectations(t)
		mockProvider.AssertExpectations(t)
	})

	t.Run("should return error for non-pending order", func(t *testing.T) {
		order := &domain.Order{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			Status:    domain.OrderStatusPaid,
		}

		mockOrderRepo.On("GetByID", ctx, "order-1").Return(order, nil).Once()

		result, err := service.CreatePayment(ctx, "order-1", "wechat", "订单支付")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("should return error for unsupported payment method", func(t *testing.T) {
		order := &domain.Order{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			Status:    domain.OrderStatusPending,
		}

		mockOrderRepo.On("GetByID", ctx, "order-1").Return(order, nil).Once()

		result, err := service.CreatePayment(ctx, "order-1", "unsupported", "订单支付")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestPaymentService_ProcessCallback(t *testing.T) {
	mockOrderRepo := new(MockPaymentOrderRepository)
	mockProvider := new(MockPaymentProvider)
	service := NewPaymentService(mockOrderRepo, nil)
	service.RegisterProvider("wechat", mockProvider)

	ctx := context.Background()

	t.Run("should process payment callback successfully", func(t *testing.T) {
		callbackBody := []byte(`{"payment_no": "PAY20240101123456", "status": "success"}`)
		signature := "valid-signature"

		callbackResult := &CallbackResult{
			PaymentNo:    "PAY20240101123456",
			ThirdPartyID: "WXP20240101123456",
			Amount:       1000,
			Status:       domain.PaymentStatusSuccess,
			PaidAt:       "2024-01-01T12:00:00Z",
		}

		payment := &domain.Payment{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			OrderID:   "order-1",
			PaymentNo: "PAY20240101123456",
			Status:    domain.PaymentStatusPending,
			Amount:    1000,
		}

		mockProvider.On("ProcessCallback", ctx, callbackBody, signature).Return(callbackResult, nil).Once()
		mockOrderRepo.On("GetPaymentByNo", ctx, "PAY20240101123456").Return(payment, nil).Once()
		mockOrderRepo.On("UpdatePayment", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil).Once()
		mockOrderRepo.On("UpdatePaymentStatus", ctx, "order-1", domain.PaymentStatusPaid, float64(1000)).Return(nil).Once()

		err := service.ProcessCallback(ctx, "wechat", callbackBody, signature)

		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
		mockProvider.AssertExpectations(t)
	})

	t.Run("should return error for invalid signature", func(t *testing.T) {
		callbackBody := []byte(`{"payment_no": "PAY20240101123456"}`)
		signature := "invalid-signature"

		mockProvider.On("ProcessCallback", ctx, callbackBody, signature).Return(nil, ErrInvalidSignature).Once()

		err := service.ProcessCallback(ctx, "wechat", callbackBody, signature)

		assert.ErrorIs(t, err, ErrInvalidSignature)
		mockProvider.AssertExpectations(t)
	})

	t.Run("should return error for unsupported provider", func(t *testing.T) {
		callbackBody := []byte(`{}`)
		signature := "signature"

		err := service.ProcessCallback(ctx, "unsupported", callbackBody, signature)

		assert.Error(t, err)
	})
}

func TestPaymentService_QueryPayment(t *testing.T) {
	mockOrderRepo := new(MockPaymentOrderRepository)
	mockProvider := new(MockPaymentProvider)
	service := NewPaymentService(mockOrderRepo, nil)
	service.RegisterProvider("wechat", mockProvider)

	ctx := context.Background()

	t.Run("should query and update payment status", func(t *testing.T) {
		payment := &domain.Payment{
			BaseModel:     domain.BaseModel{ID: uuid.New()},
			PaymentNo:     "PAY20240101123456",
			PaymentMethod: "wechat",
			Status:        domain.PaymentStatusPending,
			Amount:        1000,
		}

		queryResult := &PaymentQueryResult{
			Status:       domain.PaymentStatusSuccess,
			Amount:       1000,
			ThirdPartyID: "WXP20240101123456",
			PaidAt:       "2024-01-01T12:00:00Z",
		}

		mockOrderRepo.On("GetPaymentByID", ctx, "payment-1").Return(payment, nil).Once()
		mockProvider.On("QueryPayment", ctx, "PAY20240101123456").Return(queryResult, nil).Once()
		mockOrderRepo.On("UpdatePaymentStatus", ctx, payment.OrderID, domain.PaymentStatusPaid, float64(1000)).Return(nil).Once()
		mockOrderRepo.On("UpdatePayment", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil).Once()

		result, err := service.QueryPayment(ctx, "payment-1")

		assert.NoError(t, err)
		assert.Equal(t, domain.PaymentStatusSuccess, result.Status)
		mockOrderRepo.AssertExpectations(t)
		mockProvider.AssertExpectations(t)
	})

	t.Run("should return error when payment not found", func(t *testing.T) {
		mockOrderRepo.On("GetPaymentByID", ctx, "non-existent").Return(nil, errors.New("not found")).Once()

		result, err := service.QueryPayment(ctx, "non-existent")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestPaymentService_Refund(t *testing.T) {
	mockOrderRepo := new(MockPaymentOrderRepository)
	mockProvider := new(MockPaymentProvider)
	service := NewPaymentService(mockOrderRepo, nil)
	service.RegisterProvider("wechat", mockProvider)

	ctx := context.Background()

	t.Run("should process refund for successful payment", func(t *testing.T) {
		payment := &domain.Payment{
			BaseModel:     domain.BaseModel{ID: uuid.New()},
			OrderID:       "order-1",
			PaymentNo:     "PAY20240101123456",
			PaymentMethod: "wechat",
			Status:        domain.PaymentStatusSuccess,
			Amount:        1000,
			Currency:      "CNY",
		}

		refundResult := &RefundResult{
			RefundNo:     "REF20240101123456",
			ThirdPartyID: "WXP_REF_123",
			Status:       "SUCCESS",
		}

		mockOrderRepo.On("GetPaymentByID", ctx, "payment-1").Return(payment, nil).Once()
		mockProvider.On("Refund", ctx, payment, float64(500), "Customer request").Return(refundResult, nil).Once()
		mockOrderRepo.On("UpdatePayment", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil).Once()
		mockOrderRepo.On("UpdateStatus", ctx, "order-1", domain.OrderStatusRefunded).Return(nil).Once()

		err := service.Refund(ctx, "payment-1", 500, "Customer request")

		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
		mockProvider.AssertExpectations(t)
	})

	t.Run("should return error for non-successful payment", func(t *testing.T) {
		payment := &domain.Payment{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			Status:    domain.PaymentStatusPending,
		}

		mockOrderRepo.On("GetPaymentByID", ctx, "payment-1").Return(payment, nil).Once()

		err := service.Refund(ctx, "payment-1", 500, "Customer request")

		assert.Error(t, err)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestPaymentService_GetPaymentByOrder(t *testing.T) {
	mockOrderRepo := new(MockPaymentOrderRepository)
	service := NewPaymentService(mockOrderRepo, nil)

	ctx := context.Background()

	t.Run("should return payment by order ID", func(t *testing.T) {
		payment := &domain.Payment{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			OrderID:   "order-1",
			PaymentNo: "PAY20240101123456",
		}

		order := &domain.Order{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			Payments:  []domain.Payment{*payment},
		}

		mockOrderRepo.On("GetOrderWithDetails", ctx, "order-1").Return(order, nil).Once()

		result, err := service.GetPaymentByOrder(ctx, "order-1")

		assert.NoError(t, err)
		assert.Equal(t, payment.ID, result.ID)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("should return error when no payments found", func(t *testing.T) {
		order := &domain.Order{
			BaseModel: domain.BaseModel{ID: uuid.New()},
			Payments:  []domain.Payment{},
		}

		mockOrderRepo.On("GetOrderWithDetails", ctx, "order-1").Return(order, nil).Once()

		result, err := service.GetPaymentByOrder(ctx, "order-1")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
	})
}
