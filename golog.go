package golog

import (
	"fmt"
)

var (
	errorHandler = func(err error) {
		fmt.Println(fmt.Sprintf("golog: could not write: %s\n", err))
	}
)

// Message is a log entry message
type Message = string

// ErrorHandler is a function which handle logging error in order to avoid returning it
type ErrorHandler func(error)

// DefaultErrorHandler returns the default error handler
func DefaultErrorHandler() ErrorHandler {
	return errorHandler
}
