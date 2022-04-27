package golog_test

import (
	"context"
	"strings"
	"testing"

	. "github.com/damianopetrungaro/golog"
)

func TestStackTraceDecorator_Decorate(t *testing.T) {
	const fieldName = "stacktrace"

	w := &FakeWriter{}
	logger := New(w).WithDecorator(NewStackTraceDecorator(fieldName, 2))

	wantStack := []string{
		"golog/decorator_test.go:24",
		"testing/testing.go:1439",
	}

	ctx := context.Background()

	logger.With(Fields{String("hello", "world")}).Info(ctx, "An info message")
	for _, f := range w.Entry.Fields() {
		if f.Key() != fieldName {
			continue
		}
		stack, ok := f.Value().([]string)
		if !ok {
			break
		}

		for i, trace := range stack {
			if strings.HasSuffix(trace, "golog/decorator_test.go:24") {
				stack[i] = "golog/decorator_test.go:24"
			}
			if strings.HasSuffix(trace, "testing/testing.go:1439") {
				stack[i] = "testing/testing.go:1439"
			}
		}

		if strings.Join(stack, "##") != strings.Join(wantStack, "##") {
			t.Error("could not match trace")
			t.Errorf("got: %v", stack)
			t.Errorf("want: %v", wantStack)
		}

		return
	}
	t.Error("could not find stacktrace field")
}