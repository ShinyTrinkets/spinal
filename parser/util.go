package parser

import (
	"os"
	"time"

	"github.com/immortal/xtime"
)

// isFile: helper that returns true if the path is a regular file
func isFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); !m.IsDir() && m.IsRegular() && m&400 != 0 {
		return true
	}
	return false
}

// fileStats: helper that returns file stats (creation and modif times)
func fileStats(fname string) (time.Time, time.Time, error) {
	var c time.Time
	var m time.Time

	fi, err := os.Stat(fname)
	if err != nil {
		// File stats error
		return c, m, err
	}

	c = xtime.Get(fi).Ctime()
	m = xtime.Get(fi).Mtime()
	return c, m, nil
}
