package opentelemetry

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/damianopetrungaro/golog"
)

func TestNewProductionLogger(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = stdout
	})

	lvl := golog.DEBUG
	logger, flusher := NewProductionLogger(lvl)
	logger.Info(context.Background(), "hello")
	if err := flusher.Flush(); err != nil {
		t.Fatalf("could not flush: %s", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("could not close: %s", err)
	}
	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("could not read: %s", err)
	}
	if !bytes.Contains(output, []byte(`"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000"`)) {
		t.Fatalf("could not match trace and span keys in the log")
	}
}

func TestNewDevelopmentLogger(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = stdout
	})

	lvl := golog.DEBUG
	logger, flusher := NewDevelopmentLogger(lvl)
	logger.Info(context.Background(), "hello")
	if err := flusher.Flush(); err != nil {
		t.Fatalf("could not flush: %s", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("could not close: %s", err)
	}
	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("could not read: %s", err)
	}
	if !bytes.Contains(output, []byte(`trace_id="00000000000000000000000000000000" span_id="0000000000000000"`)) {
		t.Fatalf("could not match trace and span keys in the log")
	}
}
