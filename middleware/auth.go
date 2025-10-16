package middleware

import (
	"net/http"
	"strings"
	"udo-golang/helpers"

	"github.com/gin-gonic/gin"
)

func IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"success": false,
				"message": "No Authorization header provided",
			})
			c.Abort()
			return
		}

		updatedToken := clientToken

		if strings.HasPrefix(clientToken, "Bearer") {
			updatedToken = strings.Split(clientToken, " ")[1]
		}

		claims, err := helpers.ValidateToken(updatedToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"success": false,
				"message": err,
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("id", claims.ID)
		c.Set("isAdmin", claims.IsAdmin)

		c.Next()
	}
}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"success": false,
				"message": "No Authorization header provided",
			})
			c.Abort()
			return
		}

		updatedToken := clientToken

		if strings.HasPrefix(clientToken, "Bearer") {
			updatedToken = strings.Split(clientToken, " ")[1]
		}

		claims, err := helpers.ValidateToken(updatedToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"success": false,
				"message": err,
			})
			c.Abort()
			return
		}

		if !claims.IsAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"success": false,
				"message": "You don't have the permission to access this data",
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("id", claims.ID)
		c.Set("isAdmin", claims.IsAdmin)

		c.Next()
	}
}
