package handlers

import (
	"net/http"
	"os"
	"shelf/db/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func init() {
	database.LoadInitializers()
	database.ConnectToDb()
	initGoogleOAuthConfig()
}

func initGoogleOAuthConfig() {
	redirectURL := os.Getenv("GOOGLE_AUTH_REDIRECT_URL")
	if redirectURL == "" {
		panic("Google OAuth redirect URL not configured")
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     os.Getenv("GOOGLE_AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_AUTH_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin(c *gin.Context) {
	if googleOauthConfig == nil {
		initGoogleOAuthConfig()
	}

	clientID := os.Getenv("GOOGLE_AUTH_CLIENT_ID")
	if clientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth client ID not configured"})
		return
	}

	url := googleOauthConfig.AuthCodeURL(
		"random-state",
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)

	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}
