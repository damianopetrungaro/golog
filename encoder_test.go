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
			wantLog: fmt.Sprintln(`{"level":"INFO","message":"string message","name":"golog"}`),
		},
		"entry with int": {
			entry:   NewStdEntry(context.Background(), WARN, "int message", Fields{Int("number", 101)}),
			wantLog: fmt.Sprintln(`{"level":"WARN","message":"int message","number":101}`),
		},
		"entry with string and int": {
			entry:   NewStdEntry(context.Background(), ERROR, "string and int message", Fields{String("name", "golog"), Int("number", 101)}),
			wantLog: fmt.Sprintln(`{"level":"ERROR","message":"string and int message","name":"golog","number":101}`),
		},
		"entry with array of booleans and array of floats": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of booleans and array of floats", Fields{Bools("booleans", []bool{true, false, false, true}), Float64s("numbers", []float64{1, 1.34, 74.343})}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"array of booleans and array of floats","booleans":[true,false,false,true],"numbers":[1.0000000000,1.3400000000,74.3430000000]}`),
		},
		"entry with boolean and error": {
			entry:   NewStdEntry(context.Background(), DEBUG, "boolean and error", Fields{Bool("boolean", true), Err(fmt.Errorf("ops!"))}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"boolean and error","boolean":true,"error":"ops!"}`),
		},
		"entry with array of strings and nil error": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of strings and nil error", Fields{Strings("strings", []string{"a", "b", "c"}), Err(nil)}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"array of strings and nil error","strings":["a","b","c"],"error":"<nil>"}`),
		},
		"entry with array of errors and uint": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of errors and uint", Fields{Errs([]error{fmt.Errorf("ops 1"), nil, fmt.Errorf("ops 2")}), Uint("uint", 12)}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"array of errors and uint","errors":["ops 1","<nil>","ops 2"],"uint":12}`),
		},
		"entry with array of int and array of uint": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of int and array of uint", Fields{Ints("ints", []int{-10, 5, 10}), Uints("uints", []uint{0, 10, 20})}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"array of int and array of uint","ints":[-10,5,10],"uints":[0,10,20]}`),
		},
		"entry with float64 and float32": {
			entry:   NewStdEntry(context.Background(), DEBUG, "float64 and float32", Fields{Float64("float64", 12.19), Float32("float32", 21.0101)}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"float64 and float32","float64":12.1900000000,"float32":21.0100994110}`),
		},
		"entry with array of float32": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of float32", Fields{Float32s("float32s", []float32{1.1, 2.2, 3.3})}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"array of float32","float32s":[1.1000000238,2.2000000477,3.2999999523]}`),
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

func TestTextEncoder_Encode(t *testing.T) {
	cfg := DefaultTextConfig()

	tests := map[string]struct {
		entry   Entry
		wantLog string
	}{
		"simple entry": {
			entry:   NewStdEntry(context.Background(), DEBUG, "message", nil),
			wantLog: fmt.Sprintln(`level="DEBUG" message="message"`),
		},
		"entry with string": {
			entry:   NewStdEntry(context.Background(), INFO, "string message", Fields{String("name", "golog")}),
			wantLog: fmt.Sprintln(`level="INFO" message="string message" name="golog"`),
		},
		"entry with int": {
			entry:   NewStdEntry(context.Background(), WARN, "int message", Fields{Int("number", 101)}),
			wantLog: fmt.Sprintln(`level="WARN" message="int message" number=101`),
		},
		"entry with string and int": {
			entry:   NewStdEntry(context.Background(), ERROR, "string and int message", Fields{String("name", "golog"), Int("number", 101)}),
			wantLog: fmt.Sprintln(`level="ERROR" message="string and int message" name="golog" number=101`),
		},
		"entry with array of booleans and array of floats": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of booleans and array of floats", Fields{Bools("booleans", []bool{true, false, false, true}), Float64s("numbers", []float64{1, 1.34, 74.343})}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="array of booleans and array of floats" booleans=[true,false,false,true] numbers=[1.0000000000,1.3400000000,74.3430000000]`),
		},
		"entry with boolean and error": {
			entry:   NewStdEntry(context.Background(), DEBUG, "boolean and error", Fields{Bool("boolean", true), Err(fmt.Errorf("ops!"))}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="boolean and error" boolean=true error="ops!"`),
		},
		"entry with array of strings and nil error": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of strings and nil error", Fields{Strings("strings", []string{"a", "b", "c"}), Err(nil)}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="array of strings and nil error" strings=["a","b","c"] error="<nil>"`),
		},
		"entry with array of errors and uint": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of errors and uint", Fields{Errs([]error{fmt.Errorf("ops 1"), nil, fmt.Errorf("ops 2")}), Uint("uint", 12)}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="array of errors and uint" errors=["ops 1","<nil>","ops 2"] uint=12`),
		},
		"entry with array of int and array of uint": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of int and array of uint", Fields{Ints("ints", []int{-10, 5, 10}), Uints("uints", []uint{0, 10, 20})}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="array of int and array of uint" ints=[-10,5,10] uints=[0,10,20]`),
		},
		"entry with float64 and float32": {
			entry:   NewStdEntry(context.Background(), DEBUG, "float64 and float32", Fields{Float64("float64", 12.19), Float32("float32", 21.0101)}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="float64 and float32" float64=12.1900000000 float32=21.0100994110`),
		},
		"entry with array of float32": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of float32", Fields{Float32s("float32s", []float32{1.1, 2.2, 3.3})}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="array of float32" float32s=[1.1000000238,2.2000000477,3.2999999523]`),
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			enc := NewTextEncoder(cfg)
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
