package routes

import (
	"udo-golang/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("auth/signup", controllers.Signup())
	incomingRoutes.POST("auth/register", controllers.RegisterWithOtp())
	incomingRoutes.GET("auth/google/callback", controllers.GoogleSignUpandSignIn())

	incomingRoutes.POST("auth/verify-account", controllers.VerifyAccount())
	incomingRoutes.POST("auth/resend-otp", controllers.ResendOtp())
	incomingRoutes.POST("auth/login", controllers.Login())
}
