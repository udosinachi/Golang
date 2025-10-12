package controllers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
	"udo-golang/database"
	"udo-golang/helpers"
	"udo-golang/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout and defer cancel immediately
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		// Bind incoming JSON to the user model
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "hasError": true})
			return
		}

		user.Email = strings.ToLower(user.Email)

		// Check if the email already exists in the collection
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Printf("Error checking for email existence: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error occurred while checking for the email", "hasError": true})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "This email already exists", "hasError": true})
			return
		}

		user.ID = primitive.NewObjectID()
		user.IsAdmin = false
		user.IsVerified = false
		user.CreatedAt = time.Now()

		// Generate tokens using the email and the hex representation of the ObjectID
		token, refreshToken, err := helpers.GenerateAllTokens(user.Email, user.ID.Hex())
		if err != nil {
			log.Printf("Error generating tokens for user %s: %v", user.ID.Hex(), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		// Hash the user's password before storing it
		hashedPass, hashErr := helpers.HashPassword(user.Password)

		if hashErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": hashErr.Error(), "hasError": true})
			return
		}

		user.Password = hashedPass

		if validationErr := user.ValidateUser(); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": validationErr.Error(), "hasError": true})
			return
		}

		// Insert the user document into the database
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			log.Printf("Error inserting user: %v, %v", insertErr, user)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User item was not created", "hasError": true})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Registration successful",
			"data":         user,
			"token":        token,
			"refreshToken": refreshToken,
			"hasError":     false,
			"insertId":     resultInsertionNumber.InsertedID,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		user.Email = strings.ToLower(user.Email)
		user.LastLogin = time.Now()

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "email or password is incorrect", "error": err.Error(), "hasError": true})
			return
		}

		passwordIsValid, msg := helpers.VerifyPassword(user.Password, foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": msg, "hasError": true})
			return
		}

		if foundUser.Email == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found", "hasError": true})
			return
		}
		token, refreshToken, err := helpers.GenerateAllTokens(foundUser.Email, foundUser.ID.Hex())

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "hasError": true})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "data": foundUser, "token": token, "refreshToken": refreshToken, "hasError": false})
	}
}
