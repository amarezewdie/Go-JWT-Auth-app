package main

import (
	"go-jwt-mysql/config"
	"go-jwt-mysql/routes"
)

func main() {

	// Initialize database
	config.ConnectDB()

	// Setup routes
	router := routes.SetupRoutes()

	// Start server
	router.Run(":8080")
}
