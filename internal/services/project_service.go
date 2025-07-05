package services

import (
	"errors"
	"strings"
	"time"

	"github.com/SteelyBretty/consultant-time-tracker/internal/database"
	"github.com/SteelyBretty/consultant-time-tracker/internal/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProjectService struct{}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectCodeExists    = errors.New("project code already exists")
	ErrInvalidProjectStatus = errors.New("invalid project status")
)

func (s *ProjectService) CreateProject(userID, clientID uuid.UUID, name, code, description string, billableRate float64, currency string, startDate time.Time, endDate *time.Time) (*models.Project, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	var client models.Client
	if err := database.DB.Where("id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
		return nil, errors.New("client not found or access denied")
	}

	var existing models.Project
	if err := database.DB.Where("code = ? AND user_id = ?", code, userID).First(&existing).Error; err == nil {
		return nil, ErrProjectCodeExists
	}

	project := &models.Project{
		Name:         name,
		Code:         code,
		Description:  description,
		Status:       models.ProjectStatusActive,
		BillableRate: billableRate,
		Currency:     currency,
		StartDate:    datatypes.Date(startDate),
		IsActive:     true,
		ClientID:     clientID,
		UserID:       userID,
	}

	if endDate != nil {
		endDateValue := datatypes.Date(*endDate)
		project.EndDate = &endDateValue
	}

	if err := database.DB.Create(project).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Client").First(project, project.ID).Error; err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetProjectByID(userID, projectID uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := database.DB.Preload("Client").Where("id = ? AND user_id = ?", projectID, userID).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return &project, nil
}

func (s *ProjectService) ListProjects(userID uuid.UUID, clientID *uuid.UUID, status *models.ProjectStatus, isActive *bool, search string, offset, limit int) ([]*models.Project, int64, error) {
	var projects []*models.Project
	var total int64

	query := database.DB.Model(&models.Project{}).Preload("Client").Where("user_id = ?", userID)

	if clientID != nil {
		query = query.Where("client_id = ?", *clientID)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("start_date DESC, name ASC").Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (s *ProjectService) UpdateProject(userID, projectID uuid.UUID, updates map[string]interface{}) (*models.Project, error) {
	var project models.Project

	if err := database.DB.Where("id = ? AND user_id = ?", projectID, userID).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	if code, ok := updates["code"].(string); ok {
		code = strings.ToUpper(strings.TrimSpace(code))
		updates["code"] = code

		var existing models.Project
		if err := database.DB.Where("code = ? AND user_id = ? AND id != ?", code, userID, projectID).First(&existing).Error; err == nil {
			return nil, ErrProjectCodeExists
		}
	}

	if status, ok := updates["status"].(string); ok {
		projectStatus := models.ProjectStatus(status)
		if projectStatus != models.ProjectStatusActive &&
			projectStatus != models.ProjectStatusOnHold &&
			projectStatus != models.ProjectStatusCompleted &&
			projectStatus != models.ProjectStatusCancelled {
			return nil, ErrInvalidProjectStatus
		}
		updates["status"] = projectStatus
	}

	if err := database.DB.Model(&project).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Client").First(&project, project.ID).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *ProjectService) DeleteProject(userID, projectID uuid.UUID) error {
	result := database.DB.Where("id = ? AND user_id = ?", projectID, userID).Delete(&models.Project{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrProjectNotFound
	}
	return nil
}

func (s *ProjectService) GetProjectsByClient(userID, clientID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	err := database.DB.Where("client_id = ? AND user_id = ?", clientID, userID).Order("start_date DESC").Find(&projects).Error
	return projects, err
}

func (s *ProjectService) UpdateProjectStatus(userID, projectID uuid.UUID, status models.ProjectStatus) (*models.Project, error) {
	updates := map[string]interface{}{
		"status": status,
	}
	return s.UpdateProject(userID, projectID, updates)
}
