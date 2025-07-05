package services

import (
	"errors"
	"time"

	"github.com/SteelyBretty/consultant-time-tracker/internal/database"
	"github.com/SteelyBretty/consultant-time-tracker/internal/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TimeEntryService struct{}

func NewTimeEntryService() *TimeEntryService {
	return &TimeEntryService{}
}

var (
	ErrTimeEntryNotFound = errors.New("time entry not found")
	ErrTimeEntryExists   = errors.New("time entry already exists for this date")
	ErrExceedsAllocation = errors.New("time entry exceeds weekly allocation")
	ErrNoAllocation      = errors.New("no allocation found for this week")
	ErrDateInFuture      = errors.New("cannot log time for future dates")
)

func (s *TimeEntryService) getWeekStart(date time.Time) time.Time {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return date.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
}

func (s *TimeEntryService) CreateTimeEntry(userID, projectID uuid.UUID, date time.Time, hours float64, description string, isBillable bool) (*models.TimeEntry, error) {
	if date.After(time.Now()) {
		return nil, ErrDateInFuture
	}

	var project models.Project
	if err := database.DB.Where("id = ? AND user_id = ?", projectID, userID).First(&project).Error; err != nil {
		return nil, errors.New("project not found or access denied")
	}

	var existing models.TimeEntry
	err := database.DB.Where("project_id = ? AND user_id = ? AND date = ?",
		projectID, userID, datatypes.Date(date)).First(&existing).Error
	if err == nil {
		return nil, ErrTimeEntryExists
	}

	timeEntry := &models.TimeEntry{
		ProjectID:   projectID,
		UserID:      userID,
		Date:        datatypes.Date(date),
		Hours:       hours,
		Description: description,
		IsBillable:  isBillable,
	}

	if err := database.DB.Create(timeEntry).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Project.Client").First(timeEntry, timeEntry.ID).Error; err != nil {
		return nil, err
	}

	return timeEntry, nil
}

func (s *TimeEntryService) GetTimeEntry(userID, timeEntryID uuid.UUID) (*models.TimeEntry, error) {
	var timeEntry models.TimeEntry
	err := database.DB.Preload("Project.Client").Where("id = ? AND user_id = ?", timeEntryID, userID).First(&timeEntry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTimeEntryNotFound
		}
		return nil, err
	}
	return &timeEntry, nil
}

func (s *TimeEntryService) ListTimeEntries(userID uuid.UUID, projectID *uuid.UUID, startDate, endDate *time.Time, isBillable *bool, offset, limit int) ([]*models.TimeEntry, int64, error) {
	var timeEntries []*models.TimeEntry
	var total int64

	query := database.DB.Model(&models.TimeEntry{}).Preload("Project.Client").Where("user_id = ?", userID)

	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	}

	if startDate != nil {
		query = query.Where("date >= ?", datatypes.Date(*startDate))
	}

	if endDate != nil {
		query = query.Where("date <= ?", datatypes.Date(*endDate))
	}

	if isBillable != nil {
		query = query.Where("is_billable = ?", *isBillable)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("date DESC, created_at DESC").Find(&timeEntries).Error; err != nil {
		return nil, 0, err
	}

	return timeEntries, total, nil
}

func (s *TimeEntryService) GetDayEntries(userID uuid.UUID, date time.Time) ([]*models.TimeEntry, float64, error) {
	var entries []*models.TimeEntry
	err := database.DB.Preload("Project.Client").
		Where("user_id = ? AND date = ?", userID, datatypes.Date(date)).
		Order("created_at ASC").
		Find(&entries).Error
	if err != nil {
		return nil, 0, err
	}

	var totalHours float64
	for _, entry := range entries {
		totalHours += entry.Hours
	}

	return entries, totalHours, nil
}

func (s *TimeEntryService) GetWeekEntries(userID uuid.UUID, weekStarting time.Time) ([]*models.TimeEntry, map[string]float64, error) {
	weekStart := s.getWeekStart(weekStarting)
	weekEnd := weekStart.AddDate(0, 0, 6)

	var entries []*models.TimeEntry
	err := database.DB.Preload("Project.Client").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, datatypes.Date(weekStart), datatypes.Date(weekEnd)).
		Order("date ASC, created_at ASC").
		Find(&entries).Error
	if err != nil {
		return nil, nil, err
	}

	dailyTotals := make(map[string]float64)
	for _, entry := range entries {
		dateStr := time.Time(entry.Date).Format("2006-01-02")
		dailyTotals[dateStr] += entry.Hours
	}

	return entries, dailyTotals, nil
}

func (s *TimeEntryService) UpdateTimeEntry(userID, timeEntryID uuid.UUID, hours float64, description string, isBillable bool) (*models.TimeEntry, error) {
	var timeEntry models.TimeEntry

	if err := database.DB.Where("id = ? AND user_id = ?", timeEntryID, userID).First(&timeEntry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTimeEntryNotFound
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"hours":       hours,
		"description": description,
		"is_billable": isBillable,
	}

	if err := database.DB.Model(&timeEntry).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Project.Client").First(&timeEntry, timeEntry.ID).Error; err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

func (s *TimeEntryService) DeleteTimeEntry(userID, timeEntryID uuid.UUID) error {
	result := database.DB.Where("id = ? AND user_id = ?", timeEntryID, userID).Delete(&models.TimeEntry{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTimeEntryNotFound
	}
	return nil
}

func (s *TimeEntryService) GetProjectWeekComparison(userID, projectID uuid.UUID, weekStarting time.Time) (float64, float64, error) {
	weekStart := s.getWeekStart(weekStarting)
	weekEnd := weekStart.AddDate(0, 0, 6)

	var allocation models.Allocation
	err := database.DB.Where("project_id = ? AND user_id = ? AND week_starting = ?",
		projectID, userID, weekStart).First(&allocation).Error

	allocatedHours := float64(0)
	if err == nil {
		allocatedHours = allocation.Hours
	}

	var actualHours float64
	err = database.DB.Model(&models.TimeEntry{}).
		Where("project_id = ? AND user_id = ? AND date >= ? AND date <= ?",
			projectID, userID, datatypes.Date(weekStart), datatypes.Date(weekEnd)).
		Select("COALESCE(SUM(hours), 0)").
		Scan(&actualHours).Error

	return allocatedHours, actualHours, err
}

func (s *TimeEntryService) GetWeekSummary(userID uuid.UUID, weekStarting time.Time) (map[uuid.UUID]map[string]float64, error) {
	weekStart := s.getWeekStart(weekStarting)
	weekEnd := weekStart.AddDate(0, 0, 6)

	type ProjectWeekData struct {
		ProjectID      uuid.UUID
		AllocatedHours float64
		ActualHours    float64
	}

	var allocations []models.Allocation
	database.DB.Where("user_id = ? AND week_starting = ?", userID, weekStart).Find(&allocations)

	summary := make(map[uuid.UUID]map[string]float64)

	for _, alloc := range allocations {
		summary[alloc.ProjectID] = map[string]float64{
			"allocated": alloc.Hours,
			"actual":    0,
		}
	}

	var entries []models.TimeEntry
	database.DB.Where("user_id = ? AND date >= ? AND date <= ?",
		userID, datatypes.Date(weekStart), datatypes.Date(weekEnd)).Find(&entries)

	for _, entry := range entries {
		if _, exists := summary[entry.ProjectID]; !exists {
			summary[entry.ProjectID] = map[string]float64{
				"allocated": 0,
				"actual":    0,
			}
		}
		summary[entry.ProjectID]["actual"] += entry.Hours
	}

	return summary, nil
}
