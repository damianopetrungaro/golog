package sentry

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"

	"github.com/damianopetrungaro/golog"
)

type Writer struct {
	Encoder                   golog.Encoder
	Hub                       *sentry.Hub
	ErrHandler                golog.ErrorHandler
	DefaultLevel              golog.Level
	CaptureExceptionFromLevel golog.Level
}

func (w *Writer) WriteEntry(e golog.Entry) {
	wTo, err := w.Encoder.Encode(e)
	if err != nil {
		w.ErrHandler(fmt.Errorf("%w: sentry writer on encoding: %s", golog.ErrEntryNotWritten, err))
		return
	}

	buf := &bytes.Buffer{}
	if _, err := wTo.WriteTo(buf); err != nil {
		w.ErrHandler(fmt.Errorf("%w: sentry writer on write to: %s", golog.ErrEntryNotWritten, err))
		return
	}

	hub := w.Hub.Clone()
	if e.Level() >= w.CaptureExceptionFromLevel {
		hub.CaptureException(errors.New(buf.String()))
		return
	}

	hub.CaptureMessage(buf.String())

}

func (w *Writer) Write(msg []byte) (int, error) {
	hub := w.Hub.Clone()
	if w.DefaultLevel >= w.CaptureExceptionFromLevel {
		hub.CaptureException(errors.New(string(msg)))
		return len(msg), nil
	}

	hub.CaptureMessage(string(msg))
	return len(msg), nil
}
