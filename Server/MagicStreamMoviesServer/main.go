package main

import (
	"fmt"

	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, MagicStreamMovies!")
	})

	routes.SetupUnprotectedRoutes(router)
	routes.SetupProtectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
