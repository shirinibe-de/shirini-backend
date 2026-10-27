package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrTeamNotFound = errors.New("team not found")

type Team struct {
	ID        string
	Name      string
	CreatedBy int
}

type TeamMembership struct {
	ID     string
	TeamID string
	UserID int
}

// FindTeamByID returns nil, nil when the team does not exist.
func FindTeamByID(teamID string) (*Team, error) {
	team := &Team{}

	err := DB.QueryRow(`
		SELECT id, name, created_by
		FROM teams
		WHERE id = $1
	`, teamID).Scan(&team.ID, &team.Name, &team.CreatedBy)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return team, nil
}

// FindTeamMembership returns nil, nil when the user is not a member.
func FindTeamMembership(teamID string, userID int) (*TeamMembership, error) {
	membership := &TeamMembership{}

	err := DB.QueryRow(`
		SELECT id, team_id, user_id
		FROM team_memberships
		WHERE team_id = $1 AND user_id = $2
	`, teamID, userID).Scan(
		&membership.ID,
		&membership.TeamID,
		&membership.UserID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return membership, nil
}

func CreateTeam(name, currency, inviteToken string, creatorID int) (*Team, error) {
	tx, err := DB.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("begin team transaction: %w", err)
	}
	defer tx.Rollback()

	team := &Team{}

	err = tx.QueryRow(`
        INSERT INTO teams (name, currency, invite_token, created_by)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, created_by
    `, name, currency, inviteToken, creatorID).Scan(
		&team.ID,
		&team.Name,
		&team.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("create team: %w", err)
	}

	_, err = tx.Exec(`
        INSERT INTO team_memberships (team_id, user_id)
        VALUES ($1, $2)
    `, team.ID, creatorID)
	if err != nil {
		return nil, fmt.Errorf("add creator membership: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit team transaction: %w", err)
	}

	return team, nil
}
