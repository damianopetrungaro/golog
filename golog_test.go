package golog_test

import (
	"context"
	"testing"

	. "github.com/damianopetrungaro/golog"
)

// FakeErrorHandler used for internal testing purposes
type FakeErrorHandler struct {
	Err error
}

func (fe *FakeErrorHandler) Handle(err error) {
	fe.Err = err
}

func TestDebug(t *testing.T) {
	msg := "Hello"
	w := setLoggerHelper(t)
	Debug(context.Background(), msg)
	if w.Entry.Message() != msg {
		t.Error("could not match error message")
	}
}

func TestCheckDebug(t *testing.T) {
	msg := "Hello"
	w := setCheckLoggerHelper(t)
	if c, ok := CheckDebug(context.Background(), msg); ok {
		c.Log(nil)
	}
	if w.Entry.Message() != msg {
		t.Error("could not match error message")
	}
}

func TestInfo(t *testing.T) {
	msg := "Hello"
	w := setLoggerHelper(t)
	Info(context.Background(), msg)
	if w.Entry.Message() != msg {
		t.Error("could not match error message")
	}
}

func TestCheckInfo(t *testing.T) {
	msg := "Hello"
	w := setCheckLoggerHelper(t)
	if c, ok := CheckInfo(context.Background(), msg); ok {
		c.Log(nil)
	}
	if w.Entry.Message() != msg {
		t.Error("could not match error message")
	}
}

func TestWarning(t *testing.T) {
	msg := "Hello"
	w := setLoggerHelper(t)
	Warning(context.Background(), msg)
	if w.Entry.Message() != msg {
		t.Error("could not match error message")
	}
}

func TestCheckWarning(t *testing.T) {
	msg := "Hello"
	w := setCheckLoggerHelper(t)
	if c, ok := CheckWarning(context.Background(), msg); ok {
		c.Log(nil)
	}
	if w.Entry.Message() != msg {
		t.Error("could not match error message")
	}
}

func TestError(t *testing.T) {
	msg := "Hello"
	w := setLoggerHelper(t)
	Error(context.Background(), msg)
	if w.Entry.Message() != msg {
		t.Error("could not match error message")
	}
}

func TestCheckError(t *testing.T) {
	msg := "Hello"
	w := setCheckLoggerHelper(t)
	if c, ok := CheckError(context.Background(), msg); ok {
		c.Log(nil)
	}
	if w.Entry.Message() != msg {
		t.Error("could not match error message")
	}
}

func TestFatal(t *testing.T) {
	msg := "Hello"
	w := setLoggerHelper(t)
	defer func() {
		if w.Entry.Message() != msg {
			t.Error("could not match error message")
		}
	}()
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	Fatal(context.Background(), msg)
}

func TestCheckFatal(t *testing.T) {
	msg := "Hello"
	w := setCheckLoggerHelper(t)
	defer func() {
		if w.Entry.Message() != msg {
			t.Error("could not match error message")
		}
	}()
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	if c, ok := CheckFatal(context.Background(), msg); ok {
		c.Log(nil)
	}
}

func TestWith(t *testing.T) {
	l := New(&FakeWriter{})
	l2 := l.With(Fields{Bool("k", true)}).(StdLogger)
	if len(l2.Fields) == len(l.Fields) {
		t.Error("could match fields length")
	}
	if l2.Fields[0].Value() != true {
		t.Error("could not match field")
	}
}

func setLoggerHelper(t *testing.T) *FakeWriter {
	t.Helper()
	w := &FakeWriter{}
	SetLogger(New(w))
	return w
}

func setCheckLoggerHelper(t *testing.T) *FakeWriter {
	t.Helper()
	w := &FakeWriter{}
	SetCheckLogger(New(w))
	return w
}
