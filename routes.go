package server

import (
	"net/http"

	"github.com/fabienbellanger/echo-boilerplate/db"
	"github.com/fabienbellanger/echo-boilerplate/delivery/user"
	"github.com/fabienbellanger/echo-boilerplate/entities"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Routes construct all server routes
func Routes(e *echo.Echo, db *db.DB, logger *zap.Logger) {
	webRoutes(e, logger)
	apiRoutes(e, db, logger)
}

// Web routes
func webRoutes(e *echo.Echo, logger *zap.Logger) {
	g := e.Group("")

	g.GET("/health-check", func(c echo.Context) error {
		// return echo.NewHTTPError(http.StatusUnauthorized, nil)
		return c.String(http.StatusOK, "OK")
	})
}

// Api routes
func apiRoutes(e *echo.Echo, db *db.DB, logger *zap.Logger) {
	v1 := e.Group("/api/v1")

	// TODO: Login => Improve
	authGroup := v1.Group("")
	auth := user.New(authGroup, nil)
	authGroup.POST("/login", auth.Login)

	// Protected routes
	// ----------------
	jwtConfig := middleware.JWTConfig{
		ContextKey:    "user",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		SigningMethod: viper.GetString("JWT_ALGO"),
		Claims:        &entities.Claims{},
		SigningKey:    []byte(viper.GetString("JWT_SECRET")),
	}
	v1.Use(middleware.JWTWithConfig(jwtConfig))

	// User
	userRoutes := v1.Group("/users")
	user := user.New(userRoutes, nil)
	user.Routes()
}
