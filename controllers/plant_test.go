package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"plant-reminder/dto"
	"plant-reminder/models"
	"testing"

	"github.com/gin-gonic/gin"
)

// MockPlantService is a mock implementation of PlantService for testing
type MockPlantService struct {
	CreatePlantFunc func(*dto.PlantCreateRequest, int64) (*dto.PlantResponse, error)
	GetPlantFunc    func(int64, int64) (*dto.PlantResponse, error)
	GetPlantsFunc   func(int64) ([]dto.PlantResponse, error)
	UpdatePlantFunc func(*dto.PlantUpdateRequest, int64, int64) error
	DeletePlantFunc func(int64, int64) error
}

func (m *MockPlantService) CreatePlant(req *dto.PlantCreateRequest, userID int64) (*dto.PlantResponse, error) {
	if m.CreatePlantFunc != nil {
		return m.CreatePlantFunc(req, userID)
	}
	return nil, nil
}

func (m *MockPlantService) GetPlant(plantID, userID int64) (*dto.PlantResponse, error) {
	if m.GetPlantFunc != nil {
		return m.GetPlantFunc(plantID, userID)
	}
	return nil, nil
}

func (m *MockPlantService) GetPlants(userID int64) ([]dto.PlantResponse, error) {
	if m.GetPlantsFunc != nil {
		return m.GetPlantsFunc(userID)
	}
	return nil, nil
}

func (m *MockPlantService) UpdatePlant(req *dto.PlantUpdateRequest, plantID, userID int64) error {
	if m.UpdatePlantFunc != nil {
		return m.UpdatePlantFunc(req, plantID, userID)
	}
	return nil
}

func (m *MockPlantService) DeletePlant(userID, plantID int64) error {
	if m.DeletePlantFunc != nil {
		return m.DeletePlantFunc(userID, plantID)
	}
	return nil
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func setupPlantController(mockService *MockPlantService) (*PlantController, *gin.Engine) {
	router := setupTestRouter()
	controller := &PlantController{plantService: mockService}
	return controller, router
}

func TestPlantController_AddPlant_Success(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	// Setup mock response
	expectedResponse := &dto.PlantResponse{
		ID:        1,
		Name:      "Test Plant",
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}

	mockService.CreatePlantFunc = func(req *dto.PlantCreateRequest, userID int64) (*dto.PlantResponse, error) {
		if req.Name != "Test Plant" {
			t.Errorf("Expected name 'Test Plant', got %s", req.Name)
		}
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return expectedResponse, nil
	}

	// Setup route
	router.POST("/plants", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.AddPlant(c)
	})

	// Create request body
	requestBody := dto.PlantCreateRequest{
		Name:      "Test Plant",
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}
	jsonData, _ := json.Marshal(requestBody)

	// Create request
	req, _ := http.NewRequest("POST", "/plants", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	plant, ok := response["plant"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected plant in response")
	}

	if plant["name"] != "Test Plant" {
		t.Errorf("Expected plant name 'Test Plant', got %s", plant["name"])
	}
}

func TestPlantController_AddPlant_InvalidJSON(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	router.POST("/plants", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.AddPlant(c)
	})

	req, _ := http.NewRequest("POST", "/plants", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPlantController_AddPlant_ValidationError(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	router.POST("/plants", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.AddPlant(c)
	})

	// Empty name should fail validation
	requestBody := dto.PlantCreateRequest{
		Name:      "", // Empty name
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/plants", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPlantController_AddPlant_ServiceError(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	mockService.CreatePlantFunc = func(req *dto.PlantCreateRequest, userID int64) (*dto.PlantResponse, error) {
		return nil, errors.New("service error")
	}

	router.POST("/plants", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.AddPlant(c)
	})

	requestBody := dto.PlantCreateRequest{
		Name:      "Test Plant",
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/plants", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestPlantController_GetPlant_Success(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	expectedResponse := &dto.PlantResponse{
		ID:        1,
		Name:      "Test Plant",
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}

	mockService.GetPlantFunc = func(plantID, userID int64) (*dto.PlantResponse, error) {
		if plantID != 1 {
			t.Errorf("Expected plantID 1, got %d", plantID)
		}
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return expectedResponse, nil
	}

	router.GET("/plants/:id", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.GetPlant(c)
	})

	req, _ := http.NewRequest("GET", "/plants/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	plant, ok := response["plant"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected plant in response")
	}

	if plant["name"] != "Test Plant" {
		t.Errorf("Expected plant name 'Test Plant', got %s", plant["name"])
	}
}

func TestPlantController_GetPlant_InvalidID(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	router.GET("/plants/:id", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.GetPlant(c)
	})

	req, _ := http.NewRequest("GET", "/plants/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPlantController_GetPlant_NotFound(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	mockService.GetPlantFunc = func(plantID, userID int64) (*dto.PlantResponse, error) {
		return nil, errors.New("plant not found")
	}

	router.GET("/plants/:id", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.GetPlant(c)
	})

	req, _ := http.NewRequest("GET", "/plants/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestPlantController_GetPlants_Success(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	expectedPlants := []dto.PlantResponse{
		{ID: 1, Name: "Plant 1", TagColor: "green", PlantIcon: models.BigPlant},
		{ID: 2, Name: "Plant 2", TagColor: "blue", PlantIcon: models.SmallPlant},
	}

	mockService.GetPlantsFunc = func(userID int64) ([]dto.PlantResponse, error) {
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return expectedPlants, nil
	}

	router.GET("/plants", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.GetPlants(c)
	})

	req, _ := http.NewRequest("GET", "/plants", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	plants, ok := response["plants"].([]interface{})
	if !ok {
		t.Fatal("Expected plants array in response")
	}

	if len(plants) != 2 {
		t.Errorf("Expected 2 plants, got %d", len(plants))
	}
}

func TestPlantController_UpdatePlant_Success(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	expectedResponse := &dto.PlantResponse{
		ID:        1,
		Name:      "Updated Plant",
		TagColor:  "blue",
		PlantIcon: models.SmallPlant,
	}

	mockService.UpdatePlantFunc = func(req *dto.PlantUpdateRequest, plantID, userID int64) error {
		if req.Name != "Updated Plant" {
			t.Errorf("Expected name 'Updated Plant', got %s", req.Name)
		}
		if plantID != 1 {
			t.Errorf("Expected plantID 1, got %d", plantID)
		}
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return nil
	}

	mockService.GetPlantFunc = func(plantID, userID int64) (*dto.PlantResponse, error) {
		return expectedResponse, nil
	}

	router.PUT("/plants/:id", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.UpdatePlant(c)
	})

	requestBody := dto.PlantUpdateRequest{
		Name:      "Updated Plant",
		TagColor:  "blue",
		PlantIcon: models.SmallPlant,
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/plants/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestPlantController_DeletePlant_Success(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	mockService.DeletePlantFunc = func(userID, plantID int64) error {
		if plantID != 1 {
			t.Errorf("Expected plantID 1, got %d", plantID)
		}
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return nil
	}

	router.DELETE("/plants/:id", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.DeletePlant(c)
	})

	req, _ := http.NewRequest("DELETE", "/plants/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestPlantController_DeletePlant_InvalidID(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	router.DELETE("/plants/:id", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.DeletePlant(c)
	})

	req, _ := http.NewRequest("DELETE", "/plants/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPlantController_DeletePlant_ServiceError(t *testing.T) {
	mockService := &MockPlantService{}
	controller, router := setupPlantController(mockService)

	mockService.DeletePlantFunc = func(userID, plantID int64) error {
		return errors.New("delete failed")
	}

	router.DELETE("/plants/:id", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.DeletePlant(c)
	})

	req, _ := http.NewRequest("DELETE", "/plants/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
