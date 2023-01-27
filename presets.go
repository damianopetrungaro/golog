package golog

import (
	"bufio"
	"context"
	"os"
	"time"
)

// NewProductionLogger returns a pre-configured logger for production environment
func NewProductionLogger(lvl Level) (StdLogger, Flusher) {
	return NewWithEncoder(lvl, NewJsonEncoder(DefaultJsonConfig()))
}

// NewDevelopmentLogger returns a pre-configured logger for development environment
func NewDevelopmentLogger(lvl Level) (StdLogger, Flusher) {
	textConfig := DefaultTextConfig()
	textConfig.LevelFormatter = ColoredLevelFormatter()
	return NewWithEncoder(lvl, NewTextEncoder(textConfig))
}

// NewWithEncoder returns a preset logger with a customer Encoder
func NewWithEncoder(lvl Level, enc Encoder) (StdLogger, Flusher) {
	w := NewBufWriter(
		enc,
		bufio.NewWriter(os.Stdout),
		DefaultErrorHandler(),
		lvl,
	)

	return newLogger(lvl, w)
}

func newLogger(lvl Level, wf Writer) (StdLogger, Flusher) {
	log := New(
		wf,
		NewLevelCheckerOption(lvl),
		NewTimestampDecoratorOption("timestamp", time.RFC3339Nano),
		NewStackTraceDecoratorOption("stacktrace", 5),
	)

	// force flushing data to the disk every 5 seconds
	go func() {
		flusher := NewTickFlusher(wf, 5*time.Second)
		if err := flusher.Flush(); err != nil {
			log.With(Err(err)).Warn(context.Background(), "entries were not flushed")
		}
	}()

	// Set the created logger as main one for the golog package as well
	SetLogger(log)
	SetCheckLogger(log)

	return log, wf
}
