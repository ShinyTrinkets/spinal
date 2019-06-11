package http

import (
	"net/http"
	"net/url"

	"github.com/ShinyTrinkets/overseer.go"
	"github.com/labstack/echo"
)

// OverseerEndpoint enables Overseer endpoints
func OverseerEndpoint(srv *echo.Echo, ovr *overseer.Overseer) {
	// Get proc by ID
	// URL encoded characters in the ID are supported ("/" = "%2F")
	srv.GET("/proc/:id", func(c echo.Context) error {
		id, err := url.PathUnescape(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID format")
		}
		if ovr.HasProc(id) {
			return c.JSON(http.StatusOK, ovr.ToJSON(id))
		}
		return c.String(http.StatusBadRequest, "Invalid proc ID")
	})
	// List all procs
	srv.GET("/procs", func(c echo.Context) error {
		return c.JSON(http.StatusOK, ovr.ListAll())
	})
}
