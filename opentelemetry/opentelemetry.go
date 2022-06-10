package opentelemetry

import (
	"go.opentelemetry.io/otel/trace"

	"github.com/damianopetrungaro/golog"
)

var decorator golog.DecoratorFunc = func(e golog.Entry) golog.Entry {
	span := trace.SpanFromContext(e.Context()).SpanContext()

	return e.With(
		TraceID(span),
		SpanID(span),
	)
}

func TraceDecoratorOption() golog.Option {
	return golog.OptionFunc(func(l golog.StdLogger) golog.StdLogger {
		return l.WithDecorator(decorator)
	})
}

func TraceDecorator() golog.Decorator {
	return decorator
}

func TraceID(span trace.SpanContext) golog.Field {
	const k = "trace_id"
	return golog.String(k, span.TraceID().String())
}

func SpanID(span trace.SpanContext) golog.Field {
	const k = "span_id"
	return golog.String(k, span.SpanID().String())
}
