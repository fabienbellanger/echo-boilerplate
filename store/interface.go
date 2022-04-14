package store

import "github.com/fabienbellanger/echo-boilerplate/entities"

// UserStorer interface
type UserStorer interface {
	Login(username, password string) (entities.User, error)
	Register(user *entities.User) error
	GetAllUsers() ([]entities.User, error)
	GetUser(id string) (entities.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, userForm *entities.UserForm) (entities.User, error)
}
