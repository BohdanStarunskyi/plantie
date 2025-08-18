package models

import (
	"errors"
	"plant-reminder/config"
)

type PlantIcon string

const (
	BananaPlant  PlantIcon = "bananaPlant"
	BigCactus    PlantIcon = "bigCactus"
	BigPlant     PlantIcon = "bigPlant"
	BigRose      PlantIcon = "bigRose"
	ChilliPlant  PlantIcon = "chilliPlant"
	Daisy        PlantIcon = "daisy"
	FlowerBed    PlantIcon = "flowerBed"
	Flower       PlantIcon = "flower"
	LeafyPlant   PlantIcon = "leafyPlant"
	MediumPlant  PlantIcon = "mediumPlant"
	RedTulip     PlantIcon = "redTulip"
	SeaweedPlant PlantIcon = "seaweedPlant"
	ShortPlant   PlantIcon = "shortPlant"
	SkinnyPlant  PlantIcon = "skinnyPlant"
	SmallCactus  PlantIcon = "smallCactus"
	SmallPlant   PlantIcon = "smallPlant"
	SmallRose    PlantIcon = "smallRose"
	SpikyPlant   PlantIcon = "spikyPlant"
	TallPlant    PlantIcon = "tallPlant"
	ThreeFlowers PlantIcon = "threeFlowers"
	TwoFlowers   PlantIcon = "twoFlowers"
	TwoPlants    PlantIcon = "twoPlants"
	WhiteFlower  PlantIcon = "whiteFlower"
	YellowTulip  PlantIcon = "yellowTulip"
)

var validPlantIcons = []PlantIcon{
	BananaPlant, BigCactus, BigPlant, BigRose, ChilliPlant, Daisy,
	FlowerBed, Flower, LeafyPlant, MediumPlant, RedTulip, SeaweedPlant,
	ShortPlant, SkinnyPlant, SmallCactus, SmallPlant, SmallRose, SpikyPlant,
	TallPlant, ThreeFlowers, TwoFlowers, TwoPlants, WhiteFlower, YellowTulip,
}

func (pi PlantIcon) IsValid() bool {
	for _, valid := range validPlantIcons {
		if pi == valid {
			return true
		}
	}
	return false
}

type Plant struct {
	ID        int64      `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name" validate:"required"`
	Note      string     `json:"note"`
	TagColor  string     `json:"tagColor" validate:"required"`
	UserID    int64      `json:"-"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-" validate:"-"`
	Reminders []Reminder `gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE" json:"reminders,omitempty"`
	PlantIcon PlantIcon  `json:"plantIcon" validate:"required"`
}

func (p *Plant) Validate() error {
	if p.UserID == 0 {
		return errors.New("user ID must be set")
	}
	if !p.PlantIcon.IsValid() {
		return errors.New("invalid PlantIcon value")
	}
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.TagColor == "" {
		return errors.New("tagColor is required")
	}
	return nil
}

func (p *Plant) Save() error {
	if err := p.Validate(); err != nil {
		return err
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

	if err := p.Validate(); err != nil {
		return err
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
