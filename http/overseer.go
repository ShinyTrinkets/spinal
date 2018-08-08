package http

import (
	"net/http"

	"github.com/ShinyTrinkets/overseer.go"
	"github.com/labstack/echo"
)

// Enable Overseer endpoints
func HttpOverseer(srv *echo.Echo, ovr *overseer.Overseer) {
	// Get proc by ID
	srv.GET("/proc/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(http.StatusOK, ovr.Status(id))
	})
	// List all procs
	srv.GET("/proc", func(c echo.Context) error {
		return c.JSON(http.StatusOK, ovr.ListAll())
	})
}
