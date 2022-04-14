package store

import "github.com/fabienbellanger/echo-boilerplate/entities"

// UserStorer interface
type UserStorer interface {
	Login(username, password string) (entities.User, error)
	Register(user *entities.User) error
}
