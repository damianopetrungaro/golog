# GOLOG

[![codecov](https://codecov.io/gh/damianopetrungaro/golog/branch/main/graph/badge.svg?token=5ESXFZo2j2)](https://codecov.io/gh/damianopetrungaro/golog)

Golog is an opinionated Go logger 
with simple APIs and configurable behavior.

## Why another logger?

Golog is designed to address mainly two issues:

#### Reduce the amount of PII (personally identifiable information) data in logs
Golog exposes APIs which does not allow to simply introduce a struct or a map as part of the log fields.

This design pushes the consumers of this library to care about PII data and
pushes to reduce as much as possible the amount of data which can be logged.

It is possible to extend the logger behavior 
for handling complex data type 
by writing custom field factory functions as shown in the customization section. 

#### Add tracing and other extra data into the logging behavior
Golog expects to have a context passed down to the logging API.

The `context.Context` in Go is usually the holder for tracing information and
embedding one of the decorators available to the logger plugs this behavior for free
in all the places where the logger is used.

## Examples

### Logger

The `Logger` interface is implemented by the `StdLogger` type.
It allows you to write log messages.

An example of its usage may look like this:
 ```go
golog.With(golog.Fields{
	golog.Bool("key name", true),
	golog.Strings("another key name", []string{"one", "two"}),
}).Error(ctx, "log message here")
 ```

To override the default logger you can use the `SetLogger` API as shown here:
 ```go
// create a new custom logger
logger := golog.New(
    golog.NewBufWriter(
        golog.NewJsonEncoder(golog.DefaultJsonConfig()),
        bufio.NewWriter(os.StdErr),
        golog.DefaultErrorHandler(),
    ),
    golog.NewLevelCheckerOption(golog.WARN), 
	// any other option you may want to pass
)

// set the custom logger as the global one 
golog.SetLogger(logger)
 ```

### CheckLogger
The `CheckLogger` interface is implemented by the `StdLogger` type.
It allows you to write log messages allowing to set fields only if the log message will be written.

For example if the min log level set is higher than the one which will be logged, 
as shown in this example, there will be no extra data allocation as well as having a huge performance improvement::
```go
if checked, ok := golog.CheckDebug(ctx, "This is a message"); ok {
    checked.Log(golog.Fields{
        golog.Bool("key name", true),
        golog.Strings("another key name", []string{"one", "two"}),
    })
}
```

To override the default check logger you can use the `SetCheckLogger` API as shown here:
 ```go
// create a new custom logger
logger := golog.New(
    golog.NewBufWriter(
        golog.NewJsonEncoder(golog.DefaultJsonConfig()),
        bufio.NewWriter(os.StdErr),
        golog.DefaultErrorHandler(),
    ),
    golog.NewLevelCheckerOption(golog.WARN), 
	// any other option you may want to pass
)

// set the custom check logger as the global one 
golog.SetCheckLogger(logger)
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

    return e.With(golog.Fields{
        golog.String("span_id", span.SpanID.String()),
        golog.String("trace_id", span.TraceID.String()),
    })
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

Out of the box are provided some checkers 
for skipping log with level lower than an expected one.

Example usage:
```go
var logger golog.Logger = golog.New(
    // other arguments here
    golog.NewLevelCheckerOption(golog.INFO),
)
```

### Custom field type
Logging complex data structure is not intentionally supported out of the box, 
Golog expects you to create a Fields factory function.

An example may be something like this:
```go
// The complex data structure to log
type User struct {
	ID              string
	Email           string
	Password        string
	ReferenceCode   string
}

// The factory function to create fields out of the complex data structure
func NewUserFields(u User) golog.Fields {
    return golog.Fields{
        golog.String("user_id", u.ID),
        golog.String("reference_code", u.ReferenceCode),
    }
}
```
And its usage would look like this
```go
// Example API usage
golog.With(NewUserFields(u)).Error("an error occurred")
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
BenchmarkLogger/golog-12                 1266547               929.9 ns/op          2826 B/op         26 allocs/op
BenchmarkLogger/zap-12                   1000000              1066 ns/op            2836 B/op         20 allocs/op
BenchmarkLogger/logrus-12                 344604              3395 ns/op            6168 B/op         69 allocs/op
BenchmarkLogger/golog.Check-12          56982846                20.10 ns/op           64 B/op          1 allocs/op
BenchmarkLogger/zap.Check-12            1000000000               0.9662 ns/op          0 B/op          0 allocs/op
PASS
ok      github.com/damianopetrungaro/golog/benchmarks/logger    6.781s
```

Considering the nature of the logger and the design it has, the performances are really high.

In the future there may be a support for an even faster and zero allocations version of the logger, 
but the APIs exposed won't be matching the current one and there will be a different interface provided for that purpose.

# Note

Golog doesn't hande key deduplication.

Meaning that 
```go
golog.With(golog.Fields{
    golog.String("hello", "world"), 
    golog.String("hello", "another world"),
}).Info(ctx, "no deduplication")
```

will print 
```json
{"level":"INFO","message":"no deduplication","hello":"world","hello":"another world"}
```