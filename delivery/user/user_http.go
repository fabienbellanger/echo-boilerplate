package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/fabienbellanger/echo-boilerplate/entities"
	"github.com/fabienbellanger/echo-boilerplate/store"
	"github.com/fabienbellanger/echo-boilerplate/utils"
	"github.com/google/uuid"
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
	u.group.GET("", u.getAll())
	u.group.GET("/stream", u.stream())
	u.group.GET("/:id", u.getOne())
	u.group.PUT("/:id", u.update())
	u.group.DELETE("/:id", u.delete())
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

// register creates a new user
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

// getAll lists all users
func (u UserHandler) getAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		users, err := u.store.GetAllUsers()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, users)
	}
}

// getOne returns the user
func (u UserHandler) getOne() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad ID")
		}

		user, err := u.store.GetUser(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error when retrieving user")
		}
		if user.ID == "" {
			return echo.NewHTTPError(http.StatusNotFound, "No user found")
		}

		return c.JSON(http.StatusOK, user)
	}
}

// delete deletes the user
func (u UserHandler) delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad ID")
		}

		err := u.store.DeleteUser(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error when deleting user")
		}

		return c.NoContent(http.StatusOK)
	}
}

// update updates user information.
func (u UserHandler) update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad ID")
		}

		user := new(entities.UserForm)
		if err := c.Bind(user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Bad data")
		}

		updateErrors := utils.ValidateStruct(*user)
		if updateErrors != nil {
			return echo.NewHTTPError(http.StatusBadRequest, updateErrors)
		}

		updatedUser, err := u.store.UpdateUser(id, user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error when updating user")
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

// stream is an example of users list sent with a stream
func (u UserHandler) stream() echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := c.Response()
		resp.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		resp.WriteHeader(http.StatusOK)

		resp.Write([]byte("["))
		enc := json.NewEncoder(resp)
		n := 100_000
		for i := 0; i < n; i++ {
			user := entities.User{
				ID:        uuid.New().String(),
				Username:  "My Username",
				Password:  ",kkjkjkjkjknnqfjkkjdnfsjklqblk",
				Lastname:  "My Lastname",
				Firstname: "My Firstname",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := enc.Encode(user); err != nil {
				continue
			}

			if i < n-1 {
				resp.Write([]byte(","))
			}

			resp.Flush()
		}
		resp.Write([]byte("]"))

		return nil
	}
}
