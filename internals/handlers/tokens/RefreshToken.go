package tokens

import (
	"fmt"
	"net/http"
	"os"
	"shelf/db/database"
	"shelf/db/models"
	"shelf/internals/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RefreshHandler(c *gin.Context) {
	oldToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no_refresh_token"})
		return
	}

	claims, err := helpers.VerifyRefreshToken(oldToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_refresh_token"})
		return
	}

	var user models.User
	if err := database.DB.Preload("Profile").
		Where("id = ?", claims.Subject).
		First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_found"})
		return
	}

	var storedToken models.RefreshToken
	now := time.Now()
	gracePeriod := now.Add(-10 * time.Second)

	result := database.DB.Where("token = ? AND (revoked = false OR (revoked = true AND revoked_at > ?))", oldToken, gracePeriod).First(&storedToken)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session_expired_or_revoked"})
		return
	}

	userIDStr := fmt.Sprintf("%s", user.ID)
	newAT, newRT, err := helpers.GenerateSystemTokens(userIDStr, user.Profile.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token_generation_failed"})
		return
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&storedToken).Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
		}).Error; err != nil {
			return err
		}
		return tx.Create(&models.RefreshToken{
			Token:     newRT,
			UserID:    userIDStr,
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
			Email:     user.Profile.Email,
		}).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rotation_failed"})
		return
	}

	domain := os.Getenv("BACKEND_AUTH_DOMAIN")

	c.SetCookie("access_token", newAT, 900, "/", domain, true, true)
	c.SetCookie("refresh_token", newRT, 604800, "/", domain, true, true)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
