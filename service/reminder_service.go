package service

import (
	"errors"
	"fmt"
	"plant-reminder/config"
	"plant-reminder/constants"
	"plant-reminder/dto"
	"plant-reminder/models"
	"plant-reminder/utils"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

var scheduler *gocron.Scheduler

type ReminderService struct {
	plantService *PlantService
}

func NewReminderService() *ReminderService {
	return &ReminderService{
		plantService: NewPlantService(),
	}
}

func (s *ReminderService) CreateReminder(reminderRequest *dto.ReminderCreateRequest, userID int64) (*dto.ReminderResponse, error) {
	// Check if plant exists and belongs to user
	_, err := s.plantService.GetPlant(reminderRequest.PlantID, userID)
	if err != nil {
		return nil, errors.New("plant doesn't exist")
	}

	reminder := reminderRequest.ToModel(userID)

	var existing models.Reminder
	err = config.DB.
		Where("plant_id = ? AND time_of_day = ? AND repeat = ?", reminder.PlantID, reminder.TimeOfDay, reminder.Repeat).
		First(&existing).Error

	if err == nil && existing.ID != reminder.ID {
		return nil, errors.New("reminder with same time and repeat period already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing reminders: %w", err)
	}

	if err := s.calculateNextTriggerTime(reminder); err != nil {
		return nil, err
	}

	result := config.DB.Create(reminder)
	if result.Error != nil {
		return nil, result.Error
	}

	response := (&dto.ReminderResponse{}).FromModel(reminder)
	return response, nil
}

func (s *ReminderService) UpdateReminder(reminderRequest *dto.ReminderUpdateRequest, userID int64) error {
	if reminderRequest.ID == 0 {
		return errors.New("reminder ID must be set")
	}

	existingReminder, err := s.getReminder(reminderRequest.ID, userID, reminderRequest.PlantID)
	if err != nil {
		return err
	}

	if existingReminder.UserID != userID {
		return errors.New("not enough rights")
	}

	reminder := reminderRequest.ToModel(userID)
	if err := s.calculateNextTriggerTime(reminder); err != nil {
		return err
	}

	result := config.DB.Model(&existingReminder).Updates(reminder)
	return result.Error
}

func (s *ReminderService) DeleteReminder(reminderID, userID int64) error {
	var reminder models.Reminder
	result := config.DB.Where("id = ? AND user_id = ?", reminderID, userID).First(&reminder)
	if result.Error != nil {
		return result.Error
	}
	result = config.DB.Delete(&reminder)
	return result.Error
}

func (s *ReminderService) GetPlantReminders(plantID int64, userID int64) ([]dto.ReminderResponse, error) {
	var reminders []models.Reminder
	if userID == 0 {
		return nil, errors.New("userID must be set")
	}
	if plantID == 0 {
		return nil, errors.New("plantID must be set")
	}
	result := config.DB.Where("plant_id = ? AND user_id = ?", plantID, userID).Find(&reminders)
	if result.Error != nil {
		return nil, result.Error
	}

	return dto.FromRemindersModel(reminders), nil
}

func (s *ReminderService) GetUserReminders(userID int64) ([]dto.ReminderResponse, error) {
	var reminders []models.Reminder
	if userID == 0 {
		return nil, errors.New("userID must be set")
	}
	result := config.DB.Preload("Plant").Where("user_id = ?", userID).Find(&reminders)
	if result.Error != nil {
		return nil, result.Error
	}

	return dto.FromRemindersModel(reminders), nil
}

func (s *ReminderService) GetAllReminders(userID int64) ([]dto.ReminderResponse, error) {
	return s.GetUserReminders(userID)
}

func (s *ReminderService) SetReminders() error {
	if scheduler != nil && scheduler.IsRunning() {
		return nil
	}
	scheduler = gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(1).Minutes().Do(func() {
		errChan := make(chan error)
		go s.checkReminders(errChan)
		for err := range errChan {
			if err != nil {
				fmt.Println("error during checking reminders", err)
			}
		}
	})
	if err != nil {
		return err
	}
	scheduler.StartAsync()
	return nil
}

func (s *ReminderService) getReminder(reminderID int64, userID int64, plantID int64) (models.Reminder, error) {
	var reminder models.Reminder
	if reminderID == 0 {
		return models.Reminder{}, errors.New("reminderID must be set")
	}
	result := config.DB.Where("id = ? AND user_id = ? AND plant_id = ?", reminderID, userID, plantID).First(&reminder)
	return reminder, result.Error
}

func (s *ReminderService) GetReminder(reminderID int64, userID int64) (*dto.ReminderResponse, error) {
	var reminder models.Reminder
	if reminderID == 0 {
		return nil, errors.New("reminderID must be set")
	}
	result := config.DB.Where("id = ? AND user_id = ?", reminderID, userID).First(&reminder)
	if result.Error != nil {
		return nil, result.Error
	}

	response := (&dto.ReminderResponse{}).FromModel(&reminder)
	return response, nil
}

func (s *ReminderService) calculateNextTriggerTime(reminder *models.Reminder) error {
	now := time.Now()
	t, err := time.Parse("15:04", reminder.TimeOfDay)
	if err != nil {
		return fmt.Errorf("invalid time format, expected HH:mm: %w", err)
	}

	nextTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
	if nextTime.Before(now) {
		switch reminder.Repeat {
		case constants.RepeatDaily:
			nextTime = nextTime.Add(24 * time.Hour)
		case constants.RepeatWeekly:
			nextTime = nextTime.Add(7 * 24 * time.Hour)
		case constants.RepeatMonthly:
			nextTime = nextTime.AddDate(0, 1, 0)
		}
	}

	reminder.NextTriggerTime = nextTime
	return nil
}

func (s *ReminderService) checkReminders(ch chan error) {
	defer close(ch)
	var reminders []models.Reminder
	err := config.DB.
		Where("next_trigger_time <= ?", time.Now()).
		Find(&reminders).Error

	if err != nil {
		ch <- err
		return
	}

	var wg sync.WaitGroup
	for _, value := range reminders {
		wg.Add(1)
		go func(pid int64) {
			defer wg.Done()
			s.sendNotifications(pid)
		}(value.PlantID)
	}
	wg.Wait()

	updatedReminders := utils.Map(reminders, func(item models.Reminder) models.Reminder {
		switch item.Repeat {
		case constants.RepeatDaily:
			item.NextTriggerTime = item.NextTriggerTime.Add(24 * time.Hour)
		case constants.RepeatWeekly:
			item.NextTriggerTime = item.NextTriggerTime.Add(7 * 24 * time.Hour)
		case constants.RepeatMonthly:
			item.NextTriggerTime = item.NextTriggerTime.AddDate(0, 1, 0)
		}
		return item
	})

	var result *gorm.DB
	if len(updatedReminders) != 0 {
		result = config.DB.Save(&updatedReminders)
	}

	if result != nil {
		ch <- result.Error
	} else {
		ch <- nil
	}
}

func (s *ReminderService) sendNotifications(plantID int64) {
	var plant models.Plant
	config.DB.Where("id = ?", plantID).First(&plant)
	var user models.User
	config.DB.Where("id = ?", plant.UserID).First(&user)
	if user.PushToken != "" {
		utils.SendMessage(user.PushToken, plant.Name)
	}
}
