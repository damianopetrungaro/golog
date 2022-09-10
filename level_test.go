package golog_test

import (
	"errors"
	"testing"

	. "github.com/damianopetrungaro/golog/v2"
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
