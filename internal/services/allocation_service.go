package services

import (
	"errors"
	"time"

	"github.com/SteelyBretty/consultant-time-tracker/internal/database"
	"github.com/SteelyBretty/consultant-time-tracker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AllocationService struct{}

func NewAllocationService() *AllocationService {
	return &AllocationService{}
}

var (
	ErrAllocationNotFound = errors.New("allocation not found")
	ErrAllocationExists   = errors.New("allocation already exists for this week")
	ErrInvalidWeekStart   = errors.New("week must start on Monday")
	ErrProjectNotActive   = errors.New("project is not active")
)

func (s *AllocationService) getWeekStart(date time.Time) time.Time {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return date.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
}

func (s *AllocationService) CreateAllocation(userID, projectID uuid.UUID, weekStarting time.Time, hours float64, notes string) (*models.Allocation, error) {
	weekStart := s.getWeekStart(weekStarting)
	if !weekStart.Equal(weekStarting) {
		return nil, ErrInvalidWeekStart
	}

	var project models.Project
	if err := database.DB.Where("id = ? AND user_id = ?", projectID, userID).First(&project).Error; err != nil {
		return nil, errors.New("project not found or access denied")
	}

	if !project.IsActive || project.Status != models.ProjectStatusActive {
		return nil, ErrProjectNotActive
	}

	var existing models.Allocation
	err := database.DB.Where("project_id = ? AND user_id = ? AND week_starting = ?",
		projectID, userID, weekStart).First(&existing).Error
	if err == nil {
		return nil, ErrAllocationExists
	}

	allocation := &models.Allocation{
		ProjectID:    projectID,
		UserID:       userID,
		WeekStarting: weekStart,
		Hours:        hours,
		Notes:        notes,
	}

	if err := database.DB.Create(allocation).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Project.Client").First(allocation, allocation.ID).Error; err != nil {
		return nil, err
	}

	return allocation, nil
}

func (s *AllocationService) GetAllocation(userID, allocationID uuid.UUID) (*models.Allocation, error) {
	var allocation models.Allocation
	err := database.DB.Preload("Project.Client").Where("id = ? AND user_id = ?", allocationID, userID).First(&allocation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAllocationNotFound
		}
		return nil, err
	}
	return &allocation, nil
}

func (s *AllocationService) ListAllocations(userID uuid.UUID, projectID *uuid.UUID, startDate, endDate *time.Time, offset, limit int) ([]*models.Allocation, int64, error) {
	var allocations []*models.Allocation
	var total int64

	query := database.DB.Model(&models.Allocation{}).Preload("Project.Client").Where("user_id = ?", userID)

	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	}

	if startDate != nil {
		weekStart := s.getWeekStart(*startDate)
		query = query.Where("week_starting >= ?", weekStart)
	}

	if endDate != nil {
		weekEnd := s.getWeekStart(*endDate)
		query = query.Where("week_starting <= ?", weekEnd)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("week_starting DESC, created_at DESC").Find(&allocations).Error; err != nil {
		return nil, 0, err
	}

	return allocations, total, nil
}

func (s *AllocationService) GetWeekAllocations(userID uuid.UUID, weekStarting time.Time) ([]*models.Allocation, float64, error) {
	weekStart := s.getWeekStart(weekStarting)

	var allocations []*models.Allocation
	err := database.DB.Preload("Project.Client").
		Where("user_id = ? AND week_starting = ?", userID, weekStart).
		Order("created_at ASC").
		Find(&allocations).Error
	if err != nil {
		return nil, 0, err
	}

	var totalHours float64
	for _, allocation := range allocations {
		totalHours += allocation.Hours
	}

	return allocations, totalHours, nil
}

func (s *AllocationService) UpdateAllocation(userID, allocationID uuid.UUID, hours float64, notes string) (*models.Allocation, error) {
	var allocation models.Allocation

	if err := database.DB.Where("id = ? AND user_id = ?", allocationID, userID).First(&allocation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAllocationNotFound
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"hours": hours,
		"notes": notes,
	}

	if err := database.DB.Model(&allocation).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Project.Client").First(&allocation, allocation.ID).Error; err != nil {
		return nil, err
	}

	return &allocation, nil
}

func (s *AllocationService) DeleteAllocation(userID, allocationID uuid.UUID) error {
	result := database.DB.Where("id = ? AND user_id = ?", allocationID, userID).Delete(&models.Allocation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAllocationNotFound
	}
	return nil
}

func (s *AllocationService) CopyWeekAllocations(userID uuid.UUID, fromWeek, toWeek time.Time) ([]*models.Allocation, error) {
	fromWeekStart := s.getWeekStart(fromWeek)
	toWeekStart := s.getWeekStart(toWeek)

	if !toWeekStart.Equal(toWeek) {
		return nil, ErrInvalidWeekStart
	}

	var sourceAllocations []*models.Allocation
	err := database.DB.Where("user_id = ? AND week_starting = ?", userID, fromWeekStart).Find(&sourceAllocations).Error
	if err != nil {
		return nil, err
	}

	var newAllocations []*models.Allocation
	for _, source := range sourceAllocations {
		var existing models.Allocation
		err := database.DB.Where("project_id = ? AND user_id = ? AND week_starting = ?",
			source.ProjectID, userID, toWeekStart).First(&existing).Error
		if err == nil {
			continue
		}

		newAllocation := &models.Allocation{
			ProjectID:    source.ProjectID,
			UserID:       userID,
			WeekStarting: toWeekStart,
			Hours:        source.Hours,
			Notes:        "Copied from week of " + fromWeekStart.Format("2006-01-02"),
		}

		if err := database.DB.Create(newAllocation).Error; err == nil {
			newAllocations = append(newAllocations, newAllocation)
		}
	}

	for _, allocation := range newAllocations {
		database.DB.Preload("Project.Client").First(allocation, allocation.ID)
	}

	return newAllocations, nil
}
