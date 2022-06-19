package http

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	util "github.com/ShinyTrinkets/spinal/util"
	"github.com/labstack/echo"
)

const validLogExt = ".log"

// LogsEndpoint enables log read/write endpoints
func LogsEndpoint(srv *echo.Echo, logDir string) {
	if logDir == "" {
		logDir = "logs"
	} else {
		logDir = strings.TrimSuffix(logDir, "/")
	}

	// List all logs
	srv.GET("/logs", func(c echo.Context) error {
		files, err := ioutil.ReadDir(logDir)
		if err != nil {
			return c.String(http.StatusBadRequest, "Cannot list logs!")
		}
		logsList := []string{}
		for _, file := range files {
			name := file.Name()
			if util.IsFile(logDir+"/"+name) && filepath.Ext(name) == validLogExt {
				logsList = append(logsList, name)
			}
		}
		return c.JSON(http.StatusOK, logsList)
	})

	// Read from a log
	srv.GET("/log", func(c echo.Context) error {
		name := c.QueryParam("name")
		if name == "" {
			return c.String(http.StatusBadRequest, "Name cannot be empty!")
		}
		text, err := ioutil.ReadFile(logDir + "/" + name)
		if err != nil {
			return c.String(http.StatusBadRequest, "Cannot read log file!")
		}
		return c.String(http.StatusBadRequest, string(text))
	})

	// Append to a log
	srv.POST("/log", func(c echo.Context) error {
		// err := ioutil.WriteFile(logFile, []byte(message), 0644)
		// if err != nil {
		// file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)
		//     if err != nil {
		//         log.Println(err)
		//     }
		// defer file.Close()
		//     if _, err := file.WriteString("log message"); err != nil {
		//         log.Println(err)
		//     }
		return c.String(http.StatusBadRequest, "OK")
	})
}
