package handler

import (
	"backend/internal/domain"
	"backend/internal/pagination"
	"backend/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Cruise Service
type MockCruiseService struct {
	mock.Mock
}

func (m *MockCruiseService) Create(ctx context.Context, req service.CreateCruiseRequest) (*domain.Cruise, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cruise), args.Error(1)
}

func (m *MockCruiseService) GetByID(ctx context.Context, id string) (*domain.Cruise, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cruise), args.Error(1)
}

func (m *MockCruiseService) GetByCode(ctx context.Context, code string) (*domain.Cruise, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cruise), args.Error(1)
}

func (m *MockCruiseService) List(ctx context.Context, req service.ListCruisesRequest) (*pagination.Result, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pagination.Result), args.Error(1)
}

func (m *MockCruiseService) Update(ctx context.Context, id string, req service.UpdateCruiseRequest) (*domain.Cruise, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cruise), args.Error(1)
}

func (m *MockCruiseService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCruiseService) Restore(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCruiseService) UpdateStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockCruiseService) UpdateSortWeight(ctx context.Context, id string, sortWeight int) error {
	args := m.Called(ctx, id, sortWeight)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestCruiseHandler_GetByID(t *testing.T) {
	router := setupTestRouter()
	mockService := new(MockCruiseService)
	handler := NewCruiseHandler(mockService)

	router.GET("/cruises/:id", handler.GetByID)

	t.Run("should get cruise by ID successfully", func(t *testing.T) {
		expectedCruise := &domain.Cruise{
			BaseModel: domain.BaseModel{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			NameCN:            "测试邮轮",
			Code:              "TEST001",
			GrossTonnage:      100000,
			PassengerCapacity: 5000,
			Status:            domain.CruiseStatusActive,
		}

		mockService.On("GetByID", mock.Anything, "cruise-1").Return(expectedCruise, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/cruises/cruise-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
	})

	t.Run("should return 404 when cruise not found", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, "non-existent").Return(nil, service.ErrCruiseNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/cruises/non-existent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCruiseHandler_List(t *testing.T) {
	router := setupTestRouter()
	mockService := new(MockCruiseService)
	handler := NewCruiseHandler(mockService)

	router.GET("/cruises", handler.List)

	t.Run("should list cruises successfully", func(t *testing.T) {
		expectedResult := &pagination.Result{
			Data: []*domain.Cruise{
				{
					BaseModel: domain.BaseModel{ID: uuid.New()},
					NameCN:    "邮轮1",
				},
				{
					BaseModel: domain.BaseModel{ID: uuid.New()},
					NameCN:    "邮轮2",
				},
			},
			Pagination: pagination.Paginator{
				Page:     1,
				PageSize: 10,
				Total:    2,
				Pages:    1,
			},
		}

		mockService.On("List", mock.Anything, mock.AnythingOfType("service.ListCruisesRequest")).Return(expectedResult, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/cruises?page=1&page_size=10", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
	})
}

func TestCruiseHandler_Create(t *testing.T) {
	router := setupTestRouter()
	mockService := new(MockCruiseService)
	handler := NewCruiseHandler(mockService)

	router.POST("/cruises", handler.Create)

	t.Run("should create cruise successfully", func(t *testing.T) {
		reqBody := service.CreateCruiseRequest{
			CompanyID:    "company-1",
			NameCN:       "新邮轮",
			Code:         "NEW001",
			GrossTonnage: 100000,
		}

		expectedCruise := &domain.Cruise{
			BaseModel: domain.BaseModel{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			CompanyID:    reqBody.CompanyID,
			NameCN:       reqBody.NameCN,
			Code:         reqBody.Code,
			GrossTonnage: reqBody.GrossTonnage,
			Status:       domain.CruiseStatusActive,
		}

		mockService.On("Create", mock.Anything, reqBody).Return(expectedCruise, nil).Once()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/cruises", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["data"])
	})

	t.Run("should return 400 for invalid request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/cruises", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCruiseHandler_Update(t *testing.T) {
	router := setupTestRouter()
	mockService := new(MockCruiseService)
	handler := NewCruiseHandler(mockService)

	router.PUT("/cruises/:id", handler.Update)

	t.Run("should update cruise successfully", func(t *testing.T) {
		reqBody := service.UpdateCruiseRequest{
			NameCN: "更新的邮轮名称",
		}

		expectedCruise := &domain.Cruise{
			BaseModel: domain.BaseModel{
				ID:        uuid.New(),
				UpdatedAt: time.Now(),
			},
			NameCN: "更新的邮轮名称",
			Code:   "TEST001",
		}

		mockService.On("Update", mock.Anything, "cruise-1", reqBody).Return(expectedCruise, nil).Once()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/cruises/cruise-1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
	})

	t.Run("should return 404 when cruise not found", func(t *testing.T) {
		reqBody := service.UpdateCruiseRequest{NameCN: "新名称"}
		mockService.On("Update", mock.Anything, "non-existent", reqBody).Return(nil, service.ErrCruiseNotFound).Once()

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/cruises/non-existent", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCruiseHandler_Delete(t *testing.T) {
	router := setupTestRouter()
	mockService := new(MockCruiseService)
	handler := NewCruiseHandler(mockService)

	router.DELETE("/cruises/:id", handler.Delete)

	t.Run("should delete cruise successfully", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, "cruise-1").Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/cruises/cruise-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return 404 when cruise not found", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, "non-existent").Return(service.ErrCruiseNotFound).Once()

		req := httptest.NewRequest(http.MethodDelete, "/cruises/non-existent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
