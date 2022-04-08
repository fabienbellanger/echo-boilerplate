package server

import (
	"net/http"

	"github.com/fabienbellanger/echo-boilerplate/db"
	"github.com/fabienbellanger/echo-boilerplate/delivery/user"
	"github.com/labstack/echo/v4"
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

	// TODO: Protected routes

	// User
	userRoutes := v1.Group("/users")
	user := user.New(userRoutes, nil)
	user.Routes()
}
