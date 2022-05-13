package golog_test

import (
	"testing"

	. "github.com/damianopetrungaro/golog"
)

func TestNoopCheckedLogger_Log(t *testing.T) {
	checkedLogger := &NoopCheckedLogger{}
	checkedLogger.Log(String("test", "noop"))
}

func TestStdCheckedLogger_Log(t *testing.T) {
	t.Run("fatal entry", func(t *testing.T) {
		flds := Fields{Bool("key", true)}
		var e Entry = fatalEntry
		w := &FakeWriter{}
		defer func() {
			EntryMatcher(t, w.Entry, e.With(flds...))
		}()
		defer func() {
			// nolint:staticcheck
			if r := recover(); r != nil {
			}
		}()
		checkedLogger := &StdCheckedLogger{Entry: fatalEntry, Writer: w}
		checkedLogger.Log(flds...)
	})

	t.Run("non fatal entry", func(t *testing.T) {
		flds := Fields{Bool("key", true)}
		w := &FakeWriter{}
		var e Entry = debugEntry
		e = e.With(flds...)
		checkedLogger := &StdCheckedLogger{Entry: debugEntry, Writer: w}
		checkedLogger.Log(flds...)
		EntryMatcher(t, w.Entry, e)
	})
}
