package http

import (
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog/log"
)

// NewServer sets up a new HTTP server
func NewServer(port string) *echo.Echo {
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
	log.Info().Msgf("HTTP server start on '%s'", srv.Server.Addr)
	if err := gracehttp.Serve(srv.Server); err != nil {
		log.Fatal().Err(err).Msg("HTTP server error")
	} else {
		log.Info().Msg("HTTP server shutdown")
	}
}
