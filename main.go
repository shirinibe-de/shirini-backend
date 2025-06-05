// cmd/main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/shirinibe-de/shirini-backend/config"
	"github.com/shirinibe-de/shirini-backend/internal/router"
	"github.com/shirinibe-de/shirini-backend/pkg/db"
)

func main() {
	// Load env variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize DB
	err = db.Init(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Fiber app
	app := fiber.New()

	// Setup routes
	router.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":8080"))
}
