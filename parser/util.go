package parser

import (
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/immortal/xtime"
)

func containsListStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func normalizeMapIgnore(obj interface{}, ignore []string) interface{} {
	switch x := obj.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			s := k.(string)
			if !containsListStr(ignore, s) {
				m2[s] = normalizeMapIgnore(v, ignore)
			}
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = normalizeMapIgnore(v, ignore)
		}
	}
	return obj
}

// getTagsByName returns the tag IDs for all fields in the struct
func getTagsByName(obj interface{}, tag string) (tags []string) {
	val := reflect.ValueOf(obj)
	// if pointer, get the underlying element≤
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// we can't access the value of unexported fields
		if field.PkgPath != "" {
			continue
		}
		tag := field.Tag.Get(tag)
		if tag != "" {
			tags = append(tags, strings.Split(tag, ",")[0])
		}
	}

	return tags
}

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
