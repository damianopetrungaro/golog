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
			Encoder(nil), // TODO: add a real encoder
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
