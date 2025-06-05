package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shirinibe-de/shirini-backend/internal/handler"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/health", handler.HealthCheck)

	// Auth
	api.Post("/auth/google", handler.GoogleLogin)

	// Team
	api.Post("/teams", handler.CreateTeam)
	api.Post("/join/:token", handler.JoinTeam)

	// Claim & Vote
	api.Post("/claims", handler.CreateClaim)
	api.Post("/votes", handler.VoteOnClaim)
}
