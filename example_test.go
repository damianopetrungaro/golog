package golog_test

import (
	"bufio"
	"context"
	"os"

	"github.com/damianopetrungaro/golog"
)

func ExampleLogger() {

	w := golog.NewBufWriter(
		golog.NewJsonEncoder(golog.DefaultJsonConfig()),
		bufio.NewWriter(os.Stdout),
		golog.DefaultErrorHandler(),
		golog.INFO,
	)
	defer w.Flush()

	var logger golog.Logger = golog.New(
		w,
		golog.NewLevelCheckerOption(golog.INFO),
	)

	golog.SetLogger(logger)

	logger = golog.With(golog.Fields{golog.String("hello", "world")})
	logger.Error(context.Background(), "an error message")
	logger.Error(context.Background(), "another error message")

	// Output:
	// {"level":"ERROR","message":"an error message","fields":{"hello":"world"}}
	// {"level":"ERROR","message":"another error message","fields":{"hello":"world"}}
	//
}
