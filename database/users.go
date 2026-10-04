package database

import (
	"database/sql"
	"errors"
)

type User struct {
	ID        int
	GoogleID  string
	Email     string
	Name      string
	AvatarURL string
}

func FindUserByGoogleID(googleID string) (*User, error) {

	row := DB.QueryRow(
		`
		SELECT id, google_id, email, name, avatar_url
		FROM users
		WHERE google_id = $1
		`,
		googleID,
	)

	user := &User{}

	err := row.Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&user.AvatarURL,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}


func CreateUser(user *User) (*User, error) {

	err := DB.QueryRow(
		`
		INSERT INTO users
		(
			google_id,
			email,
			name,
			avatar_url
		)
		VALUES
		($1,$2,$3,$4)
		RETURNING id
		`,
		user.GoogleID,
		user.Email,
		user.Name,
		user.AvatarURL,
	).Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	return user, nil
}