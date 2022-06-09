package test

import (
	"bufio"
	"io"

	"github.com/damianopetrungaro/golog"
)

func NewNullLogger() golog.Logger {
	return newNullStdLogger()
}

func NewNullCheckLogger() golog.CheckLogger {
	return newNullStdLogger()
}

func newNullStdLogger() golog.StdLogger {
	return golog.New(
		golog.NewBufWriter(
			golog.NewJsonEncoder(golog.DefaultJsonConfig()),
			bufio.NewWriter(io.Discard),
			golog.DefaultErrorHandler(),
			golog.DEBUG,
		),
	)
}
