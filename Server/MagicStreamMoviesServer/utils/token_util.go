package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/database"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY string = os.Getenv("SECRET_REFRESH_KEY")

// Generate BOTH Access Token & Refresh Token
func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	// Create Access Token WITH Secret Key
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	// Create Refresh Token WITH Secret Refresh Key
	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

// Update Access & Refresh Token information from DB (updating "users" Collection)
func UpdateAllTokens(userId, token, refreshToken string, client *mongo.Client, c *gin.Context) (err error) {
	// Setup Timeout for Request
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    updatedAt,
		},
	}
	// bson.M{} inside means "filtering"
	var userCollection *mongo.Collection = database.OpenCollection("users", client)
	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	if err != nil {
		return err
	}
	// all good
	return nil
}

// Get Access Token String from Request
func GetAccessToken(c *gin.Context) (string, error) {
	// // Unsafe way: without Cookies
	// authHeader := c.Request.Header.Get("Authorization")
	// if authHeader == "" {
	// 	return "", errors.New("Authorization header is required")
	// }
	// // Header has an Authorization field e.g. "Bearer xxxyyy"
	// tokenString := authHeader[len("Bearer "):]
	// if tokenString == "" {
	// 	return "", errors.New("Bearer token is required")
	// }

	// Safe way with Cookie
	tokenString, err := c.Cookie("access_token")
	if err != nil {
		return "", err
	}

	// token good
	return tokenString, nil
}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}
	// 1 - Parse the token string => token object
	// 2 - Decode claims into "claims" variable
	// 3 - The callback at the end is to verify the JWT's signature
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the algorithm used is HMAC
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	// If token has expired (from info in claims)
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	// all good => return claims info
	return claims, nil
}

func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return "", errors.New("userId doesn't exists in this context")
	}
	id, ok := userId.(string)
	if !ok {
		return "", errors.New("Unable to retrieve userId")
	}
	// all good
	return id, nil
}

func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("role does not exists in this context")
	}
	memberRole, ok := role.(string)
	if !ok {
		return "", errors.New("unable to retrieve userId")
	}
	// all good
	return memberRole, nil
}

func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}
	// Parse information into "claims" (check Signature with the Callback)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_REFRESH_KEY), nil
	})
	if err != nil {
		return nil, err
	}
	// Check if Token was created with HMAC method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}
	// Check if claims expired or not
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
