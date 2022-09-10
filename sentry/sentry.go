package sentry

import (
	"context"

	"github.com/getsentry/sentry-go"

	"github.com/damianopetrungaro/golog"
)

type Writer struct {
	Hub          *sentry.Hub
	DefaultLevel golog.Level
}

func (w *Writer) WriteEntry(e golog.Entry) {
	ev := sentry.NewEvent()
	ev.Level = toSentryLevel(e.Level())
	ev.Message = e.Message()
	for _, f := range e.Fields() {
		ev.Extra[f.Key()] = f.Value()
	}

	hub := w.Hub.Clone()
	hub.CaptureEvent(ev)
}

func (w *Writer) Write(msg []byte) (int, error) {
	e := golog.NewStdEntry(context.Background(), w.DefaultLevel, string(msg), golog.Fields{})
	w.WriteEntry(e)

	return len(msg), nil
}

func toSentryLevel(lvl golog.Level) sentry.Level {
	switch lvl {
	case golog.DEBUG:
		return sentry.LevelDebug
	case golog.INFO:
		return sentry.LevelInfo
	case golog.WARN:
		return sentry.LevelWarning
	case golog.ERROR:
		return sentry.LevelError
	case golog.FATAL:
		return sentry.LevelFatal
	}

	return sentry.LevelError
}
