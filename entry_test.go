package golog_test

import (
	"context"
	"testing"

	. "github.com/damianopetrungaro/golog"
)

func TestStdEntry(t *testing.T) {
	ctx := context.Background()
	lvl := DEBUG
	msg := "This is a message"
	flds := Fields{Int("one", 1), Int("two", 2), Int("three", 3)}
	var e Entry = NewStdEntry(
		ctx,
		lvl,
		msg,
		flds,
	)

	if ctx != e.Context() {
		t.Error("could not match context")
		t.Errorf("got: %s", e.Context())
		t.Errorf("want: %s", ctx)
	}

	if lvl != e.Level() {
		t.Error("could not match level")
		t.Errorf("got: %s", e.Level())
		t.Errorf("want: %s", lvl)
	}

	if msg != e.Message() {
		t.Error("could not match message")
		t.Errorf("got: %s", e.Message())
		t.Errorf("want: %s", msg)
	}

	FieldMatcher(t, flds, e.Fields())

	otherFlds := Fields{Int("four", 4), Int("five", 5)}
	e = e.With(otherFlds)
	FieldMatcher(t, append(flds, otherFlds...), e.Fields())
}

func EntryMatcher(t *testing.T, x, y Entry) {
	t.Helper()
	if x == nil && y == nil {
		return
	}

	stdX, ok := x.(StdEntry)
	if !ok {
		t.Error("x is not a stdEntry")
		return
	}
	stdY, ok := y.(StdEntry)
	if !ok {
		t.Error("y is not a stdEntry")
		return
	}

	if stdX.Ctx != stdY.Ctx {
		t.Error("could not match context")
		t.Errorf("x: %v", stdX.Ctx)
		t.Errorf("y: %v", stdY.Ctx)
	}
	if stdX.Lvl != stdY.Lvl {
		t.Error("could not match level")
		t.Errorf("x: %s", stdX.Lvl)
		t.Errorf("y: %s", stdY.Lvl)
	}
	if stdX.Msg != stdY.Msg {
		t.Error("could not match message")
		t.Errorf("x: %s", stdX.Msg)
		t.Errorf("y: %s", stdY.Msg)
	}

	FieldMatcher(t, stdX.Fields(), stdY.Fields())
}
