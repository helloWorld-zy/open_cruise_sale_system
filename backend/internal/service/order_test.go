package service

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/repository"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock Order Repository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	args := m.Called(ctx, orderNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) List(ctx context.Context, filters repository.OrderFilters, paginator *pagination.Paginator) ([]*domain.Order, error) {
	args := m.Called(ctx, filters, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) Count(ctx context.Context, filters repository.OrderFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockOrderRepository) ListByUser(ctx context.Context, userID string, paginator *pagination.Paginator) ([]*domain.Order, error) {
	args := m.Called(ctx, userID, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) ListByStatus(ctx context.Context, status string, paginator *pagination.Paginator) ([]*domain.Order, error) {
	args := m.Called(ctx, status, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) ListByVoyage(ctx context.Context, voyageID string) ([]*domain.Order, error) {
	args := m.Called(ctx, voyageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdatePaymentStatus(ctx context.Context, id string, paymentStatus string, paidAmount float64) error {
	args := m.Called(ctx, id, paymentStatus, paidAmount)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) CreateOrderItem(ctx context.Context, item *domain.OrderItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderItemByID(ctx context.Context, id string) (*domain.OrderItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OrderItem), args.Error(1)
}

func (m *MockOrderRepository) ListOrderItemsByOrder(ctx context.Context, orderID string) ([]*domain.OrderItem, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.OrderItem), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrderItem(ctx context.Context, item *domain.OrderItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateOrderItemStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) DeleteOrderItem(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) CreatePassenger(ctx context.Context, passenger *domain.Passenger) error {
	args := m.Called(ctx, passenger)
	return args.Error(0)
}

func (m *MockOrderRepository) GetPassengerByID(ctx context.Context, id string) (*domain.Passenger, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Passenger), args.Error(1)
}

func (m *MockOrderRepository) ListPassengersByOrder(ctx context.Context, orderID string) ([]*domain.Passenger, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Passenger), args.Error(1)
}

func (m *MockOrderRepository) ListPassengersByOrderItem(ctx context.Context, orderItemID string) ([]*domain.Passenger, error) {
	args := m.Called(ctx, orderItemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Passenger), args.Error(1)
}

func (m *MockOrderRepository) UpdatePassenger(ctx context.Context, passenger *domain.Passenger) error {
	args := m.Called(ctx, passenger)
	return args.Error(0)
}

func (m *MockOrderRepository) DeletePassenger(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) BatchCreatePassengers(ctx context.Context, passengers []*domain.Passenger) error {
	args := m.Called(ctx, passengers)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderWithDetails(ctx context.Context, id string) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockOrderRepository) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockOrderRepository) GetPaymentByNo(ctx context.Context, paymentNo string) (*domain.Payment, error) {
	args := m.Called(ctx, paymentNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockOrderRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

// Mock Voyage Repository
type MockVoyageRepository struct {
	mock.Mock
}

func (m *MockVoyageRepository) Create(ctx context.Context, voyage *domain.Voyage) error {
	args := m.Called(ctx, voyage)
	return args.Error(0)
}

func (m *MockVoyageRepository) GetByID(ctx context.Context, id string) (*domain.Voyage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Voyage), args.Error(1)
}

func (m *MockVoyageRepository) GetByVoyageNumber(ctx context.Context, voyageNumber string) (*domain.Voyage, error) {
	args := m.Called(ctx, voyageNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Voyage), args.Error(1)
}

func (m *MockVoyageRepository) List(ctx context.Context, filters repository.VoyageFilters, paginator *pagination.Paginator) ([]*domain.Voyage, error) {
	args := m.Called(ctx, filters, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Voyage), args.Error(1)
}

func (m *MockVoyageRepository) Count(ctx context.Context, filters repository.VoyageFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockVoyageRepository) ListByCruise(ctx context.Context, cruiseID string) ([]*domain.Voyage, error) {
	args := m.Called(ctx, cruiseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Voyage), args.Error(1)
}

func (m *MockVoyageRepository) ListByRoute(ctx context.Context, routeID string) ([]*domain.Voyage, error) {
	args := m.Called(ctx, routeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Voyage), args.Error(1)
}

func (m *MockVoyageRepository) Update(ctx context.Context, voyage *domain.Voyage) error {
	args := m.Called(ctx, voyage)
	return args.Error(0)
}

func (m *MockVoyageRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockVoyageRepository) UpdateBookingStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockVoyageRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock Cabin Repository
type MockCabinRepository struct {
	mock.Mock
}

func (m *MockCabinRepository) Create(ctx context.Context, cabin *domain.Cabin) error {
	args := m.Called(ctx, cabin)
	return args.Error(0)
}

func (m *MockCabinRepository) GetByID(ctx context.Context, id string) (*domain.Cabin, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cabin), args.Error(1)
}

func (m *MockCabinRepository) GetByCabinNumber(ctx context.Context, voyageID, cabinNumber string) (*domain.Cabin, error) {
	args := m.Called(ctx, voyageID, cabinNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cabin), args.Error(1)
}

func (m *MockCabinRepository) List(ctx context.Context, filters repository.CabinFilters, paginator *pagination.Paginator) ([]*domain.Cabin, error) {
	args := m.Called(ctx, filters, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Cabin), args.Error(1)
}

func (m *MockCabinRepository) Count(ctx context.Context, filters repository.CabinFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCabinRepository) ListByVoyage(ctx context.Context, voyageID string) ([]*domain.Cabin, error) {
	args := m.Called(ctx, voyageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Cabin), args.Error(1)
}

func (m *MockCabinRepository) ListByVoyageAndType(ctx context.Context, voyageID, cabinTypeID string) ([]*domain.Cabin, error) {
	args := m.Called(ctx, voyageID, cabinTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Cabin), args.Error(1)
}

func (m *MockCabinRepository) ListAvailableByVoyage(ctx context.Context, voyageID string) ([]*domain.Cabin, error) {
	args := m.Called(ctx, voyageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Cabin), args.Error(1)
}

func (m *MockCabinRepository) Update(ctx context.Context, cabin *domain.Cabin) error {
	args := m.Called(ctx, cabin)
	return args.Error(0)
}

func (m *MockCabinRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockCabinRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCabinRepository) BatchCreate(ctx context.Context, cabins []*domain.Cabin) error {
	args := m.Called(ctx, cabins)
	return args.Error(0)
}

// Mock Inventory Repository
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) GetInventory(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinInventory, error) {
	args := m.Called(ctx, voyageID, cabinTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CabinInventory), args.Error(1)
}

func (m *MockInventoryRepository) CreateInventory(ctx context.Context, inventory *domain.CabinInventory) error {
	args := m.Called(ctx, inventory)
	return args.Error(0)
}

func (m *MockInventoryRepository) LockCabin(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	args := m.Called(ctx, voyageID, cabinTypeID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) UnlockCabin(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	args := m.Called(ctx, voyageID, cabinTypeID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) ConfirmBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	args := m.Called(ctx, voyageID, cabinTypeID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) CancelBooking(ctx context.Context, voyageID, cabinTypeID string, quantity int) error {
	args := m.Called(ctx, voyageID, cabinTypeID, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) UpdateInventory(ctx context.Context, inventory *domain.CabinInventory) error {
	args := m.Called(ctx, inventory)
	return args.Error(0)
}

func (m *MockInventoryRepository) ListInventoryByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinInventory, error) {
	args := m.Called(ctx, voyageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CabinInventory), args.Error(1)
}

// Mock Price Repository
type MockPriceRepository struct {
	mock.Mock
}

func (m *MockPriceRepository) Create(ctx context.Context, price *domain.CabinPrice) error {
	args := m.Called(ctx, price)
	return args.Error(0)
}

func (m *MockPriceRepository) GetByID(ctx context.Context, id string) (*domain.CabinPrice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CabinPrice), args.Error(1)
}

func (m *MockPriceRepository) GetByVoyageAndCabinType(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error) {
	args := m.Called(ctx, voyageID, cabinTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CabinPrice), args.Error(1)
}

func (m *MockPriceRepository) List(ctx context.Context, filters repository.PriceFilters, paginator *pagination.Paginator) ([]*domain.CabinPrice, error) {
	args := m.Called(ctx, filters, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CabinPrice), args.Error(1)
}

func (m *MockPriceRepository) Count(ctx context.Context, filters repository.PriceFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPriceRepository) ListByVoyage(ctx context.Context, voyageID string) ([]*domain.CabinPrice, error) {
	args := m.Called(ctx, voyageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CabinPrice), args.Error(1)
}

func (m *MockPriceRepository) ListByVoyageWithTypes(ctx context.Context, voyageID string, cabinTypeIDs []string) ([]*domain.CabinPrice, error) {
	args := m.Called(ctx, voyageID, cabinTypeIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CabinPrice), args.Error(1)
}

func (m *MockPriceRepository) GetCurrentPrice(ctx context.Context, voyageID, cabinTypeID string) (*domain.CabinPrice, error) {
	args := m.Called(ctx, voyageID, cabinTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CabinPrice), args.Error(1)
}

func (m *MockPriceRepository) Update(ctx context.Context, price *domain.CabinPrice) error {
	args := m.Called(ctx, price)
	return args.Error(0)
}

func (m *MockPriceRepository) UpdatePrice(ctx context.Context, id string, adultPrice, childPrice, infantPrice float64) error {
	args := m.Called(ctx, id, adultPrice, childPrice, infantPrice)
	return args.Error(0)
}

func (m *MockPriceRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPriceRepository) BatchCreate(ctx context.Context, prices []*domain.CabinPrice) error {
	args := m.Called(ctx, prices)
	return args.Error(0)
}

// Mock Order State Service
type MockOrderStateService struct {
	mock.Mock
}

func (m *MockOrderStateService) CanPay(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderStateService) CanConfirm(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderStateService) CanComplete(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderStateService) CanCancel(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderStateService) CanRefund(order *domain.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderStateService) TransitionToPaid(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderStateService) TransitionToConfirmed(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderStateService) TransitionToCompleted(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderStateService) TransitionToCancelled(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderStateService) TransitionToRefunded(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func TestOrderService_GetByID(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockVoyageRepo := new(MockVoyageRepository)
	mockCabinRepo := new(MockCabinRepository)
	mockPriceRepo := new(MockPriceRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	service := NewOrderService(mockOrderRepo, mockVoyageRepo, mockCabinRepo, mockPriceRepo, mockInventoryRepo)
	ctx := context.Background()

	t.Run("should return order by ID", func(t *testing.T) {
		order := &domain.Order{
			ID:          "order-1",
			OrderNumber: "ORD202401010001",
			Status:      domain.OrderStatusPending,
		}

		mockOrderRepo.On("GetByID", ctx, "order-1").Return(order, nil).Once()

		result, err := service.GetByID(ctx, "order-1")

		assert.NoError(t, err)
		assert.Equal(t, order, result)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		mockOrderRepo.On("GetByID", ctx, "non-existent").Return(nil, gorm.ErrRecordNotFound).Once()

		result, err := service.GetByID(ctx, "non-existent")

		assert.ErrorIs(t, err, ErrOrderNotFound)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_Cancel(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockVoyageRepo := new(MockVoyageRepository)
	mockCabinRepo := new(MockCabinRepository)
	mockPriceRepo := new(MockPriceRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	service := NewOrderService(mockOrderRepo, mockVoyageRepo, mockCabinRepo, mockPriceRepo, mockInventoryRepo)
	ctx := context.Background()

	t.Run("should cancel pending order successfully", func(t *testing.T) {
		order := &domain.Order{
			ID:     "order-1",
			Status: domain.OrderStatusPending,
		}

		mockOrderRepo.On("GetByID", ctx, "order-1").Return(order, nil).Once()
		mockOrderRepo.On("UpdateStatus", ctx, "order-1", domain.OrderStatusCancelled).Return(nil).Once()

		err := service.Cancel(ctx, "order-1")

		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		mockOrderRepo.On("GetByID", ctx, "non-existent").Return(nil, gorm.ErrRecordNotFound).Once()

		err := service.Cancel(ctx, "non-existent")

		assert.ErrorIs(t, err, ErrOrderNotFound)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_List(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockVoyageRepo := new(MockVoyageRepository)
	mockCabinRepo := new(MockCabinRepository)
	mockPriceRepo := new(MockPriceRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	service := NewOrderService(mockOrderRepo, mockVoyageRepo, mockCabinRepo, mockPriceRepo, mockInventoryRepo)
	ctx := context.Background()

	t.Run("should return paginated orders", func(t *testing.T) {
		req := ListOrdersRequest{
			Status: "pending",
		}
		req.Paginator = pagination.Paginator{Page: 1, PageSize: 10}

		orders := []*domain.Order{
			{ID: "order-1", Status: domain.OrderStatusPending},
			{ID: "order-2", Status: domain.OrderStatusPending},
		}

		mockOrderRepo.On("Count", ctx, mock.AnythingOfType("repository.OrderFilters")).Return(int64(2), nil).Once()
		mockOrderRepo.On("List", ctx, mock.AnythingOfType("repository.OrderFilters"), mock.AnythingOfType("*pagination.Paginator")).Return(orders, nil).Once()

		result, err := service.List(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), result.Pagination.Total)
		assert.Len(t, result.Data, 2)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_Update(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockVoyageRepo := new(MockVoyageRepository)
	mockCabinRepo := new(MockCabinRepository)
	mockPriceRepo := new(MockPriceRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	service := NewOrderService(mockOrderRepo, mockVoyageRepo, mockCabinRepo, mockPriceRepo, mockInventoryRepo)
	ctx := context.Background()

	t.Run("should update pending order", func(t *testing.T) {
		order := &domain.Order{
			ID:           "order-1",
			Status:       domain.OrderStatusPending,
			ContactName:  "Old Name",
			ContactPhone: "1234567890",
		}

		req := UpdateOrderRequest{
			ContactName:  "New Name",
			ContactPhone: "0987654321",
		}

		mockOrderRepo.On("GetByID", ctx, "order-1").Return(order, nil).Once()
		mockOrderRepo.On("Update", ctx, mock.AnythingOfType("*domain.Order")).Return(nil).Once()

		result, err := service.Update(ctx, "order-1", req)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", result.ContactName)
		assert.Equal(t, "0987654321", result.ContactPhone)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("should return error for non-pending order", func(t *testing.T) {
		order := &domain.Order{
			ID:     "order-1",
			Status: domain.OrderStatusPaid,
		}

		req := UpdateOrderRequest{
			ContactName: "New Name",
		}

		mockOrderRepo.On("GetByID", ctx, "order-1").Return(order, nil).Once()

		result, err := service.Update(ctx, "order-1", req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_CalculateTotal(t *testing.T) {
	mockOrderRepo := new(MockOrderRepository)
	mockVoyageRepo := new(MockVoyageRepository)
	mockCabinRepo := new(MockCabinRepository)
	mockPriceRepo := new(MockPriceRepository)
	mockInventoryRepo := new(MockInventoryRepository)

	service := NewOrderService(mockOrderRepo, mockVoyageRepo, mockCabinRepo, mockPriceRepo, mockInventoryRepo)
	ctx := context.Background()

	t.Run("should calculate total for items", func(t *testing.T) {
		items := []OrderItemRequest{
			{
				CabinTypeID: "cabin-type-1",
				AdultCount:  2,
				ChildCount:  1,
			},
		}

		price := &domain.CabinPrice{
			AdultPrice:  1000,
			ChildPrice:  500,
			InfantPrice: 0,
			PortFee:     100,
			ServiceFee:  50,
		}

		mockPriceRepo.On("GetCurrentPrice", ctx, "", "cabin-type-1").Return(price, nil).Once()

		result, err := service.CalculateTotal(ctx, items)

		assert.NoError(t, err)
		assert.Equal(t, 2500.0, result.Subtotal)  // 2*1000 + 1*500
		assert.Equal(t, 300.0, result.PortFee)    // 3*100
		assert.Equal(t, 150.0, result.ServiceFee) // 3*50
		assert.Equal(t, 2950.0, result.TotalAmount)
		mockPriceRepo.AssertExpectations(t)
	})
}

func TestOrderStateMachine(t *testing.T) {
	t.Run("should allow valid transitions", func(t *testing.T) {
		sm := NewOrderStateMachine(domain.OrderStatusPending)

		assert.True(t, sm.CanTransition(domain.OrderStatusPaid))
		assert.True(t, sm.CanTransition(domain.OrderStatusCancelled))
		assert.False(t, sm.CanTransition(domain.OrderStatusConfirmed))
	})

	t.Run("should not allow invalid transitions", func(t *testing.T) {
		sm := NewOrderStateMachine(domain.OrderStatusCompleted)

		assert.False(t, sm.CanTransition(domain.OrderStatusPending))
		assert.False(t, sm.CanTransition(domain.OrderStatusPaid))
	})

	t.Run("should perform transition", func(t *testing.T) {
		sm := NewOrderStateMachine(domain.OrderStatusPending)

		err := sm.Transition(domain.OrderStatusPaid)

		assert.NoError(t, err)
		assert.Equal(t, domain.OrderStatusPaid, sm.CurrentState())
	})

	t.Run("should return error for invalid transition", func(t *testing.T) {
		sm := NewOrderStateMachine(domain.OrderStatusPending)

		err := sm.Transition(domain.OrderStatusConfirmed)

		assert.Error(t, err)
		assert.Equal(t, domain.OrderStatusPending, sm.CurrentState())
	})
}
