package golog_test

import (
	"context"
	"testing"

	. "github.com/damianopetrungaro/golog"
)

func TestStdLogger(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		w         Writer
		log       func(StdLogger, context.Context, Message)
		wantEntry Entry
	}{
		"debug logger must write an debug entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: debugEntry,
		},
		"debug logger must write an info entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: infoEntry,
		},
		"debug logger must write an warn entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Warn(ctx, msg) },
			wantEntry: warnEntry,
		},
		"debug logger must write an error entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: errorEntry,
		},
		"debug logger must write an fatal entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
		"info logger must write an debug entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: debugEntry,
		},
		"info logger must write an info entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: infoEntry,
		},
		"info logger must write an warn entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Warn(ctx, msg) },
			wantEntry: warnEntry,
		},
		"info logger must write an error entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: errorEntry,
		},
		"info logger must write an fatal entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
		"warn logger must write an debug entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: debugEntry,
		},
		"warn logger must write an info entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: infoEntry,
		},
		"warn logger must write an warn entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Warn(ctx, msg) },
			wantEntry: warnEntry,
		},
		"warn logger must write an error entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: errorEntry,
		},
		"warn logger must write an fatal entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
		"error logger must write an debug entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: debugEntry,
		},
		"error logger must write an info entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: infoEntry,
		},
		"error logger must write an warn entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Warn(ctx, msg) },
			wantEntry: warnEntry,
		},
		"error logger must write an error entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: errorEntry,
		},
		"error logger must write an fatal entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
		"fatal logger must write an debug entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: debugEntry,
		},
		"fatal logger must write an info entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: infoEntry,
		},
		"fatal logger must write an warn entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Warn(ctx, msg) },
			wantEntry: warnEntry,
		},
		"fatal logger must write an error entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: errorEntry,
		},
		"fatal logger must write an fatal entry": {
			log:       func(l StdLogger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			w := &FakeWriter{}

			defer func() {
				EntryMatcher(t, test.wantEntry, w.Entry)
			}()
			defer func() {
				// nolint:staticcheck
				if r := recover(); r != nil {
				}
			}()

			test.log(New(w), ctx, msg)
		})
	}
}

func TestStdLogger_With(t *testing.T) {
	var l Logger = New(&FakeWriter{}, OptionFunc(func(l StdLogger) StdLogger {
		return l.WithDecorator(DecoratorFunc(func(e Entry) Entry { return e }))
	}))

	flds := Fields{String("a", "A"), String("b", "B")}
	l = l.With(flds...)
	if len(l.(StdLogger).Decorators) != 1 {
		t.Fatal("could not match decorators")
	}
	FieldMatcher(t, flds, l.(StdLogger).Fields)

	otherFlds := Fields{String("c", "C"), String("d", "D")}
	l = l.With(otherFlds...)
	if len(l.(StdLogger).Decorators) != 1 {
		t.Fatal("could not match decorators")
	}
	FieldMatcher(t, append(flds, otherFlds...), l.(StdLogger).Fields)
}
