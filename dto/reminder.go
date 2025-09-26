package dto

import (
	"plant-reminder/constants"
	"plant-reminder/models"
	"time"

	"github.com/go-playground/validator"
)

type ReminderCreateRequest struct {
	RepeatType constants.RepeatType `json:"repeatType"`
	TimeOfDay  string               `json:"timeOfDay" validate:"required,len=5"`
	DayOfWeek  *int16               `json:"dayOfWeek" validate:"omitempty,min=0,max=6"`
	DayOfMonth *int16               `json:"dayOfMonth" validate:"omitempty,min=1,max=31"`
}

type ReminderUpdateRequest struct {
	ID         int64                `json:"id" validate:"required"`
	RepeatType constants.RepeatType `json:"repeatType"`
	TimeOfDay  string               `json:"timeOfDay" validate:"required,len=5"`
	DayOfWeek  *int16               `json:"dayOfWeek" validate:"omitempty,min=0,max=6"`
	DayOfMonth *int16               `json:"dayOfMonth" validate:"omitempty,min=1,max=31"`
}

type ReminderResponse struct {
	ID              int64                `json:"id"`
	Repeat          constants.RepeatType `json:"repeatType"`
	TimeOfDay       string               `json:"timeOfDay"`
	NextTriggerTime time.Time            `json:"nextTriggerTime"`
	Plant           *PlantResponse       `json:"plant,omitempty"`
	DayOfWeek       *int16               `json:"dayOfWeek"`
	DayOfMonth      *int16               `json:"dayOfMonth"`
}

func (r *ReminderCreateRequest) ToModel(userID int64) *models.Reminder {
	return &models.Reminder{
		Repeat:     r.RepeatType,
		TimeOfDay:  r.TimeOfDay,
		UserID:     userID,
		DayOfMonth: r.DayOfMonth,
		DayOfWeek:  r.DayOfWeek,
	}
}

func (r *ReminderUpdateRequest) ToModel(userID int64, plantId int64) *models.Reminder {
	return &models.Reminder{
		ID:         r.ID,
		PlantID:    plantId,
		Repeat:     r.RepeatType,
		TimeOfDay:  r.TimeOfDay,
		UserID:     userID,
		DayOfMonth: r.DayOfMonth,
		DayOfWeek:  r.DayOfWeek,
	}
}

func (r *ReminderResponse) FromModel(reminder *models.Reminder) *ReminderResponse {
	response := &ReminderResponse{
		ID:              reminder.ID,
		Repeat:          reminder.Repeat,
		TimeOfDay:       reminder.TimeOfDay,
		NextTriggerTime: reminder.NextTriggerTime,
		DayOfMonth:      reminder.DayOfMonth,
		DayOfWeek:       reminder.DayOfWeek,
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

var validate = validator.New()

func (r *ReminderCreateRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}
	return constants.ValidateReminderFields(r.RepeatType, r.DayOfWeek, r.DayOfMonth)
}

func (r *ReminderUpdateRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}
	return constants.ValidateReminderFields(r.RepeatType, r.DayOfWeek, r.DayOfMonth)
}
