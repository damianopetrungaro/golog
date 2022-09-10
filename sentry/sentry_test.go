package sentry

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/damianopetrungaro/golog"
)

var (
	msg        = `This is a message for sentry`
	debugEntry = golog.NewStdEntry(context.Background(), golog.DEBUG, msg, nil)
	errorEntry = golog.NewStdEntry(context.Background(), golog.ERROR, msg, nil)
)

type fakeTransport struct {
	Level   string
	Message string
	Extra   map[string]interface{}
}

func (t *fakeTransport) Configure(sentry.ClientOptions) {}

func (t *fakeTransport) SendEvent(ev *sentry.Event) {
	t.Level = string(ev.Level)
	t.Message = ev.Message
	t.Extra = ev.Extra
}

func (t *fakeTransport) Flush(time.Duration) bool {
	return true
}

func TestWriter_WriteEntry(t *testing.T) {
	transport := &fakeTransport{}

	hub := getHubHelper(t, transport)

	w := &Writer{
		Hub:          hub,
		DefaultLevel: golog.INFO,
	}

	if _, err := w.Write([]byte(`message`)); err != nil {
		t.Fatalf("could not write log message: %s", err)
	}

	if transport.Message != `message` {
		t.Error("could not match message")
		t.Errorf("got: %s", transport.Message)
		t.Errorf("want: %s", debugEntry.Message())
	}

	if transport.Level != strings.ToLower(w.DefaultLevel.String()) {
		t.Error("could not match level")
		t.Errorf("got: %s", transport.Level)
		t.Errorf("want: %s", debugEntry.Level().String())
	}

	if len(transport.Extra) > 0 {
		t.Error("could not match extra")
		t.Errorf("got: %v", transport.Extra)
	}

	w.WriteEntry(errorEntry.With(golog.String("extra_key", "extra_value")))
	if transport.Message != errorEntry.Message() {
		t.Error("could not match message")
		t.Errorf("got: %s", transport.Message)
		t.Errorf("want: %s", errorEntry.Message())
	}

	if transport.Level != strings.ToLower(errorEntry.Level().String()) {
		t.Error("could not match level")
		t.Errorf("got: %s", transport.Level)
		t.Errorf("want: %s", errorEntry.Level().String())
	}

	if len(transport.Extra) != 1 {
		t.Error("could not match extra")
		t.Errorf("got: %v", transport.Extra)
	}

	if transport.Extra["extra_key"] != "extra_value" {
		t.Error("could not match extra key")
		t.Errorf("got: %v", transport.Extra["extra_key"])
	}
}

func TestWriter_Write(t *testing.T) {
	transport := &fakeTransport{}

	hub := getHubHelper(t, transport)

	w := &Writer{
		Hub:          hub,
		DefaultLevel: golog.INFO,
	}

	w.WriteEntry(debugEntry)
	if transport.Message != debugEntry.Message() {
		t.Error("could not match message")
		t.Errorf("got: %s", transport.Message)
		t.Errorf("want: %s", debugEntry.Message())
	}

	if transport.Level != strings.ToLower(debugEntry.Level().String()) {
		t.Error("could not match level")
		t.Errorf("got: %s", transport.Level)
		t.Errorf("want: %s", debugEntry.Level().String())
	}

	if len(transport.Extra) > 0 {
		t.Error("could not match extra")
		t.Errorf("got: %v", transport.Extra)
	}
}

func Test_toSentryLevel(t *testing.T) {
	tests := []struct {
		lvl  golog.Level
		want sentry.Level
	}{
		{lvl: golog.DEBUG, want: sentry.LevelDebug},
		{lvl: golog.INFO, want: sentry.LevelInfo},
		{lvl: golog.WARN, want: sentry.LevelWarning},
		{lvl: golog.ERROR, want: sentry.LevelError},
		{lvl: golog.FATAL, want: sentry.LevelFatal},
		{lvl: golog.Level(99), want: sentry.LevelError},
	}

	for _, test := range tests {
		if got := toSentryLevel(test.lvl); got != test.want {
			t.Error("could not match level")
			t.Errorf("got: %s", got)
			t.Errorf("want: %s", test.want)
		}
	}
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
