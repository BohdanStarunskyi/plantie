package dto

import (
	"plant-reminder/constants"
	"plant-reminder/models"
	"time"
)

type ReminderCreateRequest struct {
	PlantID   int64                `json:"plantId" validate:"required"`
	Repeat    constants.RepeatType `json:"repeatType" validate:"required"`
	TimeOfDay string               `json:"timeOfDay" validate:"required,len=5"`
}

type ReminderUpdateRequest struct {
	ID        int64                `json:"id" validate:"required"`
	PlantID   int64                `json:"plantId" validate:"required"`
	Repeat    constants.RepeatType `json:"repeatType" validate:"required"`
	TimeOfDay string               `json:"timeOfDay" validate:"required,len=5"`
}

type ReminderResponse struct {
	ID              int64                `json:"id"`
	Repeat          constants.RepeatType `json:"repeatType"`
	TimeOfDay       string               `json:"timeOfDay"`
	NextTriggerTime time.Time            `json:"nextTriggerTime"`
	Plant           *PlantResponse       `json:"plant,omitempty"`
}

func (r *ReminderCreateRequest) ToModel(userID int64) *models.Reminder {
	return &models.Reminder{
		PlantID:   r.PlantID,
		Repeat:    r.Repeat,
		TimeOfDay: r.TimeOfDay,
		UserID:    userID,
	}
}

func (r *ReminderUpdateRequest) ToModel(userID int64) *models.Reminder {
	return &models.Reminder{
		ID:        r.ID,
		PlantID:   r.PlantID,
		Repeat:    r.Repeat,
		TimeOfDay: r.TimeOfDay,
		UserID:    userID,
	}
}

func (r *ReminderResponse) FromModel(reminder *models.Reminder) *ReminderResponse {
	response := &ReminderResponse{
		ID:              reminder.ID,
		Repeat:          reminder.Repeat,
		TimeOfDay:       reminder.TimeOfDay,
		NextTriggerTime: reminder.NextTriggerTime,
	}

	if reminder.Plant != nil {
		response.Plant = (&PlantResponse{}).FromModel(reminder.Plant)
	}

	return response
}

func FromRemindersModel(reminders []models.Reminder) []ReminderResponse {
	responses := make([]ReminderResponse, len(reminders))
	for i, reminder := range reminders {
		responses[i] = *(&ReminderResponse{}).FromModel(&reminder)
	}
	return responses
}
