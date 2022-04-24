package golog

import (
	"bytes"
	"io"
	"strconv"
	"time"
)

var (
	defaultJsonConfig = JsonConfig{
		LevelKeyName:        "level",
		TimestampKeyName:    "timestamp",
		TimestampLayout:     time.RFC3339Nano,
		MessageKeyName:      "message",
		FieldsKeyName:       "fields",
		EnableStackTrace:    false,
		StackTraceFieldName: "stacktrace",
	}
)

// Encoder transforms an entry into io.WriterTo which holds the encoded content
type Encoder interface {
	Encode(Entry) (io.WriterTo, error)
}

// JsonConfig is a configuration for JsonEncoder
type JsonConfig struct {
	LevelKeyName        string
	TimestampKeyName    string
	TimestampLayout     string
	MessageKeyName      string
	FieldsKeyName       string
	EnableStackTrace    bool
	StackTraceFieldName string
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
	addElemQuoted(w, j.Config.LevelKeyName, e.Level().String())
	w.WriteString(`,`)
	addElemQuoted(w, j.Config.TimestampKeyName, time.Now().Format(j.Config.TimestampLayout))
	w.WriteString(`,`)
	addElemQuoted(w, j.Config.MessageKeyName, e.Message())
	j.encodeFields(e.Fields(), w)
	j.encodeStackTrace(nil)
	w.WriteString(`}\n`)
	return w, nil
}

func (j JsonEncoder) encodeFields(flds Fields, w *bytes.Buffer) {
	if len(flds) == 0 {
		return
	}

	w.WriteString(`,"`)
	w.WriteString(j.Config.FieldsKeyName)
	w.WriteString(`":{`)
	for i, f := range flds {
		j.encodeField(f, w)
		if i != len(flds)-1 {
			w.WriteString(`,`)
		}
	}
	w.WriteString(`}`)
}

func (j JsonEncoder) encodeStackTrace(w *bytes.Buffer) {
	if !j.Config.EnableStackTrace {
		return
	}

	// TODO
}

func (j JsonEncoder) encodeField(f Field, w *bytes.Buffer) {
	switch t := f.Value().(type) {
	case bool:
		addElem(w, f.Key(), strconv.FormatBool(t))
	case []bool:
		addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElem(w, strconv.FormatBool(val), i != len(t)-1)
			}
		})
	case string:
		addElemQuoted(w, f.Key(), t)
	case error:
		switch t {
		case nil:
			addElemQuoted(w, f.Key(), "null")
		default:
			addElemQuoted(w, f.Key(), t.Error())
		}
	case []string:
		addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElemQuoted(w, val, i != len(t)-1)
			}
		})
	case []error:
		addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, val := range t {
				switch t {
				case nil:
					addElemQuoted(w, f.Key(), "null")
				default:
					addElemQuoted(w, f.Key(), val.Error())
				}
				if i != len(t)-1 {
					w.WriteString(`,`)
				}
			}
		})
	case uint:
		addElem(w, f.Key(), strconv.Itoa(int(t)))
	case int:
		addElem(w, f.Key(), strconv.Itoa(t))
	case []uint:
		addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElem(w, strconv.Itoa(int(val)), i != len(t)-1)
			}
		})
	case []int:
		addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElem(w, strconv.Itoa(val), i != len(t)-1)
			}
		})
	case float64:
		addElem(w, f.Key(), strconv.FormatFloat(t, 'f', 10, 64))
	case float32:
		addElem(w, f.Key(), strconv.FormatFloat(float64(t), 'f', 10, 32))
	case []float64:
		addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElem(w, strconv.FormatFloat(val, 'f', 10, 64), i != len(t)-1)
			}
		})
	case []float32:
		addElements(w, f.Key(), func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElem(w, strconv.FormatFloat(float64(val), 'f', 10, 32), i != len(t)-1)
			}
		})
	}
}

func addElem(w *bytes.Buffer, k string, val string) {
	w.WriteString(`"`)
	w.WriteString(k)
	w.WriteString(`":`)
	w.WriteString(val)
}

func addElemQuoted(w *bytes.Buffer, k string, val string) {
	w.WriteString(`"`)
	w.WriteString(k)
	w.WriteString(`":"`)
	w.WriteString(val)
	w.WriteString(`"`)
}

func addArrayElem(w *bytes.Buffer, val string, hasNext bool) {
	w.WriteString(val)
	if hasNext {
		w.WriteString(`,`)
	}
}

func addArrayElemQuoted(w *bytes.Buffer, val string, hasNext bool) {
	w.WriteString(`"`)
	w.WriteString(val)
	w.WriteString(`"`)
	if hasNext {
		w.WriteString(`,`)
	}
}

func addElements(w *bytes.Buffer, k string, fn func(w *bytes.Buffer)) {
	w.WriteString(`"`)
	w.WriteString(k)
	w.WriteString(`":[`)
	fn(w)
	w.WriteString(`]`)
}
