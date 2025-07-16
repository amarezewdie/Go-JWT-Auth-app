package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-jwt-mysql/config"
	"go-jwt-mysql/models"
	"go-jwt-mysql/utils"
)

type UserController struct{}

type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := models.GetAllUsers(config.DB)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, users)
}

func (uc *UserController) GetUser(c *gin.Context) {
	userID := c.GetInt("userID")

	user, err := models.GetUserByID(config.DB, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	userID := c.GetInt("userID")

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Name:  input.Name,
		Email: input.Email,
	}

	err := models.UpdateUser(config.DB, userID, &user)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	updatedUser, err := models.GetUserByID(config.DB, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch updated user")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, updatedUser)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.GetInt("userID")

	err := models.DeleteUser(config.DB, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}