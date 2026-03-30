package handlers

import (
	"net/http"
	"os"
	"shelf/db/database"
	"shelf/db/models"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	userEmail := c.GetString("userEmail")
	if userEmail == "" {
		c.JSON(401, gin.H{"error": "user email not found"})
		return
	}

	var profile models.Profile
	err := database.DB.Where("email = ?", userEmail).First(&profile).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or already logged out"})
		return
	}

	accessToken, err := c.Cookie("access_token")
	refreshToken, err := c.Cookie("refresh_token")

	if err == nil {
		database.DB.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})
	}

	if err != nil || accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
		return
	}

	domain := os.Getenv("BACKEND_AUTH_DOMAIN")
	c.SetCookie("access_token", "", -1, "/", domain, true, true)
	c.SetCookie("refresh_token", "", -1, "/", domain, true, true)

	revokeURL := "https://oauth2.googleapis.com/revoke?token=" + accessToken
	resp, err := http.Post(revokeURL, "application/x-www-form-urlencoded", nil)
	if err == nil {
		resp.Body.Close()
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
