package test

import (
	"context"
	"sync"

	"github.com/damianopetrungaro/golog"
)

var (
	_ golog.Writer = &InMemWriter{}
)

// NewNullLogger returns a null logger
// useful for passing a logger as dependency
func NewNullLogger() golog.Logger {
	return newNullStdLogger()
}

// NewFakeLogger returns a logger
// that writes entries to the passed inmemory writer
// useful to check whether the logs are written as expected
func NewFakeLogger(w *InMemWriter) golog.Logger {
	return golog.New(w)
}

// NewNullCheckLogger returns a null check logger
// useful for passing a logger as dependency
func NewNullCheckLogger() golog.CheckLogger {
	return newNullStdLogger()
}

// NewFakeCheckLogger returns a check logger
// that writes entries to the passed inmemory writer
// useful to check whether the logs are written as expected
func NewFakeCheckLogger(w *InMemWriter) golog.Logger {
	return golog.New(w)
}

func newNullStdLogger() golog.StdLogger {
	return golog.New(NewInMemWriter())
}

// InMemWriter is a Writer which keep in memory all the entries
type InMemWriter struct {
	Entries []golog.Entry
	mu      sync.Mutex
}

// NewInMemWriter returns a InMemWriter
func NewInMemWriter() *InMemWriter {
	return &InMemWriter{}
}

// WriteEntry writes an Entry to the slice
func (w *InMemWriter) WriteEntry(e golog.Entry) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Entries = append(w.Entries, e)
}

// Write writes an Entry to the slice
func (w *InMemWriter) Write(msg []byte) (int, error) {
	w.WriteEntry(golog.NewStdEntry(context.Background(), golog.DEBUG, string(msg), golog.Fields{}))
	return len(msg), nil
}

// Flush flushes the data
func (w *InMemWriter) Flush() error {
	return nil
}
