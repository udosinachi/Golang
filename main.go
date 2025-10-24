package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	http "udo-golang/internal/adapters/http"
	mongoadpt "udo-golang/internal/adapters/mongo"
	userRepo "udo-golang/internal/adapters/mongo/repositories/user"
	authService "udo-golang/internal/services/auth"
	userService "udo-golang/internal/services/user"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	godotenv.Load(".env")
	uri := os.Getenv("MONGODB_ATLAS_URI")
	if uri == "" {
		log.Fatal("MONGODB_ATLAS_URI not found in environment")
	}

	_, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal(err, "err plenty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	mc, err := mongoadpt.Connect(ctx, uri, "test")
	if err != nil {
		log.Fatal("mongo connect:", err)
	}

	UserRepository := userRepo.NewUserRepository(mc.DB)
	UserService := userService.NewServer(UserRepository)
	AuthService := authService.NewService(UserRepository, "rwewfyuieowoo")

	httpServer := http.NewServer(UserService, AuthService)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Run(); err != nil {
			log.Panicln(err)
		}
	}()

	log.Println("Server up")
	wg.Wait()

	// if err := godotenv.Load(".env"); err != nil {
	// 	log.Println("Warning: .env file not found, using system environment variables.")
	// }

	// gin.SetMode(gin.ReleaseMode)

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	log.Fatal("PORT is not set in the environment")
	// }

	// router := gin.Default()
	// router.Use(middleware.CORSMiddleware())

	// // Public Routes
	// routes.AuthRoutes(router)

	// // Private Routes
	// routes.UserRoutes(router)

	// fmt.Println("🚀 Server is running on port:", port)

	// if err := router.Run(":" + port); err != nil {
	// 	log.Fatal("Failed to start server:", err)
	// }
}
