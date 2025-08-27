package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"plant-reminder/dto"
	"testing"

	"github.com/gin-gonic/gin"
)

// MockUserService is a mock implementation of UserService for testing
type MockUserService struct {
	CreateUserFunc   func(*dto.UserCreateRequest) (*dto.AuthResponse, error)
	VerifyUserFunc   func(string, string) (*dto.AuthResponse, error)
	SetPushTokenFunc func(string, string) error
	DeleteUserFunc   func(int64) error
	GetUserFunc      func(int64) (*dto.UserResponse, error)
}

func (m *MockUserService) CreateUser(req *dto.UserCreateRequest) (*dto.AuthResponse, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(req)
	}
	return nil, nil
}

func (m *MockUserService) VerifyUser(email, password string) (*dto.AuthResponse, error) {
	if m.VerifyUserFunc != nil {
		return m.VerifyUserFunc(email, password)
	}
	return nil, nil
}

func (m *MockUserService) SetPushToken(userID, token string) error {
	if m.SetPushTokenFunc != nil {
		return m.SetPushTokenFunc(userID, token)
	}
	return nil
}

func (m *MockUserService) DeleteUser(userID int64) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(userID)
	}
	return nil
}

func (m *MockUserService) GetUser(userID int64) (*dto.UserResponse, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(userID)
	}
	return nil, nil
}

func setupUserController(mockService *MockUserService) (*UserController, *gin.Engine) {
	router := setupTestRouter()
	controller := &UserController{userService: mockService}
	return controller, router
}

func TestUserController_SignUp_Success(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	expectedResponse := &dto.AuthResponse{
		User: dto.UserResponse{
			ID:    1,
			Email: "test@example.com",
			Name:  "testuser",
		},
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	mockService.CreateUserFunc = func(req *dto.UserCreateRequest) (*dto.AuthResponse, error) {
		if req.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got %s", req.Email)
		}
		if req.Name != "testuser" {
			t.Errorf("Expected name 'testuser', got %s", req.Name)
		}
		return expectedResponse, nil
	}

	router.POST("/signup", controller.SignUp)

	requestBody := dto.UserCreateRequest{
		Email:    "test@example.com",
		Name:     "testuser",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response dto.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.User.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", response.User.Email)
	}
}

func TestUserController_SignUp_InvalidJSON(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	router.POST("/signup", controller.SignUp)

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserController_SignUp_ValidationError(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	router.POST("/signup", controller.SignUp)

	// Empty email should fail validation
	requestBody := dto.UserCreateRequest{
		Email:    "", // Empty email
		Name:     "testuser",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserController_SignUp_ServiceError(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	mockService.CreateUserFunc = func(req *dto.UserCreateRequest) (*dto.AuthResponse, error) {
		return nil, errors.New("user already exists")
	}

	router.POST("/signup", controller.SignUp)

	requestBody := dto.UserCreateRequest{
		Email:    "test@example.com",
		Name:     "testuser",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserController_Login_Success(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	expectedResponse := &dto.AuthResponse{
		User: dto.UserResponse{
			ID:    1,
			Email: "test@example.com",
			Name:  "testuser",
		},
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	mockService.VerifyUserFunc = func(email, password string) (*dto.AuthResponse, error) {
		if email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got %s", email)
		}
		if password != "password123" {
			t.Errorf("Expected password 'password123', got %s", password)
		}
		return expectedResponse, nil
	}

	router.POST("/login", controller.Login)

	requestBody := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response dto.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.User.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", response.User.Email)
	}
}

func TestUserController_Login_InvalidCredentials(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	mockService.VerifyUserFunc = func(email, password string) (*dto.AuthResponse, error) {
		return nil, errors.New("invalid credentials")
	}

	router.POST("/login", controller.Login)

	requestBody := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUserController_SetPushToken_Success(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	mockService.SetPushTokenFunc = func(userID, token string) error {
		if userID != "123" {
			t.Errorf("Expected userID '123', got %s", userID)
		}
		if token != "push_token_123" {
			t.Errorf("Expected token 'push_token_123', got %s", token)
		}
		return nil
	}

	router.POST("/push-token", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.SetPushToken(c)
	})

	requestBody := dto.PushTokenRequest{
		Token: "push_token_123",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/push-token", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserController_DeleteUser_Success(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	mockService.DeleteUserFunc = func(userID int64) error {
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return nil
	}

	router.DELETE("/user", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.DeleteUser(c)
	})

	req, _ := http.NewRequest("DELETE", "/user", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUserController_GetMyProfile_Success(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	expectedUser := &dto.UserResponse{
		ID:    123,
		Email: "test@example.com",
		Name:  "testuser",
	}

	mockService.GetUserFunc = func(userID int64) (*dto.UserResponse, error) {
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return expectedUser, nil
	}

	router.GET("/profile", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.GetMyProfile(c)
	})

	req, _ := http.NewRequest("GET", "/profile", nil)
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

	user, ok := response["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected user in response")
	}

	if user["email"] != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", user["email"])
	}
}

func TestUserController_GetMyProfile_UserNotFound(t *testing.T) {
	mockService := &MockUserService{}
	controller, router := setupUserController(mockService)

	mockService.GetUserFunc = func(userID int64) (*dto.UserResponse, error) {
		return nil, nil // User not found
	}

	router.GET("/profile", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.GetMyProfile(c)
	})

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
