package dto

import "plant-reminder/models"

type PlantCreateRequest struct {
	Name      string           `json:"name" validate:"required"`
	Note      string           `json:"note"`
	TagColor  string           `json:"tagColor" validate:"required"`
	PlantIcon models.PlantIcon `json:"plantIcon" validate:"required"`
}

type PlantUpdateRequest struct {
	Name      string           `json:"name" validate:"required"`
	Note      string           `json:"note"`
	TagColor  string           `json:"tagColor" validate:"required"`
	PlantIcon models.PlantIcon `json:"plantIcon" validate:"required"`
}

type PlantResponse struct {
	ID        int64              `json:"id"`
	Name      string             `json:"name"`
	Note      string             `json:"note"`
	TagColor  string             `json:"tagColor"`
	PlantIcon models.PlantIcon   `json:"plantIcon"`
	Reminders []ReminderResponse `json:"reminders,omitempty"`
}

func (r *PlantCreateRequest) ToModel(userID int64) *models.Plant {
	return &models.Plant{
		Name:      r.Name,
		Note:      r.Note,
		TagColor:  r.TagColor,
		UserID:    userID,
		PlantIcon: r.PlantIcon,
	}
}

func (r *PlantUpdateRequest) ToModel(userID int64) *models.Plant {
	return &models.Plant{
		Name:      r.Name,
		Note:      r.Note,
		TagColor:  r.TagColor,
		UserID:    userID,
		PlantIcon: r.PlantIcon,
	}
}

func (r *PlantResponse) FromModel(plant *models.Plant) *PlantResponse {
	response := &PlantResponse{
		ID:        plant.ID,
		Name:      plant.Name,
		Note:      plant.Note,
		TagColor:  plant.TagColor,
		PlantIcon: plant.PlantIcon,
	}

	if plant.Reminders != nil {
		response.Reminders = make([]ReminderResponse, len(plant.Reminders))
		for i, reminder := range plant.Reminders {
			response.Reminders[i].FromModel(&reminder)
		}
	}

	return response
}

func FromPlantsModel(plants []models.Plant) []PlantResponse {
	responses := make([]PlantResponse, len(plants))
	for i, plant := range plants {
		responses[i] = *(&PlantResponse{}).FromModel(&plant)
	}
	return responses
}
