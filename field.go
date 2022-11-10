package golog

import (
	"time"
)

// FieldMapper an interface which returns Fields to add to an Entry
type FieldMapper interface {
	ToFields() Fields
}

// Fields is a slice of Field
type Fields []Field

// Field is a key value pair representing a log entry metadata field
type Field struct {
	k string
	v any
}

// Key returns the key of the Field
func (f Field) Key() string {
	return f.k
}

// Value returns the value of the Field
func (f Field) Value() any {
	return f.v
}

// Bool creates a field containing a value of type "bool"
func Bool(k string, v bool) Field {
	return Field{k: k, v: v}
}

// Bools creates a field containing a value of type "[]bool"
func Bools(k string, v []bool) Field {
	return Field{k: k, v: v}
}

// String creates a field containing a value of type "string"
func String(k string, v string) Field {
	return Field{k: k, v: v}
}

// Strings creates a field containing a value of type "[]string"
func Strings(k string, v []string) Field {
	return Field{k: k, v: v}
}

// Byte creates a field containing a value of type "byte"
func Byte(k string, v byte) Field {
	return Field{k: k, v: v}
}

// Bytes creates a field containing a value of type "[]byte"
func Bytes(k string, v []byte) Field {
	return Field{k: k, v: v}
}

// Uint creates a field containing a value of type "uint"
func Uint(k string, v uint) Field {
	return Field{k: k, v: v}
}

// Uints creates a field containing a value of type "[]uint"
func Uints(k string, v []uint) Field {
	return Field{k: k, v: v}
}

// Uint8 creates a field containing a value of type "uint8"
func Uint8(k string, v uint8) Field {
	return Field{k: k, v: v}
}

// Uint8s creates a field containing a value of type "[]uint8"
func Uint8s(k string, v []uint8) Field {
	return Field{k: k, v: v}
}

// Uint16 creates a field containing a value of type "uint16"
func Uint16(k string, v uint16) Field {
	return Field{k: k, v: v}
}

// Uint16s creates a field containing a value of type "[]uint16"
func Uint16s(k string, v []uint16) Field {
	return Field{k: k, v: v}
}

// Uint32 creates a field containing a value of type "uint32"
func Uint32(k string, v uint32) Field {
	return Field{k: k, v: v}
}

// Uint32s creates a field containing a value of type "[]uint32"
func Uint32s(k string, v []uint32) Field {
	return Field{k: k, v: v}
}

// Uint64 creates a field containing a value of type "uint64"
func Uint64(k string, v uint64) Field {
	return Field{k: k, v: v}
}

// Uint64s creates a field containing a value of type "[]uint64"
func Uint64s(k string, v []uint64) Field {
	return Field{k: k, v: v}
}

// Int creates a field containing a value of type "int"
func Int(k string, v int) Field {
	return Field{k: k, v: v}
}

// Ints creates a field containing a value of type "[]int"
func Ints(k string, v []int) Field {
	return Field{k: k, v: v}
}

// Int8 creates a field containing a value of type "int8"
func Int8(k string, v int8) Field {
	return Field{k: k, v: v}
}

// Int8s creates a field containing a value of type "[]int8"
func Int8s(k string, v []int8) Field {
	return Field{k: k, v: v}
}

// Int16 creates a field containing a value of type "int16"
func Int16(k string, v int16) Field {
	return Field{k: k, v: v}
}

// Int16s creates a field containing a value of type "[]int16"
func Int16s(k string, v []int16) Field {
	return Field{k: k, v: v}
}

// Int32 creates a field containing a value of type "int32"
func Int32(k string, v int32) Field {
	return Field{k: k, v: v}
}

// Int32s creates a field containing a value of type "[]int32"
func Int32s(k string, v []int32) Field {
	return Field{k: k, v: v}
}

// Int64 creates a field containing a value of type "int64"
func Int64(k string, v int64) Field {
	return Field{k: k, v: v}
}

// Int64s creates a field containing a value of type "[]int64"
func Int64s(k string, v []int64) Field {
	return Field{k: k, v: v}
}

// Float64 creates a field containing a value of type "float64"
func Float64(k string, v float64) Field {
	return Field{k: k, v: v}
}

// Float64s creates a field containing a value of type "[]float64"
func Float64s(k string, v []float64) Field {
	return Field{k: k, v: v}
}

// Float32 creates a field containing a value of type "float32"
func Float32(k string, v float32) Field {
	return Field{k: k, v: v}
}

// Float32s creates a field containing a value of type "[]float32"
func Float32s(k string, v []float32) Field {
	return Field{k: k, v: v}
}

// Time creates a field containing a value of type "time.Time"
func Time(k string, v time.Time) Field {
	return Field{k: k, v: v}
}

// Times creates a field containing a value of type "[]time.Time"
func Times(k string, v []time.Time) Field {
	return Field{k: k, v: v}
}

// Mapper creates a field containing a value of type "Fields"
func Mapper(k string, v FieldMapper) Field {
	return Field{k: k, v: v}
}

// Mappers creates a field containing a value of type "Fields"
func Mappers(k string, v []FieldMapper) Field {
	return Field{k: k, v: v}
}

// Err creates a field containing a value of type ") Fi"
func Err(err error) Field {
	const k = "error"
	if err == nil {
		return String(k, "<nil>")
	}
	return String(k, err.Error())
}

// Errs creates a field containing a value of type "or"
func Errs(errs []error) Field {
	const k = "errors"

	var ss = make([]string, len(errs))
	for i, err := range errs {
		if err == nil {
			ss[i] = "<nil>"
			continue
		}
		ss[i] = err.Error()
	}

	return Field{k: k, v: ss}
}
