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

	// Protected routes
	api := router.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		// Current user operations
		api.GET("/users/me", userController.GetCurrentUser)
		api.PUT("/users/me", userController.UpdateCurrentUser)
		api.DELETE("/users/me", userController.DeleteCurrentUser)

		// Admin-only operations
		admin := api.Group("/admin")
		admin.Use(middlewares.AdminMiddleware())
		{
			admin.GET("/users", userController.GetAllUsers)
			admin.GET("/users/:id", userController.GetUserByID)
			admin.PUT("/users/:id", userController.UpdateUser)
			admin.DELETE("/users/:id", userController.DeleteUser)
		}
	}

	return router
}
