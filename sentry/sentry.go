package sentry

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"

	"github.com/damianopetrungaro/golog"
)

type Writer struct {
	Encoder    golog.Encoder
	ErrHandler golog.ErrorHandler

	// Hub is the sentry Hub used for capturing events. If nil, the
	// sentry.CurrentHub() will be used.
	//
	// If the context of the golog Entry to be logged contains a sentry.Hub, it
	// will be used instead.
	Hub *sentry.Hub

	// DefaultIsException set to true means that the when this Writer is used as
	// a writer for the stdlib `log` package, its events will be treated as
	// Sentry errors.
	DefaultCaptureException bool

	// CaptureExceptionFromLevel is the level at which an Entry will be
	// treated as Sentry error instead of an arbitrary message.
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

	hub := w.getHub(e.Context())
	if e.Level() >= w.CaptureExceptionFromLevel {
		hub.CaptureException(errors.New(buf.String()))
		return
	}

	hub.CaptureMessage(buf.String())
}

func (w *Writer) Write(msg []byte) (int, error) {
	hub := w.getHub(nil)

	if w.DefaultCaptureException {
		hub.CaptureException(errors.New(string(msg)))
		return len(msg), nil
	}

	hub.CaptureMessage(string(msg))
	return len(msg), nil
}

func (w *Writer) getHub(ctx context.Context) *sentry.Hub {
	ctxHub := sentry.GetHubFromContext(ctx)
	if ctxHub != nil {
		return ctxHub
	}

	if w.Hub != nil {
		return w.Hub
	}

	return sentry.CurrentHub()
}
