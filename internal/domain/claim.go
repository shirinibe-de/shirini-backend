package domain

import (
	"time"
)

type ClaimStatus string

const (
	Pending  ClaimStatus = "pending"
	Approved ClaimStatus = "approved"
	Rejected ClaimStatus = "rejected"
)

type Claim struct {
	ID         string      `json:"id"`
	TeamID     string      `json:"team_id"`
	ClaimedBy  string      `json:"claimed_by"`
	ClaimedFor string      `json:"claimed_for"`
	Message    string      `json:"message"`
	Status     ClaimStatus `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
}
