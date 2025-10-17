package routes

import (
	"udo-golang/controllers"
	"udo-golang/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("auth/signup", controllers.Signup())
	incomingRoutes.POST("auth/register", controllers.RegisterWithOtp())
	incomingRoutes.GET("auth/google/callback", controllers.GoogleSignUpandSignIn())

	incomingRoutes.POST("auth/verify-account", controllers.VerifyAccount())
	incomingRoutes.POST("auth/resend-otp", controllers.SendOtp())
	incomingRoutes.POST("auth/login", controllers.Login())
	incomingRoutes.POST("auth/send-reset-otp", controllers.SendOtp())
	incomingRoutes.POST("auth/reset-password", controllers.ResetPassword())
	incomingRoutes.POST("auth/change-password", middleware.IsAuthenticated(), controllers.ChangePassword())
}
