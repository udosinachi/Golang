package auth

import (
	middlewares "udo-golang/internal/adapters/http/middleware"
	authService "udo-golang/internal/services/auth"
	userService "udo-golang/internal/services/user"

	"github.com/gin-gonic/gin"
)

type Facade struct {
	service     userService.Server
	middlewares middlewares.Middlewares
	auth        authService.Service
}

func NewFacade(service userService.Server, middlewares middlewares.Middlewares, auth authService.Service) *Facade {
	return &Facade{
		service,
		middlewares,
		auth,
	}
}

func (f Facade) Register(r *gin.RouterGroup) {

	r.POST("/auth/login", f.Login)
	r.POST("/auth/signup", f.Signup)
	// r.POST("/auth/forgot", f.ForgotPassword)
	// r.POST("/auth/reset", f.ResetPassword)

	// g := r.Group("/user")
	// r.Use(f.middlewares.AuthenticateUser)
	// g.GET("/me", f.GetUser)

}
