package golog

import (
	"context"
)

// Logger is a logger able to write custom log entries with
type Logger interface {
	Debug(context.Context, Message)
	Info(context.Context, Message)
	Warning(context.Context, Message)
	Error(context.Context, Message)
	Fatal(context.Context, Message)
	With(Fields) Logger
}

// StdLogger is a representation of the standard Logger
type StdLogger struct {
	MinSeverity Level
	Writer      Writer
	Fields      Fields
	Decorators  []Decorator
}

// Decorators is a slice of Decorator
type Decorators []Decorator

// Decorator modifies an entry before it get written
type Decorator interface {
	Decorate(Entry) Entry
}

// DecoratorFunc is a handy function which implements Decorator
type DecoratorFunc func(Entry) Entry

// Decorate change the entry with custom logic and return the new modified one
func (fn DecoratorFunc) Decorate(e Entry) Entry {
	return fn(e)
}

// New returns a StdLogger which writes starting from the given Level to the given Writer
// It optionally accepts decorators
func New(minSeverity Level, w Writer, ds ...Decorator) StdLogger {
	return StdLogger{
		MinSeverity: minSeverity,
		Writer:      w,
		Decorators:  ds,
	}
}

// WithDecorator returns a new StdLogger appending the given extra Decorators
func (l StdLogger) WithDecorator(ds ...Decorator) StdLogger {
	return StdLogger{
		MinSeverity: l.MinSeverity,
		Writer:      l.Writer,
		Decorators:  append(l.Decorators, ds...),
		Fields:      l.Fields,
	}
}

// Debug writes a log with the DEBUG Level
func (l StdLogger) Debug(ctx context.Context, msg Message) {
	l.log(ctx, DEBUG, msg)
}

// Info writes a log with the INFO Level
func (l StdLogger) Info(ctx context.Context, msg Message) {
	l.log(ctx, INFO, msg)
}

// Warning writes a log with the WARN Level
func (l StdLogger) Warning(ctx context.Context, msg Message) {
	l.log(ctx, WARN, msg)
}

// Error writes a log with the ERROR Level
func (l StdLogger) Error(ctx context.Context, msg Message) {
	l.log(ctx, ERROR, msg)
}

// Fatal writes a log with the FATAL Level
// Fatal also panic with the given message
func (l StdLogger) Fatal(ctx context.Context, msg Message) {
	l.log(ctx, FATAL, msg)
	panic(msg)
}

// With returns a new Logger appending the given extra Fields
func (l StdLogger) With(fields Fields) Logger {
	return StdLogger{
		MinSeverity: l.MinSeverity,
		Writer:      l.Writer,
		Decorators:  l.Decorators,
		Fields:      append(l.Fields, fields...),
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
