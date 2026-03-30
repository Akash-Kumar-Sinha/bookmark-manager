package main

import (
	"os"
	"shelf/api/auth"
	"shelf/api/graphapi"
	"shelf/db/database"
	"shelf/internals/helpers"

	"github.com/gin-gonic/gin"
)

func init() {
	database.LoadInitializers()
	database.ConnectToDb()
	helpers.InitKeys()
}

func main() {
	r := gin.Default()

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8000"
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "You are under my protection",
		})
	})

	r.GET("/api",
		func(c *gin.Context) {
			c.Redirect(301, "/")
		})
	r.GET("/api/v1", func(c *gin.Context) {
		c.Redirect(301, "/")
	})
	r.GET("/api/v1/auth", func(c *gin.Context) {
		c.Redirect(301, "/")
	})
	authRoutes := r.Group("api/v1/auth")
	auth.AuthRoutes(authRoutes)

	graphqlRoutes := r.Group("/api/v1/bookmarks")
	graphapi.GraphRoutes(graphqlRoutes)

	r.Run("0.0.0.0:" + PORT)

}
