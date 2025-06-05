package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shirinibe-de/shirini-backend/internal/domain"
	"github.com/shirinibe-de/shirini-backend/internal/repository"
)

func CreateClaim(c *fiber.Ctx) error {
	type request struct {
		TeamID     string `json:"team_id"`
		ClaimedFor string `json:"claimed_for"`
		Message    string `json:"message"`
	}
	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID := c.Locals("user_id").(string)
	claim := &domain.Claim{
		ID:         uuid.New().String(),
		TeamID:     body.TeamID,
		ClaimedBy:  userID,
		ClaimedFor: body.ClaimedFor,
		Message:    body.Message,
		Status:     domain.Pending,
		CreatedAt:  time.Now(),
	}

	claimRepo := repository.NewClaimRepository()
	ctx := context.Background()
	err := claimRepo.Create(ctx, claim)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	return c.JSON(fiber.Map{"claim_id": claim.ID})
}

func VoteOnClaim(c *fiber.Ctx) error {
	type request struct {
		ClaimID string `json:"claim_id"`
		Vote    bool   `json:"vote"`
	}
	var body request
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID := c.Locals("user_id").(string)
	vote := &domain.Vote{
		ID:      uuid.New().String(),
		ClaimID: body.ClaimID,
		VotedBy: userID,
		Vote:    body.Vote,
		VotedAt: time.Now(),
	}

	voteRepo := repository.NewVoteRepository()
	ctx := context.Background()
	err := voteRepo.Create(ctx, vote)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	// Optionally: check majority and update status (omitted for brevity)

	return c.JSON(fiber.Map{"message": "Vote recorded"})
}
