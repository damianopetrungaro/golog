package golog

import (
	"fmt"
	"runtime"
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
	return e.With(String(td.TimestampFieldName, time.Now().Format(td.TimestampLayout)))
}

// StackTraceDecorator is a Decorator which add the log stacktrace
type StackTraceDecorator struct {
	StacktraceFieldName string
	Depth               int
}

// NewStackTraceDecorator returns a StackTraceDecorator with the given field name
func NewStackTraceDecorator(n string, depth int) StackTraceDecorator {
	return StackTraceDecorator{StacktraceFieldName: n, Depth: depth}
}

// NewStackTraceDecoratorOption returns an Option which applies a StackTraceDecorator with the given field name
func NewStackTraceDecoratorOption(n string, depth int) Option {
	return OptionFunc(func(l StdLogger) StdLogger {
		return l.WithDecorator(NewStackTraceDecorator(n, depth))
	})
}

// Decorate adds the stacktrace to the entry
func (sd StackTraceDecorator) Decorate(e Entry) Entry {
	framesToField := func(fs *runtime.Frames) Field {
		var trace []string
		for {
			f, ok := fs.Next()
			trace = append(trace, fmt.Sprintf("%v:%v", f.File, f.Line))
			if !ok {
				break
			}
		}

		return Strings(sd.StacktraceFieldName, trace)
	}

	const skip = 4
	pc := make([]uintptr, sd.Depth)
	if n := runtime.Callers(skip, pc[:]); n != 0 {
		field := framesToField(runtime.CallersFrames(pc[:]))
		return e.With(field)
	}

	return e
}
