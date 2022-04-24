package logger

import (
	"io"
)

type Syncer struct {
	err    error
	called bool
}

func (s *Syncer) SetError(err error) {
	s.err = err
}

func (s *Syncer) Sync() error {
	s.called = true
	return s.err
}

func (s *Syncer) Called() bool {
	return s.called
}

type Discarder struct{ Syncer }

func (d *Discarder) Write(b []byte) (int, error) {
	return io.Discard.Write(b)
}
