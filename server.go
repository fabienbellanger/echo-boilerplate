package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// Run web server
func Run() {
	e := echo.New()

	initConfig(e)
	initMiddlerwares(e)

	// Routes
	// ------
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Echo Boilerplate")
	})

	// Start server
	// ------------
	s := &http.Server{
		Addr: fmt.Sprintf("%s:%s", viper.GetString("APP_ADDR"), viper.GetString("APP_PORT")),
		// ReadTimeout:  time.Duration(viper.GetInt("server.readTimeout")) * time.Second,
		// WriteTimeout: time.Duration(viper.GetInt("server.writeTimeout")) * time.Second,
	}
	e.Logger.Fatal(e.StartServer(s))
}

// Initialize server configuration
func initConfig(e *echo.Echo) {
	// Startup banner
	// --------------
	e.HideBanner = viper.GetString("APP_ENV") == "production"

	// Debug mode
	// ----------
	e.Debug = viper.GetString("APP_ENV") != "production"

	// Validator
	// ---------
	// TODO: e.Validator = ...

	// HTTP Error handler
	// ------------------
	e.HTTPErrorHandler = customHTTPErrorHandler
}

// Initialize server middlewares
func initMiddlerwares(e *echo.Echo) {
	// Recover
	// -------
	e.Use(middleware.Recover())
}

// CustomHTTPErrorHandler
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)

	log.Printf("Code: %v\n", code)
	// errorPage := fmt.Sprintf("%d.html", code)
	// if err := c.File(errorPage); err != nil {
	// 	c.Logger().Error(err)
	// }
}
