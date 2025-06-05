package domain

import "time"

type Membership struct {
	UserID   string    `json:"user_id"`
	TeamID   string    `json:"team_id"`
	JoinedAt time.Time `json:"joined_at"`
}
