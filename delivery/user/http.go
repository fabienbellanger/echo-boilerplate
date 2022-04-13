package user

import (
	"net/http"

	"github.com/fabienbellanger/echo-boilerplate/entities"
	"github.com/fabienbellanger/echo-boilerplate/store"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type UserHandler struct {
	group *echo.Group
	store store.UserStorer
}

// New returns a new UserHandler
func New(g *echo.Group, user store.UserStorer) UserHandler {
	return UserHandler{
		group: g,
		store: user,
	}
}

// Routes adds users routes
func (u *UserHandler) Routes() {
	routes := []struct {
		Method  string
		Path    string
		Handler echo.HandlerFunc
	}{
		{"POST", "", u.register()},
	}

	for _, route := range routes {
		switch route.Method {
		case "GET":
			u.group.GET(route.Path, route.Handler)
		case "POST":
			u.group.POST(route.Path, route.Handler)
		}
	}
}

// Login route
func (u UserHandler) Login(c echo.Context) error {
	claims := entities.NewClaims("ID", "Username", "Lastname", "Firstname", viper.GetInt("JWT_LIFETIME"))
	token, err := claims.GenerateJWT(viper.GetString("JWT_ALGO"), viper.GetString("JWT_SECRET"))
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Login route, token: "+token)
}

// Register a new user
func (u UserHandler) register() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO:
		return c.String(http.StatusOK, "Register route")
	}
}
