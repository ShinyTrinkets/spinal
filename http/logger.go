package http

type Logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
}

type Attrs map[string]interface{}

var log Logger
