package golog_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	. "github.com/damianopetrungaro/golog"
)

var _ Encoder = &FakeEncoder{}
var ErrFakeEncoder = errors.New("an error occurred on the encoder")

// FakeEncoder used for internal testing purposes
type FakeEncoder struct {
	Entry          Entry
	ShouldFail     bool
	ShouldWriterTo io.WriterTo
}

func (fe *FakeEncoder) Encode(e Entry) (io.WriterTo, error) {
	fe.Entry = e
	if fe.ShouldFail {
		return nil, ErrFakeEncoder
	}
	return fe.ShouldWriterTo, nil
}

func TestJsonEncoder_Encode(t *testing.T) {
	cfg := DefaultJsonConfig()

	tests := map[string]struct {
		entry   Entry
		wantLog string
	}{
		"simple entry": {
			entry:   NewStdEntry(context.Background(), DEBUG, "message", nil),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"message"}`),
		},
		"entry with string": {
			entry:   NewStdEntry(context.Background(), INFO, "string message", Fields{String("name", "golog")}),
			wantLog: fmt.Sprintln(`{"level":"INFO","message":"string message","fields":{"name":"golog"}}`),
		},
		"entry with int": {
			entry:   NewStdEntry(context.Background(), WARN, "int message", Fields{Int("number", 101)}),
			wantLog: fmt.Sprintln(`{"level":"WARN","message":"int message","fields":{"number":101}}`),
		},
		"entry with string and int": {
			entry:   NewStdEntry(context.Background(), ERROR, "string and int message", Fields{String("name", "golog"), Int("number", 101)}),
			wantLog: fmt.Sprintln(`{"level":"ERROR","message":"string and int message","fields":{"name":"golog","number":101}}`),
		},
		"entry with array of booleans and array of floats": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of booleans and array of floats", Fields{Bools("booleans", []bool{true, false, false, true}), Float64s("numbers", []float64{1, 1.34, 74.343})}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"array of booleans and array of floats","fields":{"booleans":[true,false,false,true],"numbers":[1.0000000000,1.3400000000,74.3430000000]}}`),
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			enc := NewJsonEncoder(cfg)
			w, err := enc.Encode(test.entry)
			if err != nil {
				t.Errorf("could not encode: %s", err)
			}

			buf := &bytes.Buffer{}
			if _, err := w.WriteTo(buf); err != nil {
				t.Errorf("could not write: %s", err)
			}

			if got := buf.String(); got != test.wantLog {
				t.Error("could not match log")
				t.Errorf("want: %s", test.wantLog)
				t.Errorf("got: %s", got)
			}
		})
	}
}
