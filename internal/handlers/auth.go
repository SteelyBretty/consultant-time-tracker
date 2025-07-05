package handlers

import (
	"net/http"

	"github.com/SteelyBretty/consultant-time-tracker/internal/models"
	"github.com/SteelyBretty/consultant-time-tracker/internal/schemas"
	"github.com/SteelyBretty/consultant-time-tracker/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req schemas.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	user, err := h.authService.Register(req.Username, req.Email, req.Password, req.FullName)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == services.ErrUsernameExists || err == services.ErrEmailExists {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, schemas.AuthResponse{
		User: schemas.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			IsActive: user.IsActive,
		},
		Message: "User registered successfully",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req schemas.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		statusCode := http.StatusUnauthorized
		if err == services.ErrUserNotActive {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.AuthResponse{
		User: schemas.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			IsActive: user.IsActive,
		},
		Message: "Login successful",
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	currentUser := user.(*models.User)
	c.JSON(http.StatusOK, schemas.UserResponse{
		ID:       currentUser.ID,
		Username: currentUser.Username,
		Email:    currentUser.Email,
		FullName: currentUser.FullName,
		IsActive: currentUser.IsActive,
	})
}
