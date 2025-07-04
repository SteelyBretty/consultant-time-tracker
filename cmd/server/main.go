package main

import (
	"log"
	"os"

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

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "consultant-time-tracker",
			"version": "0.1.2",
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

	log.Printf("Server starting on %s:%s", host, port)
	if err := r.Run(host + ":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
