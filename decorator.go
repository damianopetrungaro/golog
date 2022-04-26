package golog

import (
	"time"
)

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

// TimestampDecorator is a Decorator which add the log timestamp
type TimestampDecorator struct {
	TimestampLayout    string
	TimestampFieldName string
}

// NewTimestampDecorator returns a TimestampDecorator with the given field name and layout
func NewTimestampDecorator(name, layout string) TimestampDecorator {
	return TimestampDecorator{TimestampFieldName: name, TimestampLayout: layout}
}

// NewTimestampDecoratorOption returns an Option which applies a TimestampDecorator with the given field name
func NewTimestampDecoratorOption(name, layout string) Option {
	return OptionFunc(func(l StdLogger) StdLogger {
		return l.WithDecorator(NewTimestampDecorator(name, layout))
	})
}

// Decorate adds the timestamp to the entry
func (td TimestampDecorator) Decorate(e Entry) Entry {
	return e.With(Fields{String(td.TimestampFieldName, time.Now().Format(td.TimestampLayout))})
}

// StackTraceDecorator is a Decorator which add the log stacktrace
type StackTraceDecorator struct {
	StacktraceFieldName string
}

// NewStackTraceDecorator returns a StackTraceDecorator with the given field name
func NewStackTraceDecorator(n string) StackTraceDecorator {
	return StackTraceDecorator{StacktraceFieldName: n}
}

// NewStackTraceDecoratorOption returns an Option which applies a StackTraceDecorator with the given field name
func NewStackTraceDecoratorOption(n string) Option {
	return OptionFunc(func(l StdLogger) StdLogger {
		return l.WithDecorator(NewStackTraceDecorator(n))
	})
}

// Decorate adds the stacktrace to the entry
func (sd StackTraceDecorator) Decorate(e Entry) Entry {
	return e.With(Fields{String(sd.StacktraceFieldName, "TODO")})
}
