package main

import (
	"log"

	"shirini-backend/database"
	"shirini-backend/handlers"
	"shirini-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found")
	}

	database.Connect()
	database.RunMigrations()
	r := gin.Default()

	r.GET("/auth/google", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback)
	r.GET("/api/v1/me", middleware.AuthMiddleware(), handlers.Me)

	r.Run(":8080")
}
