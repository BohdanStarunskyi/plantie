package service

import (
	"plant-reminder/constants"
	"plant-reminder/dto"
	"plant-reminder/models"
	"testing"
	"time"
)

func TestReminderService_CreateReminder_DTOConversion(t *testing.T) {
	reminderRequest := &dto.ReminderCreateRequest{
		RepeatType: constants.RepeatDaily,
		TimeOfDay:  "08:00",
	}

	userID := int64(123)

	reminder := reminderRequest.ToModel(userID)

	if reminder.Repeat != reminderRequest.RepeatType {
		t.Errorf("Expected Repeat %v, got %v", reminderRequest.RepeatType, reminder.Repeat)
	}

	if reminder.TimeOfDay != reminderRequest.TimeOfDay {
		t.Errorf("Expected TimeOfDay %s, got %s", reminderRequest.TimeOfDay, reminder.TimeOfDay)
	}

	if reminder.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, reminder.UserID)
	}
}

func TestReminderUpdateRequest_DTOConversion(t *testing.T) {
	reminderRequest := &dto.ReminderUpdateRequest{
		ID:         1,
		PlantID:    2,
		RepeatType: constants.RepeatWeekly,
		TimeOfDay:  "18:00",
	}

	userID := int64(123)

	reminder := reminderRequest.ToModel(userID)

	if reminder.ID != reminderRequest.ID {
		t.Errorf("Expected ID %d, got %d", reminderRequest.ID, reminder.ID)
	}

	if reminder.PlantID != reminderRequest.PlantID {
		t.Errorf("Expected PlantID %d, got %d", reminderRequest.PlantID, reminder.PlantID)
	}

	if reminder.Repeat != reminderRequest.RepeatType {
		t.Errorf("Expected Repeat %v, got %v", reminderRequest.RepeatType, reminder.Repeat)
	}

	if reminder.TimeOfDay != reminderRequest.TimeOfDay {
		t.Errorf("Expected TimeOfDay %s, got %s", reminderRequest.TimeOfDay, reminder.TimeOfDay)
	}

	if reminder.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, reminder.UserID)
	}
}

func TestReminderResponse_FromModel(t *testing.T) {
	reminder := &models.Reminder{
		ID:              1,
		PlantID:         2,
		Repeat:          constants.RepeatDaily,
		TimeOfDay:       "08:00",
		UserID:          123,
		NextTriggerTime: time.Now(),
	}

	response := (&dto.ReminderResponse{}).FromModel(reminder)

	if response.ID != reminder.ID {
		t.Errorf("Expected ID %d, got %d", reminder.ID, response.ID)
	}

	if response.Repeat != reminder.Repeat {
		t.Errorf("Expected Repeat %v, got %v", reminder.Repeat, response.Repeat)
	}

	if response.TimeOfDay != reminder.TimeOfDay {
		t.Errorf("Expected TimeOfDay %s, got %s", reminder.TimeOfDay, response.TimeOfDay)
	}

	if response.NextTriggerTime != reminder.NextTriggerTime {
		t.Errorf("Expected NextTriggerTime %v, got %v", reminder.NextTriggerTime, response.NextTriggerTime)
	}
}

func TestRepeatType_Validation(t *testing.T) {
	validTypes := []constants.RepeatType{
		constants.RepeatDaily,
		constants.RepeatWeekly,
		constants.RepeatMonthly,
	}

	for _, repeatType := range validTypes {
		if repeatType < 0 {
			t.Errorf("Invalid repeat type: %v", repeatType)
		}
	}
}

func TestRepeatType_String(t *testing.T) {
	testCases := []struct {
		repeatType constants.RepeatType
		expected   string
	}{
		{constants.RepeatDaily, "daily"},
		{constants.RepeatWeekly, "weekly"},
		{constants.RepeatMonthly, "monthly"},
	}

	for _, tc := range testCases {
		result := tc.repeatType.String()
		if result != tc.expected {
			t.Errorf("Expected %s, got %s for repeat type %v", tc.expected, result, tc.repeatType)
		}
	}
}

func TestReminderCreateRequest_Validation(t *testing.T) {
	validRequest := dto.ReminderCreateRequest{RepeatType: constants.RepeatDaily, TimeOfDay: "08:00"}

	if validRequest.TimeOfDay == "" {
		t.Error("Expected valid TimeOfDay")
	}

	if validRequest.RepeatType.String() == "" {
		t.Error("Expected valid RepeatType")
	}

	invalidRequest := dto.ReminderCreateRequest{RepeatType: constants.RepeatDaily, TimeOfDay: ""}

	if invalidRequest.TimeOfDay != "" {
		t.Error("Expected invalid TimeOfDay for invalid request")
	}
}

func TestReminderUpdateRequest_Validation(t *testing.T) {
	validRequest := dto.ReminderUpdateRequest{ID: 1, PlantID: 1, RepeatType: constants.RepeatDaily, TimeOfDay: "08:00"}

	if validRequest.ID <= 0 {
		t.Error("Expected valid ID")
	}

	if validRequest.PlantID <= 0 {
		t.Error("Expected valid PlantID")
	}

	if validRequest.TimeOfDay == "" {
		t.Error("Expected valid TimeOfDay")
	}

	invalidRequest := dto.ReminderUpdateRequest{
		ID:         0,
		PlantID:    1,
		RepeatType: constants.RepeatDaily,
		TimeOfDay:  "08:00",
	}

	if invalidRequest.ID > 0 {
		t.Error("Expected invalid ID for invalid request")
	}
}
