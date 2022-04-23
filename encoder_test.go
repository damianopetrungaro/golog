package golog_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

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
	timestampLayout := "2006-01-02"
	cfg := DefaultJsonConfig()
	cfg.TimestampLayout = timestampLayout

	tests := map[string]struct {
		entry   Entry
		wantLog string
	}{
		"simple entry": {
			entry:   NewStdEntry(context.Background(), DEBUG, "message", nil),
			wantLog: fmt.Sprintf(`{"level":"DEBUG","timestamp":"%s","message":"message"}\n`, time.Now().Format(timestampLayout)),
		},
		"entry with string": {
			entry:   NewStdEntry(context.Background(), INFO, "string message", Fields{String("name", "golog")}),
			wantLog: fmt.Sprintf(`{"level":"INFO","timestamp":"%s","message":"string message","fields":{"name":"golog"}}\n`, time.Now().Format(timestampLayout)),
		},
		"entry with int": {
			entry:   NewStdEntry(context.Background(), WARN, "int message", Fields{Int("number", 101)}),
			wantLog: fmt.Sprintf(`{"level":"WARN","timestamp":"%s","message":"int message","fields":{"number":101}}\n`, time.Now().Format(timestampLayout)),
		},
		"entry with string and int": {
			entry:   NewStdEntry(context.Background(), ERROR, "string and int message", Fields{String("name", "golog"), Int("number", 101)}),
			wantLog: fmt.Sprintf(`{"level":"ERROR","timestamp":"%s","message":"string and int message","fields":{"name":"golog","number":101}}\n`, time.Now().Format(timestampLayout)),
		},
		"entry with array of booleans and array of floats": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of booleans and array of floats", Fields{Bools("booleans", []bool{true, false, false, true}), Float64s("numbers", []float64{1, 1.34, 74.343})}),
			wantLog: fmt.Sprintf(`{"level":"DEBUG","timestamp":"%s","message":"array of booleans and array of floats","fields":{"booleans":[true,false,false,true],"numbers":[1.0000000000,1.3400000000,74.3430000000]}}\n`, time.Now().Format(timestampLayout)),
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
