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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Struct to receive the signup payload
		var input struct {
			FirstName string `json:"firstName" binding:"required"`
			LastName  string `json:"lastName" binding:"required"`
			Email     string `json:"email" binding:"required,email"`
			Password  string `json:"password" binding:"required,min=6"`
			IsAdmin   bool   `json:"isAdmin"`
		}

		// Bind JSON request body
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload", "error": err.Error(), "hasError": true})
			return
		}

		email := strings.ToLower(input.Email)

		// Check if email already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
		if err != nil {
			log.Printf("Error checking for existing user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error while checking email", "hasError": true})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "This email already exists", "hasError": true})
			return
		}

		// Hash password
		hashedPassword, err := helpers.HashPassword(input.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process password", "hasError": true})
			return
		}

		// Create new user model
		newUser := models.User{
			ID:         primitive.NewObjectID(),
			FirstName:  input.FirstName,
			LastName:   input.LastName,
			Email:      email,
			Password:   hashedPassword,
			IsAdmin:    input.IsAdmin,
			IsVerified: false,
			CreatedAt:  time.Now(),
		}

		// Generate tokens
		token, refreshToken, err := helpers.GenerateAllTokens(newUser.Email, newUser.ID.Hex())
		if err != nil {
			log.Printf("Error generating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate tokens", "hasError": true})
			return
		}

		// Insert into DB
		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "hasError": true})
			return
		}

		// Return a safe response (no password)
		response := gin.H{
			"id":           newUser.ID,
			"firstName":    newUser.FirstName,
			"lastName":     newUser.LastName,
			"email":        newUser.Email,
			"isAdmin":      newUser.IsAdmin,
			"isVerified":   newUser.IsVerified,
			"createdAt":    newUser.CreatedAt,
			"token":        token,
			"refreshToken": refreshToken,
			"insertId":     result.InsertedID,
			"hasError":     false,
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Registration successful",
			"data":    response,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var loginRequest struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload", "error": err.Error(), "hasError": true})
			return
		}

		email := strings.ToLower(loginRequest.Email)

		// Find user
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		if err != nil {
			log.Printf("Login error (find): %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password", "hasError": true})
			return
		}

		// Verify password
		passwordIsValid, msg := helpers.VerifyPassword(loginRequest.Password, foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": msg, "hasError": true})
			return
		}

		// Generate tokens
		token, refreshToken, err := helpers.GenerateAllTokens(foundUser.Email, foundUser.ID.Hex())
		if err != nil {
			log.Printf("Login error (token generation): %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate authentication tokens", "hasError": true})
			return
		}

		// Update last login timestamp
		foundUser.LastLogin = time.Now()
		_, updateErr := userCollection.UpdateOne(
			ctx,
			bson.M{"_id": foundUser.ID},
			bson.M{"$set": bson.M{"last_login": foundUser.LastLogin}},
		)
		if updateErr != nil {
			log.Printf("Failed to update last login: %v", updateErr)
		}

		// Prepare safe response
		response := gin.H{
			"id":           foundUser.ID,
			"firstName":    foundUser.FirstName,
			"lastName":     foundUser.LastName,
			"email":        foundUser.Email,
			"isAdmin":      foundUser.IsAdmin,
			"isVerified":   foundUser.IsVerified,
			"lastLogin":    foundUser.LastLogin,
			"token":        token,
			"refreshToken": refreshToken,
			"hasError":     false,
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"data":    response,
		})
	}
}
