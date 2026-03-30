package auth

import (
	"shelf/internals/handlers"
	"shelf/internals/handlers/tokens"
	"shelf/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) {
	router.GET("/oauth/google/login", handlers.GoogleLogin)
	router.GET("/oauth/callback", handlers.GoogleCallback)

	router.POST("/oauth/refresh", tokens.RefreshHandler)

	router.POST("/logout", middlewares.AuthSession(), handlers.Logout)
	router.GET("/me", middlewares.AuthSession(), handlers.GetCurrentUser)
}
