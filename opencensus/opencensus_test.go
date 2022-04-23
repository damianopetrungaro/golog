package opencensus_test

import (
	"context"
	"testing"

	"go.opencensus.io/trace"

	. "github.com/damianopetrungaro/golog"
	. "github.com/damianopetrungaro/golog/opencensus"
)

func TestTraceDecorator(t *testing.T) {
	t.Run("context with tracing", func(t *testing.T) {
		ctx, span := trace.StartSpan(context.Background(), "-")
		var e Entry = NewStdEntry(ctx, DEBUG, "", nil)

		flds := TraceDecorator().Decorate(e).(StdEntry).Fields()

		if len(flds) != 2 {
			t.Fatal("could not match fields")
		}
		if flds[0].Key() != "trace_id" {
			t.Error("could not match trace key")
		}
		if flds[1].Key() != "span_id" {
			t.Error("could not match span key")
		}
		if flds[0].Value() != span.SpanContext().TraceID.String() {
			t.Error("could not match trace value")
		}
		if flds[1].Value() != span.SpanContext().SpanID.String() {
			t.Error("could not match span value")
		}
	})

	t.Run("context with no tracing", func(t *testing.T) {
		var e Entry = NewStdEntry(context.Background(), DEBUG, "", nil)
		flds := TraceDecorator().Decorate(e).(StdEntry).Fields()

		if len(flds) != 2 {
			t.Fatal("could not match fields")
		}
		if flds[0].Key() != "trace_id" {
			t.Error("could not match trace key")
		}
		if flds[1].Key() != "span_id" {
			t.Error("could not match span key")
		}
		if flds[0].Value() != "00000000000000000000000000000000" {
			t.Error("could not match trace value")
		}
		if flds[1].Value() != "0000000000000000" {
			t.Error("could not match span value")
		}
	})
}
