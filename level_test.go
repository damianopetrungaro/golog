package golog_test

import (
	"errors"
	"testing"

	. "github.com/damianopetrungaro/golog"
)

func TestParseLevel(t *testing.T) {
	tests := map[string]struct {
		given   string
		want    Level
		wantErr error
	}{
		"parse DEBUG level must return DEBUG":      {given: "DEBUG", want: DEBUG, wantErr: nil},
		"parse INFO level must return INFO":        {given: "INFO", want: INFO, wantErr: nil},
		"parse WARN level must return WARN":        {given: "WARN", want: WARN, wantErr: nil},
		"parse WARNING level must return WARN":     {given: "WARNING", want: WARN, wantErr: nil},
		"parse ERROR level must return ERROR":      {given: "ERROR", want: ERROR, wantErr: nil},
		"parse FATAL level must return FATAL":      {given: "FATAL", want: FATAL, wantErr: nil},
		"parse UNKNOWN level must return an error": {given: "UNKNOWN", want: 0, wantErr: ErrLevelNotParsed},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			got, err := ParseLevel(test.given)

			if got != test.want {
				t.Error("could not match parsed level")
				t.Errorf("got: %s", got)
				t.Errorf("want: %s", test.want)
			}

			if !errors.Is(err, test.wantErr) {
				t.Error("could not match error")
				t.Errorf("got: %v", err)
				t.Errorf("want: %v", test.wantErr)
			}
		})
	}
}

func TestLevel_String(t *testing.T) {
	tests := map[string]struct {
		given Level
		want  string
	}{
		"a DEBUG level must return DEBUG":             {given: DEBUG, want: "DEBUG"},
		"a INFO level must return INFO":               {given: INFO, want: "INFO"},
		"a WARN level must return WARN":               {given: WARN, want: "WARN"},
		"a ERROR level must return ERROR":             {given: ERROR, want: "ERROR"},
		"a FATAL level must return FATAL":             {given: FATAL, want: "FATAL"},
		"a UNKNOWN level must return an empty string": {given: 0, want: ""},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			if got := test.given.String(); got != test.want {
				t.Error("could not match string level")
				t.Errorf("got: %s", got)
				t.Errorf("want: %s", test.want)
			}
		})
	}
}

func TestLevel_ColouredString(t *testing.T) {
	tests := map[string]struct {
		given Level
		want  string
	}{
		"a DEBUG level must return a coloured DEBUG":  {given: DEBUG, want: "\033[32mDEBUG\033[0m"},
		"a INFO level must return a coloured INFO":    {given: INFO, want: "\033[34mINFO\033[0m"},
		"a WARN level must return a coloured WARN":    {given: WARN, want: "\033[33mWARN\033[0m"},
		"a ERROR level must return a coloured ERROR":  {given: ERROR, want: "\033[31mERROR\033[0m"},
		"a FATAL level must return a coloured FATAL":  {given: FATAL, want: "\033[41mFATAL\033[0m"},
		"a UNKNOWN level must return an empty string": {given: 0, want: ""},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			if got := test.given.ColouredString(); got != test.want {
				t.Error("could not match coloured string level")
				t.Errorf("got: %s", got)
				t.Errorf("want: %s", test.want)
			}
		})
	}
}
