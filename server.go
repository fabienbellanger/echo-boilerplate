package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fabienbellanger/echo-boilerplate/db"
	"github.com/fabienbellanger/echo-boilerplate/delivery/pprof"
	"github.com/fabienbellanger/echo-boilerplate/utils"
	"github.com/fabienbellanger/goutils"
	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// Run web server
func Run(logger *zap.Logger, db *db.DB) {
	e := echo.New()

	initConfig(e)
	initMiddlerwares(e, logger)

	// Routes
	// ------
	Routes(e, db, logger)

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
func initMiddlerwares(e *echo.Echo, logger *zap.Logger) {
	// Recover
	// -------
	e.Use(middleware.Recover())

	// Logger
	// ------
	// e.Use(middleware.Logger())
	e.Use(zapLogger(logger))

	// Request ID
	// ----------
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	// CORS
	// ----
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     viper.GetStringSlice("CORS_ALLOW_ORIGINS"),
		AllowHeaders:     viper.GetStringSlice("CORS_EXPOSE_HEADERS"),
		AllowMethods:     viper.GetStringSlice("CORS_ALLOW_METHODS"),
		AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
		ExposeHeaders:    viper.GetStringSlice("CORS_EXPOSE_HEADERS"),
		MaxAge:           int(12 * time.Hour),
	}))

	// Basic Auth
	// ----------
	basicAuthConfig := middleware.BasicAuthConfig{
		Validator: func(username, password string, c echo.Context) (bool, error) {
			basicAuthUsername := viper.GetString("SERVER_BASICAUTH_USERNAME")
			basicAuthPassword := viper.GetString("SERVER_BASICAUTH_PASSWORD")

			if username == basicAuthUsername && password == basicAuthPassword {
				return true, nil
			}
			return false, nil
		},
	}
	protectedGroup := e.Group("/private")
	protectedGroup.Use(middleware.BasicAuthWithConfig(basicAuthConfig))

	// Pprof
	// -----
	if viper.GetBool("SERVER_PROMETHEUS") {
		pprof := pprof.New(protectedGroup)
		pprof.Routes()
	}

	// Prometheus
	// ----------
	if viper.GetBool("SERVER_PROMETHEUS") {
		p := prometheus.NewPrometheus(viper.GetString("APP_NAME"), nil)
		p.Use(e)
	}

	// Rate Limiter
	// ------------
	// TODO: https://echo.labstack.com/middleware/rate-limiter/
	if viper.GetBool("LIMITER_ENABLE") {
		rateLimiterConfig := middleware.RateLimiterConfig{
			Skipper: func(c echo.Context) bool {
				excludedIP := viper.GetStringSlice("LIMITER_EXCLUDE_IP")
				if len(excludedIP) == 0 {
					return false
				}
				return goutils.StringInSlice(c.RealIP(), excludedIP)
			},
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{
					Rate:      rate.Limit(viper.GetFloat64("LIMITER_LIMIT")), // req/sec
					Burst:     viper.GetInt("LIMITER_BURST"),
					ExpiresIn: 60 * time.Second, // time.Duration(viper.GetInt("LIMITER_EXPIRATION")) * time.Second,
				},
			),
			IdentifierExtractor: func(c echo.Context) (string, error) {
				return c.RealIP(), nil
			},
			ErrorHandler: func(context echo.Context, err error) error {
				return context.JSON(http.StatusForbidden, nil)
			},
			DenyHandler: func(context echo.Context, identifier string, err error) error {
				return context.JSON(http.StatusTooManyRequests, nil)
			},
		}
		e.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))
		// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
	}

	// Secure
	// ------
	e.Use(middleware.Secure())
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
