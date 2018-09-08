package http

import (
	"net/http"

	"github.com/ShinyTrinkets/spinal/logger"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Global log instance
var log logger.Logger

// NewServer sets up a new HTTP server
func NewServer(port string) *echo.Echo {
	log = logger.NewLogger("HttpServer")

	srv := echo.New()
	srv.Server.Addr = port
	srv.Pre(middleware.RemoveTrailingSlash())

	srv.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "There's nothing here, stranger")
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
