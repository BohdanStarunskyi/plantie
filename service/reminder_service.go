package service

import (
	"errors"
	"fmt"
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
	db           *gorm.DB
}

type ReminderServiceInterface interface {
	CreateReminder(reminderRequest *dto.ReminderCreateRequest, plantId int64, userID int64) (*dto.ReminderResponse, error)
	GetReminder(reminderID int64, userID int64) (*dto.ReminderResponse, error)
	GetPlantReminders(plantID int64, userID int64) ([]dto.ReminderResponse, error)
	GetUserReminders(userID int64) ([]dto.ReminderResponse, error)
	UpdateReminder(reminder *dto.ReminderUpdateRequest, userID int64, plantId int64) error
	DeleteReminder(reminderID int64, userID int64) error
	TestReminder(userId int64) error
}

func NewReminderService(ps *PlantService, db *gorm.DB) *ReminderService {
	return &ReminderService{
		plantService: ps,
		db:           db,
	}
}

func (s *ReminderService) CreateReminder(reminderRequest *dto.ReminderCreateRequest, plantId int64, userID int64) (*dto.ReminderResponse, error) {
	_, err := s.plantService.GetPlant(plantId, userID)
	if err != nil {
		return nil, errors.New("plant doesn't exist")
	}

	reminder := reminderRequest.ToModel(userID)

	var existing models.Reminder
	err = s.db.
		Where("plant_id = ? AND time_of_day = ? AND repeat = ?", plantId, reminder.TimeOfDay, reminder.Repeat).
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

	reminder.PlantID = plantId

	result := s.db.Create(reminder)
	if result.Error != nil {
		return nil, result.Error
	}

	response := (&dto.ReminderResponse{}).FromModel(reminder)
	return response, nil
}

func (s *ReminderService) UpdateReminder(reminderRequest *dto.ReminderUpdateRequest, userID int64, plantID int64) error {
	if reminderRequest.ID == 0 {
		return errors.New("reminder ID must be set")
	}

	existingReminder, err := s.getReminder(reminderRequest.ID)
	if err != nil {
		return err
	}

	if existingReminder.UserID != userID {
		return errors.New("not enough rights")
	}

	reminder := reminderRequest.ToModel(userID, plantID)
	if err := s.calculateNextTriggerTime(reminder); err != nil {
		return err
	}

	result := s.db.Model(&existingReminder).Updates(reminder)
	return result.Error
}

func (s *ReminderService) DeleteReminder(reminderID, userID int64) error {
	var reminder models.Reminder
	result := s.db.Where("id = ? AND user_id = ?", reminderID, userID).First(&reminder)
	if result.Error != nil {
		return result.Error
	}
	result = s.db.Delete(&reminder)
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
	result := s.db.Where("plant_id = ? AND user_id = ?", plantID, userID).Find(&reminders)
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
	result := s.db.Preload("Plant").Where("user_id = ?", userID).Find(&reminders)
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

func (s *ReminderService) getReminder(reminderID int64) (models.Reminder, error) {
	var reminder models.Reminder
	if reminderID == 0 {
		return models.Reminder{}, errors.New("reminderID must be set")
	}
	result := s.db.Where("id = ?", reminderID).First(&reminder)
	return reminder, result.Error
}

func (s *ReminderService) GetReminder(reminderID int64, userID int64) (*dto.ReminderResponse, error) {
	var reminder models.Reminder
	if reminderID == 0 {
		return nil, errors.New("reminderID must be set")
	}
	result := s.db.Where("id = ? AND user_id = ?", reminderID, userID).First(&reminder)
	if result.Error != nil {
		return nil, result.Error
	}

	response := (&dto.ReminderResponse{}).FromModel(&reminder)
	return response, nil
}

func (s *ReminderService) TestReminder(userID int64) error {
	var user models.User
	s.db.Where("id = ?", userID).First(&user)
	if user.PushToken == "" {
		return errors.New("user doesn't have push token")
	}
	utils.SendMessage(user.PushToken, "test")
	return nil
}

func (s *ReminderService) calculateNextTriggerTime(reminder *models.Reminder) error {
	now := time.Now()

	t, err := time.Parse("15:04", reminder.TimeOfDay)
	if err != nil {
		return fmt.Errorf("invalid time format, expected HH:mm: %w", err)
	}

	nextTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		t.Hour(), t.Minute(), 0, 0, time.Local,
	)

	switch reminder.Repeat {
	case constants.RepeatDaily:
		if nextTime.Before(now) {
			nextTime = nextTime.Add(24 * time.Hour)
		}

	case constants.RepeatWeekly:
		if reminder.DayOfWeek == nil {
			return errors.New("weekly reminder requires dayOfWeek")
		}
		targetWeekday := time.Weekday(*reminder.DayOfWeek)

		daysUntil := (int(targetWeekday) - int(now.Weekday()) + 7) % 7
		if daysUntil == 0 && nextTime.Before(now) {
			daysUntil = 7
		}
		nextTime = nextTime.AddDate(0, 0, daysUntil)

	case constants.RepeatMonthly:
		if reminder.DayOfMonth == nil {
			return errors.New("monthly reminder requires dayOfMonth")
		}
		day := int(*reminder.DayOfMonth)

		daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
		if day > daysInMonth {
			day = daysInMonth
		}

		nextTime = time.Date(now.Year(), now.Month(), day, t.Hour(), t.Minute(), 0, 0, time.Local)
		if nextTime.Before(now) {
			nextMonth := now.AddDate(0, 1, 0)
			daysInNextMonth := time.Date(nextMonth.Year(), nextMonth.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
			if day > daysInNextMonth {
				day = daysInNextMonth
			}
			nextTime = time.Date(nextMonth.Year(), nextMonth.Month(), day, t.Hour(), t.Minute(), 0, 0, time.Local)
		}

	default:
		return fmt.Errorf("unsupported repeat type: %s", reminder.Repeat)
	}

	reminder.NextTriggerTime = nextTime
	return nil
}

func (s *ReminderService) checkReminders(ch chan error) {
	defer close(ch)
	var reminders []models.Reminder
	err := s.db.
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

	var updatedReminders []models.Reminder
	for _, r := range reminders {
		if err := s.calculateNextTriggerTime(&r); err != nil {
			fmt.Println("failed to recalc nextTriggerTime:", err)
			continue
		}
		updatedReminders = append(updatedReminders, r)
	}

	var result *gorm.DB
	if len(updatedReminders) != 0 {
		result = s.db.Save(&updatedReminders)
	}

	if result != nil {
		ch <- result.Error
	} else {
		ch <- nil
	}
}

func (s *ReminderService) sendNotifications(plantID int64) {
	var plant models.Plant
	s.db.Where("id = ?", plantID).First(&plant)
	var user models.User
	s.db.Where("id = ?", plant.UserID).First(&user)
	if user.PushToken != "" {
		utils.SendMessage(user.PushToken, plant.Name)
	}
}
