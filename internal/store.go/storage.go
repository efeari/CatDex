package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		create(ctx context.Context, tx *sql.Tx, user *User) error
		CreateAndInvite(ctx context.Context,
			user *User,
			token string,
			invitationExp time.Duration) error
		GetByID(ctx context.Context, uuid uuid.UUID) (*User, error)
		DeleteByID(ctx context.Context, uuid uuid.UUID) error
		UpdateByID(ctx context.Context, user *User) error
		Activate(ctx context.Context, token string) error
	}
	Cats interface {
		Create(ctx context.Context, cat *Cat) error
		GetByID(ctx context.Context, uuid uuid.UUID) (*Cat, error)
		DeleteByID(ctx context.Context, uuid uuid.UUID) error
		UpdateByID(ctx context.Context, cat *Cat) error
		GetGlobalFeed(ctx context.Context, fq PaginatedFeedQuery) ([]CatFeed, error)
	}
}

// Creates a new storage layer
func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UsersStore{db},
		Cats:  &CatsStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
