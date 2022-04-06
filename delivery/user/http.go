package user

import (
	"github.com/fabienbellanger/echo-boilerplate/store"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	store store.UserStorer
}

// New returns a new UserHandler
func New(user store.UserStorer) UserHandler {
	return UserHandler{store: user}
}

// Register a new user
func (u UserHandler) Register(c echo.Context) error {
	// TODO:
	return nil
}
