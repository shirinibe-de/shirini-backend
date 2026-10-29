package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrTeamNotFound = errors.New("team not found")
var ErrNotTeamMember = errors.New("user is not a team member")

type Team struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Currency    string    `json:"currency"`
	InviteToken string    `json:"invite_token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	MemberCount int       `json:"member_count,omitempty"`
}

func generateInviteToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate invite token: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func ListTeamsForUser(userID int) ([]Team, error) {
	rows, err := DB.Query(`
		SELECT t.id, t.name, t.currency, t.created_at
		FROM teams t
		JOIN team_memberships tm ON tm.team_id = t.id
		WHERE tm.user_id = $1
		ORDER BY t.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}
	defer rows.Close()

	teams := make([]Team, 0)

	for rows.Next() {
		var team Team

		if err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.Currency,
			&team.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan team: %w", err)
		}

		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate teams: %w", err)
	}

	return teams, nil
}

func CreateTeam(name, currency string, userID int) (*Team, error) {
	name = strings.TrimSpace(name)
	currency = strings.TrimSpace(currency)

	if name == "" {
		return nil, errors.New("team name is required")
	}

	if currency == "" {
		return nil, errors.New("team currency is required")
	}

	inviteToken, err := generateInviteToken()
	if err != nil {
		return nil, err
	}

	tx, err := DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin create team transaction: %w", err)
	}
	defer tx.Rollback()

	team := &Team{}

	err = tx.QueryRow(`
		INSERT INTO teams (name, currency, invite_token, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, currency, invite_token, created_at
	`, name, currency, inviteToken, userID).Scan(
		&team.ID,
		&team.Name,
		&team.Currency,
		&team.InviteToken,
		&team.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert team: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO team_memberships (team_id, user_id)
		VALUES ($1, $2)
	`, team.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("add creator as team member: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create team transaction: %w", err)
	}

	return team, nil
}

func GetTeamDetailsForUser(teamID string, userID int) (*Team, error) {
	var team Team

	err := DB.QueryRow(`
		SELECT id, name, currency, invite_token, created_at
		FROM teams
		WHERE id = $1
	`, teamID).Scan(
		&team.ID,
		&team.Name,
		&team.Currency,
		&team.InviteToken,
		&team.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTeamNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get team: %w", err)
	}

	var isMember bool

	err = DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM team_memberships
			WHERE team_id = $1 AND user_id = $2
		)
	`, teamID, userID).Scan(&isMember)
	if err != nil {
		return nil, fmt.Errorf("check team membership: %w", err)
	}

	if !isMember {
		return nil, ErrNotTeamMember
	}

	err = DB.QueryRow(`
		SELECT COUNT(*)
		FROM team_memberships
		WHERE team_id = $1
	`, teamID).Scan(&team.MemberCount)
	if err != nil {
		return nil, fmt.Errorf("count team members: %w", err)
	}

	return &team, nil
}

type TeamMember struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	JoinedAt  time.Time `json:"joined_at"`
}

func GetTeamMembersForUser(teamID string, userID int) ([]TeamMember, error) {
	// First, verify that the requesting user belongs to this team.
	var isMember bool
	err := DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM team_memberships
			WHERE team_id = $1 AND user_id = $2
		)
	`, teamID, userID).Scan(&isMember)

	if err != nil {
		return nil, fmt.Errorf("check team membership: %w", err)
	}
	if !isMember {
		return nil, ErrNotTeamMember
	}

	rows, err := DB.Query(`
		SELECT u.id, u.email, COALESCE(u.name, ''),
		       COALESCE(u.avatar_url, ''), tm.created_at
		FROM team_memberships tm
		JOIN users u ON u.id = tm.user_id
		WHERE tm.team_id = $1
		ORDER BY tm.created_at ASC
	`, teamID)
	if err != nil {
		return nil, fmt.Errorf("get team members: %w", err)
	}
	defer rows.Close()

	members := make([]TeamMember, 0)
	for rows.Next() {
		var member TeamMember
		if err := rows.Scan(
			&member.ID,
			&member.Email,
			&member.Name,
			&member.AvatarURL,
			&member.JoinedAt,
		); err != nil {
			return nil, fmt.Errorf("scan team member: %w", err)
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate team members: %w", err)
	}

	return members, nil
}
