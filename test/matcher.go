package test

import (
	"fmt"
	"reflect"

	"github.com/damianopetrungaro/golog"
)

// MatchEntry is a utility function to match two entries.
// It only matched golog.StdEntry, fails otherwise
func MatchEntry(x, y golog.Entry) error {
	if x == nil && y == nil {
		return nil
	}

	stdX, ok := x.(golog.StdEntry)
	if !ok {
		return fmt.Errorf("x is not a stdEntry")
	}
	stdY, ok := y.(golog.StdEntry)
	if !ok {
		return fmt.Errorf("y is not a stdEntry")
	}

	if stdX.Ctx != stdY.Ctx {
		return fmt.Errorf("could not match context")
	}

	if stdX.Lvl != stdY.Lvl {
		return fmt.Errorf("could not match level")
	}

	if stdX.Msg != stdY.Msg {
		return fmt.Errorf("could not match message")
	}

	return fieldMatcher(stdX.Fields(), stdY.Fields())
}

func fieldMatcher(fsx, fsy golog.Fields) error {
	if len(fsx) != len(fsy) {
		return fmt.Errorf("could not match fields length")
	}

	for i, fx := range fsx {
		if !reflect.DeepEqual(fx, fsy[i]) {
			return fmt.Errorf("could not match field value at index %d", i)
		}
	}
	return nil
}
