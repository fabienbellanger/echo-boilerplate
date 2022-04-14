package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/fabienbellanger/echo-boilerplate/entities"
	"github.com/fabienbellanger/echo-boilerplate/store"
	"github.com/fabienbellanger/echo-boilerplate/utils"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type userLogin struct {
	entities.User
	Token     string `json:"token" xml:"token" form:"token"`
	ExpiresAt string `json:"expires_at" xml:"expires_at" form:"expires_at"`
}

type userAuth struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
}

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
	u.group.POST("", u.register())
}

// Login route
func (u *UserHandler) Login(c echo.Context) error {
	ua := new(userAuth)
	if err := c.Bind(ua); err != nil {
		return err
	}

	loginErrors := utils.ValidateStruct(*ua)
	if loginErrors != nil {
		return echo.NewHTTPError(http.StatusBadRequest, loginErrors)
	}

	user, err := u.store.Login(ua.Username, ua.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, nil)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Error during authentication")
	}

	claims := entities.NewClaims(user.ID, ua.Username, user.Lastname, user.Firstname, viper.GetInt("JWT_LIFETIME"))
	token, err := claims.GenerateJWT(viper.GetString("JWT_ALGO"), viper.GetString("JWT_SECRET"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, userLogin{
		User:      user,
		Token:     token,
		ExpiresAt: time.Unix(claims.ExpiresAt, 0).Format("2006-01-02T15:04:05.000Z"),
	})
}

// Register a new user
func (u UserHandler) register() echo.HandlerFunc {
	return func(c echo.Context) error {
		uf := new(entities.UserForm)
		if err := c.Bind(uf); err != nil {
			return err
		}

		if uf.Firstname == "" || uf.Lastname == "" || uf.Username == "" || uf.Password == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad Parameters")
		}

		user := entities.User{
			Lastname:  uf.Lastname,
			Firstname: uf.Firstname,
			Password:  uf.Password,
			Username:  uf.Username,
		}

		if err := u.store.Register(&user); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error during user creation")
		}

		return c.JSON(http.StatusOK, user)
	}
}
