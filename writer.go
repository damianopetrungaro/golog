package golog

import (
	"bufio"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrEntriesNotFlushed = errors.New("golog: could not flush entries")
	ErrEntryNotWritten   = errors.New("golog: could not write entry")
)

// Writer is an Entry writer
type Writer interface {
	Write(Entry)
}

// Flusher ensure that the data that a Writer may hold is written
type Flusher interface {
	Flush() error
}

// BufWriter is a Writer which holds a buffer behind the scene to reduce sys calls
type BufWriter struct {
	Encoder    Encoder
	Writer     *bufio.Writer
	ErrHandler ErrorHandler
	mu         sync.Mutex
}

// NewBufWriter returns a BufWriter
func NewBufWriter(
	enc Encoder,
	w *bufio.Writer,
	errHandler ErrorHandler,
) *BufWriter {
	return &BufWriter{
		Encoder:    enc,
		Writer:     w,
		ErrHandler: errHandler,
	}
}

// Write writes an Entry to the buffer
func (w *BufWriter) Write(e Entry) {
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

// Flush forces the data in the buffer to be written
func (w *BufWriter) Flush() error {
	if err := w.Writer.Flush(); err != nil {
		return fmt.Errorf("%w: buf writer on flush: %s", ErrEntriesNotFlushed, err)
	}

	return nil
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
