package routes

import (
	controller "github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUnprotectedRoutes(router *gin.Engine) {
	router.GET("/movies", controller.GetMovies())
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
}
