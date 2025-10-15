package routes

import (
	"udo-golang/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("users", controllers.GetAllUsers())
}
