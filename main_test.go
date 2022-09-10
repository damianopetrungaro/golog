package golog_test

import (
	"context"

	. "github.com/damianopetrungaro/golog/v2"
)

type any = interface{}

var (
	msg        = "This is a log message"
	debugEntry = NewStdEntry(context.Background(), DEBUG, msg, nil)
	infoEntry  = NewStdEntry(context.Background(), INFO, msg, nil)
	warnEntry  = NewStdEntry(context.Background(), WARN, msg, nil)
	errorEntry = NewStdEntry(context.Background(), ERROR, msg, nil)
	fatalEntry = NewStdEntry(context.Background(), FATAL, msg, nil)
)
