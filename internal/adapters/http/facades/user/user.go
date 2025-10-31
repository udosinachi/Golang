package user

import (
	middlewares "udo-golang/internal/adapters/http/middleware"
	userService "udo-golang/internal/services/user"

	"github.com/gin-gonic/gin"
)

type Facade struct {
	service     userService.Server
	middlewares middlewares.Middlewares
}

func NewFacade(service userService.Server, middlewares middlewares.Middlewares) *Facade {
	return &Facade{
		service,
		middlewares,
	}
}

func (f Facade) Register(r *gin.RouterGroup) {
	r.GET("/users", f.GetAllUsers)
	r.GET("/users/:id", f.GetUsersById)
	r.DELETE("/delete-user/:id", f.DeleteUserById)
	r.PUT("/update-user/:id", f.UpdateUserById)
}
