package golog

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
)

func TestNewProductionLogger(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = stdout
	})

	lvl := DEBUG
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
	if !bytes.Contains(output, []byte(`"level":"INFO"`)) {
		t.Fatalf("could not match level key in the logs")
	}
	if !bytes.Contains(output, []byte(`"stacktrace":[`)) {
		t.Fatalf("could not match stactrace key in the logs")
	}
	if !bytes.Contains(output, []byte(`"timestamp":"`)) {
		t.Fatalf("could not match timestamp key in the logs")
	}
}

func TestNewDevelopmentLogger(t *testing.T) {
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = stdout
	})

	lvl := DEBUG
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
	if !bytes.Contains(output, []byte(`level=`)) {
		t.Fatalf("could not match level key in the logs")
	}
	if !bytes.Contains(output, []byte(`stacktrace=[`)) {
		t.Fatalf("could not match stacktrace key in the logs")
	}
	if !bytes.Contains(output, []byte(`timestamp="`)) {
		t.Fatalf("could not match timestamp key in the logs")
	}
}
