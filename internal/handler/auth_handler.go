package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/shirinibe-de/shirini-backend/internal/domain"
	"github.com/shirinibe-de/shirini-backend/internal/repository"
)

func GoogleLogin(c *fiber.Ctx) error {
	// Normally, you would extract and verify the Google ID token here.
	// For now, we'll skip token verification and use dummy user info.

	// Suppose we parsed user info from a verified token:
	email := "user@example.com"
	name := "Google User"
	avatar := "https://example.com/avatar.jpg"

	userRepo := repository.NewUserRepository()
	ctx := context.Background()
	// Check if user exists
	existing, err := userRepo.GetByEmail(ctx, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
	}
	if existing == nil {
		// Create new user
		newUser := &domain.User{
			ID:        uuid.New().String(),
			Email:     email,
			Name:      name,
			AvatarURL: avatar,
			CreatedAt: time.Now(),
		}
		err = userRepo.Create(ctx, newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB error"})
		}
		existing = newUser
	}

	// Generate session/jwt (omitted for brevity)
	return c.JSON(fiber.Map{"user": existing})
}
