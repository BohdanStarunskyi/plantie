package service

import (
	"errors"
	"plant-reminder/dto"
	"plant-reminder/models"
	"plant-reminder/utils"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

type UserServiceInterface interface {
	CreateUser(userRequest *dto.UserCreateRequest) (*dto.AuthResponse, error)
	VerifyUser(email, password string) (*dto.AuthResponse, error)
	SetPushToken(userID, token string) error
	DeleteUser(userID int64) error
	GetUser(userID int64) (*dto.UserResponse, error)
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(userRequest *dto.UserCreateRequest) (*dto.AuthResponse, error) {
	exists, _ := s.userExists(userRequest.Email)
	if exists {
		return nil, errors.New("user with such email already exists")
	}

	user := userRequest.ToModel()
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.CreationDate = time.Now()
	user.Password = hashedPassword

	result := s.db.Create(user)
	if result.Error != nil {
		return nil, errors.New("error while writing to database")
	}

	accessToken, err := utils.SignPayload(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.SignRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	userResponse := (&dto.UserResponse{}).FromModel(user)
	authResponse := &dto.AuthResponse{
		User:         *userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authResponse, nil
}

func (s *UserService) VerifyUser(email string, password string) (*dto.AuthResponse, error) {
	var user models.User
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if err := utils.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("wrong credentials")
	}

	accessToken, err := utils.SignPayload(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.SignRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	userResponse := (&dto.UserResponse{}).FromModel(&user)
	authResponse := &dto.AuthResponse{
		User:         *userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authResponse, nil
}

func (s *UserService) GetUser(id int64) (*dto.UserResponse, error) {
	var user models.User
	result := s.db.Where("id = ?", id).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	user.Password = ""

	userResponse := (&dto.UserResponse{}).FromModel(&user)
	return userResponse, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	if user.ID == 0 {
		return errors.New("user ID must be set")
	}

	updates := map[string]interface{}{}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.Name != "" {
		updates["name"] = user.Name
	}
	if user.PushToken != "" {
		updates["push_token"] = user.PushToken
	}

	result := s.db.Model(user).Updates(updates)
	return result.Error
}

func (s *UserService) SetPushToken(userID string, token string) error {
	result := s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("push_token", token)
	return result.Error
}

func (s *UserService) DeleteUser(userID int64) error {
	var user models.User
	result := s.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return result.Error
	}

	result = s.db.Delete(&user)
	return result.Error
}

func (s *UserService) userExists(email string) (bool, error) {
	var user models.User
	result := s.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.Error == nil, result.Error
}
