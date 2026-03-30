package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"shelf/db/database"
	"shelf/db/models"
	"shelf/internals/helpers"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	userInfo, err := helpers.GetUserInfoFromGoogle(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	if idToken, ok := token.Extra("id_token").(string); ok {
		if err := ValidateGoogleIDToken(idToken, userInfo.ID); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
			return
		}
	}

	user, err := findOrCreateUser(userInfo)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User processing failed"})
		return
	}

	userIDStr := fmt.Sprintf("%s", user.ProfileID)
	fmt.Print("profileId")
	accessToken, refreshToken, err := helpers.GenerateSystemTokens(userIDStr, user.Profile.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate system tokens"})
		return
	}

	err = database.DB.Create(&models.RefreshToken{
		Token:     refreshToken,
		UserID:    userIDStr,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Email:     user.Profile.Email,
	}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	backendDomain := os.Getenv("BACKEND_AUTH_DOMAIN")

	c.SetCookie("access_token", accessToken, 900, "/", backendDomain, true, true)

	c.SetCookie("refresh_token", refreshToken, 604800, "/", backendDomain, true, true)

	frontendURL := os.Getenv("FRONTEND_DOMAIN")
	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/workspace")
}

func ValidateGoogleIDToken(idToken, expectedSubject string) error {
	clientID := os.Getenv("GOOGLE_AUTH_CLIENT_ID")
	payload, err := idtoken.Validate(context.Background(), idToken, clientID)
	if err != nil {
		return fmt.Errorf("failed to validate ID token: %v", err)
	}

	if payload.Subject != expectedSubject {
		return fmt.Errorf("subject mismatch: expected %s, got %s", expectedSubject, payload.Subject)
	}

	return nil
}

func findOrCreateUser(userInfo *helpers.GoogleUserInfo) (*models.User, error) {
	var user models.User
	var profile models.Profile

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Preload("Profile").Where("provider = ? AND provider_client_id = ?", "google", userInfo.ID).First(&user).Error
	if err != nil {
		profileErr := tx.Where("email = ?", userInfo.Email).First(&profile).Error

		if profileErr != nil {
			addr, err := mail.ParseAddress(userInfo.Email)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("invalid email format: %v", err)
			}
			if addr.Address != userInfo.Email {
				tx.Rollback()
				return nil, fmt.Errorf("email address mismatch: expected %s, got %s", userInfo.Email, addr.Address)
			}
			username := strings.SplitN(addr.Address, "@", 2)[0]
			firstName, middleName, lastName := splitName(userInfo.Name)

			middleNameStr := ""
			if middleName != nil {
				middleNameStr = *middleName
			}

			profile = models.Profile{
				Email:      userInfo.Email,
				FirstName:  firstName,
				MiddleName: middleNameStr,
				Username:   username,
				LastName:   lastName,
				Avatar:     userInfo.Picture,
			}

			if err := tx.Create(&profile).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create profile: %v", err)
			}
		}

		user = models.User{
			ProfileID:        profile.ID,
			Provider:         "google",
			ProviderClientID: userInfo.ID,
			LastLogin:        time.Now(),
		}

		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create user: %v", err)
		}

		user.Profile = profile
	} else {

		if err := tx.Model(&user).Update("last_login", time.Now()).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update last login: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &user, nil
}

func splitName(fullName string) (string, *string, string) {
	var firstName, middleName, lastName string
	nameParts := strings.Fields(fullName)

	if len(nameParts) == 1 {
		firstName = nameParts[0]
	} else if len(nameParts) == 2 {
		firstName = nameParts[0]
		lastName = nameParts[1]
	} else if len(nameParts) > 2 {
		firstName = nameParts[0]
		lastName = nameParts[len(nameParts)-1]
		middle := strings.Join(nameParts[1:len(nameParts)-1], " ")
		middleName = middle
	}

	if middleName == "" {
		return firstName, nil, lastName
	}
	return firstName, &middleName, lastName
}
