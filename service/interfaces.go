package service

import "plant-reminder/dto"

type PlantServiceInterface interface {
	CreatePlant(plantRequest *dto.PlantCreateRequest, userID int64) (*dto.PlantResponse, error)
	GetPlant(plantID int64, userID int64) (*dto.PlantResponse, error)
	GetPlants(userID int64) ([]dto.PlantResponse, error)
	UpdatePlant(plant *dto.PlantUpdateRequest, plantId int64, userID int64) error
	DeletePlant(userID int64, plantID int64) error
}

type UserServiceInterface interface {
	CreateUser(userRequest *dto.UserCreateRequest) (*dto.AuthResponse, error)
	VerifyUser(email, password string) (*dto.AuthResponse, error)
	SetPushToken(userID, token string) error
	DeleteUser(userID int64) error
	GetUser(userID int64) (*dto.UserResponse, error)
}

type ReminderServiceInterface interface {
	CreateReminder(reminderRequest *dto.ReminderCreateRequest, plantId int64, userID int64) (*dto.ReminderResponse, error)
	GetReminder(reminderID int64, userID int64) (*dto.ReminderResponse, error)
	GetPlantReminders(plantID int64, userID int64) ([]dto.ReminderResponse, error)
	GetUserReminders(userID int64) ([]dto.ReminderResponse, error)
	UpdateReminder(reminder *dto.ReminderUpdateRequest, userID int64) error
	DeleteReminder(reminderID int64, userID int64) error
}
