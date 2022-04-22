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
