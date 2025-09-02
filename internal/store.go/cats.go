package store

import (
	"context"
	"database/sql"
)

type Cat struct {
	ID            int64  `json:"id" db:"id"`
	Name          string `json:"name" db:"name"` // user_given_name if not a stray cat
	UserGivenName string `json:"user_given_name" db:"user_given_name"`
	Description   string `json:"description" db:"description"`
	Location      string `json:"location" db:"location"`
	PhotoPath     string `json:"-" db:"photo_path"` // Hidden from JSON
	PhotoURL      string `json:"photo_url" db:"-"`  // Computed field, not in DB
	UserID        int64  `json:"user_id" db:"user_id"`
	CreatedAt     string `json:"created_at" db:"created_at"`
	LastSeen      string `json:"last_seen" db:"last_seen"`
}

type CatsStore struct {
	db *sql.DB
}

func (s *CatsStore) Create(ctx context.Context, cat *Cat) error {
	query := `
	INSERT INTO cats (name, description, location, photo_path, user_id)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, last_seen
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		cat.Name,
		cat.Description,
		cat.Location,
		cat.PhotoPath,
		cat.UserID,
	).Scan(
		&cat.ID,
		&cat.CreatedAt,
		&cat.LastSeen,
	)
	if err != nil {
		return err
	}
	return nil
}
