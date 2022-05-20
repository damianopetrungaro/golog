package golog

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"
)

type any = interface{}

var (
	logger Logger = New(
		writer,
		levelChecker,
		timestampDecorator,
		stackTraceDecorator,
	)

	checkLogger CheckLogger = New(
		writer,
		levelChecker,
		timestampDecorator,
		stackTraceDecorator,
	)

	writer Writer = NewBufWriter(
		NewJsonEncoder(DefaultJsonConfig()),
		bufio.NewWriter(os.Stdout),
		DefaultErrorHandler(),
		INFO,
	)

	errorHandler = func(err error) {
		fmt.Printf("golog: could not write: %s\n", err)
	}

	levelChecker        = NewLevelCheckerOption(INFO)
	timestampDecorator  = NewTimestampDecoratorOption("timestamp", time.RFC3339Nano)
	stackTraceDecorator = NewStackTraceDecoratorOption("stacktrace", 5)
)

// Message is a log entry message
type Message = string

// ErrorHandler is a function which handle logging error in order to avoid returning it
type ErrorHandler func(error)

// DefaultErrorHandler returns the default error handler
func DefaultErrorHandler() ErrorHandler {
	return errorHandler
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

// SetCheckLogger overrides the base CheckLogger
func SetCheckLogger(l CheckLogger) {
	checkLogger = l
}

// Debug calls the base Logger's Debug method
func Debug(ctx context.Context, msg string) {
	logger.Debug(ctx, msg)
}

// CheckDebug calls the base Logger's CheckDebug method
func CheckDebug(ctx context.Context, msg string) (CheckedLogger, bool) {
	return checkLogger.CheckDebug(ctx, msg)
}

// Info calls the base Logger's Info method
func Info(ctx context.Context, msg string) {
	logger.Info(ctx, msg)
}

// CheckInfo calls the base Logger's CheckInfo method
func CheckInfo(ctx context.Context, msg string) (CheckedLogger, bool) {
	return checkLogger.CheckInfo(ctx, msg)
}

// Warning calls the base Logger's Warning method
func Warning(ctx context.Context, msg string) {
	logger.Warn(ctx, msg)
}

// CheckWarning calls the base Logger's CheckWarning method
func CheckWarning(ctx context.Context, msg string) (CheckedLogger, bool) {
	return checkLogger.CheckWarn(ctx, msg)
}

// Error calls the base Logger's Error method
func Error(ctx context.Context, msg string) {
	logger.Error(ctx, msg)
}

// CheckError calls the base Logger's CheckError method
func CheckError(ctx context.Context, msg string) (CheckedLogger, bool) {
	return checkLogger.CheckError(ctx, msg)
}

// Fatal calls the base Logger's Fatal method
func Fatal(ctx context.Context, msg string) {
	logger.Fatal(ctx, msg)
}

// CheckFatal calls the base Logger's CheckFatal method
func CheckFatal(ctx context.Context, msg string) (CheckedLogger, bool) {
	return checkLogger.CheckFatal(ctx, msg)
}

// With calls the base Logger's With method
func With(fields ...Field) Logger {
	return logger.With(fields...)
}
