package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"go-jwt-mysql/config"
	"go-jwt-mysql/models"
	"go-jwt-mysql/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

type UserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// For regular users to get their own profile
func (uc *UserController) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := models.GetUserByID(config.DB, userID.(int))
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, user)
}

// For admin to get all users
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := models.GetAllUsers(config.DB)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, users)
}

// For admin to get specific user by ID
func (uc *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := models.GetUserByID(config.DB, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, user)
}

// For regular users to update their own profile
func (uc *UserController) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	err := models.UpdateUser(config.DB, userID.(int), &user) // Added type assertion
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	updatedUser, err := models.GetUserByID(config.DB, userID.(int)) // Added type assertion
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch updated user")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, updatedUser)
}

// For admin to update any user
func (uc *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	err = models.UpdateUser(config.DB, userID, &user)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %v", err))
		return
	}

	updatedUser, err := models.GetUserByID(config.DB, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch updated user")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, updatedUser)
}

// For regular users to delete their own account
func (uc *UserController) DeleteCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err := models.DeleteUser(config.DB, userID.(int)) // Added type assertion
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// For admin to delete any user
func (uc *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = models.DeleteUser(config.DB, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}
