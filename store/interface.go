package store

import "github.com/fabienbellanger/echo-boilerplate/entities"

// UserStorer interface
type UserStorer interface {
	Register(entities.User) (entities.User, error)
}
