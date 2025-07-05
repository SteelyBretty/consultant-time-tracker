package main

import (
	"log"
	"os"

	"github.com/SteelyBretty/consultant-time-tracker/internal/api"
	"github.com/SteelyBretty/consultant-time-tracker/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	if err := database.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	if err := database.CreateIndexes(); err != nil {
		log.Fatal("Failed to create indexes:", err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := database.DB.DB()
		dbStatus := "healthy"
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "unhealthy"
		}

		c.JSON(200, gin.H{
			"status":   "healthy",
			"service":  "consultant-time-tracker",
			"version":  "0.1.0",
			"database": dbStatus,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Consultant Time Tracker API",
			"endpoints": gin.H{
				"health":  "/health",
				"api":     "/api/v1",
				"graphql": "/graphql",
			},
		})
	})

	api.SetupRoutes(r)

	log.Printf("Server starting on %s:%s", host, port)
	if err := r.Run(host + ":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
