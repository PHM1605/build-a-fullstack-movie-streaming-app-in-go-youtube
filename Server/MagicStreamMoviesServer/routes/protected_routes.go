package routes

import (
	controller "github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/controllers"
	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	// Register the Middleware Handler first
	router.Use(middleware.AuthMiddleWare())
	// Register normal Handlers
	router.GET("/movie/:imdb_id", controller.GetMovie())
	router.POST("/addmovie", controller.AddMovie())
	router.GET("/recommendedmovies", controller.GetRecommendedMovies())
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate())
}
