package datadog

import (
	"testing"
)

func TestTraceIDFromHexFormat(t *testing.T) {
	tests := map[string]struct {
		traceID string
		want    string
	}{
		"correct trace id": {
			traceID: "0ef42e66522e5a9183241106963acf99",
			want:    "9449696638118055833",
		},
		"too short trace id": {
			traceID: "0123",
			want:    "",
		},
		"invalid trace id": {
			traceID: "xyzxyzxyzxyzxyzxyzxyzxyzxyzxyz",
			want:    "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := TraceIDFromHexFormat(test.traceID)

			if result != test.want {
				t.Errorf("expected result trace id to be %q got %q", test.want, result)
			}
		})
	}
}

func TestSpanIDFromHexFormat(t *testing.T) {
	tests := map[string]struct {
		spanID string
		want   string
	}{
		"correct span id": {
			spanID: "0102040810203040",
			want:   "72624976668143680",
		},
		"too short span id": {
			spanID: "0123",
			want:   "",
		},
		"invalid span id": {
			spanID: "xyzxyzxyzxyzxyzxyzxyzxyzxyzxyz",
			want:   "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := SpanIDFromHexFormat(test.spanID)

			if result != test.want {
				t.Errorf("expected result span id to be %q got %q", test.want, result)
			}
		})
	}
}
