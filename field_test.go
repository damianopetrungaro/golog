package golog_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/damianopetrungaro/golog"
)

func Test_Bool(t *testing.T) {
	k := "key name"
	v := true
	f := Bool(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Bools(t *testing.T) {
	k := "key name"
	v := []bool{true, false, false, true}
	f := Bools(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_String(t *testing.T) {
	k := "key name"
	v := "value name"
	f := String(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Strings(t *testing.T) {
	k := "key name"
	v := []string{"name one", "name two", "name three"}
	f := Strings(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Uint(t *testing.T) {
	k := "key name"
	v := uint(101)
	f := Uint(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Uints(t *testing.T) {
	k := "key name"
	v := []uint{10, 0, 202}
	f := Uints(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Int(t *testing.T) {
	k := "key name"
	v := 12
	f := Int(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Ints(t *testing.T) {
	k := "key name"
	v := []int{1, 2, 3, 4, 5}
	f := Ints(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Float64(t *testing.T) {
	k := "key name"
	v := 3.39
	f := Float64(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Float64s(t *testing.T) {
	k := "key name"
	v := []float64{1.12, 21.12, 3.419}
	f := Float64s(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Float32(t *testing.T) {
	k := "key name"
	v := float32(9.01)
	f := Float32(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Float32s(t *testing.T) {
	k := "key name"
	v := []float32{1.12, 21.12, 3.419}
	f := Float32s(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Time(t *testing.T) {
	k := "key name"
	v := time.Now()
	f := Time(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Times(t *testing.T) {
	k := "key name"
	v := []time.Time{time.Now(), time.Now(), time.Now()}
	f := Times(k, v)
	testFieldHelper(t, k, v, f)
}

func Test_Err(t *testing.T) {
	t.Run("test nil error", func(t *testing.T) {
		k := "error"
		errMsg := "<nil>"
		f := Err(nil)
		testFieldHelper(t, k, errMsg, f)
	})

	t.Run("test actual error", func(t *testing.T) {
		k := "error"
		errMsg := "an error"
		v := fmt.Errorf(errMsg)
		f := Err(v)
		testFieldHelper(t, k, errMsg, f)
	})
}

func Test_Errs(t *testing.T) {
	k := "errors"
	errMsgs := []string{"an error", "<nil>", "another error"}
	v := []error{fmt.Errorf("an error"), nil, fmt.Errorf("another error")}
	f := Errs(v)
	testFieldHelper(t, k, errMsgs, f)
}

// tests Key and Value methods as well implicitly
func testFieldHelper(t *testing.T, k string, v any, f Field) {
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

func FieldMatcher(t *testing.T, xs, ys Fields) {
	t.Helper()
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
