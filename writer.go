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
	Flusher
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

// LeveledWriterOptionFunc is a handy function which implements attach a Writer for a given Level in a LeveledWriter
type LeveledWriterOptionFunc func(*LeveledWriter)

// DefaultLeveledWriterOptionFunc implements LeveledWriterOptionFunc
func DefaultLeveledWriterOptionFunc(lvl Level, w Writer) LeveledWriterOptionFunc {
	return func(mux *LeveledWriter) {
		mux.LevelWriter[lvl] = w
	}
}

// LeveledWriter is a Writer which based on the log level will write to a writer
// It also uses a Default one for the Write method
// as well as supporting the case when the Writer is not found in the Level map
type LeveledWriter struct {
	Default     Writer
	LevelWriter map[Level]Writer
}

// NewLeveledWriter returns a LeveledWriter
func NewLeveledWriter(
	defaultWriter Writer,
	fns ...LeveledWriterOptionFunc,
) *LeveledWriter {
	w := &LeveledWriter{Default: defaultWriter, LevelWriter: map[Level]Writer{}}
	for _, fn := range fns {
		fn(w)
	}

	return w
}

// WriteEntry writes an Entry to the related Writer
// If not found, then fallback on the Default
func (m *LeveledWriter) WriteEntry(e Entry) {
	w, ok := m.LevelWriter[e.Level()]
	if !ok {
		m.Default.WriteEntry(e)
		return
	}

	w.WriteEntry(e)
}

// Write calls the Default Write method
func (m *LeveledWriter) Write(msg []byte) (int, error) {
	return m.Default.Write(msg)
}

// Flush flushes the data
func (m *LeveledWriter) Flush() error {
	return m.Default.Flush()
}

// DeduplicatorWriter is a Writer which deduplicate fields with the same name
type DeduplicatorWriter struct {
	Default Writer
}

// NewDeduplicatorWriter returns a DeduplicatorWriter
func NewDeduplicatorWriter(
	defaultWriter Writer,
) *DeduplicatorWriter {
	w := &DeduplicatorWriter{Default: defaultWriter}

	return w
}

// WriteEntry writes an Entry to the related Writer
// If not found, then fallback on the Default
func (m *DeduplicatorWriter) WriteEntry(e Entry) {
	var flds Fields
	counter := map[string]int{}

	for _, f := range e.Fields() {
		k := f.Key()
		c, exists := counter[k]
		if exists {
			c++
			k = fmt.Sprintf("%s_%d", k, c)
		}

		flds = append(flds, Field{k: k, v: f.Value()})
		counter[f.Key()] = c
	}

	m.Default.WriteEntry(StdEntry{
		Ctx:  e.Context(),
		Lvl:  e.Level(),
		Msg:  e.Message(),
		Flds: flds,
	})
}

// Write calls the Default Write method
func (m *DeduplicatorWriter) Write(msg []byte) (int, error) {
	return m.Default.Write(msg)
}

// Flush flushes the data
func (m *DeduplicatorWriter) Flush() error {
	return m.Default.Flush()
}

// MultiWriter is a Writer which based on the log level will write to a writer
// It also uses a Default one for the Write method
// as well as supporting the case when the Writer is not found in the Level map
type MultiWriter struct {
	Writers []Writer
}

// NewMultiWriter returns a MultiWriter
func NewMultiWriter(
	ws ...Writer,
) *MultiWriter {
	return &MultiWriter{Writers: ws}
}

// WriteEntry writes an Entry to the related Writer
// If not found, then fallback on the Default
func (m *MultiWriter) WriteEntry(e Entry) {
	wg := &sync.WaitGroup{}
	wg.Add(len(m.Writers))
	for _, w := range m.Writers {
		go func(w Writer) {
			w.WriteEntry(e)
			wg.Done()
		}(w)
	}

	wg.Wait()
}

// Write calls the Default Write method
func (m *MultiWriter) Write(msg []byte) (int, error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(m.Writers))

	var err error
	for _, w := range m.Writers {
		go func(w Writer) {
			if _, werr := w.Write(msg); werr != nil {
				err = werr
			}

			wg.Done()
		}(w)
	}

	wg.Wait()
	return len(msg), err
}

// Flush flushes the data to the related Writer
func (m *MultiWriter) Flush() error {
	wg := &sync.WaitGroup{}
	wg.Add(len(m.Writers))

	var err error
	for _, w := range m.Writers {
		go func(w Writer) {
			if ferr := w.Flush(); ferr != nil {
				err = ferr
			}

			wg.Done()
		}(w)
	}

	wg.Wait()
	return err
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
	for range f.Ticker.C {
		if err := f.Flusher.Flush(); err != nil {
			return fmt.Errorf("%w: tick flusher on flush: %s", ErrEntriesNotFlushed, err)
		}
	}

	return nil
}

// Close stops the time.Ticker and flushes the data
func (f *TickFlusher) Close() error {
	f.Ticker.Stop()
	if err := f.Flusher.Flush(); err != nil {
		return fmt.Errorf("%w: tick flusher on close: %s", ErrEntriesNotFlushed, err)
	}

	return nil
}
