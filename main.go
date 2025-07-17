package main

import (
	"go-jwt-mysql/config"
	"go-jwt-mysql/models"
	"go-jwt-mysql/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load .env file: %v", err)
	}

	// Initialize database
	config.ConnectDB()
	config.CreateUserTable()
	// Initialize Admin user if not exists
	if err := models.InitAdminUser(config.DB); err != nil {
		log.Fatalf("❌ Failed to initialize admin user: %v", err)
	}

	// Setup routes
	router := routes.SetupRoutes()

	// Start server
	router.Run(":8080")
}
