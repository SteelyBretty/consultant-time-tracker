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

type AllocationHandler struct {
	allocationService *services.AllocationService
}

func NewAllocationHandler() *AllocationHandler {
	return &AllocationHandler{
		allocationService: services.NewAllocationService(),
	}
}

func (h *AllocationHandler) CreateAllocation(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	var req schemas.CreateAllocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	weekStarting, _ := time.Parse("2006-01-02", req.WeekStarting)

	allocation, err := h.allocationService.CreateAllocation(
		userID, req.ProjectID, weekStarting, req.Hours, req.Notes,
	)
	if err != nil {
		switch err {
		case services.ErrInvalidWeekStart:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Week must start on Monday"})
		case services.ErrAllocationExists:
			c.JSON(http.StatusConflict, gin.H{"error": "Allocation already exists for this week"})
		case services.ErrProjectNotActive:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Project is not active"})
		default:
			if err.Error() == "project not found or access denied" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create allocation"})
			}
		}
		return
	}

	c.JSON(http.StatusCreated, h.mapAllocationToResponse(allocation))
}

func (h *AllocationHandler) GetAllocation(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	allocationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid allocation ID"})
		return
	}

	allocation, err := h.allocationService.GetAllocation(userID, allocationID)
	if err != nil {
		if err == services.ErrAllocationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Allocation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch allocation"})
		return
	}

	c.JSON(http.StatusOK, h.mapAllocationToResponse(allocation))
}

func (h *AllocationHandler) ListAllocations(c *gin.Context) {
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

	if limit > 100 {
		limit = 100
	}

	allocations, total, err := h.allocationService.ListAllocations(userID, projectID, startDate, endDate, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch allocations"})
		return
	}

	response := schemas.AllocationListResponse{
		Allocations: make([]schemas.AllocationResponse, len(allocations)),
		Total:       total,
		Offset:      offset,
		Limit:       limit,
	}

	for i, allocation := range allocations {
		response.Allocations[i] = *h.mapAllocationToResponse(allocation)
	}

	c.JSON(http.StatusOK, response)
}

func (h *AllocationHandler) GetWeekAllocations(c *gin.Context) {
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

	allocations, totalHours, err := h.allocationService.GetWeekAllocations(userID, week)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch week allocations"})
		return
	}

	response := schemas.WeekAllocationResponse{
		WeekStarting: week.Format("2006-01-02"),
		TotalHours:   totalHours,
		Allocations:  make([]schemas.AllocationResponse, len(allocations)),
	}

	for i, allocation := range allocations {
		response.Allocations[i] = *h.mapAllocationToResponse(allocation)
	}

	c.JSON(http.StatusOK, response)
}

func (h *AllocationHandler) UpdateAllocation(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	allocationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid allocation ID"})
		return
	}

	var req schemas.UpdateAllocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	allocation, err := h.allocationService.UpdateAllocation(userID, allocationID, req.Hours, req.Notes)
	if err != nil {
		if err == services.ErrAllocationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Allocation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update allocation"})
		return
	}

	c.JSON(http.StatusOK, h.mapAllocationToResponse(allocation))
}

func (h *AllocationHandler) DeleteAllocation(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	allocationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid allocation ID"})
		return
	}

	if err := h.allocationService.DeleteAllocation(userID, allocationID); err != nil {
		if err == services.ErrAllocationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Allocation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete allocation"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *AllocationHandler) CopyWeekAllocations(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	var req schemas.CopyAllocationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	fromWeek, _ := time.Parse("2006-01-02", req.FromWeek)
	toWeek, _ := time.Parse("2006-01-02", req.ToWeek)

	allocations, err := h.allocationService.CopyWeekAllocations(userID, fromWeek, toWeek)
	if err != nil {
		if err == services.ErrInvalidWeekStart {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Target week must start on Monday"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy allocations"})
		return
	}

	response := schemas.CopyAllocationsResponse{
		CopiedCount: len(allocations),
		Allocations: make([]schemas.AllocationResponse, len(allocations)),
	}

	for i, allocation := range allocations {
		response.Allocations[i] = *h.mapAllocationToResponse(allocation)
	}

	c.JSON(http.StatusOK, response)
}

func (h *AllocationHandler) mapAllocationToResponse(allocation *models.Allocation) *schemas.AllocationResponse {
	response := &schemas.AllocationResponse{
		ID:           allocation.ID,
		ProjectID:    allocation.ProjectID,
		WeekStarting: allocation.WeekStarting.Format("2006-01-02"),
		Hours:        allocation.Hours,
		Notes:        allocation.Notes,
		CreatedAt:    allocation.CreatedAt,
		UpdatedAt:    allocation.UpdatedAt,
	}

	if allocation.Project.ID != uuid.Nil {
		response.Project = &schemas.ProjectSummary{
			ID:           allocation.Project.ID,
			Name:         allocation.Project.Name,
			Code:         allocation.Project.Code,
			BillableRate: allocation.Project.BillableRate,
			Currency:     allocation.Project.Currency,
			Client: schemas.ClientSummary{
				ID:   allocation.Project.Client.ID,
				Name: allocation.Project.Client.Name,
				Code: allocation.Project.Client.Code,
			},
		}
	}

	return response
}
