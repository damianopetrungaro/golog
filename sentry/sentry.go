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

	if e.Level() >= w.CaptureExceptionFromLevel {
		w.Hub.CaptureException(errors.New(buf.String()))
		return
	}

	w.Hub.CaptureMessage(buf.String())

}

func (w *Writer) Write(msg []byte) (int, error) {
	if w.DefaultLevel >= w.CaptureExceptionFromLevel {
		w.Hub.CaptureException(errors.New(string(msg)))
		return len(msg), nil
	}

	w.Hub.CaptureMessage(string(msg))
	return len(msg), nil
}
