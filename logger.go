package golog

import (
	"context"
	"fmt"
)

// Logger is a logger able to write custom log Message with Fields
type Logger interface {
	Debug(context.Context, Message)
	Info(context.Context, Message)
	Warn(context.Context, Message)
	Error(context.Context, Message)
	Fatal(context.Context, Message)
	With(...Field) Logger
}

// CheckLogger is a logger able to check if a message should be written
type CheckLogger interface {
	CheckDebug(context.Context, Message) (CheckedLogger, bool)
	CheckInfo(context.Context, Message) (CheckedLogger, bool)
	CheckWarn(context.Context, Message) (CheckedLogger, bool)
	CheckError(context.Context, Message) (CheckedLogger, bool)
	CheckFatal(context.Context, Message) (CheckedLogger, bool)
}

// CheckedLogger logs an already checked log
type CheckedLogger interface {
	Log(...Field)
}

// NoopCheckedLogger is a nil-like CheckedLogger
type NoopCheckedLogger struct{}

// Log does nothing
func (n NoopCheckedLogger) Log(_ ...Field) {}

// StdCheckedLogger is a CheckedLogger which will write when called
type StdCheckedLogger struct {
	Writer Writer
	Entry  Entry
}

// Log writes a log with the given Fields
// Log panics with the message if the Level is FATAL
func (l StdCheckedLogger) Log(flds ...Field) {
	l.Writer.WriteEntry(l.Entry.With(flds...))
	if l.Entry.Level() != FATAL {
		return
	}

	msg := l.Entry.Message()
	if err := l.Writer.Flush(); err != nil {
		msg = fmt.Sprintf("%s: %s", err, msg)
	}

	panic(msg)
}
