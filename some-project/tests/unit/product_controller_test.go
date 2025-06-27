package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/domain"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/controller"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/usecase"
)

type MockproductUseCase struct {
	mock.Mock
}

func (m *MockproductUseCase) Create(ctx context.Context, entity *domain.product) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockproductUseCase) Get(ctx context.Context, id string) (*domain.product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.product), args.Error(1)
}

func (m *MockproductUseCase) Update(ctx context.Context, entity *domain.product) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockproductUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockproductUseCase) List(ctx context.Context) ([]*domain.product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.product), args.Error(1)
}

func setupproductTest() (*gin.Engine, *MockproductUseCase, *controller.productController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	mockUseCase := &MockproductUseCase{}
	controller := controller.NewproductController(mockUseCase)
	
	api := router.Group("/api/v1")
	{
		api.POST("/products", controller.Create)
		api.GET("/products/:id", controller.Get)
		api.PUT("/products/:id", controller.Update)
		api.DELETE("/products/:id", controller.Delete)
		api.GET("/products", controller.List)
	}
	
	return router, mockUseCase, controller
}

func TestproductController_Create(t *testing.T) {
	router, mockUseCase, _ := setupproductTest()
	
	entity := &domain.product{
		
		ID: 123,
		
		Name: "test-value",
		
		Description: "test-value",
		
		Price: 123.45,
		
		CreatedAt: "test-value",
		
	}
	
	mockUseCase.On("Create", mock.Anything, mock.AnythingOfType("*domain.product"))
		.Return(nil)
	
	body, _ := json.Marshal(entity)
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestproductController_Get(t *testing.T) {
	router, mockUseCase, _ := setupproductTest()
	
	entity := &domain.product{
		ID: "test-id",
		
		ID: 123,
		
		Name: "test-value",
		
		Description: "test-value",
		
		Price: 123.45,
		
		CreatedAt: "test-value",
		
	}
	
	mockUseCase.On("Get", mock.Anything, "test-id").Return(entity, nil)
	
	req, _ := http.NewRequest("GET", "/api/v1/products/test-id", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestproductController_Update(t *testing.T) {
	router, mockUseCase, _ := setupproductTest()
	
	entity := &domain.product{
		ID: "test-id",
		
		ID: 123,
		
		Name: "test-value",
		
		Description: "test-value",
		
		Price: 123.45,
		
		CreatedAt: "test-value",
		
	}
	
	mockUseCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.product"))
		.Return(nil)
	
	body, _ := json.Marshal(entity)
	req, _ := http.NewRequest("PUT", "/api/v1/products/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestproductController_Delete(t *testing.T) {
	router, mockUseCase, _ := setupproductTest()
	
	mockUseCase.On("Delete", mock.Anything, "test-id").Return(nil)
	
	req, _ := http.NewRequest("DELETE", "/api/v1/products/test-id", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestproductController_List(t *testing.T) {
	router, mockUseCase, _ := setupproductTest()
	
	entities := []*domain.product{
		{
			ID: "test-id-1",
			
			ID: 123,
			
			Name: "test-value",
			
			Description: "test-value",
			
			Price: 123.45,
			
			CreatedAt: "test-value",
			
		},
		{
			ID: "test-id-2",
			
			ID: 123,
			
			Name: "test-value",
			
			Description: "test-value",
			
			Price: 123.45,
			
			CreatedAt: "test-value",
			
		},
	}
	
	mockUseCase.On("List", mock.Anything).Return(entities, nil)
	
	req, _ := http.NewRequest("GET", "/api/v1/products", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}
