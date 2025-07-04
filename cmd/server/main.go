package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Consultant Time Tracker starting on port %s...\n", port)
	log.Fatal("Not implemented yet")
}
