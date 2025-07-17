package controllers

import (
	"log"
	"net/http"
	"os"
	"strings"

	"go-jwt-mysql/config"
	"go-jwt-mysql/models"
	"go-jwt-mysql/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user already exists
	_, err := models.GetUserByEmail(config.DB, input.Email)
	if err == nil {
		utils.RespondWithError(c, http.StatusConflict, "Email already exists")
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password, // raw password here, hashing done inside model
		Role:     "user",
	}

	// Set admin role if email matches env
	if input.Email == os.Getenv("ADMIN_EMAIL") {
		user.Role = "admin"
	}

	if err := models.CreateUser(config.DB, &user); err != nil {
		log.Printf("User creation error: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	user.Password = ""
	utils.RespondWithJSON(c, http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Login handles user login
func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Trim whitespace from inputs
	input.Email = strings.TrimSpace(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	user, err := models.GetUserByEmail(config.DB, input.Email)
	if err != nil {
		log.Printf("Login failed for %s: %v", input.Email, err)
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		log.Printf("Password mismatch for %s: %v", input.Email, err)
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.Name, user.ID, user.Role)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Successful response
	user.Password = "" // Don't return password
	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}
