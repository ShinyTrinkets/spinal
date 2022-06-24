package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	config "github.com/ShinyTrinkets/spinal/config"
	util "github.com/ShinyTrinkets/spinal/util"
	"github.com/labstack/echo"
)

type LogEntry struct {
	Level uint   `json:"level"`
	Time  uint   `json:"time"`
	Msg   string `json:"msg"`
	Pid   uint   `json:"pid,omitempty"`
}

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
	srv.GET("/log/:id", func(c echo.Context) error {
		id, err := url.PathUnescape(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}
		text, err := ioutil.ReadFile(cfg.LogDir + "/" + id + cfg.LogExt)
		if err != nil {
			return c.String(http.StatusBadRequest, "Cannot read log file!")
		}
		return c.String(http.StatusBadRequest, string(text))
	})

	// Append to a log; file ext is added automatically
	srv.POST("/log/:id", func(c echo.Context) error {
		// Using the pino & pino-pretty log format
		// https://github.com/pinojs/pino-pretty
		id, err := url.PathUnescape(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid ID")
		}
		msg := strings.Trim(c.QueryParam("msg"), " ")
		if msg == "" {
			return c.String(http.StatusBadRequest, "Message cannot be empty!")
		}
		lvl, err := strconv.ParseUint(c.QueryParam("lvl"), 10, 16)
		if err != nil {
			return c.String(http.StatusBadRequest,
				fmt.Sprintf("Invalid Level value! Error: %v\n", err))
		}
		pid, err := strconv.ParseUint(c.QueryParam("pid"), 10, 16)
		if err != nil {
			fmt.Printf("Invalid PID value! Error: %v\n", err)
		}

		logFile := cfg.LogDir + "/" + id + cfg.LogExt
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

		ts := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)
		le := LogEntry{Level: uint(lvl), Time: uint(ts), Msg: msg}
		if pid != 0 {
			le.Pid = uint(pid)
		}
		line, _ := json.Marshal(le)
		if _, err := file.WriteString(string(line) + "\n"); err != nil {
			return c.String(http.StatusBadRequest, "Cannot append into log file!")
		}
		return c.String(http.StatusBadRequest, "OK")
	})
}
