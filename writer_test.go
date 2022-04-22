package golog_test

import (
	. "github.com/damianopetrungaro/golog"
)

var _ Writer = &FakeWriter{}

type FakeWriter struct {
	Entry Entry
}

func (fw *FakeWriter) Write(e Entry) {
	fw.Entry = e
}
