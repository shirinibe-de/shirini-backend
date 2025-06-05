package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shirinibe-de/shirini-backend/internal/domain"
	"github.com/shirinibe-de/shirini-backend/internal/repository"
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
