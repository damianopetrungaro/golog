package golog_test

import (
	"errors"
	"io"

	. "github.com/damianopetrungaro/golog"
)

var _ Encoder = &FakeEncoder{}
var ErrFakeEncoder = errors.New("an error occurred on the encoder")

// FakeEncoder used for internal testing purposes
type FakeEncoder struct {
	Entry          Entry
	ShouldFail     bool
	ShouldWriterTo io.WriterTo
}

func (fe *FakeEncoder) Encode(e Entry) (io.WriterTo, error) {
	fe.Entry = e
	if fe.ShouldFail {
		return nil, ErrFakeEncoder
	}
	return fe.ShouldWriterTo, nil
}
