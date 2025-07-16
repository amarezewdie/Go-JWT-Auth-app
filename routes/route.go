package routes

import (
	"go-jwt-mysql/controllers"
	"go-jwt-mysql/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	  router.Use(cors.Default()) 

	authController := controllers.AuthController{}
	userController := controllers.UserController{}

	// Auth routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// User routes (protected)
	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		api.GET("/users", userController.GetUsers)
		api.GET("/users/me", userController.GetUser)
		api.PUT("/users/me", userController.UpdateUser)
		api.DELETE("/users/me", userController.DeleteUser)
	}

	return router
}