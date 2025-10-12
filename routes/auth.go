package routes

import (
	"udo-golang/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("auth/signup", controllers.Signup())
	incomingRoutes.POST("auth/login", controllers.Login())
}
