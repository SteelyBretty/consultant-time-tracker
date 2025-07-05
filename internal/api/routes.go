package api

import (
	"github.com/SteelyBretty/consultant-time-tracker/internal/handlers"
	"github.com/SteelyBretty/consultant-time-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	authHandler := handlers.NewAuthHandler()
	clientHandler := handlers.NewClientHandler()
	projectHandler := handlers.NewProjectHandler()
	allocationHandler := handlers.NewAllocationHandler()

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.BasicAuth(), authHandler.GetCurrentUser)
		}

		protected := api.Group("/")
		protected.Use(middleware.BasicAuth())
		{
			clients := protected.Group("/clients")
			{
				clients.POST("", clientHandler.CreateClient)
				clients.GET("", clientHandler.ListClients)
				clients.GET("/:id", clientHandler.GetClient)
				clients.PUT("/:id", clientHandler.UpdateClient)
				clients.DELETE("/:id", clientHandler.DeleteClient)
			}

			projects := protected.Group("/projects")
			{
				projects.POST("", projectHandler.CreateProject)
				projects.GET("", projectHandler.ListProjects)
				projects.GET("/:id", projectHandler.GetProject)
				projects.PUT("/:id", projectHandler.UpdateProject)
				projects.DELETE("/:id", projectHandler.DeleteProject)
			}

			allocations := protected.Group("/allocations")
			{
				allocations.POST("", allocationHandler.CreateAllocation)
				allocations.GET("", allocationHandler.ListAllocations)
				allocations.GET("/week", allocationHandler.GetWeekAllocations)
				allocations.GET("/:id", allocationHandler.GetAllocation)
				allocations.PUT("/:id", allocationHandler.UpdateAllocation)
				allocations.DELETE("/:id", allocationHandler.DeleteAllocation)
				allocations.POST("/copy", allocationHandler.CopyWeekAllocations)
			}
		}
	}
}
