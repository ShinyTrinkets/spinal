package http

import (
	"net/http"
	"net/url"

	logr "github.com/ShinyTrinkets/meta-logger"
	"github.com/ShinyTrinkets/spinal/state"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	// Logger is a type alias
	Logger = logr.Logger
	// DefaultLogger is a type alias
	DefaultLogger = logr.DefaultLogger
)

// Global log instance
var log Logger

// NewServer sets up a new HTTP server
func NewServer(port string) *echo.Echo {
	if logr.NewLogger == nil {
		// When the logger is not defined, use the basic logger
		logr.NewLogger = func(name string) Logger {
			return &DefaultLogger{Name: name}
		}
	}
	// Setup the logs by calling user's provided log builder
	log = logr.NewLogger("HttpServer")

	srv := echo.New()
	srv.Server.Addr = port
	srv.Pre(middleware.RemoveTrailingSlash())

	srv.Static("/static", "static")

	srv.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "The Spinal server is running")
	})

	// Get state lvl1 by ID
	// URL encoded characters in the ID are supported ("/" = "%2F")
	srv.GET("/state/:id", func(c echo.Context) error {
		id, err := url.PathUnescape(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID format")
		}
		if state.HasLevel1(id) {
			return c.JSON(http.StatusOK, state.GetLevel1(id))
		}
		return c.String(http.StatusBadRequest, "Invalid state ID")
	})
	// Get app state
	srv.GET("/state", func(c echo.Context) error {
		return c.JSON(http.StatusOK, state.GetState())
	})

	return srv
}

// Serve listens and serves
func Serve(srv *echo.Echo) {
	log.Info("HTTP server start on '%s'", srv.Server.Addr)
	if err := gracehttp.Serve(srv.Server); err != nil {
		log.Error("HTTP server error: %s", err)
	} else {
		log.Info("HTTP server shutdown")
	}
}
