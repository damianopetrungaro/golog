package golog

import (
	"context"
)

// StdLogger is a representation of the standard Logger
type StdLogger struct {
	MinSeverity MinSeverity
	Writer      Writer
	Fields      Fields
	Decorators  Decorators
	Checkers    Checkers
}

// New returns a StdLogger which writes starting from the given Level to the given Writer
func New(minSev MinSeverity, w Writer, options ...Option) StdLogger {
	l := StdLogger{
		MinSeverity: minSev,
		Writer:      w,
	}

	for _, o := range options {
		l = o.Apply(l)
	}

	return l
}

// WithDecorator returns a new StdLogger appending the given extra Decorators
func (l StdLogger) WithDecorator(ds ...Decorator) StdLogger {
	return StdLogger{
		MinSeverity: l.MinSeverity,
		Writer:      l.Writer,
		Fields:      l.Fields,
		Decorators:  append(l.Decorators, ds...),
		Checkers:    l.Checkers,
	}
}

// WithCheckers returns a new StdLogger appending the given extra Checkers
func (l StdLogger) WithCheckers(cs ...Checker) StdLogger {
	return StdLogger{
		MinSeverity: l.MinSeverity,
		Writer:      l.Writer,
		Fields:      l.Fields,
		Decorators:  l.Decorators,
		Checkers:    append(l.Checkers, cs...),
	}
}

// Debug writes a log with the DEBUG Level
func (l StdLogger) Debug(ctx context.Context, msg Message) {
	l.log(ctx, DEBUG, msg)
}

// CheckDebug returns a CheckedLogger and a guard
// When the guard is true and the CheckDebug is called a log with the DEBUG Level is written
func (l StdLogger) CheckDebug(ctx context.Context, msg Message) (CheckedLogger, bool) {
	return l.check(ctx, DEBUG, msg)
}

// Info writes a log with the INFO Level
func (l StdLogger) Info(ctx context.Context, msg Message) {
	l.log(ctx, INFO, msg)
}

// CheckInfo returns a CheckedLogger and a guard
// When the guard is true and the CheckInfo is called a log with the INFO Level is written
func (l StdLogger) CheckInfo(ctx context.Context, msg Message) (CheckedLogger, bool) {
	return l.check(ctx, INFO, msg)
}

// Warning writes a log with the WARN Level
func (l StdLogger) Warning(ctx context.Context, msg Message) {
	l.log(ctx, WARN, msg)
}

// CheckWarning returns a CheckedLogger and a guard
// When the guard is true and the CheckWarning is called a log with the WARN Level is written
func (l StdLogger) CheckWarning(ctx context.Context, msg Message) (CheckedLogger, bool) {
	return l.check(ctx, WARN, msg)
}

// Error writes a log with the ERROR Level
func (l StdLogger) Error(ctx context.Context, msg Message) {
	l.log(ctx, ERROR, msg)
}

// CheckError returns a CheckedLogger and a guard
// When the guard is true and the CheckError is called a log with the ERROR Level is written
func (l StdLogger) CheckError(ctx context.Context, msg Message) (CheckedLogger, bool) {
	return l.check(ctx, ERROR, msg)
}

// Fatal writes a log with the FATAL Level
// Fatal also panic with the given message
func (l StdLogger) Fatal(ctx context.Context, msg Message) {
	l.log(ctx, FATAL, msg)
	panic(msg)
}

// CheckFatal returns a CheckedLogger and a guard
// When the guard is true and the with is called a log with the FATAL Level is written
// CheckFatal will also panic with the given message
func (l StdLogger) CheckFatal(ctx context.Context, msg Message) (CheckedLogger, bool) {
	return l.check(ctx, ERROR, msg)
}

// With returns a new Logger appending the given extra Fields
func (l StdLogger) With(fields Fields) Logger {
	return StdLogger{
		MinSeverity: l.MinSeverity,
		Writer:      l.Writer,
		Fields:      append(l.Fields, fields...),
		Decorators:  l.Decorators,
		Checkers:    l.Checkers,
	}
}

func (l StdLogger) log(ctx context.Context, lvl Level, msg Message) {
	if lvl < l.MinSeverity {
		return
	}

	var e Entry = NewStdEntry(ctx, lvl, msg, l.Fields)
	for _, d := range l.Decorators {
		e = d.Decorate(e)
	}

	l.Writer.Write(e)
}

func (l StdLogger) check(ctx context.Context, lvl Level, msg Message) (CheckedLogger, bool) {

	var e Entry = NewStdEntry(ctx, lvl, msg, l.Fields)
	for _, c := range l.Checkers {
		if !c.Check(e) {
			return NoopCheckedLogger{}, false
		}
	}

	return StdCheckedLogger{Writer: l.Writer, Entry: e}, true
}
