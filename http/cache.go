package http

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ShinyTrinkets/spinal/kvstore"
	"github.com/labstack/echo"
)

// CacheEndpoint is a key-value cache store
func CacheEndpoint(srv *echo.Echo) {
	// List all stores
	srv.GET("/kv", func(c echo.Context) error {
		kvList := kvstore.List()
		return c.JSON(http.StatusOK, kvList)
	})

	srv.GET("/kv/:id/:key", func(c echo.Context) error {
		id, err := url.PathUnescape(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}
		key, err := url.PathUnescape(c.Param("key"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid key")
		}
		kv := kvstore.Store(id)
		data, _ := kv.Get(key)
		return c.JSON(http.StatusOK, data)
	})

	srv.POST("/kv/:id/:key", func(c echo.Context) error {
		id, err := url.PathUnescape(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}
		key, err := url.PathUnescape(c.Param("key"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid key")
		}
		data := c.QueryParam("data")
		if data == "" {
			return c.String(http.StatusBadRequest, "Data cannot be empty!")
		}
		var i interface{}
		json.Unmarshal([]byte(data), &i)
		kv := kvstore.Store(id)
		kv.Set(key, i, -1)
		return c.String(http.StatusOK, "OK")
	})
}
