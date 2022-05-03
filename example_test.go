package golog_test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/damianopetrungaro/golog"
)

func ExampleLogger() {
	ctx := context.Background()
	t, err := time.Parse("2006", "2021")
	if err != nil {
		fmt.Println(err)
		return
	}

	w := golog.NewBufWriter(
		golog.NewJsonEncoder(golog.DefaultJsonConfig()),
		bufio.NewWriter(os.Stdout),
		golog.DefaultErrorHandler(),
		golog.DEBUG,
	)
	defer w.Flush()

	golog.SetLogger(golog.New(w, golog.NewLevelCheckerOption(golog.DEBUG)))

	logger := golog.With(golog.String("hello", "world"))
	logger.Error(ctx, "an error message")
	logger.Error(ctx, "another error message")
	loggerWithErr := logger.With(golog.Err(fmt.Errorf("error: ops!")))
	logger.Info(ctx, "an info message")
	loggerWithErr.Warn(ctx, "a warning message")
	loggerWithErr.With(golog.Mapper("user", user{ID: "uuid", Reference: 123, Birthdate: t})).Error(ctx, "a warning message")

	// Output:
	// {"level":"ERROR","message":"an error message","hello":"world"}
	// {"level":"ERROR","message":"another error message","hello":"world"}
	// {"level":"INFO","message":"an info message","hello":"world"}
	// {"level":"WARN","message":"a warning message","hello":"world","error":"error: ops!"}
	// {"level":"ERROR","message":"a warning message","hello":"world","error":"error: ops!","user":{"id":"uuid","ref":123,"birthdate":"2021-01-01T00:00:00Z"}}
	//
}
