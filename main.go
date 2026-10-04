package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"shirini-backend/config"
	"shirini-backend/database"
	"shirini-backend/handlers"
	"shirini-backend/middleware"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	database.Connect(cfg.Database.URL)

	database.RunMigrations()

	r := gin.Default()

	r.GET(
		"/auth/google",
		handlers.GoogleLogin(cfg),
	)

	r.GET(
		"/auth/google/callback",
		handlers.GoogleCallback(cfg),
	)

	r.GET(
		"/api/v1/me",
		middleware.AuthMiddleware(cfg),
		handlers.Me,
	)

	log.Println("Server running on port", cfg.Server.Port)

	err = r.Run(":" + cfg.Server.Port)

	if err != nil {
		log.Fatal(err)
	}
}