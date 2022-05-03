package sentry_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/damianopetrungaro/golog"
	. "github.com/damianopetrungaro/golog/sentry"
)

var (
	msg        = `This is a message for sentry`
	debugEntry = golog.NewStdEntry(context.Background(), golog.DEBUG, msg, nil)
	errorEntry = golog.NewStdEntry(context.Background(), golog.ERROR, msg, nil)
)

type fakeTransport struct {
	Message   string
	Excpetion []sentry.Exception
}

func (t *fakeTransport) Configure(sentry.ClientOptions) {}

func (t *fakeTransport) SendEvent(ev *sentry.Event) {
	t.Message = ev.Message
	t.Excpetion = ev.Exception
}

func (t *fakeTransport) Flush(time.Duration) bool {
	return true
}

// FakeEncoder used for internal testing purposes
type FakeEncoder struct {
	Entry          golog.Entry
	ShouldFail     bool
	ShouldWriterTo io.WriterTo
}

func (fe *FakeEncoder) Encode(e golog.Entry) (io.WriterTo, error) {
	fe.Entry = e
	return fe.ShouldWriterTo, nil
}

func TestWriter(t *testing.T) {
	t.Run("capture message", func(t *testing.T) {
		data := []byte(`This is the data written`)
		writerTo := bytes.NewBuffer(data)
		transport := &fakeTransport{}

		hub := getHubHelper(t, transport)
		enc := &FakeEncoder{ShouldWriterTo: writerTo}

		w := &Writer{
			Encoder:                   enc,
			Hub:                       hub,
			ErrHandler:                golog.DefaultErrorHandler(),
			DefaultLevel:              golog.INFO,
			CaptureExceptionFromLevel: golog.WARN,
		}

		w.WriteEntry(debugEntry)
		if transport.Message != string(data) {
			t.Error("could not match message")
			t.Errorf("got: %s", transport.Message)
			t.Errorf("want: %s", data)
		}

		if len(transport.Excpetion) > 0 {
			t.Error("could not match exception")
			t.Errorf("got: %v", transport.Excpetion)
		}
	})

	t.Run("capture message", func(t *testing.T) {
		data := []byte(`This is the data written`)
		writerTo := bytes.NewBuffer(data)
		transport := &fakeTransport{}

		hub := getHubHelper(t, transport)
		enc := &FakeEncoder{ShouldWriterTo: writerTo}

		w := &Writer{
			Encoder:                   enc,
			Hub:                       hub,
			ErrHandler:                golog.DefaultErrorHandler(),
			DefaultLevel:              golog.INFO,
			CaptureExceptionFromLevel: golog.WARN,
		}

		w.WriteEntry(errorEntry)
		if transport.Excpetion[0].Value != string(data) {
			t.Error("could not match exception")
			t.Errorf("got: %s", transport.Message)
			t.Errorf("want: %s", data)
		}

		if transport.Message != "" {
			t.Error("could not match message")
			t.Errorf("got: %v", transport.Excpetion)
		}
	})
}

func getHubHelper(t *testing.T, transport sentry.Transport) *sentry.Hub {
	t.Helper()

	client, err := sentry.NewClient(sentry.ClientOptions{Transport: transport})
	if err != nil {
		t.Fatalf("could not create client: %s", err)
	}

	scope := sentry.NewScope()
	hub := sentry.NewHub(client, scope)
	return hub
}
