package dto

import (
	"plant-reminder/models"
	"time"
)

type UserCreateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
}

type UserUpdateRequest struct {
	Name string `json:"name" validate:"omitempty,min=2,max=100"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID           int64           `json:"id"`
	Email        string          `json:"email"`
	Name         string          `json:"name"`
	CreationDate time.Time       `json:"createdAt"`
	Plants       []PlantResponse `json:"plants,omitempty"`
}

type PushTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

func (r *UserCreateRequest) ToModel() *models.User {
	return &models.User{
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
	}
}

func (r *UserResponse) FromModel(user *models.User) *UserResponse {
	response := &UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		CreationDate: user.CreationDate,
	}

	if user.Plants != nil {
		response.Plants = FromPlantsModel(user.Plants)
	}

	return response
}
