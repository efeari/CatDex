package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Users interface {
		Create(ctx context.Context, user *User) error
	}
	Cats interface {
		Create(ctx context.Context, cat *Cat) error
	}
}

// Creates a new storage layer
func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UsersStore{db},
		Cats:  &CatsStore{db},
	}
}
