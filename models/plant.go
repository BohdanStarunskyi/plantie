package models

import (
	"errors"
	"plant-reminder/config"
)

type Plant struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Name     string `json:"name" validate:"required"`
	Note     string `json:"note"`
	TagColor string `json:"tagColor" validate:"required"`
	UserID   int64  `json:"-"`
	User     User   `gorm:"foreignKey:UserID" json:"-" validate:"-"`
}

func (p *Plant) Save() error {
	if p.UserID == 0 {
		return errors.New("user ID must be set before saving plant")
	}
	result := config.DB.Save(&p)

	return result.Error
}

func (p *Plant) Update(userID int64) error {
	if p.ID == 0 {
		return errors.New("plant ID must be set")
	}

	plant, err := GetPlant(p.ID, userID)
	if err != nil {
		return err
	}

	if plant.UserID != userID {
		return errors.New("not enough rights")
	}

	result := config.DB.Model(&plant).Updates(p)
	return result.Error
}

func DeletePlant(userID int64, plantID int64) error {
	plant, err := GetPlant(plantID, userID)
	if err != nil {
		return err
	}
	if plant.UserID != userID {
		return errors.New("not enough rights")
	}
	result := config.DB.Delete(&plant)
	return result.Error
}

func GetPlant(plantID int64, userID int64) (Plant, error) {
	var plant Plant
	if plantID == 0 {
		return Plant{}, errors.New("plantID must be set")
	}
	result := config.DB.Where("id = ? AND user_id = ?", plantID, userID).First(&plant)
	return plant, result.Error
}
func GetPlants(userID int64) ([]Plant, error) {
	var plants []Plant
	if userID == 0 {
		return nil, errors.New("userID must be set")
	}
	result := config.DB.Where("user_id = ?", userID).Find(&plants)
	return plants, result.Error
}
