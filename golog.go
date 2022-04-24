package golog

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

var (
	logger Logger = New(
		INFO,
		NewBufWriter(
			NewJsonEncoder(DefaultJsonConfig()),
			bufio.NewWriter(os.Stdout),
			DefaultErrorHandler(),
		),
	)

	errorHandler = func(err error) {
		fmt.Println(fmt.Sprintf("golog: could not write: %s\n", err))
	}
)

// MinSeverity is the min log Level which can be written
type MinSeverity = Level

// Message is a log entry message
type Message = string

// ErrorHandler is a function which handle logging error in order to avoid returning it
type ErrorHandler func(error)

// DefaultErrorHandler returns the default error handler
func DefaultErrorHandler() ErrorHandler {
	return errorHandler
}

// Decorators is a slice of Decorator
type Decorators []Decorator

// Decorator modifies an entry before it get written
type Decorator interface {
	Decorate(Entry) Entry
}

// DecoratorFunc is a handy function which implements Decorator
type DecoratorFunc func(Entry) Entry

// Decorate changes the entry with custom logic and return the new modified one
func (fn DecoratorFunc) Decorate(e Entry) Entry {
	return fn(e)
}

// Checkers is a slice of Checker
type Checkers []Checker

// Checker modifies an entry before it get written
type Checker interface {
	Check(Entry) bool
}

// CheckerFunc is a handy function which implements Checker
type CheckerFunc func(Entry) bool

// Check checks if an entry should proceed to be written when using a CheckLogger
func (fn CheckerFunc) Check(e Entry) bool {
	return fn(e)
}

// Option modifies a StdLogger and returns the modified one
type Option interface {
	Apply(StdLogger) StdLogger
}

// OptionFunc is a handy function which implements Option
type OptionFunc func(StdLogger) StdLogger

// Apply change a StdLogger with custom logic and return the new modified one
func (fn OptionFunc) Apply(l StdLogger) StdLogger {
	return fn(l)
}

// SetLogger overrides the base Logger
func SetLogger(l Logger) {
	logger = l
}

// Debug calls the base Logger's Debug method
func Debug(ctx context.Context, msg string) {
	logger.Debug(ctx, msg)
}

// Info calls the base Logger's Info method
func Info(ctx context.Context, msg string) {
	logger.Info(ctx, msg)
}

// Warning calls the base Logger's Warning method
func Warning(ctx context.Context, msg string) {
	logger.Warning(ctx, msg)
}

// Error calls the base Logger's Error method
func Error(ctx context.Context, msg string) {
	logger.Error(ctx, msg)
}

// Fatal calls the base Logger's Fatal method
func Fatal(ctx context.Context, msg string) {
	logger.Fatal(ctx, msg)
}

// With calls the base Logger's With method
func With(fields Fields) Logger {
	return logger.With(fields)
}
