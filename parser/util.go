package parser

import (
	"reflect"
	"strings"
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
	switch t := obj.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range t {
			s := k.(string)
			if !containsListStr(ignore, s) {
				m[s] = normalizeMapIgnore(v, ignore)
			}
		}
		return m
	case map[string]interface{}:
		m := map[string]interface{}{}
		for k, v := range t {
			if !containsListStr(ignore, k) {
				m[k] = normalizeMapIgnore(v, ignore)
			}
		}
		return m
	case []interface{}:
		for i, v := range t {
			t[i] = normalizeMapIgnore(v, ignore)
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
