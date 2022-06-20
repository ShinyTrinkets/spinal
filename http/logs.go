package http

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	config "github.com/ShinyTrinkets/spinal/config"
	util "github.com/ShinyTrinkets/spinal/util"
	"github.com/labstack/echo"
)

// LogsEndpoint enables log read/write endpoints
func LogsEndpoint(srv *echo.Echo, cfg *config.SpinalConfig) {
	// List all logs
	srv.GET("/logs", func(c echo.Context) error {
		files, err := ioutil.ReadDir(cfg.LogDir)
		if err != nil {
			return c.String(http.StatusBadRequest, "Cannot list logs!")
		}
		logsList := []string{}
		for _, file := range files {
			name := file.Name()
			if util.IsFile(cfg.LogDir+"/"+name) && filepath.Ext(name) == cfg.LogExt {
				logsList = append(logsList, name)
			}
		}
		return c.JSON(http.StatusOK, logsList)
	})

	// Read from a log; file ext is added automatically
	srv.GET("/log", func(c echo.Context) error {
		name := strings.Trim(c.QueryParam("name"), " ")
		if name == "" {
			return c.String(http.StatusBadRequest, "Name cannot be empty!")
		}
		text, err := ioutil.ReadFile(cfg.LogDir + "/" + name + cfg.LogExt)
		if err != nil {
			return c.String(http.StatusBadRequest, "Cannot read log file!")
		}
		return c.String(http.StatusBadRequest, string(text))
	})

	// Append to a log; file ext is added automatically
	srv.POST("/log", func(c echo.Context) error {
		name := strings.Trim(c.QueryParam("name"), " ")
		if name == "" {
			return c.String(http.StatusBadRequest, "Name cannot be empty!")
		}
		msg := strings.Trim(c.QueryParam("msg"), " ")
		if msg == "" {
			return c.String(http.StatusBadRequest, "Message cannot be empty!")
		}
		logFile := cfg.LogDir + "/" + name + cfg.LogExt
		// create if it doesn't exist
		if !util.IsFile(logFile) {
			err := ioutil.WriteFile(logFile, []byte(""), 0644)
			if err != nil {
				return c.String(http.StatusBadRequest, "Cannot create log file!")
			}
		}
		// start appending to the log
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return c.String(http.StatusBadRequest, "Cannot open log file for append!")
		}
		defer file.Close()
		if _, err := file.WriteString(msg + "\n"); err != nil {
			return c.String(http.StatusBadRequest, "Cannot append into log file!")
		}
		return c.String(http.StatusBadRequest, "OK")
	})
}
