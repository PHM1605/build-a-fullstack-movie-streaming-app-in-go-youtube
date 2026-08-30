package middleware

import (
	"net/http"

	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Access Token raw string
		token, err := utils.GetAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort() // NEW: for middleware we handle Error like this
			return
		}
		// token is empty string
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			c.Abort() // NEW: for middleware we handle Error like this
			return
		}
		// Validate token; decode token into claims
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Put some info to request's context
		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)
		// Continue to Handler after Middleware
		c.Next()
	}
}
