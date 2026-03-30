package handlers

import (
	"net/http"
	"shelf/db/database"
	"shelf/db/models"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) {

	userEmail := c.GetString("userEmail")
	if userEmail == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User email not found in context"})
		return
	}

	var profile models.Profile
	err := database.DB.Where("email = ?", userEmail).First(&profile).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or token invalid"})
		return
	}
	c.JSON(
		http.StatusOK, gin.H{
			"profile": profile,
		},
	)
}
