package util

import (
	"os"
	"time"

	"github.com/immortal/xtime"
)

// IsFile - helper that returns true if the path is a regular file
func IsFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); !m.IsDir() && m.IsRegular() && m&400 != 0 {
		return true
	}
	return false
}

// IsDir - helper that returns true if the path is a dir
func IsDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); m.IsDir() && m&400 != 0 {
		return true
	}
	return false
}

// FileStats - helper that returns file stats (creation and modif times)
func FileStats(fname string) (time.Time, time.Time, error) {
	var t time.Time

	fi, err := os.Stat(fname)
	if err != nil {
		// File stats error
		return t, t, err
	}

	x := xtime.Get(fi)
	return x.Ctime(), x.Mtime(), nil
}
