// File: internal/handler/team_handler.go
package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shirinibe-de/shirini-backend/internal/domain"
	"github.com/shirinibe-de/shirini-backend/internal/repository"
	"github.com/shirinibe-de/shirini-backend/pkg/db"
)

func CreateTeam(c *fiber.Ctx) error {
	type request struct {
		Name string `json:"name"`
	}
	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID := c.Locals("user_id").(string) // Assume middleware sets this

	// Create new team
	token := uuid.New().String()
	team := &domain.Team{
		ID:        uuid.New().String(),
		Name:      body.Name,
		JoinToken: token,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	teamRepo := repository.NewTeamRepository()
	ctx := context.Background()
	err := teamRepo.Create(ctx, team)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	// Add creator as member
	membershipRepo := repository.NewMembershipRepository()
	membership := &domain.Membership{
		UserID:   userID,
		TeamID:   team.ID,
		JoinedAt: time.Now(),
	}
	err = membershipRepo.Create(ctx, membership)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	return c.JSON(fiber.Map{"join_url": "/join/" + token})
}

func JoinTeam(c *fiber.Ctx) error {
	token := c.Params("token")
	teamRepo := repository.NewTeamRepository()
	ctx := context.Background()
	team, err := teamRepo.GetByJoinToken(ctx, token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}
	if team == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Team not found"})
	}

	userID := c.Locals("user_id").(string)
	membershipRepo := repository.NewMembershipRepository()
	membership := &domain.Membership{
		UserID:   userID,
		TeamID:   team.ID,
		JoinedAt: time.Now(),
	}
	err = membershipRepo.Create(ctx, membership)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	return c.JSON(fiber.Map{"message": "Joined team successfully"})
}

func ListTeams(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	d := db.GetPool()
	query := `
    SELECT t.id, t.name, t.join_token, t.created_by, t.created_at
    FROM teams t
    JOIN team_members tm ON t.id = tm.team_id
    WHERE tm.user_id = $1
    `
	rows, err := d.Query(context.Background(), query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}
	defer rows.Close()

	type TeamView struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		JoinToken string    `json:"join_token"`
		CreatedBy string    `json:"created_by"`
		CreatedAt time.Time `json:"created_at"`
	}

	var teams []TeamView
	for rows.Next() {
		var t TeamView
		if err := rows.Scan(&t.ID, &t.Name, &t.JoinToken, &t.CreatedBy, &t.CreatedAt); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB scan error"})
		}
		teams = append(teams, t)
	}

	return c.JSON(teams)
}
