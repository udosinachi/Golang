package middlewares

import (
	"strings"
	"udo-golang/internal/adapters/http/common"
	"udo-golang/internal/common/token"
	userService "udo-golang/internal/services/user"

	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	user userService.Server
}

func NewMiddleware(user userService.Server) *Middlewares { return &Middlewares{user} }

func (m *Middlewares) Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, token, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func (m *Middlewares) AuthenticateUser(c *gin.Context) {
	clientToken := c.Request.Header.Get("Authorization")
	if clientToken == "" {
		common.SendUnauthorized(c, "Token not provided")
		c.Abort()
		return
	}

	updatedToken := clientToken

	if strings.HasPrefix(clientToken, "Bearer") {
		updatedToken = strings.Split(clientToken, " ")[1]
	}

	claims, err := token.ValidateToken(updatedToken)
	if err != "" {
		common.SendUnauthorized(c, err)

		c.Abort()
		return
	}

	if !claims.IsAdmin {
		common.SendUnauthorized(c, "You don't have the permission to access this data")
		c.Abort()
		return
	}

	c.Set("email", claims.Email)
	c.Set("uid", claims.ID)
	c.Set("isAdmin", claims.IsAdmin)

	c.Next()
}
