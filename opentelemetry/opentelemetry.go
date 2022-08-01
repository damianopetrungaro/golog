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

func CustomTraceDecoratorOption(
	traceIDKey string,
	traceIDTransformer func(traceID string) string,
	spanIDKey string,
	spanIDTransformer func(spanID string) string,
) golog.Option {
	return golog.OptionFunc(func(l golog.StdLogger) golog.StdLogger {
		d := CustomTraceDecorator(traceIDKey, traceIDTransformer, spanIDKey, spanIDTransformer)
		return l.WithDecorator(d)
	})
}

func CustomTraceDecorator(
	traceIDKey string,
	traceIDTransformer func(traceID string) string,
	spanIDKey string,
	spanIDTransformer func(spanID string) string,
) golog.DecoratorFunc {
	return func(e golog.Entry) golog.Entry {
		span := trace.SpanFromContext(e.Context()).SpanContext()

		return e.With(
			golog.String(traceIDKey, traceIDTransformer(span.TraceID().String())),
			golog.String(spanIDKey, spanIDTransformer(span.SpanID().String())),
		)
	}
}
