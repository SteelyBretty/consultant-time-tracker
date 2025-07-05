package services

import (
	"errors"

	"github.com/SteelyBretty/consultant-time-tracker/internal/database"
	"github.com/SteelyBretty/consultant-time-tracker/internal/models"
	"github.com/google/uuid"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotActive      = errors.New("user account is not active")
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
)

func (s *AuthService) Register(username, email, password, fullName string) (*models.User, error) {
	var existingUser models.User

	if err := database.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, ErrUsernameExists
	}

	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, ErrEmailExists
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
		FullName: fullName,
		IsActive: true,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (*models.User, error) {
	var user models.User

	if err := database.DB.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) ValidateCredentials(username, password string) (*models.User, error) {
	return s.Login(username, password)
}
