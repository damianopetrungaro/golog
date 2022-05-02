package golog_test

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/damianopetrungaro/golog"
)

func ExampleLogger() {

	w := golog.NewBufWriter(
		golog.NewJsonEncoder(golog.DefaultJsonConfig()),
		bufio.NewWriter(os.Stdout),
		golog.DefaultErrorHandler(),
		golog.DEBUG,
	)
	defer w.Flush()

	var logger golog.Logger = golog.New(
		w,
		golog.NewLevelCheckerOption(golog.DEBUG),
	)

	golog.SetLogger(logger)

	ctx := context.Background()

	logger = golog.With(golog.String("hello", "world"))
	logger.Error(ctx, "an error message")
	logger.Error(ctx, "another error message")
	loggerWithErr := logger.With(golog.Err(fmt.Errorf("error: ops!")))
	logger.Info(ctx, "an info message")
	loggerWithErr.Warning(ctx, "a warning message")

	// Output:
	// {"level":"ERROR","message":"an error message","hello":"world"}
	// {"level":"ERROR","message":"another error message","hello":"world"}
	// {"level":"INFO","message":"an info message","hello":"world"}
	// {"level":"WARN","message":"a warning message","hello":"world","error":"error: ops!"}
	//
}
