package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/SteelyBretty/consultant-time-tracker/internal/middleware"
	"github.com/SteelyBretty/consultant-time-tracker/internal/models"
	"github.com/SteelyBretty/consultant-time-tracker/internal/schemas"
	"github.com/SteelyBretty/consultant-time-tracker/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TimeEntryHandler struct {
	timeEntryService *services.TimeEntryService
	projectService   *services.ProjectService
}

func NewTimeEntryHandler() *TimeEntryHandler {
	return &TimeEntryHandler{
		timeEntryService: services.NewTimeEntryService(),
		projectService:   services.NewProjectService(),
	}
}

func (h *TimeEntryHandler) CreateTimeEntry(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	var req schemas.CreateTimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	date, _ := time.Parse("2006-01-02", req.Date)

	timeEntry, err := h.timeEntryService.CreateTimeEntry(
		userID, req.ProjectID, date, req.Hours, req.Description, req.IsBillable,
	)
	if err != nil {
		switch err {
		case services.ErrTimeEntryExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Time entry already exists for this date"})
		case services.ErrDateInFuture:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot log time for future dates"})
		default:
			if err.Error() == "project not found or access denied" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create time entry"})
			}
		}
		return
	}

	c.JSON(http.StatusCreated, h.mapTimeEntryToResponse(timeEntry))
}

func (h *TimeEntryHandler) GetTimeEntry(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	timeEntryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time entry ID"})
		return
	}

	timeEntry, err := h.timeEntryService.GetTimeEntry(userID, timeEntryID)
	if err != nil {
		if err == services.ErrTimeEntryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Time entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch time entry"})
		return
	}

	c.JSON(http.StatusOK, h.mapTimeEntryToResponse(timeEntry))
}

func (h *TimeEntryHandler) ListTimeEntries(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var projectID *uuid.UUID
	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		if id, err := uuid.Parse(projectIDStr); err == nil {
			projectID = &id
		}
	}

	var startDate, endDate *time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		if date, err := time.Parse("2006-01-02", startStr); err == nil {
			startDate = &date
		}
	}
	if endStr := c.Query("end_date"); endStr != "" {
		if date, err := time.Parse("2006-01-02", endStr); err == nil {
			endDate = &date
		}
	}

	var isBillable *bool
	if billableStr := c.Query("is_billable"); billableStr != "" {
		billable := billableStr == "true"
		isBillable = &billable
	}

	if limit > 100 {
		limit = 100
	}

	timeEntries, total, err := h.timeEntryService.ListTimeEntries(userID, projectID, startDate, endDate, isBillable, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch time entries"})
		return
	}

	response := schemas.TimeEntryListResponse{
		TimeEntries: make([]schemas.TimeEntryResponse, len(timeEntries)),
		Total:       total,
		Offset:      offset,
		Limit:       limit,
	}

	for i, entry := range timeEntries {
		response.TimeEntries[i] = *h.mapTimeEntryToResponse(entry)
	}

	c.JSON(http.StatusOK, response)
}

func (h *TimeEntryHandler) GetDayEntries(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date parameter is required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	entries, totalHours, err := h.timeEntryService.GetDayEntries(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch day entries"})
		return
	}

	response := schemas.DayEntriesResponse{
		Date:        dateStr,
		TotalHours:  totalHours,
		TimeEntries: make([]schemas.TimeEntryResponse, len(entries)),
	}

	for i, entry := range entries {
		response.TimeEntries[i] = *h.mapTimeEntryToResponse(entry)
	}

	c.JSON(http.StatusOK, response)
}

func (h *TimeEntryHandler) GetWeekEntries(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	weekStr := c.Query("week")
	if weekStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Week parameter is required"})
		return
	}

	week, err := time.Parse("2006-01-02", weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week format, use YYYY-MM-DD"})
		return
	}

	entries, dailyTotals, err := h.timeEntryService.GetWeekEntries(userID, week)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch week entries"})
		return
	}

	weekTotal := float64(0)
	for _, total := range dailyTotals {
		weekTotal += total
	}

	response := schemas.WeekEntriesResponse{
		WeekStarting: week.Format("2006-01-02"),
		TimeEntries:  make([]schemas.TimeEntryResponse, len(entries)),
		DailyTotals:  dailyTotals,
		WeekTotal:    weekTotal,
	}

	for i, entry := range entries {
		response.TimeEntries[i] = *h.mapTimeEntryToResponse(entry)
	}

	c.JSON(http.StatusOK, response)
}

func (h *TimeEntryHandler) GetProjectWeekComparison(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	weekStr := c.Query("week")
	if weekStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Week parameter is required"})
		return
	}

	week, err := time.Parse("2006-01-02", weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week format"})
		return
	}

	allocatedHours, actualHours, err := h.timeEntryService.GetProjectWeekComparison(userID, projectID, week)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comparison"})
		return
	}

	variance := actualHours - allocatedHours
	variancePercent := float64(0)
	if allocatedHours > 0 {
		variancePercent = (variance / allocatedHours) * 100
	}

	response := schemas.ProjectWeekComparisonResponse{
		ProjectID:       projectID,
		WeekStarting:    week.Format("2006-01-02"),
		AllocatedHours:  allocatedHours,
		ActualHours:     actualHours,
		Variance:        variance,
		VariancePercent: variancePercent,
	}

	c.JSON(http.StatusOK, response)
}

func (h *TimeEntryHandler) GetWeekSummary(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	weekStr := c.Query("week")
	if weekStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Week parameter is required"})
		return
	}

	week, err := time.Parse("2006-01-02", weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week format"})
		return
	}

	summary, err := h.timeEntryService.GetWeekSummary(userID, week)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get week summary"})
		return
	}

	response := schemas.WeekSummaryResponse{
		WeekStarting: week.Format("2006-01-02"),
		Projects:     []schemas.ProjectWeekSummary{},
	}

	for projectID, hours := range summary {
		project, _ := h.projectService.GetProjectByID(userID, projectID)

		projectSummary := schemas.ProjectWeekSummary{
			ProjectID:      projectID,
			AllocatedHours: hours["allocated"],
			ActualHours:    hours["actual"],
			Variance:       hours["actual"] - hours["allocated"],
		}

		if project != nil {
			projectSummary.Project = &schemas.ProjectSummary{
				ID:           project.ID,
				Name:         project.Name,
				Code:         project.Code,
				BillableRate: project.BillableRate,
				Currency:     project.Currency,
				Client: schemas.ClientSummary{
					ID:   project.Client.ID,
					Name: project.Client.Name,
					Code: project.Client.Code,
				},
			}
		}

		response.Projects = append(response.Projects, projectSummary)
		response.TotalAllocated += hours["allocated"]
		response.TotalActual += hours["actual"]
	}

	c.JSON(http.StatusOK, response)
}

func (h *TimeEntryHandler) UpdateTimeEntry(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	timeEntryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time entry ID"})
		return
	}

	var req schemas.UpdateTimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	timeEntry, err := h.timeEntryService.UpdateTimeEntry(userID, timeEntryID, req.Hours, req.Description, req.IsBillable)
	if err != nil {
		if err == services.ErrTimeEntryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Time entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update time entry"})
		return
	}

	c.JSON(http.StatusOK, h.mapTimeEntryToResponse(timeEntry))
}

func (h *TimeEntryHandler) DeleteTimeEntry(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	timeEntryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time entry ID"})
		return
	}

	if err := h.timeEntryService.DeleteTimeEntry(userID, timeEntryID); err != nil {
		if err == services.ErrTimeEntryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Time entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete time entry"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *TimeEntryHandler) mapTimeEntryToResponse(entry *models.TimeEntry) *schemas.TimeEntryResponse {
	response := &schemas.TimeEntryResponse{
		ID:          entry.ID,
		ProjectID:   entry.ProjectID,
		Date:        time.Time(entry.Date).Format("2006-01-02"),
		Hours:       entry.Hours,
		Description: entry.Description,
		IsBillable:  entry.IsBillable,
		CreatedAt:   entry.CreatedAt,
		UpdatedAt:   entry.UpdatedAt,
	}

	if entry.Project.ID != uuid.Nil {
		response.Project = &schemas.ProjectSummary{
			ID:           entry.Project.ID,
			Name:         entry.Project.Name,
			Code:         entry.Project.Code,
			BillableRate: entry.Project.BillableRate,
			Currency:     entry.Project.Currency,
			Client: schemas.ClientSummary{
				ID:   entry.Project.Client.ID,
				Name: entry.Project.Client.Name,
				Code: entry.Project.Client.Code,
			},
		}
	}

	return response
}
