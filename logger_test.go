package golog_test

import (
	"context"
	"testing"

	. "github.com/damianopetrungaro/golog"
)

func TestStdLogger(t *testing.T) {
	ctx := context.Background()
	msg := "This is a log message"

	debugEntry := NewStdEntry(ctx, DEBUG, msg, nil)
	infoEntry := NewStdEntry(ctx, INFO, msg, nil)
	warnEntry := NewStdEntry(ctx, WARN, msg, nil)
	errorEntry := NewStdEntry(ctx, ERROR, msg, nil)
	fatalEntry := NewStdEntry(ctx, FATAL, msg, nil)

	tests := map[string]struct {
		lvl       Level
		w         Writer
		log       func(Logger, context.Context, Message)
		wantEntry Entry
	}{
		"debug logger must write an debug entry": {
			lvl:       DEBUG,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: debugEntry,
		},
		"debug logger must write an info entry": {
			lvl:       DEBUG,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: infoEntry,
		},
		"debug logger must write an warn entry": {
			lvl:       DEBUG,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Warning(ctx, msg) },
			wantEntry: warnEntry,
		},
		"debug logger must write an error entry": {
			lvl:       DEBUG,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: errorEntry,
		},
		"debug logger must write an fatal entry": {
			lvl:       DEBUG,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
		"info logger must write an debug entry": {
			lvl:       INFO,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: nil,
		},
		"info logger must write an info entry": {
			lvl:       INFO,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: infoEntry,
		},
		"info logger must write an warn entry": {
			lvl:       INFO,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Warning(ctx, msg) },
			wantEntry: warnEntry,
		},
		"info logger must write an error entry": {
			lvl:       INFO,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: errorEntry,
		},
		"info logger must write an fatal entry": {
			lvl:       INFO,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
		"warn logger must write an debug entry": {
			lvl:       WARN,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: nil,
		},
		"warn logger must write an info entry": {
			lvl:       WARN,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: nil,
		},
		"warn logger must write an warn entry": {
			lvl:       WARN,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Warning(ctx, msg) },
			wantEntry: warnEntry,
		},
		"warn logger must write an error entry": {
			lvl:       WARN,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: errorEntry,
		},
		"warn logger must write an fatal entry": {
			lvl:       WARN,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
		"error logger must write an debug entry": {
			lvl:       ERROR,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: nil,
		},
		"error logger must write an info entry": {
			lvl:       ERROR,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: nil,
		},
		"error logger must write an warn entry": {
			lvl:       ERROR,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Warning(ctx, msg) },
			wantEntry: nil,
		},
		"error logger must write an error entry": {
			lvl:       ERROR,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: errorEntry,
		},
		"error logger must write an fatal entry": {
			lvl:       ERROR,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
		"fatal logger must write an debug entry": {
			lvl:       FATAL,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Debug(ctx, msg) },
			wantEntry: nil,
		},
		"fatal logger must write an info entry": {
			lvl:       FATAL,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Info(ctx, msg) },
			wantEntry: nil,
		},
		"fatal logger must write an warn entry": {
			lvl:       FATAL,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Warning(ctx, msg) },
			wantEntry: nil,
		},
		"fatal logger must write an error entry": {
			lvl:       FATAL,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Error(ctx, msg) },
			wantEntry: nil,
		},
		"fatal logger must write an fatal entry": {
			lvl:       FATAL,
			log:       func(l Logger, ctx context.Context, msg Message) { l.Fatal(ctx, msg) },
			wantEntry: fatalEntry,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			w := &FakeWriter{}

			defer func() {
				EntryMatcher(t, test.wantEntry, w.Entry)
			}()
			defer func() {
				if r := recover(); r != nil {
				}
			}()

			test.log(New(test.lvl, w), ctx, msg)
		})
	}
}
