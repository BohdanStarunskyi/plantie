package models

import (
	"plant-reminder/constants"
	"time"
)

type Reminder struct {
	ID              int64 `gorm:"primaryKey"`
	PlantID         int64
	Repeat          constants.RepeatType `gorm:"type:smallint"`
	TimeOfDay       string               `validate:"required,len=5"`
	NextTriggerTime time.Time
	UserID          int64
	Plant           *Plant `gorm:"foreignKey:PlantID;constraint:OnDelete:CASCADE"`
}
