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
	cfg := DefaultJsonConfig()
	date, err := time.Parse("2006-01-02", "2000-12-25")
	if err != nil {
		t.Fatalf("could not create time.Time: %s", err)
	}

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
		"entry with string with special chars": {
			entry: NewStdEntry(context.Background(), INFO, "string message", Fields{String("name", `golog
is	 working!"\`)}),
			wantLog: fmt.Sprintln(`{"level":"INFO","message":"string message","name":"golog\nis\t working!\"\\"}`),
		},
		"entry with int, int8, int16, int32, int64": {
			entry:   NewStdEntry(context.Background(), WARN, "int, int8, int16, int32, int64 message", Fields{Int("int", 100), Int8("int8", 101), Int16("int16", 102), Int32("int32", 103), Int64("int64", 104)}),
			wantLog: fmt.Sprintln(`{"level":"WARN","message":"int, int8, int16, int32, int64 message","int":100,"int8":101,"int16":102,"int32":103,"int64":104}`),
		},
		"entry with uint, uint8, uint16, uint32, uint64": {
			entry:   NewStdEntry(context.Background(), WARN, "uint, uint8, uint16, uint32, uint64 message", Fields{Uint("uint", 100), Uint8("uint8", 101), Uint16("uint16", 102), Uint32("uint32", 103), Uint64("uint64", 104)}),
			wantLog: fmt.Sprintln(`{"level":"WARN","message":"uint, uint8, uint16, uint32, uint64 message","uint":100,"uint8":101,"uint16":102,"uint32":103,"uint64":104}`),
		},
		"entry with arrays of int, int8, int16, int32, int64": {
			entry:   NewStdEntry(context.Background(), WARN, "arrays of int, int8, int16, int32, int64 message", Fields{Ints("ints", []int{100, 110}), Int8s("int8s", []int8{101, 111}), Int16s("int16s", []int16{102, 112}), Int32s("int32s", []int32{103, 113}), Int64s("int64s", []int64{104, 114})}),
			wantLog: fmt.Sprintln(`{"level":"WARN","message":"arrays of int, int8, int16, int32, int64 message","ints":[100,110],"int8s":[101,111],"int16s":[102,112],"int32s":[103,113],"int64s":[104,114]}`),
		},
		"entry with arrays of uint, uint8, uint16, uint32, uint64": {
			entry:   NewStdEntry(context.Background(), WARN, "arrays of uint, uint8, uint16, uint32, uint64 message", Fields{Uints("uints", []uint{100, 110}), Uint8s("uint8s", []uint8{101, 111}), Uint16s("uint16s", []uint16{102, 112}), Uint32s("uint32s", []uint32{103, 113}), Uint64s("uint64s", []uint64{104, 114})}),
			wantLog: fmt.Sprintln(`{"level":"WARN","message":"arrays of uint, uint8, uint16, uint32, uint64 message","uints":[100,110],"uint8s":[101,111],"uint16s":[102,112],"uint32s":[103,113],"uint64s":[104,114]}`),
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
		"entry with float64 and float32": {
			entry:   NewStdEntry(context.Background(), DEBUG, "float64 and float32", Fields{Float64("float64", 12.19), Float32("float32", 21.0101)}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"float64 and float32","float64":12.1900000000,"float32":21.0100994110}`),
		},
		"entry with array of float32": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of float32", Fields{Float32s("float32s", []float32{1.1, 2.2, 3.3})}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"array of float32","float32s":[1.1000000238,2.2000000477,3.2999999523]}`),
		},
		"entry with array of time": {
			entry:   NewStdEntry(context.Background(), DEBUG, "time and an array of time", Fields{Time("25 Dec", date), Times("26/27 Dec", []time.Time{date.AddDate(0, 0, 1), date.AddDate(0, 0, 2)})}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"time and an array of time","25 Dec":"2000-12-25T00:00:00Z","26/27 Dec":["2000-12-26T00:00:00Z","2000-12-27T00:00:00Z"]}`),
		},
		"entry with mapper": {
			entry:   NewStdEntry(context.Background(), DEBUG, "mapper", Fields{Mapper("user", user{ID: "1", Reference: 321, Birthdate: date})}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"mapper","user":{"id":"1","ref":321,"birthdate":"2000-12-25T00:00:00Z"}}`),
		},
		"entry with array of mapper": {
			entry:   NewStdEntry(context.Background(), DEBUG, "mapper", Fields{Mappers("users", []FieldMapper{user{ID: "1", Reference: 321, Birthdate: date}, user{ID: "2", Reference: 123, Birthdate: date}})}),
			wantLog: fmt.Sprintln(`{"level":"DEBUG","message":"mapper","users":[{"id":"1","ref":321,"birthdate":"2000-12-25T00:00:00Z"},{"id":"2","ref":123,"birthdate":"2000-12-25T00:00:00Z"}]}`),
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
	date, err := time.Parse("2006-01-02", "2000-12-25")
	if err != nil {
		t.Fatalf("could not create time.Time: %s", err)
	}

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
		"entry with int, int8, int16, int32, int64": {
			entry:   NewStdEntry(context.Background(), WARN, "int, int8, int16, int32, int64 message", Fields{Int("int", 100), Int8("int8", 101), Int16("int16", 102), Int32("int32", 103), Int64("int64", 104)}),
			wantLog: fmt.Sprintln(`level="WARN" message="int, int8, int16, int32, int64 message" int=100 int8=101 int16=102 int32=103 int64=104`),
		},
		"entry with uint, uint8, uint16, uint32, uint64": {
			entry:   NewStdEntry(context.Background(), WARN, "uint, uint8, uint16, uint32, uint64 message", Fields{Uint("uint", 100), Uint8("uint8", 101), Uint16("uint16", 102), Uint32("uint32", 103), Uint64("uint64", 104)}),
			wantLog: fmt.Sprintln(`level="WARN" message="uint, uint8, uint16, uint32, uint64 message" uint=100 uint8=101 uint16=102 uint32=103 uint64=104`),
		},
		"entry with arrays of int, int8, int16, int32, int64": {
			entry:   NewStdEntry(context.Background(), WARN, "arrays of int, int8, int16, int32, int64 message", Fields{Ints("ints", []int{100, 110}), Int8s("int8s", []int8{101, 111}), Int16s("int16s", []int16{102, 112}), Int32s("int32s", []int32{103, 113}), Int64s("int64s", []int64{104, 114})}),
			wantLog: fmt.Sprintln(`level="WARN" message="arrays of int, int8, int16, int32, int64 message" ints=[100,110] int8s=[101,111] int16s=[102,112] int32s=[103,113] int64s=[104,114]`),
		},
		"entry with arrays of uint, uint8, uint16, uint32, uint64": {
			entry:   NewStdEntry(context.Background(), WARN, "arrays of uint, uint8, uint16, uint32, uint64 message", Fields{Uints("uints", []uint{100, 110}), Uint8s("uint8s", []uint8{101, 111}), Uint16s("uint16s", []uint16{102, 112}), Uint32s("uint32s", []uint32{103, 113}), Uint64s("uint64s", []uint64{104, 114})}),
			wantLog: fmt.Sprintln(`level="WARN" message="arrays of uint, uint8, uint16, uint32, uint64 message" uints=[100,110] uint8s=[101,111] uint16s=[102,112] uint32s=[103,113] uint64s=[104,114]`),
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
		"entry with float64 and float32": {
			entry:   NewStdEntry(context.Background(), DEBUG, "float64 and float32", Fields{Float64("float64", 12.19), Float32("float32", 21.0101)}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="float64 and float32" float64=12.1900000000 float32=21.0100994110`),
		},
		"entry with array of float32": {
			entry:   NewStdEntry(context.Background(), DEBUG, "array of float32", Fields{Float32s("float32s", []float32{1.1, 2.2, 3.3})}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="array of float32" float32s=[1.1000000238,2.2000000477,3.2999999523]`),
		},
		"entry with array of time": {
			entry:   NewStdEntry(context.Background(), DEBUG, "time and an array of time", Fields{Time("25 Dec", date), Times("26/27 Dec", []time.Time{date.AddDate(0, 0, 1), date.AddDate(0, 0, 2)})}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="time and an array of time" 25 Dec="2000-12-25T00:00:00Z" 26/27 Dec=["2000-12-26T00:00:00Z","2000-12-27T00:00:00Z"]`),
		},
		"entry with mapper": {
			entry:   NewStdEntry(context.Background(), DEBUG, "mapper", Fields{Mapper("user", user{ID: "1", Reference: 321, Birthdate: date})}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="mapper" user=[id="1" ref=321 birthdate="2000-12-25T00:00:00Z"]`),
		},
		"entry with array of mapper": {
			entry:   NewStdEntry(context.Background(), DEBUG, "mapper", Fields{Mappers("users", []FieldMapper{user{ID: "1", Reference: 321, Birthdate: date}, user{ID: "2", Reference: 123, Birthdate: date}})}),
			wantLog: fmt.Sprintln(`level="DEBUG" message="mapper" users=[id="1" ref=321 birthdate="2000-12-25T00:00:00Z",id="2" ref=123 birthdate="2000-12-25T00:00:00Z"]`),
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

func TestTextEncoder_EncodeWithCustomLevelFormatter(t *testing.T) {
	cfg := DefaultTextConfig()
	cfg.LevelFormatter = func(l Level) string {
		var colour string

		switch l {
		case DEBUG:
			colour = COLOUR_GREEN
		case INFO:
			colour = COLOUR_BLUE
		case WARN:
			colour = COLOUR_YELLOW
		case ERROR:
			colour = COLOUR_RED
		case FATAL:
			colour = COLOUR_REDBG
		default:
			return l.String()
		}

		return colour + l.String() + COLOUR_RESET
	}

	unknownLevel, _ := ParseLevel("unknown")

	tests := map[string]struct {
		entry   Entry
		wantLog string
	}{
		"debug entry": {
			entry:   NewStdEntry(context.Background(), DEBUG, "message", nil),
			wantLog: fmt.Sprintf(`level="%sDEBUG%s" message="message"%s`, COLOUR_GREEN, COLOUR_RESET, "\n"),
		},
		"info entry": {
			entry:   NewStdEntry(context.Background(), INFO, "message", nil),
			wantLog: fmt.Sprintf(`level="%sINFO%s" message="message"%s`, COLOUR_BLUE, COLOUR_RESET, "\n"),
		},
		"warn entry": {
			entry:   NewStdEntry(context.Background(), WARN, "message", nil),
			wantLog: fmt.Sprintf(`level="%sWARN%s" message="message"%s`, COLOUR_YELLOW, COLOUR_RESET, "\n"),
		},
		"error entry": {
			entry:   NewStdEntry(context.Background(), ERROR, "message", nil),
			wantLog: fmt.Sprintf(`level="%sERROR%s" message="message"%s`, COLOUR_RED, COLOUR_RESET, "\n"),
		},
		"fatal entry": {
			entry:   NewStdEntry(context.Background(), FATAL, "message", nil),
			wantLog: fmt.Sprintf(`level="%sFATAL%s" message="message"%s`, COLOUR_REDBG, COLOUR_RESET, "\n"),
		},
		"default entry": {
			entry:   NewStdEntry(context.Background(), unknownLevel, "message", nil),
			wantLog: fmt.Sprintf(`level="" message="message"%s`, "\n"),
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
