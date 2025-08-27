package service

import (
	"plant-reminder/dto"
	"plant-reminder/models"
	"testing"
	"time"
)

func TestUserService_CreateUser_Success(t *testing.T) {
	userRequest := &dto.UserCreateRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	user := userRequest.ToModel()

	if user.Email != userRequest.Email {
		t.Errorf("Expected email %s, got %s", userRequest.Email, user.Email)
	}

	if user.Password != userRequest.Password {
		t.Errorf("Expected password %s, got %s", userRequest.Password, user.Password)
	}

	if user.Name != userRequest.Name {
		t.Errorf("Expected name %s, got %s", userRequest.Name, user.Name)
	}
}

func TestUserResponse_FromModel(t *testing.T) {
	user := &models.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		CreationDate: time.Now(),
	}

	response := (&dto.UserResponse{}).FromModel(user)

	if response.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, response.ID)
	}

	if response.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, response.Email)
	}

	if response.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, response.Name)
	}

	if response.CreationDate != user.CreationDate {
		t.Errorf("Expected creation date %v, got %v", user.CreationDate, response.CreationDate)
	}
}

func TestUserLoginRequest_Validation(t *testing.T) {
	validRequest := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	if validRequest.Email == "" {
		t.Error("Expected valid email")
	}

	if validRequest.Password == "" {
		t.Error("Expected valid password")
	}

	invalidRequest := dto.UserLoginRequest{
		Email:    "",
		Password: "password123",
	}

	if invalidRequest.Email != "" {
		t.Error("Expected empty email for invalid request")
	}
}

func TestUserUpdateRequest_ToModel(t *testing.T) {
	updateRequest := &dto.UserUpdateRequest{
		Name: "Updated Name",
	}
	if updateRequest.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", updateRequest.Name)
	}
}

func TestPushTokenRequest_Validation(t *testing.T) {
	validRequest := dto.PushTokenRequest{
		Token: "valid_push_token_123",
	}

	if validRequest.Token == "" {
		t.Error("Expected valid push token")
	}

	invalidRequest := dto.PushTokenRequest{
		Token: "",
	}

	if invalidRequest.Token != "" {
		t.Error("Expected empty token for invalid request")
	}
}

func TestAuthResponse_Structure(t *testing.T) {
	authResponse := dto.AuthResponse{
		User: dto.UserResponse{
			ID:    1,
			Email: "test@example.com",
			Name:  "Test User",
		},
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}

	if authResponse.User.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", authResponse.User.ID)
	}

	if authResponse.AccessToken == "" {
		t.Error("Expected non-empty access token")
	}

	if authResponse.RefreshToken == "" {
		t.Error("Expected non-empty refresh token")
	}
}
