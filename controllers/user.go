package controllers

import (
	"fmt"
	"net/http"
	"time"
	"udo-golang/queries"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		allUsers, err := queries.GetAllUsers()
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"success": false,
				"message": "Unable to Fetch Users",
				"error":   err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": true,
			"message": "Users Fetched Successfully",
			"data":    allUsers,
		})
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, err := queries.GetUserByID(id)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"success": false,
				"message": "Unable to Fetch this user",
			})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"success": false,
				"message": "User does not exist",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": true,
			"message": "User Fetched Successfully",
			"data":    user,
		})
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := queries.DeleteUserById(id)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"success": false,
				"message": "Unable to delete this User",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": true,
			"message": "User Deleted Successfully",
		})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, getUserErr := queries.GetUserByID(id)
		if getUserErr != nil {
			fmt.Println(getUserErr)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"success": false,
				"message": "User does not exist",
			})
			return
		}

		update := bson.M{
			"isVerified": true,
			"otp":        nil,
			"otpExpire":  nil,
			"updatedAt":  time.Now(),
		}

		err := queries.UpdateUser(id, update)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"success": false,
				"message": "Unable to Update this User",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": true,
			"message": "User Deleted Successfully",
		})
	}
}
