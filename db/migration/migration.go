package main

import (
	"log"
	"shelf/db/database"
	"shelf/db/models"
)

func init() {
	database.LoadInitializers()
	database.ConnectToDb()
}

func main() {
	if err := database.DB.AutoMigrate(&models.RefreshToken{}); err != nil {
		log.Printf("Error migrating RefreshToken: %v", err)
		panic(err)
	}

	if err := database.DB.AutoMigrate(models.Profile{}); err != nil {
		log.Printf("Error during migration of Profile: %v", err)
		panic(err)
	}
	if err := database.DB.AutoMigrate(models.User{}); err != nil {
		log.Printf("Error during migration of User: %v", err)
		panic(err)
	}

	if err := database.DB.AutoMigrate(models.Bookmark{}); err != nil {
		log.Printf("Error during migration of Bookmark: %v", err)
		panic(err)
	}

	if err := database.DB.AutoMigrate(models.Tag{}); err != nil {
		log.Printf("Error during migration of Tag: %v", err)
		panic(err)
	}

	if err := database.DB.AutoMigrate(models.BookmarkTag{}); err != nil {
		log.Printf("Error during migration of BookmarkTag: %v", err)
		panic(err)
	}
}
