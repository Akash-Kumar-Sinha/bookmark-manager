package middlewares

import (
	"fmt"
	"net/http"
	"shelf/internals/helpers"

	"github.com/gin-gonic/gin"
)

func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := c.Cookie("access_token")
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := helpers.VerifyAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token_expired"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", fmt.Sprintf("%s", claims.Email))

		c.Next()
	}
}
