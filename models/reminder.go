package models

import (
	"errors"
	"fmt"
	"plant-reminder/config"
	"plant-reminder/constants"
	"plant-reminder/utils"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

var scheduler *gocron.Scheduler

type Reminder struct {
	ID              int64                `gorm:"primaryKey" json:"id"`
	PlantID         int64                `json:"-"`
	Repeat          constants.RepeatType `gorm:"type:smallint" json:"repeatType"`
	TimeOfDay       string               `json:"timeOfDay" validate:"required,len=5"`
	NextTriggerTime time.Time            `json:"nextTriggerTime"`
	UserID          int64                `json:"-"`
	Plant           *Plant               `gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE" json:"plant,omitempty"`
}

func (r *Reminder) Save() error {
	plant, err := GetPlant(r.PlantID, r.UserID)
	if err != nil {
		return errors.New("plant doesn't exist")
	}
	if plant.UserID != r.UserID {
		return errors.New("not enough rights")
	}

	var existing Reminder
	err = config.DB.
		Where("plant_id = ? AND time_of_day = ? AND repeat = ?", r.PlantID, r.TimeOfDay, r.Repeat).
		First(&existing).Error

	if err == nil && existing.ID != r.ID {
		return errors.New("reminder with same time and repeat period already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing reminders: %w", err)
	}

	now := time.Now()
	t, err := time.Parse("15:04", r.TimeOfDay)
	if err != nil {
		return fmt.Errorf("invalid time format, expected HH:mm: %w", err)
	}

	nextTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
	if nextTime.Before(now) {
		switch r.Repeat {
		case constants.RepeatDaily:
			nextTime = nextTime.Add(24 * time.Hour)
		case constants.RepeatWeekly:
			nextTime = nextTime.Add(7 * 24 * time.Hour)
		case constants.RepeatMonthly:
			nextTime = nextTime.AddDate(0, 1, 0)
		}
	}

	r.NextTriggerTime = nextTime

	result := config.DB.Save(&r)
	return result.Error
}

func (r *Reminder) Update(userID int64) error {
	if r.ID == 0 {
		return errors.New("reminder ID must be set")
	}
	now := time.Now()
	t, err := time.Parse("15:04", r.TimeOfDay)
	if err != nil {
		return fmt.Errorf("invalid time format, expected HH:mm: %w", err)
	}
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)

	reminder, err := getReminder(r.ID, userID, r.PlantID)
	if err != nil {
		return err
	}

	if reminder.UserID != userID {
		return errors.New("not enough rights")
	}

	if nextTime.Before(now) {
		switch r.Repeat {
		case constants.RepeatDaily:
			nextTime = nextTime.Add(24 * time.Hour)
		case constants.RepeatWeekly:
			nextTime = nextTime.Add(7 * 24 * time.Hour)
		case constants.RepeatMonthly:
			nextTime = nextTime.AddDate(0, 1, 0)
		}
	}

	r.NextTriggerTime = nextTime

	result := config.DB.Model(&reminder).Updates(r)
	return result.Error
}

func getReminder(reminderID int64, userID int64, plantId int64) (Reminder, error) {
	var reminder Reminder
	if reminderID == 0 {
		return Reminder{}, errors.New("reminderID must be set")
	}
	result := config.DB.Where("id = ? AND user_id = ? AND plant_id = ?", reminderID, userID, plantId).First(&reminder)
	return reminder, result.Error
}

func (r Reminder) Delete() error {
	reminder, err := getReminder(r.ID, r.UserID, r.PlantID)
	if err != nil {
		return err
	}
	result := config.DB.Delete(&reminder)
	return result.Error
}

func GetPlantReminders(userID int64, plantID int64) ([]Reminder, error) {
	var reminders []Reminder
	if userID == 0 {
		return nil, errors.New("userID must be set")
	}
	if plantID == 0 {
		return nil, errors.New("plantID must be set")
	}
	result := config.DB.Where("plant_id = ? AND user_id = ?", plantID, userID).Find(&reminders)
	return reminders, result.Error
}

func GetAllReminders(userID int64) ([]Reminder, error) {
	var reminders []Reminder
	if userID == 0 {
		return nil, errors.New("userID must be set")
	}
	result := config.DB.Preload("Plant").Where("user_id = ?", userID).Find(&reminders)
	return reminders, result.Error
}

func SetReminders() error {
	if scheduler != nil && scheduler.IsRunning() {
		return nil
	}
	scheduler = gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(1).Minutes().Do(func() {
		errChan := make(chan error)
		go checkReminders(errChan)
		for err := range errChan {
			if err != nil {
				fmt.Println("error during checking reinders", err)
			}
		}

	})
	if err != nil {
		return err
	}
	scheduler.StartAsync()
	return nil
}

func checkReminders(ch chan error) {
	defer close(ch)
	var reminders []Reminder
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
			sendNotifications(pid)
		}(value.PlantID)
	}
	wg.Wait()

	updatedReminders := utils.Map(reminders, func(item Reminder) Reminder {
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

func sendNotifications(plantID int64) {
	var plant Plant
	config.DB.Where("id = ?", plantID).First(&plant)
	var user User
	config.DB.Where("id = ?", plant.UserID).First(&user)
	if user.PushToken != "" {
		utils.SendMessage(user.PushToken, plant.Name)
	}
}
