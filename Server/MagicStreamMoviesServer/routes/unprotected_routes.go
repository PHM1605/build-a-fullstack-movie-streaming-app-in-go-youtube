package routes

import (
	controller "github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUnprotectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.GET("/movies", controller.GetMovies(client))
	router.POST("/register", controller.RegisterUser(client))
	router.POST("/login", controller.LoginUser(client))
	router.POST("/logout", controller.LogoutHandler(client))
	router.GET("/genres", controller.GetGenres(client))
	router.POST("/refresh", controller.RefreshTokenHandler(client))
}
