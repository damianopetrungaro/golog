package golog_test

import (
	"context"
	"strings"
	"testing"

	. "github.com/damianopetrungaro/golog/v2"
)

func TestStackTraceDecorator_Decorate(t *testing.T) {
	const fieldName = "stacktrace"

	w := &FakeWriter{}
	logger := New(w).WithDecorator(NewStackTraceDecorator(fieldName, 2))

	wantStack := []string{
		"golog/decorator_test.go:000",
		"testing/testing.go:111",
	}

	ctx := context.Background()

	logger.With(String("hello", "world")).Info(ctx, "An info message")
	for _, f := range w.Entry.Fields() {
		if f.Key() != fieldName {
			continue
		}
		stack, ok := f.Value().([]string)
		if !ok {
			break
		}

		for i, trace := range stack {
			if strings.Contains(trace, "golog/decorator_test.go") {
				stack[i] = "golog/decorator_test.go:000"
			}
			if strings.Contains(trace, "testing/testing.go") {
				stack[i] = "testing/testing.go:111"
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
