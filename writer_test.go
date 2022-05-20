package golog_test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/damianopetrungaro/golog"
)

var (
	_ Writer      = &FakeWriter{}
	_ io.WriterTo = &FakeIoWriterTo{}
)

var (
	ErrFakeFlusher  = errors.New("an error occurred on the flusher")
	ErrFakeWriterTo = errors.New("an error occurred on the writer to")
)

// FakeWriter used for internal testing purposes
type FakeWriter struct {
	Entry   Entry
	Message Message
}

func (fw *FakeWriter) WriteEntry(e Entry) {
	fw.Entry = e
}

func (fw *FakeWriter) Write(msg []byte) (int, error) {
	fw.Message = Message(msg)
	return len(msg), nil
}

// FakeFlusher used for internal testing purposes
type FakeFlusher struct {
	Counter    int64
	ShouldFail bool
}

func (fw *FakeFlusher) Flush() error {
	atomic.AddInt64(&fw.Counter, 1)

	if fw.ShouldFail {
		return ErrFakeFlusher
	}

	return nil
}

// FakeIoWriterTo used for internal testing purposes
type FakeIoWriterTo struct {
	ShouldFail bool
}

func (fw *FakeIoWriterTo) WriteTo(w io.Writer) (int64, error) {
	if fw.ShouldFail {
		return 0, ErrFakeWriterTo
	}

	return 0, nil
}

func TestBufWriter_WriteEntry(t *testing.T) {
	t.Run("successfully write", func(t *testing.T) {
		data := []byte(`This is the data written`)
		buf := &bytes.Buffer{}
		writerTo := bytes.NewBuffer(data)
		enc := &FakeEncoder{ShouldWriterTo: writerTo}
		errHandler := &FakeErrorHandler{}

		w := BufWriter{
			Encoder:    enc,
			Writer:     bufio.NewWriter(buf),
			ErrHandler: errHandler.Handle,
		}

		e := NewStdEntry(context.Background(), ERROR, "A log error message", Fields{Bool("test", true)})
		w.WriteEntry(e)
		_ = w.Writer.Flush()

		EntryMatcher(t, e, enc.Entry)
		if got := buf.Bytes(); !bytes.Equal(got, data) {
			t.Error("could not match data written")
			t.Errorf("got: %s", got)
			t.Errorf("want: %s", data)
		}
		if errHandler.Err != nil {
			t.Errorf("could not match error: %s", errHandler.Err)
		}
	})

	t.Run("when fails encoding errors are handled", func(t *testing.T) {
		buf := &bytes.Buffer{}
		enc := &FakeEncoder{ShouldFail: true}
		errHandler := &FakeErrorHandler{}

		w := BufWriter{
			Encoder:    enc,
			Writer:     bufio.NewWriter(buf),
			ErrHandler: errHandler.Handle,
		}

		e := NewStdEntry(context.Background(), ERROR, "A log error message", Fields{Bool("test", true)})
		w.WriteEntry(e)
		_ = w.Writer.Flush()

		EntryMatcher(t, e, enc.Entry)
		if got := buf.Bytes(); len(got) != 0 {
			t.Error("could not match data written")
			t.Errorf("got: %s", got)
			t.Errorf("want: %v", nil)
		}
		if !errors.Is(errHandler.Err, ErrEntryNotWritten) || !strings.Contains(errHandler.Err.Error(), ErrFakeEncoder.Error()) {
			t.Errorf("could not match error")
			t.Errorf("got: %s", errHandler.Err)
			t.Errorf("want: %s", ErrFakeEncoder)
		}
	})

	t.Run("when fails writing down to the buffer errors are handled", func(t *testing.T) {
		buf := &bytes.Buffer{}
		writerTo := &FakeIoWriterTo{ShouldFail: true}
		enc := &FakeEncoder{ShouldWriterTo: writerTo}
		errHandler := &FakeErrorHandler{}

		w := BufWriter{
			Encoder:    enc,
			Writer:     bufio.NewWriter(buf),
			ErrHandler: errHandler.Handle,
		}

		e := NewStdEntry(context.Background(), ERROR, "A log error message", Fields{Bool("test", true)})
		w.WriteEntry(e)
		_ = w.Writer.Flush()

		EntryMatcher(t, e, enc.Entry)
		if got := buf.Bytes(); len(got) != 0 {
			t.Error("could not match data written")
			t.Errorf("got: %s", got)
			t.Errorf("want: %v", nil)
		}
		if !errors.Is(errHandler.Err, ErrEntryNotWritten) || !strings.Contains(errHandler.Err.Error(), ErrFakeWriterTo.Error()) {
			t.Errorf("could not match error")
			t.Errorf("got: %s", errHandler.Err)
			t.Errorf("want: %s", ErrFakeEncoder)
		}
	})
}

func TestBufWriter_Write(t *testing.T) {
	data := []byte(`This is the data written`)
	buf := &bytes.Buffer{}
	writerTo := bytes.NewBuffer(data)
	enc := &FakeEncoder{ShouldWriterTo: writerTo}
	errHandler := &FakeErrorHandler{}

	w := &BufWriter{
		Encoder:         enc,
		Writer:          bufio.NewWriter(buf),
		ErrHandler:      errHandler.Handle,
		DefaultLogLevel: DEBUG,
	}
	log.SetFlags(0)
	log.SetOutput(w)
	log.Printf("%s", data)
	_ = w.Writer.Flush()

	e := NewStdEntry(context.Background(), DEBUG, fmt.Sprintf("%s\n", data), nil)
	EntryMatcher(t, e, enc.Entry)

	if got := buf.Bytes(); !bytes.Equal(got, data) {
		t.Error("could not match data written")
		t.Errorf("got: %s", got)
		t.Errorf("want: %s", data)
	}
}

func TestLeveledWriter(t *testing.T) {
	msg := []byte("A default message")

	defaultWriter := &FakeWriter{}
	debugWriter := &FakeWriter{}
	infoWriter := &FakeWriter{}
	warnWriter := &FakeWriter{}
	errorWriter := &FakeWriter{}
	fatalWriter := &FakeWriter{}

	w := NewLeveledWriter(
		defaultWriter,
		DefaultLeveledWriterOptionFunc(DEBUG, debugWriter),
		DefaultLeveledWriterOptionFunc(INFO, infoWriter),
		DefaultLeveledWriterOptionFunc(WARN, warnWriter),
		DefaultLeveledWriterOptionFunc(ERROR, errorWriter),
		DefaultLeveledWriterOptionFunc(FATAL, fatalWriter),
	)

	w.WriteEntry(debugEntry)
	w.WriteEntry(infoEntry)
	w.WriteEntry(warnEntry)
	w.WriteEntry(errorEntry)
	w.WriteEntry(fatalEntry)
	if _, err := w.Write(msg); err != nil {
		t.Errorf("could not write: %s", err)
	}

	if debugWriter.Entry.Level() != DEBUG {
		t.Errorf("could not match level in the debug writer, got: %s", debugWriter.Entry.Level())
	}

	if infoWriter.Entry.Level() != INFO {
		t.Errorf("could not match level in the info writer, got: %s", infoWriter.Entry.Level())
	}

	if warnWriter.Entry.Level() != WARN {
		t.Errorf("could not match level in the warn writer, got: %s", warnWriter.Entry.Level())
	}

	if errorWriter.Entry.Level() != ERROR {
		t.Errorf("could not match level in the error writer, got: %s", errorWriter.Entry.Level())
	}

	if fatalWriter.Entry.Level() != FATAL {
		t.Errorf("could not match level in the fatal writer, got: %s", fatalWriter.Entry.Level())
	}

	if defaultWriter.Message != Message(msg) {
		t.Errorf("could not match message in the default writer, got: %s", defaultWriter.Entry.Level())
	}
}

func TestDeduplicatorWriter(t *testing.T) {
	msg := []byte("A default message")

	defaultWriter := &FakeWriter{}

	w := NewDeduplicatorWriter(defaultWriter)

	entry := debugEntry.With(
		String("duplicate", "value_0"),
		String("duplicate", "value_1"),
		String("duplicate", "value_2"),
	)

	w.WriteEntry(entry)
	if _, err := w.Write(msg); err != nil {
		t.Errorf("could not write: %s", err)
	}

	if defaultWriter.Entry.Level() != DEBUG {
		t.Errorf("could not match level in the debug writer, got: %s", defaultWriter.Entry.Level())
	}

	wantEntry := debugEntry.With(
		String("duplicate", "value_0"),
		String("duplicate_1", "value_1"),
		String("duplicate_2", "value_2"),
	)
	EntryMatcher(t, defaultWriter.Entry, wantEntry)

	if defaultWriter.Message != Message(msg) {
		t.Errorf("could not match message in the default writer, got: %s", defaultWriter.Entry.Level())
	}
}

func TestMultiWriter(t *testing.T) {

	w1 := &FakeWriter{}
	w2 := &FakeWriter{}
	w := NewMultiWriter(w1, w2)

	w.WriteEntry(debugEntry)

	if w1.Entry.Message() != debugEntry.Message() {
		t.Errorf("could not match message first writer, got: %s", w1.Entry.Message())
	}

	if w2.Entry.Message() != debugEntry.Message() {
		t.Errorf("could not match message second writer, got: %s", w2.Entry.Message())
	}

	msg2 := []byte("Another default message")
	if _, err := w.Write(msg2); err != nil {
		t.Errorf("could not write: %s", err)
	}
	if w1.Message != Message(msg2) {
		t.Errorf("could not match message first writer, got: %s", w1.Entry.Message())
	}

	if w2.Message != Message(msg2) {
		t.Errorf("could not match message second writer, got: %s", w2.Entry.Message())
	}
}

func TestTickFlusher(t *testing.T) {
	fflush := &FakeFlusher{}
	tflush := NewTickFlusher(fflush, 100*time.Millisecond)
	go func() {
		if err := tflush.Flush(); err != nil {
			t.Errorf("could not flush")
		}
	}()
	time.Sleep(410 * time.Millisecond)
	if err := tflush.Close(); err != nil {
		t.Errorf("could not close")
	}
	if got, want := fflush.Counter, int64(5); got != want {
		t.Error("could not match expected flushes")
		t.Errorf("got: %d", got)
		t.Errorf("want: %d", want)
	}
}
