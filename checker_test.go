package golog_test

import (
	"testing"

	. "github.com/damianopetrungaro/golog/v2"
)

func TestLevelChecker_Check(t *testing.T) {
	c := NewLevelChecker(WARN)

	tests := map[string]struct {
		entry Entry
		want  bool
	}{
		"debug entry must not be reported": {
			entry: debugEntry,
			want:  false,
		},
		"info entry must not be reported": {
			entry: infoEntry,
			want:  false,
		},
		"warning entry must not be reported": {
			entry: warnEntry,
			want:  true,
		},
		"error entry must not be reported": {
			entry: errorEntry,
			want:  true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := c.Check(test.entry); test.want != got {
				t.Error("could not match check result")
				t.Errorf("want: %t", test.want)
				t.Errorf("got: %t", got)
			}
		})
	}
}
