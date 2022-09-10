package opentelemetry_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"

	"github.com/damianopetrungaro/golog/v2"
	. "github.com/damianopetrungaro/golog/v2/opentelemetry"
)

func TestTraceDecorator(t *testing.T) {
	t.Run("context with tracing", func(t *testing.T) {
		ctx, span := otel.Tracer("-").Start(context.Background(), "-")
		var e golog.Entry = golog.NewStdEntry(ctx, golog.DEBUG, "", nil)

		flds := TraceDecorator().Decorate(e).(golog.StdEntry).Fields()

		if len(flds) != 2 {
			t.Fatal("could not match fields")
		}
		if flds[0].Key() != "trace_id" {
			t.Error("could not match trace key")
		}
		if flds[1].Key() != "span_id" {
			t.Error("could not match span key")
		}
		if flds[0].Value() != span.SpanContext().TraceID().String() {
			t.Error("could not match trace value")
		}
		if flds[1].Value() != span.SpanContext().SpanID().String() {
			t.Error("could not match span value")
		}
	})

	t.Run("context with no tracing", func(t *testing.T) {
		var e golog.Entry = golog.NewStdEntry(context.Background(), golog.DEBUG, "", nil)
		flds := TraceDecorator().Decorate(e).(golog.StdEntry).Fields()

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

func TestCustomTraceDecoratorOption(t *testing.T) {
	var e golog.Entry = golog.NewStdEntry(context.Background(), golog.DEBUG, "", nil)
	d := CustomTraceDecorator(
		"test.trace_id",
		func(traceID string) string {
			return "trace_id"
		},
		"test.span_id",
		func(spanID string) string {
			return "span_id"
		},
	)
	flds := d.Decorate(e).(golog.StdEntry).Fields()

	if len(flds) != 2 {
		t.Fatal("could not match fields")
	}
	if flds[0].Key() != "test.trace_id" {
		t.Error("could not match trace key")
	}
	if flds[1].Key() != "test.span_id" {
		t.Error("could not match span key")
	}
	if flds[0].Value() != "trace_id" {
		t.Error("could not match trace value")
	}
	if flds[1].Value() != "span_id" {
		t.Error("could not match span value")
	}
}
