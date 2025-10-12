package controllers

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
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
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request payload",
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		email := strings.ToLower(input.Email)

		// Check if email already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
		if err != nil {
			log.Printf("Error checking for existing user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Server error while checking email",
				"success": false,
			})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"status":  http.StatusConflict,
				"message": "This email already exists",
				"success": false,
			})
			return
		}

		// Hash password
		hashedPassword, err := helpers.HashPassword(input.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to process password",
				"success": false,
			})
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
			IsVerified: true,
			CreatedAt:  time.Now(),
		}

		// Generate tokens
		token, refreshToken, err := helpers.GenerateAllTokens(newUser.Email, newUser.ID.Hex())
		if err != nil {
			log.Printf("Error generating tokens: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to generate tokens",
				"sucess":  false,
			})
			return
		}

		// Insert into DB
		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to create user",
				"success": false,
			})
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
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Registration successful",
			"data":    response,
			"success": true,
		})
	}
}

func RegisterWithOtp() gin.HandlerFunc {
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
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request payload",
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		email := strings.ToLower(input.Email)

		// Check if email already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
		if err != nil {
			log.Printf("Error checking for existing user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Server error while checking email",
				"success": false,
			})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"status":  http.StatusConflict,
				"message": "This email already exists",
				"success": false,
			})
			return
		}

		// Hash password
		hashedPassword, err := helpers.HashPassword(input.Password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to process password",
				"success": false,
			})
			return
		}

		n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
		otp := int(n.Int64())
		otpExpire := time.Now().Add(5 * time.Minute)
		fmt.Println(otp)
		fmt.Println(otpExpire)

		// // Create new user model
		newUser := models.User{
			ID:         primitive.NewObjectID(),
			FirstName:  input.FirstName,
			LastName:   input.LastName,
			Email:      email,
			Password:   hashedPassword,
			IsAdmin:    input.IsAdmin,
			IsVerified: false,
			CreatedAt:  time.Now(),
			Otp:        &otp,
			OtpExpire:  &otpExpire,
		}

		// Insert into DB
		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to create user",
				"success": false,
			})
			return
		}

		fmt.Println(result)

		message := "An OTP has been sent to your email address for account verification " + strconv.Itoa(otp)

		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": message,
			"success": true,
		})
	}
}

func VerifyAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Struct to receive the signup payload
		var input struct {
			Email string `json:"email" binding:"required,email"`
			Otp   int    `json:"otp" binding:"required"`
		}

		// Bind JSON request body
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request payload",
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		email := strings.ToLower(input.Email)

		// Check if email exists
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		if err != nil {
			log.Printf("Login error (find): %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "User Account does not exist",
				"sucess":  false,
			})
			return
		}

		if foundUser.Otp == nil || *foundUser.Otp != input.Otp {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid OTP",
				"success": false,
			})
			return
		}

		if foundUser.Otp != nil && *foundUser.Otp == input.Otp && foundUser.OtpExpire != nil && (*foundUser.OtpExpire).After(time.Now()) {
			foundUser.IsVerified = true
			_, updateErr := userCollection.UpdateOne(
				ctx,
				bson.M{"_id": foundUser.ID},
				bson.M{"$set": bson.M{"isVerified": true}},
			)
			if updateErr != nil {
				log.Printf("Failed to verify account: %v", updateErr)
			}
		} else if foundUser.OtpExpire != nil && !(*foundUser.OtpExpire).After(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "OTP has expired",
				"success": false,
			})
			return
		}

		response := gin.H{
			"id":         foundUser.ID,
			"firstName":  foundUser.FirstName,
			"lastName":   foundUser.LastName,
			"email":      foundUser.Email,
			"isAdmin":    foundUser.IsAdmin,
			"isVerified": foundUser.IsVerified,
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "Account Verification Completed",
			"success": true,
			"data":    response,
		})
	}
}

func ResendOtp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Struct to receive the signup payload
		var input struct {
			Email string `json:"email" binding:"required,email"`
		}

		// Bind JSON request body
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request payload",
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		email := strings.ToLower(input.Email)

		// Check if email exists
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		if err != nil {
			log.Printf("Login error (find): %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "User Account does not exist",
				"sucess":  false,
			})
			return
		}

		n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
		otp := int(n.Int64())
		otpExpire := time.Now().Add(5 * time.Minute)
		fmt.Println(otp)
		fmt.Println(otpExpire)

		update := bson.M{
			"$set": bson.M{
				"otp":       otp,
				"otpExpire": otpExpire,
			},
		}

		_, updateErr := userCollection.UpdateOne(ctx, bson.M{"_id": foundUser.ID}, update)
		if updateErr != nil {
			log.Printf("Failed to update OTP: %v", updateErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to resend OTP",
				"success": false,
			})
			return
		}

		message := "An OTP has been sent to your email address for account verification " + strconv.Itoa(otp)

		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": message,
			"success": true,
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
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request payload",
				"error":   err.Error(),
				"success": false})
			return
		}

		email := strings.ToLower(loginRequest.Email)

		// Find user
		var foundUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		if err != nil {
			log.Printf("Login error (find): %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid email or password",
				"sucess":  false,
			})
			return
		}

		passwordIsValid, msg := helpers.VerifyPassword(loginRequest.Password, foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": msg,
				"success": false,
			})
			return
		}

		if !foundUser.IsVerified {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Account Verification is Incomplete",
				"success": false,
			})
			return
		}

		// Generate tokens
		token, refreshToken, err := helpers.GenerateAllTokens(foundUser.Email, foundUser.ID.Hex())
		if err != nil {
			log.Printf("Login error (token generation): %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to generate authentication tokens",
				"success": false,
			})
			return
		}

		// Update last login timestamp
		now := time.Now()
		foundUser.LastLogin = &now
		_, updateErr := userCollection.UpdateOne(
			ctx,
			bson.M{"_id": foundUser.ID},
			bson.M{"$set": bson.M{"lastLogin": foundUser.LastLogin}},
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
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Login successful",
			"data":    response,
			"success": true,
		})
	}
}
