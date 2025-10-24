package http

import "github.com/gin-gonic/gin"

type Facade interface {
	Register(r *gin.RouterGroup)
}
