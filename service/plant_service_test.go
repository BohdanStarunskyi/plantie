package service

import (
	"plant-reminder/dto"
	"plant-reminder/models"
	"testing"
)

type MockDB struct {
	CreateFunc  func(interface{}) error
	FirstFunc   func(interface{}, ...interface{}) error
	FindFunc    func(interface{}, ...interface{}) error
	UpdatesFunc func(interface{}, interface{}) error
	DeleteFunc  func(interface{}) error
}

func (m *MockDB) Create(value interface{}) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(value)
	}
	return nil
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) error {
	if m.FirstFunc != nil {
		return m.FirstFunc(dest, conds...)
	}
	return nil
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) error {
	if m.FindFunc != nil {
		return m.FindFunc(dest, conds...)
	}
	return nil
}

func (m *MockDB) Updates(dest interface{}, values interface{}) error {
	if m.UpdatesFunc != nil {
		return m.UpdatesFunc(dest, values)
	}
	return nil
}

func (m *MockDB) Delete(value interface{}) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(value)
	}
	return nil
}

func TestPlantService_CreatePlant_Success(t *testing.T) {
	plantRequest := &dto.PlantCreateRequest{
		Name:      "Test Plant",
		Note:      "Test note",
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}

	userID := int64(123)

	plant := plantRequest.ToModel(userID)

	if plant.Name != "Test Plant" {
		t.Errorf("Expected name 'Test Plant', got %s", plant.Name)
	}

	if plant.UserID != userID {
		t.Errorf("Expected userID %d, got %d", userID, plant.UserID)
	}

	if plant.TagColor != "green" {
		t.Errorf("Expected tagColor 'green', got %s", plant.TagColor)
	}

	if plant.PlantIcon != models.BigPlant {
		t.Errorf("Expected plantIcon %s, got %s", models.BigPlant, plant.PlantIcon)
	}
}

func TestPlantService_ValidatePlant_Success(t *testing.T) {
	service := NewPlantService()

	plant := &models.Plant{
		Name:      "Test Plant",
		UserID:    123,
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}

	err := service.validatePlant(plant)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestPlantService_ValidatePlant_InvalidUserID(t *testing.T) {
	service := NewPlantService()

	plant := &models.Plant{
		Name:      "Test Plant",
		UserID:    0,
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}

	err := service.validatePlant(plant)
	if err == nil {
		t.Error("Expected error for invalid userID, got nil")
	}

	expectedError := "user ID must be set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestPlantService_ValidatePlant_InvalidPlantIcon(t *testing.T) {
	service := NewPlantService()

	plant := &models.Plant{
		Name:      "Test Plant",
		UserID:    123,
		TagColor:  "green",
		PlantIcon: models.PlantIcon("invalid"),
	}

	err := service.validatePlant(plant)
	if err == nil {
		t.Error("Expected error for invalid plant icon, got nil")
	}

	expectedError := "invalid PlantIcon value"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestPlantService_ValidatePlant_EmptyName(t *testing.T) {
	service := NewPlantService()

	plant := &models.Plant{
		Name:      "",
		UserID:    123,
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}

	err := service.validatePlant(plant)
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}

	expectedError := "name is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestPlantService_ValidatePlant_EmptyTagColor(t *testing.T) {
	service := NewPlantService()

	plant := &models.Plant{
		Name:      "Test Plant",
		UserID:    123,
		TagColor:  "",
		PlantIcon: models.BigPlant,
	}

	err := service.validatePlant(plant)
	if err == nil {
		t.Error("Expected error for empty tagColor, got nil")
	}

	expectedError := "tagColor is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestPlantService_GetPlant_InvalidPlantID(t *testing.T) {
	service := NewPlantService()

	_, err := service.GetPlant(0, 123)
	if err == nil {
		t.Error("Expected error for invalid plantID, got nil")
	}

	expectedError := "plantID must be set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestPlantService_GetPlants_InvalidUserID(t *testing.T) {
	service := NewPlantService()

	_, err := service.GetPlants(0)
	if err == nil {
		t.Error("Expected error for invalid userID, got nil")
	}

	expectedError := "userID must be set"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestPlantCreateRequest_ToModel(t *testing.T) {
	request := &dto.PlantCreateRequest{
		Name:      "Test Plant",
		Note:      "Test note",
		TagColor:  "green",
		PlantIcon: models.BigPlant,
	}

	userID := int64(123)
	plant := request.ToModel(userID)

	if plant.Name != request.Name {
		t.Errorf("Expected name %s, got %s", request.Name, plant.Name)
	}

	if plant.Note != request.Note {
		t.Errorf("Expected note %s, got %s", request.Note, plant.Note)
	}

	if plant.TagColor != request.TagColor {
		t.Errorf("Expected tagColor %s, got %s", request.TagColor, plant.TagColor)
	}

	if plant.PlantIcon != request.PlantIcon {
		t.Errorf("Expected plantIcon %s, got %s", request.PlantIcon, plant.PlantIcon)
	}

	if plant.UserID != userID {
		t.Errorf("Expected userID %d, got %d", userID, plant.UserID)
	}
}

func TestPlantUpdateRequest_ToModel(t *testing.T) {
	request := &dto.PlantUpdateRequest{
		Name:      "Updated Plant",
		Note:      "Updated note",
		TagColor:  "blue",
		PlantIcon: models.SmallPlant,
	}

	userID := int64(123)
	plant := request.ToModel(userID)

	if plant.Name != request.Name {
		t.Errorf("Expected name %s, got %s", request.Name, plant.Name)
	}

	if plant.Note != request.Note {
		t.Errorf("Expected note %s, got %s", request.Note, plant.Note)
	}

	if plant.TagColor != request.TagColor {
		t.Errorf("Expected tagColor %s, got %s", request.TagColor, plant.TagColor)
	}

	if plant.PlantIcon != request.PlantIcon {
		t.Errorf("Expected plantIcon %s, got %s", request.PlantIcon, plant.PlantIcon)
	}

	if plant.UserID != userID {
		t.Errorf("Expected userID %d, got %d", userID, plant.UserID)
	}
}

func TestPlantResponse_FromModel(t *testing.T) {
	plant := &models.Plant{
		ID:        1,
		Name:      "Test Plant",
		Note:      "Test note",
		TagColor:  "green",
		UserID:    123,
		PlantIcon: models.BigPlant,
	}

	response := (&dto.PlantResponse{}).FromModel(plant)

	if response.ID != plant.ID {
		t.Errorf("Expected ID %d, got %d", plant.ID, response.ID)
	}

	if response.Name != plant.Name {
		t.Errorf("Expected name %s, got %s", plant.Name, response.Name)
	}

	if response.Note != plant.Note {
		t.Errorf("Expected note %s, got %s", plant.Note, response.Note)
	}

	if response.TagColor != plant.TagColor {
		t.Errorf("Expected tagColor %s, got %s", plant.TagColor, response.TagColor)
	}

	if response.PlantIcon != plant.PlantIcon {
		t.Errorf("Expected plantIcon %s, got %s", plant.PlantIcon, response.PlantIcon)
	}
}
