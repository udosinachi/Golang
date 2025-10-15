package controllers

import (
	"fmt"
	"net/http"
	"udo-golang/queries"

	"github.com/gin-gonic/gin"
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
