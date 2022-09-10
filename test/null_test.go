package test

import (
	"context"
	"testing"

	"github.com/damianopetrungaro/golog/v2"
)

func TestInMemWriter(t *testing.T) {
	ctx := context.Background()

	w := NewInMemWriter()
	entryOne := golog.NewStdEntry(ctx, golog.INFO, "Message one", golog.Fields{})
	entryTwo := golog.NewStdEntry(ctx, golog.WARN, "Message two", golog.Fields{})
	w.WriteEntry(entryOne)
	w.WriteEntry(entryTwo)
	if _, err := w.Write([]byte("Message three")); err != nil {
		t.Fatalf("could not write: %s", err)
	}
	if err := MatchEntry(entryOne, w.Entries[0]); err != nil {
		t.Fatalf("could not match first entry: %s", err)
	}
	if err := MatchEntry(entryTwo, w.Entries[1]); err != nil {
		t.Fatalf("could not match second entry: %s", err)
	}
	if err := MatchEntry(golog.NewStdEntry(ctx, golog.DEBUG, "Message three", golog.Fields{}), w.Entries[2]); err != nil {
		t.Fatalf("could not match third entry: %s", err)
	}
}

func TestNewNullLogger(t *testing.T) {
	logger := NewNullLogger()
	logger.Error(context.Background(), "An error")
}

func TestNewNullCheckLogger(t *testing.T) {
	logger := NewNullCheckLogger()
	logger.CheckError(context.Background(), "An error")
}

func TestNewNullLoggerWithWriter(t *testing.T) {
	ctx := context.Background()

	w := NewInMemWriter()
	logger := NewFakeLogger(w)
	logger.Info(ctx, "Message one")
	logger.Error(ctx, "Message two")

	entryOne := golog.NewStdEntry(ctx, golog.INFO, "Message one", golog.Fields{})
	entryTwo := golog.NewStdEntry(ctx, golog.ERROR, "Message two", golog.Fields{})

	if err := MatchEntry(entryOne, w.Entries[0]); err != nil {
		t.Fatalf("could not match first entry: %s", err)
	}
	if err := MatchEntry(entryTwo, w.Entries[1]); err != nil {
		t.Fatalf("could not match second entry: %s", err)
	}
}

func TestNewNullCheckLoggerWithWriter(t *testing.T) {
	ctx := context.Background()

	w := NewInMemWriter()
	logger := NewFakeCheckLogger(w)
	logger.Info(ctx, "Message one")
	logger.Error(ctx, "Message two")

	entryOne := golog.NewStdEntry(ctx, golog.INFO, "Message one", golog.Fields{})
	entryTwo := golog.NewStdEntry(ctx, golog.ERROR, "Message two", golog.Fields{})

	if err := MatchEntry(entryOne, w.Entries[0]); err != nil {
		t.Fatalf("could not match first entry: %s", err)
	}
	if err := MatchEntry(entryTwo, w.Entries[1]); err != nil {
		t.Fatalf("could not match first entry: %s", err)
	}
}
