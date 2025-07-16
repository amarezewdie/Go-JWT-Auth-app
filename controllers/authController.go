package controllers

import (
	"net/http"

	"go-jwt-mysql/config"
	"go-jwt-mysql/models"
	"go-jwt-mysql/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user already exists
	_, err := models.GetUserByEmail(config.DB, input.Email)
	if err == nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Email already exists")
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	err = models.CreateUser(config.DB, &user)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

//login handles user login 
func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := models.GetUserByEmail(config.DB, input.Email)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	err = user.ComparePassword(input.Password)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateToken(user.Name, user.ID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
