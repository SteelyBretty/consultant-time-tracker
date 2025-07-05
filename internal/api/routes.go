package api

import (
	"github.com/SteelyBretty/consultant-time-tracker/internal/handlers"
	"github.com/SteelyBretty/consultant-time-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	authHandler := handlers.NewAuthHandler()

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.BasicAuth(), authHandler.GetCurrentUser)
		}
	}
}
