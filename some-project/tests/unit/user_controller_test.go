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

type MockuserUseCase struct {
	mock.Mock
}

func (m *MockuserUseCase) Create(ctx context.Context, entity *domain.user) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockuserUseCase) Get(ctx context.Context, id string) (*domain.user, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.user), args.Error(1)
}

func (m *MockuserUseCase) Update(ctx context.Context, entity *domain.user) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockuserUseCase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockuserUseCase) List(ctx context.Context) ([]*domain.user, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.user), args.Error(1)
}

func setupuserTest() (*gin.Engine, *MockuserUseCase, *controller.userController) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	mockUseCase := &MockuserUseCase{}
	controller := controller.NewuserController(mockUseCase)
	
	api := router.Group("/api/v1")
	{
		api.POST("/users", controller.Create)
		api.GET("/users/:id", controller.Get)
		api.PUT("/users/:id", controller.Update)
		api.DELETE("/users/:id", controller.Delete)
		api.GET("/users", controller.List)
	}
	
	return router, mockUseCase, controller
}

func TestuserController_Create(t *testing.T) {
	router, mockUseCase, _ := setupuserTest()
	
	entity := &domain.user{
		
		ID: 123,
		
		Username: "test-value",
		
		Email: "test-value",
		
		CreatedAt: "test-value",
		
	}
	
	mockUseCase.On("Create", mock.Anything, mock.AnythingOfType("*domain.user"))
		.Return(nil)
	
	body, _ := json.Marshal(entity)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestuserController_Get(t *testing.T) {
	router, mockUseCase, _ := setupuserTest()
	
	entity := &domain.user{
		ID: "test-id",
		
		ID: 123,
		
		Username: "test-value",
		
		Email: "test-value",
		
		CreatedAt: "test-value",
		
	}
	
	mockUseCase.On("Get", mock.Anything, "test-id").Return(entity, nil)
	
	req, _ := http.NewRequest("GET", "/api/v1/users/test-id", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestuserController_Update(t *testing.T) {
	router, mockUseCase, _ := setupuserTest()
	
	entity := &domain.user{
		ID: "test-id",
		
		ID: 123,
		
		Username: "test-value",
		
		Email: "test-value",
		
		CreatedAt: "test-value",
		
	}
	
	mockUseCase.On("Update", mock.Anything, mock.AnythingOfType("*domain.user"))
		.Return(nil)
	
	body, _ := json.Marshal(entity)
	req, _ := http.NewRequest("PUT", "/api/v1/users/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestuserController_Delete(t *testing.T) {
	router, mockUseCase, _ := setupuserTest()
	
	mockUseCase.On("Delete", mock.Anything, "test-id").Return(nil)
	
	req, _ := http.NewRequest("DELETE", "/api/v1/users/test-id", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestuserController_List(t *testing.T) {
	router, mockUseCase, _ := setupuserTest()
	
	entities := []*domain.user{
		{
			ID: "test-id-1",
			
			ID: 123,
			
			Username: "test-value",
			
			Email: "test-value",
			
			CreatedAt: "test-value",
			
		},
		{
			ID: "test-id-2",
			
			ID: 123,
			
			Username: "test-value",
			
			Email: "test-value",
			
			CreatedAt: "test-value",
			
		},
	}
	
	mockUseCase.On("List", mock.Anything).Return(entities, nil)
	
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}
