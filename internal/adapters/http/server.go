package http

import (
	"fmt"
	"net/http"

	"udo-golang/internal/adapters/http/common"
	"udo-golang/internal/adapters/http/facades/auth"
	middlewares "udo-golang/internal/adapters/http/middleware"
	authService "udo-golang/internal/services/auth"
	userService "udo-golang/internal/services/user"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port   int
	router *gin.Engine
}

func (s *Server) Handler() http.Handler {
	return s.router.Handler()
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}

func NewServer(userService userService.Server, authService authService.Service) *Server {
	router := gin.Default()

	middlewares := middlewares.NewMiddleware(userService)

	router.Use(middlewares.Cors)

	router.NoRoute(func(ctx *gin.Context) {
		common.SendNotFound(ctx, "The resource you were looking for does not exist")
	})

	facades := []Facade{
		auth.NewFacade(userService, *middlewares, authService),
	}

	for _, f := range facades {
		f.Register(&router.RouterGroup)
	}

	return &Server{
		8080,
		router,
	}
}
