package models

import (
	"time"
)

type User struct {
	ID           int64  `gorm:"primaryKey"`
	Email        string `validate:"required,email"`
	Password     string `validate:"required,min=6"`
	Name         string `validate:"omitempty,min=2,max=100"`
	CreationDate time.Time
	PushToken    string
	Plants       []Plant `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
