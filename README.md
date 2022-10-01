# GOLOG

[![codecov](https://codecov.io/gh/damianopetrungaro/golog/branch/main/graph/badge.svg?token=5ESXFZo2j2)](https://codecov.io/gh/damianopetrungaro/golog)

Golog is an opinionated Go logger
with simple APIs and configurable behavior.

## Why another logger?

Golog is designed to address mainly two issues:

#### Reduce the amount of PII (personally identifiable information) data in logs

Golog exposes APIs which does not allow to simply introduce a struct or a map as part of the log fields.

This design pushes the consumers of this library to care about PII data and
aim to reduce as much as possible the amount of data which can be logged.

It is possible to extend the logger behavior
for handling complex data type
by implementing an interface as shown in the "Custom field type" section.

#### Add tracing and other extra data into the logging behavior

Golog expects to have a context passed down to the logging API.

The `context.Context` in Go is usually the holder for tracing information and
embedding one of the decorators available to the logger plugs this behavior for free
in all the places where the logger is used.

## Who uses Golog?

* [Pento](https://www.pento.io) - Used in multiple production services.


## Examples

Based on your needs, you can find some presets available.

The ones highly recommended are the one that adds tracing metadata in your logs.

If you use opentelemetry, the code snippet will look like this:
```go
package main

import (
	"context"

	"github.com/damianopetrungaro/golog"
	"github.com/damianopetrungaro/golog/opentelemetry"
)

func main() {
	logger, flusher := NewLogger(golog.DEBUG) // min level
	defer flusher.Flush()
	// ....
	logger.Info(context.Background(), "Hello world")
}

func NewLogger(lvl golog.Level) (golog.StdLogger, golog.Flusher) {
	return opentelemetry.NewProductionLogger(lvl)
}

```

If you use opencensus, the code snippet will look like this:
```go
package main

import (
	"context"

	"github.com/damianopetrungaro/golog"
	"github.com/damianopetrungaro/golog/opencensus"
)

func main() {
	logger, flusher := NewLogger(golog.DEBUG) // min level
	defer flusher.Flush()
	// ....
	logger.Info(context.Background(), "Hello world")
}

func NewLogger(lvl golog.Level) (golog.StdLogger, golog.Flusher) {
	return opencensus.NewProductionLogger(lvl)
}
```

There is also a method for development purposes, wo use that use the factory function `NewDevelopemntLogger`. 

For more extensive and customized implementations, plesse continue reading the documentation!  

### Logger

The `Logger` interface is implemented by the `StdLogger` type.
It allows you to write log messages.

An example of its usage may look like this:

 ```go
w := golog.NewBufWriter(
    golog.NewJsonEncoder(golog.DefaultJsonConfig()),
    bufio.NewWriter(os.Stdout),
    golog.DefaultErrorHandler(),
    golog.INFO,
)
defer w.Flush()

logger := golog.New(w, golog.NewTimestampDecoratorOption("timestamp", time.RFC3339))
golog.SetLogger(logger)

golog.With(
	golog.Bool("key name", true),
	golog.Strings("another key name", []string{"one", "two"}),
).Error(ctx, "log message here")
 ```

which will print
```json
{"level":"ERROR","message":"log message here","key name":true,"another key name":["one","two"],"timestamp":"2022-05-20T16:16:29+02:00"}
```

### CheckLogger

The `CheckLogger` interface is implemented by the `StdLogger` type.
It allows you to write log messages allowing to set fields only if the log message will be written.

For example if the min log level set is higher than the one which will be logged,
as shown in this example, there will be no extra data allocation as well as having a huge performance improvement::

```go
w := golog.NewBufWriter(
    golog.NewJsonEncoder(golog.DefaultJsonConfig()),
    bufio.NewWriter(os.Stdout),
    golog.DefaultErrorHandler(),
    golog.INFO,
)
defer w.Flush()

logger := golog.New(w, golog.NewTimestampDecoratorOption("timestamp", time.RFC3339))
golog.SetCheckLogger(logger)

if checked, ok := golog.CheckDebug(ctx, "This is a message"); ok {
    checked.Log(
        golog.Bool("key name", true),
        golog.Strings("another key name", []string{"one", "two"}),
    )
}
```

```json
{"level":"DEBUG","message":"This is a message","timestamp":"2022-05-20T16:28:15+02:00","key name":true,"another key name":["one","two"]}
```

### Standard Library support

Golog Writer can be used by the go `log` package as well as output

```go
w := &BufWriter{
    Encoder:         enc,
    Writer:          bufio.NewWriter(buf),
    ErrHandler:      errHandler.Handle,
    DefaultLogLevel: DEBUG, //! This will be the log level used for all the logs by the stdlib logger
}

log.SetOutput(w)
log.Println("your log message here...")
```

## Customization

Golog provides multiple ways to customize behaviors

### Decorators

A decorator is a function that gets executed before a log message gets written,
allowing to inject only once a recurring logging behavior
to modify the log message.

An example may be adding a trace and span ids to the log:

```go
var customTraceDecorator golog.DecoratorFunc = func(e golog.Entry) golog.Entry {
    span := trace.FromContext(e.Context()).SpanContext()

    return e.With(
        golog.String("span_id", span.SpanID.String()),
        golog.String("trace_id", span.TraceID.String()),
    )
}

var logger golog.Logger = golog.New(
    // other arguments here
    golog.OptionFunc(func(l golog.StdLogger) golog.StdLogger {
        return l.WithDecorator(customTraceDecorator)
    }),
)
```

Out of the box are provided some decorators
for tracing purposes in the `opencensus` and `opentelemetry` packages,
PRs are welcome to add more behavior.

### Checkers

A checker is a function that gets executed before a log message gets decorated,
allowing to skip the decoration and the writing of a log entry due to custom logic.

An example may be skipping a log if the context doesn't have a value:

```go
var customCtxValueChecker golog.Checker = golog.CheckerFunc(func(e golog.Entry) bool {
    if _, ok := e.Context().Value("key").(string); !ok {
        return false
    }

    return true
})

var logger golog.Logger = golog.New(
    // other arguments here
    golog.OptionFunc(func(l golog.StdLogger) golog.StdLogger {
        return l.WithChecker(customCtxValueChecker)
    }),
)
```

Out of the box are provided some checkers:

### Min level checker

This checker will skip logging all the entry with level lower than an expected one.

Example usage:

```go
var logger golog.Logger = golog.New(
    // other arguments here
    golog.NewLevelCheckerOption(golog.INFO),
)
```

## Custom field type

Logging complex data structure is not intentionally supported out of the box,
Golog expects you to implement a FieldMapper interface.

An example may be something like this:

```go
// The complex data structure to log
type User struct {
	ID              string
	Email           string
	Password        string
	ReferenceCode   string
}

// The FieldMapper interface method to create fields out of the complex data structure
func (u User) ToFields() golog.Fields {
    return golog.Fields{
        golog.String("user_id", u.ID),
        golog.String("reference_code", u.ReferenceCode),
    }
}

//...

var u User{...} 
golog.With(golog.Mapper("user", u)).Debug(ctx, "...")
```

## Writers

Based on your need you may want to use different entry writers.

Golog provide you those implementations:

#### BufWriter

It is the standard implementation, and it can be created in this way:

```go
w := golog.NewBufWriter(
    golog.NewJsonEncoder(golog.DefaultJsonConfig()),
    bufio.NewWriter(os.Stdout),
    golog.DefaultErrorHandler(),
    golog.INFO,
)
```

#### LeveledWriter

This implementation provides you a way to use a different writer based on the log level, 
with a default writer used in case there is not an override defined for a log level

```go
var stdOutWriter golog.Writer 
var stdErrWriter golog.Writer 

w := NewLeveledWriter(
    stdOutWriter,
    golog.DefaultMuxWriterOptionFunc(golog.ERROR, stdErrWriter),
    golog.DefaultMuxWriterOptionFunc(golog.FATAL, stdErrWriter),
)
```


#### MultiWriter

This implementation simply writes an across multiple writers concurrently

```go
var w1 golog.Writer
var w2 golog.Writer
var w3 golog.Writer
w := golog.NewMultiWriter(w1, w2, w3)
```


#### DeduplicatorWriter

This implementation will deduplicate keys with the same values.

The logger is slower when using this writer, so make sure you actually need it. 

```go
w := golog.NewBufWriter(
    golog.NewJsonEncoder(golog.DefaultJsonConfig()),
    bufio.NewWriter(os.Stdout),
    golog.DefaultErrorHandler(),
    golog.DEBUG,
)
defer w.Flush()

w = golog.NewDeduplicatorWriter(w)
golog.SetLogger(golog.New(w, golog.NewLevelCheckerOption(golog.DEBUG)))

golog.With(golog.String("hello", "world"), golog.String("hello", "another world")).Error(ctx, "an error message")
```

This will print
```json
{"level":"ERROR","message":"an error message","hello":"world","hello_1":"another world"}
```


## Testing utilities

The `golog/test` package provides a mock generated using [gomock](https://github.com/golang/mock) for helping developers
to test the logger.

The `golog/test` package provides a null logger as well, 
useful when the logger need to be passed as dependency but there the test don't care about the logging logic.

## HTTP utilities

The `golog/http` utility package provide a simple and customizable API for adding some logging behavior on an HTTP
server.

```go
// ...
import (
    "net/http"
    
	"github.com/damianopetrungaro/golog"
    httplog "github.com/damianopetrungaro/golog/http"
)

// ...
var h http.Handler // the handler you want to decorate
var logger golog.Logger // your logger
httplog.NewHandler(h, logger, httplog.DefaultLogHandle()) // returns your decorated handler
```

## Presets

Golog provides configurations and encoders utilities for several log consumers 

### Datadog (DD)

```go
import (
    "github.com/damianopetrungaro/golog/datadog"
)

logger := golog.New(
    golog.NewBufWriter(
        datadog.NewJsonEncoder(),
        bufio.NewWriter(os.Stdout),
        golog.DefaultErrorHandler(),
        golog.DEBUG,
    ),
    golog.NewLevelCheckerOption(golog.WARN),
    // any other option you may want to pass
)

```

## Performances

Golog is a really fast logging solution,
with a low number of allocations as well as crazy performances.

Benchmarks comparing it to [logrus](https://github.com/sirupsen/logrus) and [zap](https://github.com/uber-go/zap)

```text
goos: darwin
goarch: amd64
pkg: github.com/damianopetrungaro/golog/benchmarks/logger
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkLogger/logrus-12				322524			3544 ns/op		6162 B/op	70 allocs/op
BenchmarkLogger/golog-12				1314696			878.1 ns/op		2841 B/op	27 allocs/op
BenchmarkLogger/golog.deduplicator-12	811656			1318 ns/op		4370 B/op	34 allocs/op
BenchmarkLogger/zap-12					1292510			1093 ns/op		2836 B/op	20 allocs/op
BenchmarkLogger/golog.Check-12			69898831		18.23 ns/op		64 B/op 	1 allocs/op
BenchmarkLogger/zap.Check-12			1000000000		0.8543 ns/op	0 B/op 		0 allocs/op
PASS
ok      github.com/damianopetrungaro/golog/benchmarks/logger    9.179s
```

Considering the nature of the logger and the design it has, the performances are really high.

In the future there may be a support for an even faster and zero allocations version of the logger,
but the APIs exposed won't be matching the current one and there will be a different interface provided for that
purpose.

[More updated benchmarks can be found on this page](https://damianopetrungaro.github.io/golog/)

# Note

Golog doesn't hande key deduplication by default.

Meaning that

```go
golog.With(
    golog.String("hello", "world"),
    golog.String("hello", "another world"),
).Info(ctx, "no deduplication")
```

will print

```json
{"level": "INFO","message":"no deduplication","hello":"world","hello":"another world"}
```

If you need deduplicated keys, please check the DeduplicatorWriter section.