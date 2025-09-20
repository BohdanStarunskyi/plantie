package service

import (
	"errors"
	"plant-reminder/dto"
	"plant-reminder/models"

	"gorm.io/gorm"
)

type PlantService struct {
	db *gorm.DB
}

type PlantServiceInterface interface {
	CreatePlant(plantRequest *dto.PlantCreateRequest, userID int64) (*dto.PlantResponse, error)
	GetPlant(plantID int64, userID int64) (*dto.PlantResponse, error)
	GetPlants(userID int64) ([]dto.PlantResponse, error)
	UpdatePlant(plant *dto.PlantUpdateRequest, plantId int64, userID int64) error
	DeletePlant(userID int64, plantID int64) error
}

func NewPlantService(db *gorm.DB) *PlantService {
	return &PlantService{
		db: db,
	}
}

func (s *PlantService) CreatePlant(plantRequest *dto.PlantCreateRequest, userID int64) (*dto.PlantResponse, error) {
	plant := plantRequest.ToModel(userID)

	if err := s.validatePlant(plant); err != nil {
		return nil, err
	}

	result := s.db.Create(plant)
	if result.Error != nil {
		return nil, result.Error
	}

	response := (&dto.PlantResponse{}).FromModel(plant)
	return response, nil
}

func (s *PlantService) UpdatePlant(plant *dto.PlantUpdateRequest, plantId int64, userID int64) error {
	var existingPlant models.Plant
	result := s.db.Where("id = ? AND user_id = ?", plantId, userID).First(&existingPlant)
	if result.Error != nil {
		return result.Error
	}

	updateModel := plant.ToModel(userID)
	updateModel.ID = plantId

	result = s.db.Model(&existingPlant).Updates(updateModel)
	return result.Error
}

func (s *PlantService) DeletePlant(userID int64, plantID int64) error {
	var plant models.Plant
	result := s.db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant)
	if result.Error != nil {
		return result.Error
	}

	result = s.db.Delete(&plant)
	return result.Error
}

func (s *PlantService) GetPlant(plantID int64, userID int64) (*dto.PlantResponse, error) {
	var plant models.Plant
	if plantID == 0 {
		return nil, errors.New("plantID must be set")
	}
	result := s.db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant)
	if result.Error != nil {
		return nil, result.Error
	}

	response := (&dto.PlantResponse{}).FromModel(&plant)
	return response, nil
}

func (s *PlantService) GetPlants(userID int64) ([]dto.PlantResponse, error) {
	var plants []models.Plant
	if userID == 0 {
		return nil, errors.New("userID must be set")
	}
	result := s.db.Where("user_id = ?", userID).Find(&plants)
	if result.Error != nil {
		return nil, result.Error
	}

	return dto.FromPlantsModel(plants), nil
}

func (s *PlantService) validatePlant(plant *models.Plant) error {
	if plant.UserID == 0 {
		return errors.New("user ID must be set")
	}
	if !plant.PlantIcon.IsValid() {
		return errors.New("invalid PlantIcon value")
	}
	if plant.Name == "" {
		return errors.New("name is required")
	}
	if plant.TagColor == "" {
		return errors.New("tagColor is required")
	}
	return nil
}
