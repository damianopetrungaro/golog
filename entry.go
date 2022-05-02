package golog

import (
	"context"
)

// Entry is a log entry
type Entry interface {
	Level() Level
	Message() Message
	Fields() Fields
	Context() context.Context
	With(...Field) Entry
}

// StdEntry is a representation of the standard log Entry
type StdEntry struct {
	Ctx  context.Context
	Lvl  Level
	Msg  Message
	Flds Fields
}

// NewStdEntry returns a StdEntry
func NewStdEntry(
	ctx context.Context,
	lvl Level,
	msg string,
	flds Fields,
) StdEntry {
	return StdEntry{
		Ctx:  ctx,
		Lvl:  lvl,
		Msg:  msg,
		Flds: flds,
	}
}

// Context returns an entry assigned to the entry
// This could be used to enrich the entry after it get created for the first time
func (e StdEntry) Context() context.Context {
	return e.Ctx
}

// Level returns the entry Level
func (e StdEntry) Level() Level {
	return e.Lvl
}

// Message returns the entry Message
func (e StdEntry) Message() Message {
	return e.Msg
}

// Fields returns the fields assigned to a log entry
func (e StdEntry) Fields() Fields {
	return e.Flds
}

// With returns a new StdEntry appending the given extra Fields
func (e StdEntry) With(flds ...Field) Entry {
	return StdEntry{
		Lvl:  e.Lvl,
		Msg:  e.Msg,
		Flds: append(e.Flds, flds...),
		Ctx:  e.Ctx,
	}
}
