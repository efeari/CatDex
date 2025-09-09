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
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, uuid uuid.UUID) (*User, error)
		DeleteByID(ctx context.Context, uuid uuid.UUID) error
		UpdateByID(ctx context.Context, user *User) error
	}
	Cats interface {
		Create(ctx context.Context, cat *Cat) error
		GetByID(ctx context.Context, uuid uuid.UUID) (*Cat, error)
		DeleteByID(ctx context.Context, uuid uuid.UUID) error
		UpdateByID(ctx context.Context, cat *Cat) error
		GetGlobalFeed(ctx context.Context) ([]CatFeed, error)
	}
}

// Creates a new storage layer
func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UsersStore{db},
		Cats:  &CatsStore{db},
	}
}
