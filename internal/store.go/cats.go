package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/efeari/catdex/internal/utils"
	"github.com/google/uuid"
)

type Cat struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"` // user_given_name if not a stray cat
	UserGivenName string    `json:"user_given_name" db:"user_given_name"`
	Description   string    `json:"description" db:"description"`
	Location      string    `json:"location" db:"location"`
	PhotoPath     string    `json:"-" db:"photo_path"` // Hidden from JSON
	PhotoURL      string    `json:"photo_url" db:"-"`  // Computed field, not in DB
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt     string    `json:"created_at" db:"created_at"`
	LastSeen      string    `json:"last_seen" db:"last_seen"`
	Version       int       `json:"version" db:"version"`
}

type CatFeed struct {
	Cat
	UserName string
}

type CatsStore struct {
	db *sql.DB
}

func (s *CatsStore) GetGlobalFeed(ctx context.Context, fq PaginatedFeedQuery) ([]CatFeed, error) {
	query :=
		`
	SELECT 
    c.id,
    c.name,
    c.description,
    c.location,
    c.photo_path,
    c.user_id,
    c.created_at,
    c.last_seen,
    c.version,
    u.username
	FROM cats c
	JOIN users u ON c.user_id = u.id
	WHERE 
    ($1 = '' OR c.name ILIKE '%' || $1 || '%')
    AND ($2 = '' OR u.username ILIKE '%' || $2 || '%')
    AND ($3 = '' OR c.description ILIKE '%' || $3 || '%')
    AND ($4 = '' OR c.location ILIKE '%' || $4 || '%')
	ORDER BY c.created_at ` + fq.Sort + `
	LIMIT $5 OFFSET $6;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		fq.Name,
		fq.Username,
		fq.Search,
		fq.Location,
		fq.Limit,
		fq.Offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// err = godotenv.Load("../../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// addr := os.Getenv("ADDR")

	var feed []CatFeed
	for rows.Next() {
		var c CatFeed
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.Location,
			&c.PhotoPath,
			&c.UserID,
			&c.CreatedAt,
			&c.LastSeen,
			&c.Version,
			&c.UserName,
		)
		if err != nil {
			return nil, err
		}

		//c.PhotoURL = fmt.Sprintf("%s/v1/photos/%s", addr, c.PhotoPath)

		feed = append(feed, c)
	}

	return feed, nil
}

func (s *CatsStore) UpdateByID(ctx context.Context, cat *Cat) error {
	query := `
	UPDATE cats
	SET name = $1,
    description = $2,
    location = $3,
    last_seen = $4,
	version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		cat.Name,
		cat.Description,
		cat.Location,
		cat.LastSeen,
		cat.ID,
		cat.Version).Scan(&cat.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *CatsStore) DeleteByID(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM cats WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, uuid)

	if err != nil {
		return err
	}

	return utils.DeleteCatPhoto(uuid.String())
}

func (s *CatsStore) GetByID(ctx context.Context, uuid uuid.UUID) (*Cat, error) {
	query := `
	SELECT * FROM cats WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, uuid)

	var cat Cat
	err := row.Scan(
		&cat.ID,
		&cat.Name,
		&cat.Description,
		&cat.Location,
		&cat.PhotoPath,
		&cat.UserID,
		&cat.CreatedAt,
		&cat.LastSeen,
		&cat.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &cat, nil
}

func (s *CatsStore) Create(ctx context.Context, cat *Cat) error {
	query := `
	INSERT INTO cats (id, name, description, location, photo_path, user_id)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at, last_seen
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		cat.ID,
		cat.Name,
		cat.Description,
		cat.Location,
		cat.PhotoPath,
		cat.UserID,
	).Scan(
		&cat.CreatedAt,
		&cat.LastSeen,
	)
	if err != nil {
		return err
	}
	return nil
}
