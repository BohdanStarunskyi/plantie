package models

import (
	"errors"
	"plant-reminder/config"
	"plant-reminder/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	Email        string    `binding:"required" gorm:"unique" json:"email"`
	Password     string    `binding:"required" json:"-"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"createdAt"`
	PushToken    string    `binding:"-" json:"-"`
}

func userExists(email string) (bool, error) {
	var user User
	result := config.DB.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.Error == nil, result.Error
}

func (u *User) CreateUser() (string, error) {
	exists, _ := userExists(u.Email)
	if exists {
		return "", errors.New("user with such email already exists")
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return "", err
	}
	u.CreationDate = time.Now()
	u.Password = hashedPassword

	result := config.DB.Create(&u)
	if result.Error != nil {
		return "", errors.New("error while writing to database")
	}
	token, err := utils.SignPayload(u.ID)

	return token, err
}

func VerifyUser(email string, password string) (string, *User, error) {
	var user User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return "", nil, result.Error
	}
	if err := utils.CheckPassword(user.Password, password); err != nil {
		return "", nil, errors.New("wrong credentials")
	}

	token, err := utils.SignPayload(user.ID)
	return token, &user, err
}

func SetPushToken(userID string, token string) error {
	result := config.DB.Model(&User{}).
		Where("id = ?", userID).
		Update("push_token", token)
	return result.Error
}
