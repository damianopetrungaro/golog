package golog

import (
	"bytes"
	"io"
	"strconv"
	"time"
)

var (
	defaultJsonConfig = JsonConfig{
		LevelKeyName:   "level",
		MessageKeyName: "message",
		TimeLayout:     time.RFC3339Nano,
	}
	defaultTextConfig = TextConfig{
		LevelKeyName:   "level",
		MessageKeyName: "message",
		TimeLayout:     time.RFC3339Nano,
	}
)

// Encoder transforms an entry into io.WriterTo which holds the encoded content
type Encoder interface {
	Encode(Entry) (io.WriterTo, error)
}

// TextConfig is a configuration for TextEncoder
type TextConfig struct {
	LevelKeyName   string
	MessageKeyName string
	TimeLayout     string
}

// TextEncoder is an encoder for text
type TextEncoder struct {
	Config TextConfig
}

// DefaultTextConfig returns a default TextConfig
func DefaultTextConfig() TextConfig {
	return defaultTextConfig
}

// NewTextEncoder returns a TextEncoder
func NewTextEncoder(cfg TextConfig) TextEncoder {
	return TextEncoder{Config: cfg}
}

// Encode encodes an entry into a text content holds into an io.WriterTo
func (t TextEncoder) Encode(e Entry) (io.WriterTo, error) {
	w := &bytes.Buffer{}
	t.addElemQuoted(w, t.Config.LevelKeyName, e.Level().String())
	w.WriteString(` `)
	t.addElemQuoted(w, t.Config.MessageKeyName, e.Message())
	t.encodeFields(e.Fields(), w)
	w.WriteByte('\n')
	return w, nil
}

func (t TextEncoder) encodeFields(flds Fields, w *bytes.Buffer) {
	if len(flds) == 0 {
		return
	}

	w.WriteString(` `)
	for i, f := range flds {
		t.encodeField(f, w)
		if i != len(flds)-1 {
			w.WriteString(` `)
		}
	}
}

func (t TextEncoder) encodeField(f Field, w *bytes.Buffer) {
	switch val := f.Value().(type) {
	case bool:
		t.addElem(w, f.Key(), strconv.FormatBool(val))
	case []bool:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.FormatBool(v), i != len(val)-1)
			}
		})
	case string:
		t.addElemQuoted(w, f.Key(), val)
	case []string:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElemQuoted(w, v, i != len(val)-1)
			}
		})
	case uint:
		t.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case uint8:
		t.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case uint16:
		t.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case uint32:
		t.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case uint64:
		t.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case int:
		t.addElem(w, f.Key(), strconv.Itoa(val))
	case int8:
		t.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case int16:
		t.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case int32:
		t.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case int64:
		t.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case []uint:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []uint8:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []uint16:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []uint32:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []uint64:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []int:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(v), i != len(val)-1)
			}
		})
	case []int8:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []int16:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []int32:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []int64:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case float64:
		t.addElem(w, f.Key(), strconv.FormatFloat(val, 'f', 10, 64))
	case float32:
		t.addElem(w, f.Key(), strconv.FormatFloat(float64(val), 'f', 10, 32))
	case []float64:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.FormatFloat(v, 'f', 10, 64), i != len(val)-1)
			}
		})
	case []float32:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElem(w, strconv.FormatFloat(float64(v), 'f', 10, 32), i != len(val)-1)
			}
		})
	case time.Time:
		t.addElemQuoted(w, f.Key(), val.Format(t.Config.TimeLayout))
	case []time.Time:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				t.addArrayElemQuoted(w, v.Format(t.Config.TimeLayout), i != len(val)-1)
			}
		})
	case FieldMapper:
		t.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, f := range val.ToFields() {
				t.encodeField(f, w)
				if i != len(val.ToFields())-1 {
					w.WriteString(` `)
				}
			}
		})
	}
}

func (t TextEncoder) addElem(w *bytes.Buffer, k string, val string) {
	w.WriteString(k)
	w.WriteString(`=`)
	w.WriteString(val)
}

func (t TextEncoder) addElemQuoted(w *bytes.Buffer, k string, val string) {
	w.WriteString(k)
	w.WriteString(`="`)
	w.WriteString(val)
	w.WriteString(`"`)
}

func (t TextEncoder) addArrayElem(w *bytes.Buffer, val string, hasNext bool) {
	w.WriteString(val)
	if hasNext {
		w.WriteString(`,`)
	}
}

func (t TextEncoder) addArrayElemQuoted(w *bytes.Buffer, val string, hasNext bool) {
	w.WriteString(`"`)
	w.WriteString(val)
	w.WriteString(`"`)
	if hasNext {
		w.WriteString(`,`)
	}
}

func (t TextEncoder) addElements(w *bytes.Buffer, k string, fn func(w *bytes.Buffer)) {
	w.WriteString(k)
	w.WriteString(`=[`)
	fn(w)
	w.WriteString(`]`)
}

// JsonConfig is a configuration for JsonEncoder
type JsonConfig struct {
	LevelKeyName   string
	MessageKeyName string
	TimeLayout     string
}

// JsonEncoder is an encoder for json
type JsonEncoder struct {
	Config JsonConfig
}

// DefaultJsonConfig returns a default JsonConfig
func DefaultJsonConfig() JsonConfig {
	return defaultJsonConfig
}

// NewJsonEncoder returns a JsonEncoder
func NewJsonEncoder(cfg JsonConfig) JsonEncoder {
	return JsonEncoder{Config: cfg}
}

// Encode encodes an entry into a json content holds into an io.WriterTo
func (j JsonEncoder) Encode(e Entry) (io.WriterTo, error) {
	w := &bytes.Buffer{}
	w.WriteString(`{`)
	j.addElemQuoted(w, j.Config.LevelKeyName, e.Level().String())
	w.WriteString(`,`)
	j.addElemQuoted(w, j.Config.MessageKeyName, e.Message())
	j.encodeFields(e.Fields(), w)
	w.WriteString(`}`)
	w.WriteByte('\n')
	return w, nil
}

func (j JsonEncoder) encodeFields(flds Fields, w *bytes.Buffer) {
	if len(flds) == 0 {
		return
	}

	w.WriteString(`,`)
	for i, f := range flds {
		j.encodeField(f, w)
		if i != len(flds)-1 {
			w.WriteString(`,`)
		}
	}
}

func (j JsonEncoder) encodeField(f Field, w *bytes.Buffer) {
	switch val := f.Value().(type) {
	case bool:
		j.addElem(w, f.Key(), strconv.FormatBool(val))
	case []bool:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.FormatBool(v), i != len(val)-1)
			}
		})
	case string:
		j.addElemQuoted(w, f.Key(), val)
	case []string:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElemQuoted(w, v, i != len(val)-1)
			}
		})
	case uint:
		j.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case uint8:
		j.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case uint16:
		j.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case uint32:
		j.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case uint64:
		j.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case int:
		j.addElem(w, f.Key(), strconv.Itoa(val))
	case int8:
		j.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case int16:
		j.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case int32:
		j.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case int64:
		j.addElem(w, f.Key(), strconv.Itoa(int(val)))
	case []uint:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []uint8:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []uint16:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []uint32:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []uint64:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []int:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(v), i != len(val)-1)
			}
		})
	case []int8:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []int16:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []int32:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case []int64:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.Itoa(int(v)), i != len(val)-1)
			}
		})
	case float64:
		j.addElem(w, f.Key(), strconv.FormatFloat(val, 'f', 10, 64))
	case float32:
		j.addElem(w, f.Key(), strconv.FormatFloat(float64(val), 'f', 10, 32))
	case []float64:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.FormatFloat(v, 'f', 10, 64), i != len(val)-1)
			}
		})
	case []float32:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElem(w, strconv.FormatFloat(float64(v), 'f', 10, 32), i != len(val)-1)
			}
		})
	case time.Time:
		j.addElemQuoted(w, f.Key(), val.Format(j.Config.TimeLayout))
	case []time.Time:
		j.addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, v := range val {
				j.addArrayElemQuoted(w, v.Format(j.Config.TimeLayout), i != len(val)-1)
			}
		})
	case FieldMapper:
		j.addObject(w, f.Key(), func(w *bytes.Buffer) {
			for i, f := range val.ToFields() {
				j.encodeField(f, w)
				if i != len(val.ToFields())-1 {
					w.WriteString(`,`)
				}
			}
		})
	}
}

func (j JsonEncoder) addElem(w *bytes.Buffer, k string, val string) {
	w.WriteString(`"`)
	w.WriteString(k)
	w.WriteString(`":`)
	w.WriteString(val)
}

func (j JsonEncoder) addElemQuoted(w *bytes.Buffer, k string, val string) {
	w.WriteString(`"`)
	w.WriteString(k)
	w.WriteString(`":"`)
	w.WriteString(val)
	w.WriteString(`"`)
}

func (j JsonEncoder) addArrayElem(w *bytes.Buffer, val string, hasNext bool) {
	w.WriteString(val)
	if hasNext {
		w.WriteString(`,`)
	}
}

func (j JsonEncoder) addArrayElemQuoted(w *bytes.Buffer, val string, hasNext bool) {
	w.WriteString(`"`)
	w.WriteString(val)
	w.WriteString(`"`)
	if hasNext {
		w.WriteString(`,`)
	}
}

func (j JsonEncoder) addElements(w *bytes.Buffer, k string, fn func(w *bytes.Buffer)) {
	w.WriteString(`"`)
	w.WriteString(k)
	w.WriteString(`":[`)
	fn(w)
	w.WriteString(`]`)
}

func (j JsonEncoder) addObject(w *bytes.Buffer, k string, fn func(w *bytes.Buffer)) {
	w.WriteString(`"`)
	w.WriteString(k)
	w.WriteString(`":{`)
	fn(w)
	w.WriteString(`}`)
}
