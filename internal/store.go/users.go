package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username" binding:"required,min=3,max=25"`
	Email     string    `json:"email" db:"email" binding:"required,email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt string    `json:"created_at" db:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) UpdateByID(ctx context.Context, user *User) error {
	return nil
}

func (s *UsersStore) DeleteByID(ctx context.Context, uuid uuid.UUID) error {
	return nil
}

func (s *UsersStore) GetByID(ctx context.Context, uuid uuid.UUID) (*User, error) {
	query := `
	SELECT * FROM users WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO cats (username, email, password)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
