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

	r.Use(middleware.CORSMiddleware(cfg))

	api := r.Group("/api/v1")

	api.GET(
		"/auth/google",
		handlers.GoogleLogin(cfg),
	)

	api.GET(
		"/auth/google/callback",
		handlers.GoogleCallback(cfg),
	)
	api.GET(
		"/auth/me",
		middleware.AuthMiddleware(cfg),
		handlers.Me(cfg),
	)

	// Domain endpoints: mock implementations mirroring the frontend
	// seed data until the real logic is built.

	// Authenticated domain endpoints.
	auth := api.Group("/", middleware.AuthMiddleware(cfg))

	// Team endpoints require both authentication and team membership.
	teams := auth.Group("/teams/:team_id")
	teams.Use(middleware.RequireTeamMembership())

	// Team collection endpoints do not require membership in a specific team.
	auth.GET("/teams", handlers.MockGetTeams)
	auth.POST("/teams", handlers.MockCreateTeam)

	// Team-specific endpoints.
	teams.GET("", handlers.MockGetTeamDetails)
	teams.GET("/members", handlers.MockGetTeamMembers)
	teams.GET("/achievement-types", handlers.MockGetAchievementTypes)
	teams.POST("/achievement-types", handlers.MockCreateAchievementType)
	teams.GET("/claims", handlers.MockGetClaims)
	teams.POST("/claims", handlers.MockCreateClaim)
	teams.GET("/debts", handlers.MockGetDebts)
	teams.POST("/settlements", handlers.MockCreateSettlement)

	// Other authenticated endpoints.
	auth.POST("/claims/:claim_id/votes", handlers.MockCastVote)
	auth.POST("/invites/:token/accept", handlers.MockAcceptInvite)

	// Unauthenticated per api.md: invite resolution happens pre-login.
	api.GET("/invites/:token", handlers.MockResolveInvite)

	log.Println("Server running on port", cfg.Server.Port)

	err = r.Run(":" + cfg.Server.Port)
	if err != nil {
		log.Fatal(err)
	}
}
