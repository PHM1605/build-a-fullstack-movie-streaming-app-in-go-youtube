package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/database"
	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/models"
	"github.com/PHM1605/build-a-fullstack-movie-streaming-app-in-go-youtube/Server/MagicStreamMoviesServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")
var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// create a Cursor to read our Collection
		cursor, err := movieCollection.Find(ctx, bson.M{}) // M: a BSON document
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		defer cursor.Close(ctx)

		// Fill "movies" from Cursor reading into Collection
		var movies []models.Movie
		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
		}

		// All movies are filled good
		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add timeout to "gin-Context"
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		// Get movie's ID from the URL i.e. "/movie/:imdb_id"
		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}
		// movie's ID exists; start filtering it out from Collection
		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		// all good
		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		// context with timeout
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()
		// parse Movie to be added
		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		// validate the "tags" in Movie Struct
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}
		// insert good Movie into DB
		// result: ID of the inserted entry
		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}
		// all good
		c.JSON(http.StatusCreated, result)
	}
}

func AdminReviewUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// to update Review, User must be "ADMIN"
		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in context"})
			return
		}
		if role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be part of the ADMIN role"})
			return
		}

		movieId := c.Param("imdb_id") // on URL
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID required"})
			return
		}
		// Define Request and Response formats of Admin Review Update part
		var req struct {
			AdminReview string `json:"admin_review"`
		}
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}
		// Parsing Request's Form Fields to Request Struct "req"
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		// get sentiment of Review from LLM
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview)
		fmt.Println(sentiment, rankVal, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking"})
			return
		}
		// update MongoDB; we setup aggregation stages
		filter := bson.M{"imdb_id": movieId} // filter out the correct Document
		// update which fields in that document
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}
		// set timeout for long Requests
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		// update MongoDB
		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}
		// if we can't find the Document to update
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}
		// Setup the Response
		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview
		c.JSON(http.StatusOK, resp)
	}
}

func GetReviewRanking(admin_review string) (string, int, error) {
	// Get list of rankings possible from DB
	rankings, err := GetRankings()
	if err != nil {
		return "", 0, err
	}

	// To store "Excellent,Good,Okay,Bad,Terrible"
	sentimentDelimited := ""
	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}
	// Remove the last ","
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")
	// Load .env
	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}
	// Read .env
	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	if OpenAiApiKey == "" {
		return "", 0, errors.New("could not read OPENAI_API_KEY")
	}
	// setup LLM to ask for sentiment
	llm, err := openai.New(openai.WithToken(OpenAiApiKey))
	if err != nil {
		return "", 0, err
	}
	// get LLM prompt's template from .env
	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")
	// setup correct prompt template
	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)
	// call LLM
	response, err := llm.Call(context.Background(), base_prompt+admin_review)
	if err != nil {
		return "", 0, err
	}
	// check ranking sent from LLM with list of rankings from DB
	rankVal := 0
	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}
	return response, rankVal, nil
}

// Get list of rankings from DB (collection "rankings")
func GetRankings() ([]models.Ranking, error) {
	var rankings []models.Ranking

	// Setup timeout for request
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Get cursor into DB (ranking Collection)
	cursor, err := rankingCollection.Find(ctx, bson.M{}) // findAll() here with empty M{}
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Parsing information from DB in Cursor to "rankings variable"
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}
	// all good
	return rankings, nil
}

func GetRecommendedMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get UserID from context of request (set in Middleware - read auth_middleware.go to understand)
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Id not found in context"})
			return
		}
		// Get list of User's favourite genres e.g. ["Scifi", "Horror"]
		favourite_genres, err := GetUsersFavouriteGenres(userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Read .env
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Warning: .env file not found")
		}
		// Get number of recommend movies to be returned
		var recommendedMovieLimitVal int64 = 5
		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")
		if recommendedMovieLimitStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
		}
		// Find and sort result based on Ranking Value in ascending order (limit to 5)
		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
		findOptions.SetLimit(recommendedMovieLimitVal)
		// if 1 of genres of a Movie is IN list of User's favourite genres
		filter := bson.M{"genre.genre_name": bson.M{"$in": favourite_genres}}
		// Set timeout for long Request
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		// Filter movies with suitable genres & sort according to "ranking_value"
		cursor, err := movieCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)

		// Parse query result
		var recommendedMovies []models.Movie
		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// all good
		c.JSON(http.StatusOK, recommendedMovies)
	}
}

// Get list of favourite genres of a User Id
func GetUsersFavouriteGenres(userId string) ([]string, error) {
	// Setup timeout for too long Requests
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// Prepare stages of MongoDB aggregation
	filter := bson.M{"user_id": userId}
	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}
	opts := options.FindOne().SetProjection(projection)
	// Perform filter & projection THEN parsing to "result"
	var result bson.M
	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		// User's favourites not exist in DB
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}
	// check if we return an Array
	favGenresArray, ok := result["favourite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("unable to retrieve favourite genres for user")
	}
	// get list of Strings only; no need full item
	var genreNames []string
	for _, item := range favGenresArray {
		// if each item is an Object e.g. {genre_id: 5, genre_name: "Thriller"}
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}
	return genreNames, nil
}
