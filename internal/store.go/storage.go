package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Storage struct {
	Users interface {
		Create(ctx context.Context, user *User) error
		GetByID(ctx context.Context, uuid uuid.UUID) (*User, error)
	}
	Cats interface {
		Create(ctx context.Context, cat *Cat) error
		GetByID(ctx context.Context, uuid uuid.UUID) (*Cat, error)
	}
}

// Creates a new storage layer
func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UsersStore{db},
		Cats:  &CatsStore{db},
	}
}
