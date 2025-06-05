// File: internal/handler/claim_handler.go
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
	if err := claimRepo.Create(ctx, claim); err != nil {
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
	if err := voteRepo.Create(ctx, vote); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}

	return c.JSON(fiber.Map{"message": "Vote recorded"})
}

func ListClaims(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	d := db.GetPool()
	query := `
    SELECT
      ac.id,
      ac.team_id,
      ac.claimed_by,
      cb.name AS claimed_by_name,
      ac.claimed_for,
      cf.name AS claimed_for_name,
      ac.message
    FROM achievement_claims ac
    JOIN team_members tm ON ac.team_id = tm.team_id
    JOIN users cb ON ac.claimed_by = cb.id
    JOIN users cf ON ac.claimed_for = cf.id
    WHERE tm.user_id = $1 AND ac.status = 'pending'
    `
	rows, err := d.Query(context.Background(), query, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}
	defer rows.Close()

	type ClaimView struct {
		ID             string `json:"id"`
		TeamID         string `json:"team_id"`
		ClaimedBy      string `json:"claimed_by"`
		ClaimedByName  string `json:"claimed_by_name"`
		ClaimedFor     string `json:"claimed_for"`
		ClaimedForName string `json:"claimed_for_name"`
		Message        string `json:"message"`
	}

	var claims []ClaimView
	for rows.Next() {
		var cView ClaimView
		if err := rows.Scan(
			&cView.ID,
			&cView.TeamID,
			&cView.ClaimedBy,
			&cView.ClaimedByName,
			&cView.ClaimedFor,
			&cView.ClaimedForName,
			&cView.Message,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB scan error"})
		}
		claims = append(claims, cView)
	}

	return c.JSON(claims)
}
