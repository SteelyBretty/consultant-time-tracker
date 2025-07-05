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
	"gorm.io/datatypes"
)

type ProjectHandler struct {
	projectService *services.ProjectService
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		projectService: services.NewProjectService(),
	}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	var req schemas.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, _ := time.Parse("2006-01-02", *req.EndDate)
		endDate = &parsedEndDate
	}

	project, err := h.projectService.CreateProject(
		userID, req.ClientID, req.Name, req.Code, req.Description,
		req.BillableRate, req.Currency, startDate, endDate,
	)
	if err != nil {
		if err == services.ErrProjectCodeExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "client not found or access denied" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, h.mapProjectToResponse(project))
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.projectService.GetProjectByID(userID, projectID)
	if err != nil {
		if err == services.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}

	c.JSON(http.StatusOK, h.mapProjectToResponse(project))
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	var clientID *uuid.UUID
	if clientIDStr := c.Query("client_id"); clientIDStr != "" {
		if id, err := uuid.Parse(clientIDStr); err == nil {
			clientID = &id
		}
	}

	var status *models.ProjectStatus
	if statusStr := c.Query("status"); statusStr != "" {
		projectStatus := models.ProjectStatus(statusStr)
		status = &projectStatus
	}

	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		active := activeStr == "true"
		isActive = &active
	}

	if limit > 100 {
		limit = 100
	}

	projects, total, err := h.projectService.ListProjects(userID, clientID, status, isActive, search, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	response := schemas.ProjectListResponse{
		Projects: make([]schemas.ProjectResponse, len(projects)),
		Total:    total,
		Offset:   offset,
		Limit:    limit,
	}

	for i, project := range projects {
		response.Projects[i] = *h.mapProjectToResponse(project)
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req schemas.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.BillableRate != nil {
		updates["billable_rate"] = *req.BillableRate
	}
	if req.Currency != nil {
		updates["currency"] = *req.Currency
	}
	if req.StartDate != nil {
		startDate, _ := time.Parse("2006-01-02", *req.StartDate)
		updates["start_date"] = datatypes.Date(startDate)
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			updates["end_date"] = nil
		} else {
			endDate, _ := time.Parse("2006-01-02", *req.EndDate)
			updates["end_date"] = datatypes.Date(endDate)
		}
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	project, err := h.projectService.UpdateProject(userID, projectID, updates)
	if err != nil {
		if err == services.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		if err == services.ErrProjectCodeExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == services.ErrInvalidProjectStatus {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, h.mapProjectToResponse(project))
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.projectService.DeleteProject(userID, projectID); err != nil {
		if err == services.ErrProjectNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ProjectHandler) mapProjectToResponse(project *models.Project) *schemas.ProjectResponse {
	response := &schemas.ProjectResponse{
		ID:           project.ID,
		Name:         project.Name,
		Code:         project.Code,
		Description:  project.Description,
		Status:       string(project.Status),
		BillableRate: project.BillableRate,
		Currency:     project.Currency,
		StartDate:    time.Time(project.StartDate).Format("2006-01-02"),
		IsActive:     project.IsActive,
		ClientID:     project.ClientID,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
	}

	if project.EndDate != nil {
		endDateStr := time.Time(*project.EndDate).Format("2006-01-02")
		response.EndDate = &endDateStr
	}

	if project.Client.ID != uuid.Nil {
		response.Client = &schemas.ClientResponse{
			ID:       project.Client.ID,
			Name:     project.Client.Name,
			Code:     project.Client.Code,
			Email:    project.Client.Email,
			IsActive: project.Client.IsActive,
		}
	}

	return response
}
