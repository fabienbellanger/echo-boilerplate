package server

import (
	"fmt"
	"net/http"

	"github.com/fabienbellanger/echo-boilerplate/utils"
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
	e.GET("/health-check", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusUnauthorized, nil)
		// return c.String(http.StatusOK, "OK")
	})

	// Start server
	// ------------
	s := &http.Server{
		Addr: fmt.Sprintf("%s:%s", viper.GetString("APP_ADDR"), viper.GetString("APP_PORT")),
		// TODO: ReadTimeout:  time.Duration(viper.GetInt("server.readTimeout")) * time.Second,
		// TODO: WriteTimeout: time.Duration(viper.GetInt("server.writeTimeout")) * time.Second,
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

	// Logger
	// ------
	e.Use(middleware.Logger())
}

// CustomHTTPErrorHandler
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var msg interface{}
	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
		msg = httpError.Message
	}

	switch code {
	case http.StatusBadRequest:
		// 400
		c.JSON(code, utils.HTTPError{Code: code, Message: "Bad Request", Details: msg})
	case http.StatusUnauthorized:
		// 401
		c.JSON(code, utils.HTTPError{Code: code, Message: "Unauthorized", Details: msg})
	case http.StatusNotFound:
		// 404
		c.JSON(code, utils.HTTPError{Code: code, Message: "Resource Not Found", Details: msg})
	case http.StatusInternalServerError:
		// 500
		c.Logger().Error(err)
		c.JSON(code, utils.HTTPError{Code: code, Message: "Internal Server Error", Details: msg})
	default:
		c.Logger().Error(err)
		c.JSON(code, utils.HTTPError{Code: code, Message: "Error", Details: msg})
	}
}
