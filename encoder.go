package golog

import (
	"bytes"
	"io"
	"strconv"
)

var (
	defaultJsonConfig = JsonConfig{
		LevelKeyName:   "level",
		MessageKeyName: "message",
	}
	defaultTextConfig = TextConfig{
		LevelKeyName:   "level",
		MessageKeyName: "message",
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
	case int:
		t.addElem(w, f.Key(), strconv.Itoa(val))
	case []uint:
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
	case int:
		j.addElem(w, f.Key(), strconv.Itoa(val))
	case []uint:
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
