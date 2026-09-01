package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/database"
	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/models"
	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// hash good
	return string(hashPassword), nil
}

func RegisterUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// try parsing User info from Request
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		// User info from Request is parsed well; now check the "validate:xxx" tags
		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		// User info from Request good; now we hash password
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password"})
		}

		// setup Timeout for request
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// make sure Email is unique
		var userCollection *mongo.Collection = database.OpenCollection("users", client)
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		// Finish the entry of User to be added to DB
		user.UserID = bson.NewObjectID().Hex()
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		// Add to DB
		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		// all good
		c.JSON(http.StatusCreated, result)
	}
}

func LoginUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// to parse User Login Info from Request
		var userLogin models.UserLogin
		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}
		// Create Timeout for request's Context
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()
		// Look into DB with User email; parsing info of the logging in User
		var foundUser models.User
		var userCollection *mongo.Collection = database.OpenCollection("users", client)
		err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		// Compare the entering Password and hashed Password (in DB)
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userLogin.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		// User logins successfully => we create Access Token and Refresh Token
		token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}
		// update both Tokens to DB
		err = utils.UpdateAllTokens(foundUser.UserID, token, refreshToken, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		// Safe way to send back Tokens "http-only Cookie"
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: token,
			Path:  "/",
			// Domain: "localhost",
			MaxAge:   86400,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "refresh_token",
			Value: refreshToken,
			Path:  "/",
			// Domain: "localhost",
			MaxAge:   86400,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		// all good => return both to Client
		c.JSON(http.StatusOK, models.UserResponse{
			UserID:    foundUser.UserID,
			FirstName: foundUser.FirstName,
			LastName:  foundUser.LastName,
			Email:     foundUser.Email,
			Role:      foundUser.Role,
			// Token:           token, // unsafe
			// RefreshToken:    refreshToken, // unsafe
			FavouriteGenres: foundUser.FavouriteGenres,
		})
	}
}

func LogoutHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var UserLogout struct {
			UserId string `json:"user_id"`
		}
		// Parse info of which User is logging out (from Request)
		err := c.ShouldBindJSON(&UserLogout)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		fmt.Println("User ID from Logout request: ", UserLogout.UserId)

		// Clear tokens in the DB
		err = utils.UpdateAllTokens(UserLogout.UserId, "", "", client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging out"})
			return
		}

		// Clear the Request's Cookies
		http.SetCookie(c.Writer, &http.Cookie{
			Name:  "access_token",
			Value: "",
			Path:  "/",
			// Domain: "localhost",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		// All good
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

func RefreshTokenHandler(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// set Timeout for Request
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			fmt.Println("error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to return refresh token"})
			return
		}

		claim, err := utils.ValidateRefreshToken(refreshToken)
		if err != nil || claim == nil {
			fmt.Println("error", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		}

		// Compare (not yet expire) Refresh Token with what being stored in DB
		var userCollection *mongo.Collection = database.OpenCollection("users", client)
		// Parse User from DB
		var user models.User
		err = userCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: claim.UserId}}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Issue new Access Token and Refresh Token
		newToken, newRefreshToken, _ := utils.GenerateAllTokens(user.Email, user.FirstName, user.LastName, user.Role, user.UserID)
		err = utils.UpdateAllTokens(user.UserID, newToken, newRefreshToken, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating tokens"})
			return
		}
		// Access Token expires in 24 hours
		c.SetCookie("access_token", newToken, 86400, "/", "localhost", true, true)
		// Refresh Token expires in 1 week
		c.SetCookie("refresh_token", newRefreshToken, 604800, "/", "localhost", true, true)

		c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
	}
}
