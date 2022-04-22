package golog

//Fields is a slice of fields
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

//Bool creates a field containing a value of type "bool"
func Bool(k string, v bool) Field {
	return Field{k: k, v: v}
}

// Bools creates a field containing a value of type "[]bool"
func Bools(k string, v []bool) Field {
	return Field{k: k, v: v}
}

//String creates a field containing a value of type "string"
func String(k string, v string) Field {
	return Field{k: k, v: v}
}

//Strings creates a field containing a value of type "[]string"
func Strings(k string, v []string) Field {
	return Field{k: k, v: v}
}

//Uint creates a field containing a value of type "uint"
func Uint(k string, v uint) Field {
	return Field{k: k, v: v}
}

//Uints creates a field containing a value of type "[]uint"
func Uints(k string, v []uint) Field {
	return Field{k: k, v: v}
}

//Int creates a field containing a value of type "int"
func Int(k string, v int) Field {
	return Field{k: k, v: v}
}

//Ints creates a field containing a value of type "[]int"
func Ints(k string, v []int) Field {
	return Field{k: k, v: v}
}

//Float64 creates a field containing a value of type "float64"
func Float64(k string, v float64) Field {
	return Field{k: k, v: v}
}

//Float64s creates a field containing a value of type "[]float64"
func Float64s(k string, v []float64) Field {
	return Field{k: k, v: v}
}

//Float32 creates a field containing a value of type "float32"
func Float32(k string, v float32) Field {
	return Field{k: k, v: v}
}

//Float32s creates a field containing a value of type "[]float32"
func Float32s(k string, v []float32) Field {
	return Field{k: k, v: v}
}

//Err creates a field containing a value of type ") Fi"
func Err(err error) Field {
	const k = "error"
	return Field{k: k, v: err}
}

//Errs creates a field containing a value of type "or"
func Errs(err []error) Field {
	const k = "errors"
	return Field{k: k, v: err}
}
