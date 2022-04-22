package golog_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/damianopetrungaro/golog"
)

func Test_Bool(t *testing.T) {
	k := "key name"
	v := true
	f := golog.Bool(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Bools(t *testing.T) {
	k := "key name"
	v := []bool{true, false, false, true}
	f := golog.Bools(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_String(t *testing.T) {
	k := "key name"
	v := "value name"
	f := golog.String(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Strings(t *testing.T) {
	k := "key name"
	v := []string{"name one", "name two", "name three"}
	f := golog.Strings(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Uint(t *testing.T) {
	k := "key name"
	v := uint(101)
	f := golog.Uint(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Uints(t *testing.T) {
	k := "key name"
	v := []uint{10, 0, 202}
	f := golog.Uints(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Int(t *testing.T) {
	k := "key name"
	v := 12
	f := golog.Int(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Ints(t *testing.T) {
	k := "key name"
	v := []int{1, 2, 3, 4, 5}
	f := golog.Ints(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Float64(t *testing.T) {
	k := "key name"
	v := 3.39
	f := golog.Float64(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Float64s(t *testing.T) {
	k := "key name"
	v := []float64{1.12, 21.12, 3.419}
	f := golog.Float64s(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Float32(t *testing.T) {
	k := "key name"
	v := float32(9.01)
	f := golog.Float32(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Float32s(t *testing.T) {
	k := "key name"
	v := []float32{1.12, 21.12, 3.419}
	f := golog.Float32s(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Err(t *testing.T) {
	k := "error"
	v := fmt.Errorf("an error")
	f := golog.Err(v)
	testFieldHelper(t, k, v, f)
}

func Test_Errs(t *testing.T) {
	k := "errors"
	v := []error{fmt.Errorf("an error"), fmt.Errorf("another error")}
	f := golog.Errs(v)
	testFieldHelper(t, k, v, f)
}

// tests Key and Value methods as well implicitly
func testFieldHelper(t *testing.T, k string, v any, f golog.Field) {
	t.Helper()
	if f.Key() != k {
		t.Error("could not match key")
		t.Errorf("got: %v", f.Key())
		t.Errorf("want: %v", k)
	}
	if !reflect.DeepEqual(f.Value(), v) {
		t.Error("could not match value")
		t.Errorf("got: %v", f.Value())
		t.Errorf("want: %v", v)
	}
}

func FieldMatcher(t *testing.T, xs, ys golog.Fields) {
	if len(xs) != len(ys) {
		t.Error("could not match fields length")
		t.Errorf("xs: %v", len(xs))
		t.Errorf("ys: %v", len(ys))
	}

	for i, x := range xs {
		if !reflect.DeepEqual(x, ys[i]) {
			t.Errorf("could not match value at index %d", i)
			t.Errorf("x: %v", x)
			t.Errorf("y: %v", ys[i])
		}
	}
}
