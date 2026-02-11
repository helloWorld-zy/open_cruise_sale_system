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

// Mock Cruise Repository
type MockCruiseRepository struct {
	mock.Mock
}

func (m *MockCruiseRepository) Create(ctx context.Context, cruise *domain.Cruise) error {
	args := m.Called(ctx, cruise)
	return args.Error(0)
}

func (m *MockCruiseRepository) GetByID(ctx context.Context, id string) (*domain.Cruise, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cruise), args.Error(1)
}

func (m *MockCruiseRepository) GetByCode(ctx context.Context, code string) (*domain.Cruise, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cruise), args.Error(1)
}

func (m *MockCruiseRepository) List(ctx context.Context, filters repository.CruiseFilters, paginator *pagination.Paginator) ([]*domain.Cruise, error) {
	args := m.Called(ctx, filters, paginator)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Cruise), args.Error(1)
}

func (m *MockCruiseRepository) Count(ctx context.Context, filters repository.CruiseFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCruiseRepository) Update(ctx context.Context, cruise *domain.Cruise) error {
	args := m.Called(ctx, cruise)
	return args.Error(0)
}

func (m *MockCruiseRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCruiseRepository) Restore(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCruiseRepository) ListByCompany(ctx context.Context, companyID string) ([]*domain.Cruise, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Cruise), args.Error(1)
}

func (m *MockCruiseRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockCruiseRepository) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	args := m.Called(ctx, id, sortWeight)
	return args.Error(0)
}

func TestCruiseService_Create(t *testing.T) {
	mockRepo := new(MockCruiseRepository)
	service := NewCruiseService(mockRepo)
	ctx := context.Background()

	t.Run("should create cruise successfully", func(t *testing.T) {
		req := CreateCruiseRequest{
			CompanyID:         "company-1",
			NameCN:            "测试邮轮",
			NameEN:            "Test Cruise",
			Code:              "TEST001",
			GrossTonnage:      100000,
			PassengerCapacity: 5000,
			Status:            domain.CruiseStatusActive,
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Cruise")).Return(nil).Once()

		cruise, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, cruise)
		assert.Equal(t, req.NameCN, cruise.NameCN)
		assert.Equal(t, req.Code, cruise.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should create cruise with default status", func(t *testing.T) {
		req := CreateCruiseRequest{
			CompanyID: "company-1",
			NameCN:    "测试邮轮",
			Code:      "TEST002",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Cruise")).Return(nil).Once()

		cruise, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, domain.CruiseStatusActive, cruise.Status)
	})

	t.Run("should return error for invalid status", func(t *testing.T) {
		req := CreateCruiseRequest{
			CompanyID: "company-1",
			NameCN:    "测试邮轮",
			Code:      "TEST003",
			Status:    "invalid_status",
		}

		cruise, err := service.Create(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCruiseData, err)
		assert.Nil(t, cruise)
	})
}

func TestCruiseService_GetByID(t *testing.T) {
	mockRepo := new(MockCruiseRepository)
	service := NewCruiseService(mockRepo)
	ctx := context.Background()

	t.Run("should get cruise by ID successfully", func(t *testing.T) {
		expectedCruise := &domain.Cruise{
			BaseModel: domain.BaseModel{
				ID:        "cruise-1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			NameCN: "测试邮轮",
			Code:   "TEST001",
		}

		mockRepo.On("GetByID", ctx, "cruise-1").Return(expectedCruise, nil).Once()

		cruise, err := service.GetByID(ctx, "cruise-1")

		assert.NoError(t, err)
		assert.NotNil(t, cruise)
		assert.Equal(t, expectedCruise.ID, cruise.ID)
	})

	t.Run("should return error when cruise not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "non-existent").Return(nil, gorm.ErrRecordNotFound).Once()

		cruise, err := service.GetByID(ctx, "non-existent")

		assert.Error(t, err)
		assert.Equal(t, ErrCruiseNotFound, err)
		assert.Nil(t, cruise)
	})
}

func TestCruiseService_List(t *testing.T) {
	mockRepo := new(MockCruiseRepository)
	service := NewCruiseService(mockRepo)
	ctx := context.Background()

	t.Run("should list cruises with pagination", func(t *testing.T) {
		req := ListCruisesRequest{
			Keyword: "test",
			Paginator: pagination.Paginator{
				Page:     1,
				PageSize: 10,
			},
		}

		expectedCruises := []*domain.Cruise{
			{
				BaseModel: domain.BaseModel{ID: "1"},
				NameCN:    "邮轮1",
			},
			{
				BaseModel: domain.BaseModel{ID: "2"},
				NameCN:    "邮轮2",
			},
		}

		mockRepo.On("Count", ctx, mock.AnythingOfType("repository.CruiseFilters")).Return(int64(2), nil).Once()
		mockRepo.On("List", ctx, mock.AnythingOfType("repository.CruiseFilters"), mock.AnythingOfType("*pagination.Paginator")).Return(expectedCruises, nil).Once()

		result, err := service.List(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(result.Data.([]*domain.Cruise)))
		assert.Equal(t, int64(2), result.Pagination.Total)
	})
}

func TestCruiseService_Update(t *testing.T) {
	mockRepo := new(MockCruiseRepository)
	service := NewCruiseService(mockRepo)
	ctx := context.Background()

	t.Run("should update cruise successfully", func(t *testing.T) {
		existingCruise := &domain.Cruise{
			BaseModel: domain.BaseModel{
				ID:        "cruise-1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			NameCN: "旧名称",
			Code:   "TEST001",
		}

		req := UpdateCruiseRequest{
			NameCN: "新名称",
		}

		mockRepo.On("GetByID", ctx, "cruise-1").Return(existingCruise, nil).Once()
		mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Cruise")).Return(nil).Once()

		cruise, err := service.Update(ctx, "cruise-1", req)

		assert.NoError(t, err)
		assert.NotNil(t, cruise)
		assert.Equal(t, "新名称", cruise.NameCN)
	})

	t.Run("should return error for non-existent cruise", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "non-existent").Return(nil, gorm.ErrRecordNotFound).Once()

		req := UpdateCruiseRequest{NameCN: "新名称"}
		cruise, err := service.Update(ctx, "non-existent", req)

		assert.Error(t, err)
		assert.Equal(t, ErrCruiseNotFound, err)
		assert.Nil(t, cruise)
	})
}

func TestCruiseService_Delete(t *testing.T) {
	mockRepo := new(MockCruiseRepository)
	service := NewCruiseService(mockRepo)
	ctx := context.Background()

	t.Run("should delete cruise successfully", func(t *testing.T) {
		existingCruise := &domain.Cruise{
			BaseModel: domain.BaseModel{ID: "cruise-1"},
			NameCN:    "测试邮轮",
		}

		mockRepo.On("GetByID", ctx, "cruise-1").Return(existingCruise, nil).Once()
		mockRepo.On("Delete", ctx, "cruise-1").Return(nil).Once()

		err := service.Delete(ctx, "cruise-1")

		assert.NoError(t, err)
	})

	t.Run("should return error for non-existent cruise", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "non-existent").Return(nil, gorm.ErrRecordNotFound).Once()

		err := service.Delete(ctx, "non-existent")

		assert.Error(t, err)
		assert.Equal(t, ErrCruiseNotFound, err)
	})
}

func TestCruiseService_UpdateStatus(t *testing.T) {
	mockRepo := new(MockCruiseRepository)
	service := NewCruiseService(mockRepo)
	ctx := context.Background()

	t.Run("should update status successfully", func(t *testing.T) {
		existingCruise := &domain.Cruise{
			BaseModel: domain.BaseModel{ID: "cruise-1"},
			Status:    domain.CruiseStatusActive,
		}

		mockRepo.On("GetByID", ctx, "cruise-1").Return(existingCruise, nil).Once()
		mockRepo.On("UpdateStatus", ctx, "cruise-1", domain.CruiseStatusInactive).Return(nil).Once()

		err := service.UpdateStatus(ctx, "cruise-1", domain.CruiseStatusInactive)

		assert.NoError(t, err)
	})

	t.Run("should return error for invalid status", func(t *testing.T) {
		err := service.UpdateStatus(ctx, "cruise-1", "invalid_status")

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCruiseData, err)
	})
}
