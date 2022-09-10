package test

import (
	"context"
	"testing"

	"github.com/damianopetrungaro/golog/v2"
)

func TestMatchEntry(t *testing.T) {
	ctx := context.Background()
	infoEntry := golog.NewStdEntry(ctx, golog.INFO, "Message one", golog.Fields{})
	errorEntry := golog.NewStdEntry(ctx, golog.ERROR, "Message two", golog.Fields{})

	tests := map[string]struct {
		x, y   golog.Entry
		errMsg string
	}{
		"entries are equal because nil": {
			x: nil,
			y: nil,
		},
		"entries are equal without fields": {
			x: infoEntry,
			y: infoEntry,
		},
		"entries are equal with fields": {
			x: errorEntry.With(golog.String("key", "value")),
			y: errorEntry.With(golog.String("key", "value")),
		},
		"entries are not equal because x is not a std entry": {
			x:      &MockEntry{},
			y:      infoEntry,
			errMsg: "x is not a stdEntry",
		},
		"entries are not equal because y is not a std entry": {
			x:      infoEntry,
			y:      &MockEntry{},
			errMsg: "y is not a stdEntry",
		},
		"entries are not equal because of different fields": {
			x:      errorEntry.With(golog.String("key", "value")),
			y:      errorEntry.With(golog.String("another_key", "value")),
			errMsg: "could not match field value at index 0",
		},
		"entries are not equal because of different level": {
			x:      golog.NewStdEntry(ctx, golog.INFO, "Message one", golog.Fields{}),
			y:      golog.NewStdEntry(ctx, golog.DEBUG, "Message one", golog.Fields{}),
			errMsg: "could not match level",
		},
		"entries are not equal because of different message": {
			x:      golog.NewStdEntry(ctx, golog.INFO, "Message one", golog.Fields{}),
			y:      golog.NewStdEntry(ctx, golog.INFO, "Message two", golog.Fields{}),
			errMsg: "could not match message",
		},
		"entries are not equal because of different fields length": {
			x:      errorEntry.With(golog.String("key", "value")),
			y:      errorEntry,
			errMsg: "could not match fields length",
		},
		"entries are not equal because of different context": {
			x:      golog.NewStdEntry(context.Background(), golog.INFO, "Message one", golog.Fields{}),
			y:      golog.NewStdEntry(context.TODO(), golog.INFO, "Message one", golog.Fields{}),
			errMsg: "could not match context",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := MatchEntry(test.x, test.y)
			switch {
			case test.errMsg != "":
				if err.Error() != test.errMsg {
					t.Fatalf("could not error message: %s", err)
				}
			default:
				if err != nil {
					t.Fatalf("could not match entry: %s", err)
				}
			}
		})
	}
}
