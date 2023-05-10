package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://admin:admin@development.xkgdtfm.mongodb.net/test"))
	if err != nil {
		errors.New("[X] Error occured while trying to connect to DB")
	}

	//Ping database to be sure we are connected
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[*] Successfully connected to Database!")

	//Gather info for database and collection
	database := client.Database("LilShorty")
	collection := database.Collection("Links")

	server := gin.Default()

	//Shorten your url to string of 4 letters and 2 numbers
	server.POST("/createShort/:url", func(ctx *gin.Context) {
		url := ctx.Param("url")
		key := generateString()

		filterLink := bson.D{
			{Key: "Link", Value: url},
		}

		filterKey := bson.D{
			{Key: "Link", Value: key},
		}

		doc := bson.D{
			{Key: "Link", Value: url},
			{Key: "Key", Value: key},
		}

		var existingLink bson.M
		errLink := collection.FindOne(context.Background(), filterLink).Decode(&existingLink)
		if errLink == nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"status": "This URL is already shortened",
			})
			return
		} else if errLink != mongo.ErrNoDocuments {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Failed to check if URL is already shortened",
			})
			return
		}

		var existingKey bson.M
		errKey := collection.FindOne(context.Background(), filterKey).Decode(&existingKey)
		if errKey == nil {
			key = generateString()
			doc = bson.D{
				{Key: "Link", Value: url},
				{Key: "Key", Value: key},
			}
			return
		} else if errKey != mongo.ErrNoDocuments {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Failed to check if key is already used",
			})
			return
		}

		_, err = collection.InsertOne(context.Background(), doc)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Failed to insert the document into the collection",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": "URL shortened, key is: " + key,
		})
	})

	//Get the link we will use to redirect users
	server.GET("/short/:key", func(ctx *gin.Context) {
		key := ctx.Param("key")

		filter := bson.D{{Key: "Key", Value: key}}
		var result bson.M
		err := collection.FindOne(context.Background(), filter).Decode(&result)
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "Shortened URL not found",
			})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status": "Failed to retrieve the URL",
			})
			return
		}

		url := result["Link"].(string)
		ctx.JSON(http.StatusOK, gin.H{
			"status": url,
		})
	})

	//Run this webserver finally
	server.Run("localhost:42069")
}

// Helper for key generation
func generateString() string {
	const (
		letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers = "0123456789"
	)

	rand.Seed(time.Now().UnixNano())

	result := make([]byte, 6)
	for i := 0; i < 4; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}

	for i := 4; i < 6; i++ {
		result[i] = numbers[rand.Intn(len(numbers))]
	}

	return string(result)
}
