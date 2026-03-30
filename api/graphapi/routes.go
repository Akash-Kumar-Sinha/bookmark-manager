package graphapi

import (
	"os"
	"shelf/db/database"
	"shelf/graph"
	"shelf/middlewares"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func GraphRoutes(r *gin.RouterGroup) {
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					DB: database.DB,
				},
			},
		),
	)

	r.POST("/query", middlewares.AuthSession(), func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	if os.Getenv("ENV") != "production" {
		h := playground.Handler("Shelf", "/api/v1/bookmarks/query")
		r.GET("/playground", middlewares.AuthSession(), func(c *gin.Context) {
			h.ServeHTTP(c.Writer, c.Request)
		})
	}
}
