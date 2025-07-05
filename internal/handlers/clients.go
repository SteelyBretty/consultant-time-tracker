package handlers

import (
	"net/http"
	"strconv"

	"github.com/SteelyBretty/consultant-time-tracker/internal/middleware"
	"github.com/SteelyBretty/consultant-time-tracker/internal/models"
	"github.com/SteelyBretty/consultant-time-tracker/internal/schemas"
	"github.com/SteelyBretty/consultant-time-tracker/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClientHandler struct {
	clientService *services.ClientService
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{
		clientService: services.NewClientService(),
	}
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	var req schemas.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	client, err := h.clientService.CreateClient(userID, req.Name, req.Code, req.Email, req.Phone, req.Address)
	if err != nil {
		if err == services.ErrClientCodeExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}

	c.JSON(http.StatusCreated, h.mapClientToResponse(client))
}

func (h *ClientHandler) GetClient(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	client, err := h.clientService.GetClientByID(userID, clientID)
	if err != nil {
		if err == services.ErrClientNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch client"})
		return
	}

	c.JSON(http.StatusOK, h.mapClientToResponse(client))
}

func (h *ClientHandler) ListClients(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		active := activeStr == "true"
		isActive = &active
	}

	if limit > 100 {
		limit = 100
	}

	clients, total, err := h.clientService.ListClients(userID, isActive, search, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clients"})
		return
	}

	response := schemas.ClientListResponse{
		Clients: make([]schemas.ClientResponse, len(clients)),
		Total:   total,
		Offset:  offset,
		Limit:   limit,
	}

	for i, client := range clients {
		response.Clients[i] = *h.mapClientToResponse(client)
	}

	c.JSON(http.StatusOK, response)
}

func (h *ClientHandler) UpdateClient(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	var req schemas.UpdateClientRequest
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
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	client, err := h.clientService.UpdateClient(userID, clientID, updates)
	if err != nil {
		if err == services.ErrClientNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		if err == services.ErrClientCodeExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client"})
		return
	}

	c.JSON(http.StatusOK, h.mapClientToResponse(client))
}

func (h *ClientHandler) DeleteClient(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	if err := h.clientService.DeleteClient(userID, clientID); err != nil {
		if err == services.ErrClientNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ClientHandler) mapClientToResponse(client *models.Client) *schemas.ClientResponse {
	response := &schemas.ClientResponse{
		ID:        client.ID,
		Name:      client.Name,
		Code:      client.Code,
		Email:     client.Email,
		Phone:     client.Phone,
		Address:   client.Address,
		IsActive:  client.IsActive,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}

	if len(client.Projects) > 0 {
		response.Projects = make([]schemas.ProjectResponse, len(client.Projects))
		for i, project := range client.Projects {
			response.Projects[i] = schemas.ProjectResponse{
				ID:           project.ID,
				Name:         project.Name,
				Code:         project.Code,
				Status:       string(project.Status),
				BillableRate: project.BillableRate,
				Currency:     project.Currency,
			}
		}
	}

	return response
}
