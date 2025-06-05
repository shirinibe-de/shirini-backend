package domain

import "time"

type Team struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	JoinToken string    `json:"join_token"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
