package http

import (
	"net/http"

	"github.com/ShinyTrinkets/overseer.go"
	"github.com/labstack/echo"
)

// OverseerEndpoint enables Overseer endpoints
func OverseerEndpoint(srv *echo.Echo, ovr *overseer.Overseer) {
	// Get proc by ID
	srv.GET("/proc/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(http.StatusOK, ovr.ToJSON(id))
	})
	// List all procs
	srv.GET("/proc", func(c echo.Context) error {
		return c.JSON(http.StatusOK, ovr.ListAll())
	})
}
