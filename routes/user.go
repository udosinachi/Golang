package routes

import (
	"udo-golang/controllers"
	"udo-golang/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("users", middleware.IsAdmin(), controllers.GetAllUsers())
	incomingRoutes.GET("users/:id", controllers.GetUser())
	incomingRoutes.DELETE("delete-user/:id", controllers.DeleteUser())
	incomingRoutes.PUT("update-user/:id", controllers.UpdateUser())
}
