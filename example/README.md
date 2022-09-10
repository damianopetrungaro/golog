# Examples

## Production ready logger usage with opentelemetry support

```go
import (
    "github.com/damianopetrungaro/golog/v2"
    "github.com/damianopetrungaro/golog/v2/opentelemetry"
)

func NewLogger(lvl golog.Level) (golog.Logger, golog.Flusher) {
	w := golog.NewBufWriter(
        golog.NewJsonEncoder(),
		bufio.NewWriter(os.Stdout),
        golog.DefaultErrorHandler(),
		lvl,
	)
	
	logger := golog.New(
		w,
		golog.NewLevelCheckerOption(lvl),
		golog.NewTimestampDecoratorOption("timestamp", time.RFC3339Nano),
		golog.NewStackTraceDecoratorOption("stacktrace", 5),
		opentelemetry.TraceDecoratorOption(),
	)

	// force flushing data to the disk every 5 seconds
	flusher := golog.NewTickFlusher(stdoutWriter, 5*time.Second)
	go func() {
		if err := flusher.Flush(); err != nil {
			logger.With(golog.Err(err)).Warn(ctx, "log entries were not flushed")
		}
	}()

	golog.SetLogger(logger)
	golog.SetCheckLogger(logger)

	return logger, flusher
}
```

## Production ready logger usage with opencensus support

```go
import (
    "github.com/damianopetrungaro/golog/v2"
    "github.com/damianopetrungaro/golog/v2/opencensus"
)

func NewLogger(lvl golog.Level) (golog.Logger, golog.Flusher) {
	w := golog.NewBufWriter(
        golog.NewJsonEncoder(),
		bufio.NewWriter(os.Stdout),
        golog.DefaultErrorHandler(),
		lvl,
	)
	
	logger := golog.New(
		w,
		golog.NewLevelCheckerOption(lvl),
		golog.NewTimestampDecoratorOption("timestamp", time.RFC3339Nano),
		golog.NewStackTraceDecoratorOption("stacktrace", 5),
        opencensus.TraceDecoratorOption(),
	)

	// force flushing data to the disk every 5 seconds
	flusher := golog.NewTickFlusher(stdoutWriter, 5*time.Second)
	go func() {
		if err := flusher.Flush(); err != nil {
			logger.With(golog.Err(err)).Warn(ctx, "log entries were not flushed")
		}
	}()

	golog.SetLogger(logger)
	golog.SetCheckLogger(logger)

	return logger, flusher
}
```

## Production ready logger usage with datadog and opentelemetry support
```go
import (
    "github.com/damianopetrungaro/golog/v2"
    "github.com/damianopetrungaro/golog/v2/datadog"
    "github.com/damianopetrungaro/golog/v2/opentelemetry"
)

func NewLogger(lvl golog.Level) (golog.Logger, golog.Flusher) {
    w := golog.NewBufWriter(
	    datadog.NewJsonEncoder(),
        bufio.NewWriter(os.Stdout),
        golog.DefaultErrorHandler(),
        lvl,
    )
    
    logger := golog.New(
        w,
        golog.NewLevelCheckerOption(lvl),
        golog.NewTimestampDecoratorOption("timestamp", time.RFC3339Nano),
        golog.NewStackTraceDecoratorOption("stacktrace", 5),
        opentelemetry.TraceDecoratorOption(),
        opentelemetry.CustomTraceDecoratorOption(
            datadog.TraceIDKey,
            datadog.TraceIDFromHexFormat,
            datadog.SpanIDKey,
            datadog.SpanIDFromHexFormat,
        ),
    )
    
    // force flushing data to the disk every 5 seconds
    flusher := golog.NewTickFlusher(stdoutWriter, 5*time.Second)
    go func() {
        if err := flusher.Flush(); err != nil {
            logger.With(golog.Err(err)).Warn(ctx, "log entries were not flushed")
        }
    }()
    
    golog.SetLogger(logger)
    golog.SetCheckLogger(logger)
    
    return logger, flusher
}
```


## Production ready logger usage with sentry and stdout support
```go
import (
    "github.com/damianopetrungaro/golog/v2"
    "github.com/damianopetrungaro/golog/v2/sentry"
    goSentry "github.com/getsentry/sentry-go"
)

func NewLogger(lvl golog.Level) (golog.Logger, golog.Flusher) {
    jsonEncoder := datadog.NewJsonEncoder()
    errHandler := golog.DefaultErrorHandler()
    
    stdoutWriter := golog.NewBufWriter(
        jsonEncoder,
        bufio.NewWriter(os.Stdout),
        errHandler,
        lvl,
    )
    
    sentryWriter := &sentry.Writer{
        Encoder:                   jsonEncoder,
        Hub:                       goSentry.NewHub(client, goSentry.NewScope()),
        ErrHandler:                errHandler,
        DefaultLevel:              lvl,
        CaptureExceptionFromLevel: lvl,
    }

    // create a writer which will write to sentry for error and fatal error
    multiWriter := golog.NewMultiWriter(stdoutWriter, sentryWriter)
    w := golog.NewLeveledWriter(
        stdoutWriter,
        golog.DefaultLeveledWriterOptionFunc(golog.ERROR, multiWriter),
        golog.DefaultLeveledWriterOptionFunc(golog.FATAL, multiWriter),
    )
    
    logger := golog.New(
        w,
        golog.NewLevelCheckerOption(lvl),
        golog.NewTimestampDecoratorOption("timestamp", time.RFC3339Nano),
        golog.NewStackTraceDecoratorOption("stacktrace", 5),
    )
    // force flushing data to the disk every 5 seconds
    flusher := golog.NewTickFlusher(stdoutWriter, 5*time.Second)
    go func() {
        if err := flusher.Flush(); err != nil {
            logger.With(golog.Err(err)).Warn(ctx, "log entries were not flushed")
        }
    }()
    
    golog.SetLogger(logger)
    golog.SetCheckLogger(logger)
    
    return logger, flusher
}
```
