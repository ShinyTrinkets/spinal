package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ShinyTrinkets/overseer"
	quote "github.com/kballard/go-shellquote"
	"github.com/labstack/echo"
)

// OverseerEndpoint enables Overseer endpoints
func OverseerEndpoint(srv *echo.Echo, ovr *overseer.Overseer) {
	// List all procs
	srv.GET("/procs", func(c echo.Context) error {
		return c.JSON(http.StatusOK, ovr.ListAll())
	})

	// Get proc by ID
	// URL encoded characters in the ID are supported ("/" = "%2F")
	srv.GET("/proc/:id", func(c echo.Context) error {
		id, err := url.PathUnescape(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID format")
		}
		if ovr.HasProc(id) {
			return c.JSON(http.StatusOK, ovr.Status(id))
		}
		return c.String(http.StatusBadRequest, "Invalid proc ID")
	})

	// Add, Supervise and Remove a process when complete
	srv.GET("/stop/:id", func(c echo.Context) error {
		id, err := url.PathUnescape(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest,
				fmt.Sprintf("No ID! Error: %v\n", err))
		}

		ovr.Stop(id)
		time.Sleep(250 * time.Millisecond)
		ovr.Remove(id)

		return c.String(http.StatusOK, "Done")
	})

	// Add, Supervise and Remove a process when complete
	srv.GET("/run/:id", func(c echo.Context) error {
		id, err := url.PathUnescape(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest,
				fmt.Sprintf("No ID! Error: %v\n", err))
		}

		exec := c.QueryParam("exec")
		if exec == "" {
			return c.String(http.StatusBadRequest, "Exec command cannot be empty!")
		}
		args, err := quote.Split(exec)
		if err != nil {
			return c.String(http.StatusBadRequest,
				fmt.Sprintf("Cannot split args! Error: %v\n", err))
		}

		delay, err := strconv.ParseUint(c.QueryParam("delay"), 10, 16)
		if err != nil {
			return c.String(http.StatusBadRequest,
				fmt.Sprintf("Invalid delay value! Error: %v\n", err))
		}
		retry, err := strconv.ParseUint(c.QueryParam("retry"), 10, 16)
		if err != nil {
			return c.String(http.StatusBadRequest,
				fmt.Sprintf("Invalid retry value! Error: %v\n", err))
		}
		cwd := c.QueryParam("cwd")

		opts := overseer.Options{Buffered: true, Streaming: false}
		if cwd != "" {
			opts.Dir = cwd
		}
		if delay > 0 {
			opts.DelayStart = uint(delay)
		}
		if retry > 0 {
			opts.RetryTimes = uint(retry)
		}

		p := ovr.Add(id, args[0], args[1:], opts)
		if p == nil {
			return c.String(http.StatusBadRequest, "Proc cannot added!")
		}
		fmt.Println(p)

		go func() {
			ovr.Supervise(id)
			time.Sleep(250 * time.Millisecond)
			ovr.Remove(id)
		}()

		return c.String(http.StatusOK, "Done")
	})
}
