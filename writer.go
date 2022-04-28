package golog

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	ErrEntriesNotFlushed = errors.New("golog: could not flush entries")
	ErrEntryNotWritten   = errors.New("golog: could not write entry")
)

// Writer is an Entry writer
type Writer interface {
	WriteEntry(Entry)
	io.Writer
}

// Flusher ensure that the data that a Writer may hold is written
type Flusher interface {
	Flush() error
}

// BufWriter is a Writer which holds a buffer behind the scene to reduce sys calls
type BufWriter struct {
	Encoder         Encoder
	Writer          *bufio.Writer
	ErrHandler      ErrorHandler
	DefaultLogLevel Level
	mu              sync.Mutex
}

// NewBufWriter returns a BufWriter
func NewBufWriter(
	enc Encoder,
	w *bufio.Writer,
	errHandler ErrorHandler,
	defaultLogLevel Level,
) *BufWriter {
	return &BufWriter{
		Encoder:         enc,
		Writer:          w,
		ErrHandler:      errHandler,
		DefaultLogLevel: defaultLogLevel,
	}
}

// WriteEntry writes an Entry to the buffer
func (w *BufWriter) WriteEntry(e Entry) {
	wTo, err := w.Encoder.Encode(e)
	if err != nil {
		w.ErrHandler(fmt.Errorf("%w: buf writer on encoding: %s", ErrEntryNotWritten, err))
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	if _, err := wTo.WriteTo(w.Writer); err != nil {
		w.ErrHandler(fmt.Errorf("%w: buf writer on write to: %s", ErrEntryNotWritten, err))
		return
	}
}

// Write writes an Entry to the buffer
// This method is implemented to add provide support the std library log package
func (w *BufWriter) Write(msg []byte) (int, error) {
	w.WriteEntry(NewStdEntry(context.Background(), w.DefaultLogLevel, string(msg), Fields{}))
	return len(msg), nil
}

// Flush forces the data in the buffer to be written
func (w *BufWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Writer.Buffered() == 0 {
		return nil
	}

	if err := w.Writer.Flush(); err != nil {
		return fmt.Errorf("%w: buf writer on flush: %s", ErrEntriesNotFlushed, err)
	}

	return nil
}

// MuxWriter is a Writer which based on the log level will write to a writer
// It also uses a Default one for the Write method
// as well as supporting the case when the Writer is not found in the Level map
type MuxWriter struct {
	Default     Writer
	LevelWriter map[Level]Writer
}

// NewMuxWriter returns a MuxWriter
func NewMuxWriter(
	defaultWriter Writer,
	debugWriter Writer,
	infoWriter Writer,
	warnWriter Writer,
	errorWriter Writer,
	fatalWriter Writer,
) *MuxWriter {
	return &MuxWriter{
		Default: defaultWriter,
		LevelWriter: map[Level]Writer{
			DEBUG: debugWriter,
			INFO:  infoWriter,
			WARN:  warnWriter,
			ERROR: errorWriter,
			FATAL: fatalWriter,
		},
	}
}

// WriteEntry writes an Entry to the related Writer
// If not found, then fallback on the Default
func (m *MuxWriter) WriteEntry(e Entry) {
	w, ok := m.LevelWriter[e.Level()]
	if !ok {
		m.Default.WriteEntry(e)
		return
	}

	w.WriteEntry(e)
}

// Write calls the Default Write method
func (m *MuxWriter) Write(msg []byte) (int, error) {
	return m.Default.Write(msg)
}

// TickFlusher is a Flusher triggered by a time.Ticker
type TickFlusher struct {
	Flusher
	Ticker *time.Ticker
}

// NewTickFlusher returns a TickFlusher
func NewTickFlusher(f Flusher, d time.Duration) *TickFlusher {
	return &TickFlusher{
		Flusher: f,
		Ticker:  time.NewTicker(d),
	}
}

// Flush forces the data in the inner Flusher to be written
func (f *TickFlusher) Flush() error {
	for {
		select {
		case _, ok := <-f.Ticker.C:
			if !ok {
				return nil
			}

			if err := f.Flusher.Flush(); err != nil {
				return fmt.Errorf("%w: tick flusher on flush: %s", ErrEntriesNotFlushed, err)
			}
		}
	}
}

// Close stops the time.Ticker and flushes the data
func (f *TickFlusher) Close() error {
	f.Ticker.Stop()
	if err := f.Flusher.Flush(); err != nil {
		return fmt.Errorf("%w: tick flusher on close: %s", ErrEntriesNotFlushed, err)
	}

	return nil
}
