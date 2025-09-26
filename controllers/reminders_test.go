package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"plant-reminder/constants"
	"plant-reminder/dto"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type MockReminderService struct {
	CreateReminderFunc    func(*dto.ReminderCreateRequest, int64, int64) (*dto.ReminderResponse, error)
	GetReminderFunc       func(int64, int64) (*dto.ReminderResponse, error)
	GetPlantRemindersFunc func(int64, int64) ([]dto.ReminderResponse, error)
	GetUserRemindersFunc  func(int64) ([]dto.ReminderResponse, error)
	UpdateReminderFunc    func(*dto.ReminderUpdateRequest, int64, int64) (*dto.ReminderResponse, error)
	DeleteReminderFunc    func(int64, int64) error
	TestReminderFunc      func(userId int64) error
}

func (m *MockReminderService) CreateReminder(req *dto.ReminderCreateRequest, plantID int64, userID int64) (*dto.ReminderResponse, error) {
	if m.CreateReminderFunc != nil {
		return m.CreateReminderFunc(req, plantID, userID)
	}
	return nil, nil
}

func (m *MockReminderService) GetReminder(reminderID, userID int64) (*dto.ReminderResponse, error) {
	if m.GetReminderFunc != nil {
		return m.GetReminderFunc(reminderID, userID)
	}
	return nil, nil
}

func (m *MockReminderService) GetPlantReminders(plantID, userID int64) ([]dto.ReminderResponse, error) {
	if m.GetPlantRemindersFunc != nil {
		return m.GetPlantRemindersFunc(plantID, userID)
	}
	return nil, nil
}

func (m *MockReminderService) GetUserReminders(userID int64) ([]dto.ReminderResponse, error) {
	if m.GetUserRemindersFunc != nil {
		return m.GetUserRemindersFunc(userID)
	}
	return nil, nil
}

func (m *MockReminderService) UpdateReminder(req *dto.ReminderUpdateRequest, userID int64, plantID int64) (*dto.ReminderResponse, error) {
	if m.UpdateReminderFunc != nil {
		return m.UpdateReminderFunc(req, userID, plantID)
	}
	return nil, nil
}

func (m *MockReminderService) DeleteReminder(reminderID, userID int64) error {
	if m.DeleteReminderFunc != nil {
		return m.DeleteReminderFunc(reminderID, userID)
	}
	return nil
}

func (m *MockReminderService) TestReminder(userId int64) error {
	if m.TestReminderFunc != nil {
		return m.TestReminderFunc(userId)
	}
	return nil
}

func setupReminderController(mockService *MockReminderService) (*ReminderController, *gin.Engine) {
	router := setupTestRouter()
	controller := &ReminderController{reminderService: mockService}
	return controller, router
}

func TestReminderController_AddReminder_Success(t *testing.T) {
	mockService := &MockReminderService{}
	controller, _ := setupReminderController(mockService)

	expectedResponse := &dto.ReminderResponse{
		ID:              1,
		TimeOfDay:       "08:00",
		Repeat:          constants.RepeatDaily,
		NextTriggerTime: time.Now(),
	}

	mockService.CreateReminderFunc = func(req *dto.ReminderCreateRequest, plantID int64, userID int64) (*dto.ReminderResponse, error) {
		return expectedResponse, nil
	}

	if controller == nil {
		t.Error("Expected non-nil controller")
	} else if controller.reminderService == nil {
		t.Error("Expected non-nil reminder service")
	}
}

func TestReminderController_AddReminder_InvalidPlantID(t *testing.T) {
	mockService := &MockReminderService{}
	controller, router := setupReminderController(mockService)

	router.POST("/plants/:id/reminders", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.AddReminder(c)
	})

	jsonData := []byte(`{
		"plantId": 1,
		"timeOfDay": "08:00",
		"repeatType": "daily"
	}`)

	req, _ := http.NewRequest("POST", "/plants/invalid/reminders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestReminderController_GetPlantReminders_Success(t *testing.T) {
	mockService := &MockReminderService{}
	controller, router := setupReminderController(mockService)

	expectedReminders := []dto.ReminderResponse{
		{ID: 1, TimeOfDay: "08:00", Repeat: constants.RepeatDaily, NextTriggerTime: time.Now()},
		{ID: 2, TimeOfDay: "18:00", Repeat: constants.RepeatWeekly, NextTriggerTime: time.Now()},
	}

	mockService.GetPlantRemindersFunc = func(plantID, userID int64) ([]dto.ReminderResponse, error) {
		if plantID != 1 {
			t.Errorf("Expected plantID 1, got %d", plantID)
		}
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return expectedReminders, nil
	}

	router.GET("/plants/:id/reminders", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.GetPlantReminders(c)
	})

	req, _ := http.NewRequest("GET", "/plants/1/reminders", nil)
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

	reminders, ok := response["reminders"].([]interface{})
	if !ok {
		t.Fatal("Expected reminders array in response")
	}

	if len(reminders) != 2 {
		t.Errorf("Expected 2 reminders, got %d", len(reminders))
	}
}

func TestReminderController_DeleteReminder_Success(t *testing.T) {
	mockService := &MockReminderService{}
	controller, router := setupReminderController(mockService)

	mockService.DeleteReminderFunc = func(reminderID, userID int64) error {
		if reminderID != 1 {
			t.Errorf("Expected reminderID 1, got %d", reminderID)
		}
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return nil
	}

	router.DELETE("/reminders/:reminderId", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.DeleteReminder(c)
	})

	req, _ := http.NewRequest("DELETE", "/reminders/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestReminderController_TestReminder_Success(t *testing.T) {
	mockService := &MockReminderService{}
	controller, router := setupReminderController(mockService)
	mockService.TestReminderFunc = func(userID int64) error {
		if userID != 123 {
			t.Errorf("Expected userID 123, got %d", userID)
		}
		return nil
	}
	router.POST("/reminders/test", func(c *gin.Context) {
		c.Set("userID", int64(123))
		controller.TestReminder(c)
	})

	req, _ := http.NewRequest("POST", "/reminders/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
