package main

import (
	"fmt"
	"log"
	"os"
	"udo-golang/middleware"
	"udo-golang/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	gin.SetMode(gin.ReleaseMode)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment")
	}

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Register routes
	routes.AuthRoutes(router)

	fmt.Println("🚀 Server is running on port:", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
