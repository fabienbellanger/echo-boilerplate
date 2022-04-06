package user

import (
	"github.com/fabienbellanger/echo-boilerplate/db"
	"github.com/fabienbellanger/echo-boilerplate/entities"
)

// UserStore ...
type UserStore struct {
	db *db.DB
}

// New returns a new UserStore
func New(db *db.DB) UserStore {
	return UserStore{db: db}
}

// Register creates a new user in database
func (u UserStore) Register(user entities.User) (entities.User, error) {
	// TODO:
	return user, nil
}
