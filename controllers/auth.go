package controllers

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
	"udo-golang/helpers"
	"udo-golang/models"
	"udo-golang/queries"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		_, err := queries.GetUserByEmail(email)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
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
		token, refreshToken, err := helpers.GenerateAllTokens(newUser.Email, newUser.ID.Hex(), newUser.IsAdmin)
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
		result, err := queries.CreateNewUser(&newUser)
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

		_, err := queries.GetUserByEmail(email)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
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
		otpString := strconv.Itoa(otp)
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
			Otp:        &otpString,
			OtpExpire:  &otpExpire,
		}

		// Insert into DB
		result, err := queries.CreateNewUser(&newUser)
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

		message := "An OTP has been sent to your email address for account verification " + otpString

		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": message,
			"success": true,
		})
	}
}

func GoogleSignUpandSignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Query("access_token")
		if accessToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Access token is required",
				"success": false,
			})
			return
		}

		userInfo, err := queries.GetGoogleUserInfo(accessToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Error fetching Google user info",
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		email, _ := userInfo["email"].(string)
		id, _ := userInfo["id"].(string)
		name, _ := userInfo["name"].(string)

		parts := strings.Fields(name)

		firstName := ""
		lastName := ""

		if len(parts) > 0 {
			firstName = parts[0]
		}
		if len(parts) > 1 {
			lastName = parts[len(parts)-1]
		}

		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Email not found in Google profile",
				"success": false,
			})
			return
		}

		signedToken, err := helpers.SignJWt(email, id, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Error generating token",
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		foundUser, err := queries.GetUserByEmail(email)
		if err == nil {
			now := time.Now()
			if err := queries.UpdateUser(foundUser.ID.Hex(), bson.M{"lastLogin": &now}); err != nil {
				log.Printf("Failed to update last login: %v", err)
			}

			response := gin.H{
				"id":         foundUser.ID,
				"firstName":  foundUser.FirstName,
				"lastName":   foundUser.LastName,
				"email":      foundUser.Email,
				"isAdmin":    foundUser.IsAdmin,
				"isVerified": foundUser.IsVerified,
				"lastLogin":  now,
				"token":      signedToken,
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "Login successful",
				"data":    response,
				"success": true,
			})
			return
		}

		newUser := models.User{
			ID:         primitive.NewObjectID(),
			FirstName:  firstName,
			LastName:   lastName,
			Email:      email,
			Password:   "",
			IsAdmin:    false,
			IsVerified: true,
			CreatedAt:  time.Now(),
		}

		result, err := queries.CreateNewUser(&newUser)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to create user",
				"success": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "User Created Successful",
			"data": gin.H{
				"id":    result.InsertedID,
				"name":  name,
				"email": email,
				"token": signedToken,
			},
			"success": true,
		})
	}
}

func VerifyAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email string `json:"email" binding:"required,email"`
			Otp   string `json:"otp" binding:"required"`
		}

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

		foundUser, err := queries.GetUserByEmail(email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "User account does not exist",
				"success": false,
			})
			return
		}

		if foundUser.Otp == nil || *foundUser.Otp != input.Otp {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid OTP",
				"success": false,
			})
			return
		}

		if foundUser.OtpExpire == nil || time.Now().After(*foundUser.OtpExpire) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "OTP has expired",
				"success": false,
			})
			return
		}

		updateData := bson.M{
			"isVerified": true,
			"otp":        nil,
			"otpExpire":  nil,
		}

		if err := queries.UpdateUser(foundUser.ID.Hex(), updateData); err != nil {
			log.Printf("Failed to verify account: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to complete account verification",
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
			"isVerified": true,
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Account verification completed successfully",
			"success": true,
			"data":    response,
		})
	}
}

func SendOtp() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		foundUser, err := queries.GetUserByEmail(email)
		if err != nil {
			log.Printf("Resend OTP error (find): %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "User Account does not exist",
				"sucess":  false,
			})
			return
		}

		n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
		otp := int(n.Int64())
		otpString := strconv.Itoa(otp)
		otpExpire := time.Now().Add(5 * time.Minute)
		fmt.Println(otp)
		fmt.Println(otpExpire)

		update := bson.M{
			"otp":       otpString,
			"otpExpire": otpExpire,
		}

		if err := queries.UpdateUser(foundUser.ID.Hex(), update); err != nil {
			log.Printf("Failed to update OTP: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to resend OTP",
				"success": false,
			})
			return
		}

		message := "An OTP has been sent to your email address for account verification " + otpString

		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": message,
			"success": true,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request payload",
				"error":   err.Error(),
				"success": false,
			})
			return
		}

		email := strings.ToLower(loginRequest.Email)

		foundUser, err := queries.GetUserByEmail(email)
		if err != nil {
			log.Printf("Login error (find): %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid email or password",
				"success": false,
			})
			return
		}

		passwordIsValid, msg := helpers.VerifyPassword(loginRequest.Password, foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": msg,
				"success": false,
			})
			return
		}

		if !foundUser.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Account verification incomplete",
				"success": false,
			})
			return
		}

		token, refreshToken, err := helpers.GenerateAllTokens(foundUser.Email, foundUser.ID.Hex(), foundUser.IsAdmin)
		if err != nil {
			log.Printf("Token generation error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to generate authentication tokens",
				"success": false,
			})
			return
		}

		now := time.Now()
		if err := queries.UpdateUser(foundUser.ID.Hex(), bson.M{"lastLogin": &now}); err != nil {
			log.Printf("Failed to update last login: %v", err)
		}

		response := gin.H{
			"id":           foundUser.ID,
			"firstName":    foundUser.FirstName,
			"lastName":     foundUser.LastName,
			"email":        foundUser.Email,
			"isAdmin":      foundUser.IsAdmin,
			"isVerified":   foundUser.IsVerified,
			"lastLogin":    now,
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

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email           string `json:"email" binding:"required,email"`
			Otp             string `json:"otp" binding:"required" `
			Password        string `json:"password" binding:"required"`
			ConfirmPassword string `json:"confirmPassword" binding:"required"`
		}

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

		foundUser, err := queries.GetUserByEmail(email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "User account does not exist",
				"success": false,
			})
			return
		}

		if foundUser.Otp == nil || *foundUser.Otp != input.Otp {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid OTP",
				"success": false,
			})
			return
		}

		if foundUser.OtpExpire == nil || time.Now().After(*foundUser.OtpExpire) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "OTP has expired",
				"success": false,
			})
			return
		}

		if len(input.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Password must be at least 6 characters",
				"success": false,
			})
			return
		}

		if input.Password != input.ConfirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Password and Confirm Password must match",
				"success": false,
			})
			return
		}

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

		updateData := bson.M{
			"isVerified": true,
			"otp":        nil,
			"otpExpire":  nil,
			"updatedAt":  time.Now(),
			"password":   hashedPassword,
		}

		if err := queries.UpdateUser(foundUser.ID.Hex(), updateData); err != nil {
			log.Printf("Failed to verify account: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to Reset Password",
				"success": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Password reset successful",
			"success": true,
		})
	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email              string `json:"email" binding:"required,email"`
			OldPassword        string `json:"oldPassword" binding:"required" `
			NewPassword        string `json:"newPassword" binding:"required"`
			ConfirmNewPassword string `json:"confirmNewPassword" binding:"required"`
		}

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

		foundUser, err := queries.GetUserByEmail(email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "User account does not exist",
				"success": false,
			})
			return
		}

		passwordIsValid, _ := helpers.VerifyPassword(input.OldPassword, foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Old Password is incorrect",
				"success": false,
			})
			return
		}

		if len(input.NewPassword) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Password must be at least 6 characters",
				"success": false,
			})
			return
		}

		if input.NewPassword != input.ConfirmNewPassword {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "New Password and Confirm New Password must match",
				"success": false,
			})
			return
		}

		hashedPassword, err := helpers.HashPassword(input.NewPassword)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to process password",
				"success": false,
			})
			return
		}

		updateData := bson.M{
			"updatedAt": time.Now(),
			"password":  hashedPassword,
		}

		if err := queries.UpdateUser(foundUser.ID.Hex(), updateData); err != nil {
			log.Printf("Failed to verify account: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to Change Password",
				"success": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Password change successful",
			"success": true,
		})
	}
}
