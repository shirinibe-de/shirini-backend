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
	auth := api.Group("/", middleware.AuthMiddleware(cfg))

	auth.GET("/teams", handlers.MockGetTeams)
	auth.POST("/teams", handlers.MockCreateTeam)
	auth.GET("/teams/:team_id", handlers.MockGetTeamDetails)
	auth.GET("/teams/:team_id/members", handlers.MockGetTeamMembers)
	auth.GET("/teams/:team_id/achievement-types", handlers.MockGetAchievementTypes)
	auth.POST("/teams/:team_id/achievement-types", handlers.MockCreateAchievementType)
	auth.GET("/teams/:team_id/claims", handlers.MockGetClaims)
	auth.POST("/teams/:team_id/claims", handlers.MockCreateClaim)
	auth.GET("/teams/:team_id/debts", handlers.MockGetDebts)
	auth.POST("/teams/:team_id/settlements", handlers.MockCreateSettlement)
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
