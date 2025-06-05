package domain

import (
	"time"
)

type Vote struct {
	ID      string    `json:"id"`
	ClaimID string    `json:"claim_id"`
	VotedBy string    `json:"voted_by"`
	Vote    bool      `json:"vote"` // true = up, false = down
	VotedAt time.Time `json:"voted_at"`
}
