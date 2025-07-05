package services

import (
	"errors"
	"strings"

	"github.com/SteelyBretty/consultant-time-tracker/internal/database"
	"github.com/SteelyBretty/consultant-time-tracker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientService struct{}

func NewClientService() *ClientService {
	return &ClientService{}
}

var (
	ErrClientNotFound   = errors.New("client not found")
	ErrClientCodeExists = errors.New("client code already exists")
)

func (s *ClientService) CreateClient(userID uuid.UUID, name, code, email, phone, address string) (*models.Client, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	var existing models.Client
	if err := database.DB.Where("code = ? AND user_id = ?", code, userID).First(&existing).Error; err == nil {
		return nil, ErrClientCodeExists
	}

	client := &models.Client{
		Name:     name,
		Code:     code,
		Email:    email,
		Phone:    phone,
		Address:  address,
		IsActive: true,
		UserID:   userID,
	}

	if err := database.DB.Create(client).Error; err != nil {
		return nil, err
	}

	return client, nil
}

func (s *ClientService) GetClientByID(userID, clientID uuid.UUID) (*models.Client, error) {
	var client models.Client
	err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (s *ClientService) GetClientByCode(userID uuid.UUID, code string) (*models.Client, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	var client models.Client
	err := database.DB.Where("code = ? AND user_id = ?", code, userID).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (s *ClientService) ListClients(userID uuid.UUID, isActive *bool, search string, offset, limit int) ([]*models.Client, int64, error) {
	var clients []*models.Client
	var total int64

	query := database.DB.Model(&models.Client{}).Where("user_id = ?", userID)

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR email LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("name ASC").Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

func (s *ClientService) UpdateClient(userID, clientID uuid.UUID, updates map[string]interface{}) (*models.Client, error) {
	var client models.Client

	if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	if code, ok := updates["code"].(string); ok {
		code = strings.ToUpper(strings.TrimSpace(code))
		updates["code"] = code

		var existing models.Client
		if err := database.DB.Where("code = ? AND user_id = ? AND id != ?", code, userID, clientID).First(&existing).Error; err == nil {
			return nil, ErrClientCodeExists
		}
	}

	if err := database.DB.Model(&client).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &client, nil
}

func (s *ClientService) DeleteClient(userID, clientID uuid.UUID) error {
	result := database.DB.Where("id = ? AND user_id = ?", clientID, userID).Delete(&models.Client{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrClientNotFound
	}
	return nil
}

func (s *ClientService) GetClientWithProjects(userID, clientID uuid.UUID) (*models.Client, error) {
	var client models.Client
	err := database.DB.Preload("Projects").Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}
