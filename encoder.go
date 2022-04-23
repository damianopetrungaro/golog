package golog

import (
	"io"
)

type Encoder interface {
	Encode(Entry) (io.WriterTo, error)
}
