package http

import (
	"net/http"

	"github.com/labstack/echo"
)

// CacheEndpoint is a key-value cache store
func CacheEndpoint(srv *echo.Echo) {
	// List all stores
	srv.GET("/kv", func(c echo.Context) error {
		kvList := []string{}
		return c.JSON(http.StatusOK, kvList)
	})
}
