package golog

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
		return l.WithDecorator(StackTraceDecorator{StacktraceFieldName: n})
	})
}

// Decorate adds the stacktrace to the entry
func (sd StackTraceDecorator) Decorate(e Entry) Entry {
	return e.With(Fields{String(sd.StacktraceFieldName, "TODO")})
}
