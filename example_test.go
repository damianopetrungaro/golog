package golog_test

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/damianopetrungaro/golog"
)

func ExampleLogger() {

	user := struct {
		ID   string
		Name string
	}{
		ID:   "123",
		Name: "John",
	}

	invalidJsonMap := map[any]any{
		1:          1,
		"2":        2,
		"metadata": map[string]string{"count": "12"},
	}

	validJsonMap := map[string]any{
		"1":        1,
		"2":        2,
		"metadata": map[string]string{"count": "12"},
	}

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
	loggerWithErr.Warn(ctx, "a warning message")

	loggerWithErrAndUser := loggerWithErr.With(golog.Any("user", user))
	loggerWithErrAndUser.With(golog.Any("unsupported map", invalidJsonMap)).Warn(ctx, "a warning message with an unsupported map")
	loggerWithErrAndUser.With(golog.Any("supported map", validJsonMap)).Warn(ctx, "a warning message with a supported map")

	// Output:
	// {"level":"ERROR","message":"an error message","hello":"world"}
	// {"level":"ERROR","message":"another error message","hello":"world"}
	// {"level":"INFO","message":"an info message","hello":"world"}
	// {"level":"WARN","message":"a warning message","hello":"world","error":"error: ops!"}
	// {"level":"WARN","message":"a warning message with an unsupported map","hello":"world","error":"error: ops!","user":{"ID":"123","Name":"John"}
	//,"unsupported map":"json: unsupported type: map[interface {}]interface {}"}
	// {"level":"WARN","message":"a warning message with a supported map","hello":"world","error":"error: ops!","user":{"ID":"123","Name":"John"}
	//,"supported map":{"1":1,"2":2,"metadata":{"count":"12"}}
	//}
	//
}
