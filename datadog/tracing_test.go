package datadog

import (
	"testing"
)

func TestTraceIDFromHexFormat(t *testing.T) {
	traceID := "0ef42e66522e5a9183241106963acf99"
	expectedTraceID := "9449696638118055833"

	result := TraceIDFromHexFormat(traceID)

	if result != expectedTraceID {
		t.Errorf("expected result trace id to be %q got %q", expectedTraceID, result)
	}
}

func TestSpanIDFromHexFormat(t *testing.T) {
	spanID := "0102040810203040"
	expectedSpanID := "72624976668143680"

	result := SpanIDFromHexFormat(spanID)

	if result != expectedSpanID {
		t.Errorf("expected result span id to be %q got %q", expectedSpanID, result)
	}
}
