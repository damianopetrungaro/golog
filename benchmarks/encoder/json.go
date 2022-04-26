package encoder

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/damianopetrungaro/golog"
	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

type JsonEncoder struct{}

func (j *JsonEncoder) StdLinEncode(e golog.Entry) (io.WriterTo, error) {
	log := map[string]any{
		"level":     e.Level().String(),
		"timestamp": time.Now().Format(`2006-01-02`),
		"message":   e.Message(),
	}

	j.mapEncodeFields(e.Fields(), log)

	w := &bytes.Buffer{}
	if err := json.NewEncoder(w).Encode(log); err != nil {
		return nil, err
	}

	w.WriteByte('\n')
	return w, nil
}

func (j *JsonEncoder) GoJsonEncode(e golog.Entry) (io.WriterTo, error) {
	log := map[string]any{
		"level":     e.Level().String(),
		"timestamp": time.Now().Format(`2006-01-02`),
		"message":   e.Message(),
	}
	j.mapEncodeFields(e.Fields(), log)

	w := &bytes.Buffer{}
	if err := gojson.NewEncoder(w).Encode(log); err != nil {
		return nil, err
	}

	w.WriteByte('\n')
	return w, nil
}

func (j *JsonEncoder) JsoniterEncode(e golog.Entry) (io.WriterTo, error) {
	log := map[string]any{
		"level":     e.Level().String(),
		"timestamp": time.Now().Format(`2006-01-02`),
		"message":   e.Message(),
	}
	j.mapEncodeFields(e.Fields(), log)

	raw, err := jsoniter.Marshal(log)
	if err != nil {
		return nil, err
	}

	w := &bytes.Buffer{}
	w.Write(raw)
	w.WriteByte('\n')
	return w, nil
}

func (j *JsonEncoder) mapEncodeFields(flds golog.Fields, log map[string]any) {
	if len(flds) == 0 {
		return
	}
	fields := map[string]any{}
	for _, f := range flds {
		fields[f.Key()] = f.Key()
	}
	log["fields"] = fields
}

func (j *JsonEncoder) ManualMapEncode(lvl golog.Level, msg golog.Message, flds map[string]any) (io.WriterTo, error) {
	w := &bytes.Buffer{}
	w.WriteString(`{`)
	addElemQuoted(w, "level", lvl.String())
	w.WriteString(`,`)
	addElemQuoted(w, "timestamp", time.Now().Format(`2006-01-02`))
	w.WriteString(`,`)
	addElemQuoted(w, "message", msg)
	j.manualEncodeMapFields(flds, w)
	w.WriteString(`}`)
	w.WriteByte('\n')
	return w, nil
}

func (j *JsonEncoder) ManualEncode(e golog.Entry) (io.WriterTo, error) {
	w := &bytes.Buffer{}
	w.WriteString(`{`)
	addElemQuoted(w, "level", e.Level().String())
	w.WriteString(`,`)
	addElemQuoted(w, "timestamp", time.Now().Format(`2006-01-02`))
	w.WriteString(`,`)
	addElemQuoted(w, "message", e.Message())
	j.manualEncodeFields(e.Fields(), w)
	w.WriteString(`}`)
	w.WriteByte('\n')
	return w, nil
}

func (j *JsonEncoder) manualEncodeFields(flds golog.Fields, w *bytes.Buffer) {
	if len(flds) == 0 {
		return
	}

	w.WriteString(`,"fields":{`)
	for i, f := range flds {
		j.encodeField(f, w)
		if i != len(flds)-1 {
			w.WriteString(`,`)
		}
	}
	w.WriteString(`}`)
}

func (j *JsonEncoder) manualEncodeMapFields(flds map[string]any, w *bytes.Buffer) {
	if len(flds) == 0 {
		return
	}
	total := len(flds) - 1
	i := 0
	w.WriteString(`,`)
	for k, v := range flds {
		j.encodeFieldFromMap(k, v, w)
		if i != total {
			w.WriteString(`,`)
		}
		i++
	}
}

func (j *JsonEncoder) encodeField(f golog.Field, w *bytes.Buffer) {
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

func (j *JsonEncoder) encodeFieldFromMap(k string, v any, w *bytes.Buffer) {
	switch t := v.(type) {
	case bool:
		addElem(w, k, strconv.FormatBool(t))
	case []bool:
		addElements(w, k, func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElem(w, strconv.FormatBool(val), i != len(t)-1)
			}
		})
	case string:
		addElemQuoted(w, k, t)
	case error:
		switch t {
		case nil:
			addElemQuoted(w, k, "null")
		default:
			addElemQuoted(w, k, t.Error())
		}
	case []string:
		addElements(w, k, func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElemQuoted(w, val, i != len(t)-1)
			}
		})
	case []error:
		addElements(w, k, func(w *bytes.Buffer) {
			for i, val := range t {
				switch t {
				case nil:
					addElemQuoted(w, k, "null")
				default:
					addElemQuoted(w, k, val.Error())
				}
				if i != len(t)-1 {
					w.WriteString(`,`)
				}
			}
		})
	case uint:
		addElem(w, k, strconv.Itoa(int(t)))
	case int:
		addElem(w, k, strconv.Itoa(t))
	case []uint:
		addElements(w, k, func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElem(w, strconv.Itoa(int(val)), i != len(t)-1)
			}
		})
	case []int:
		addElements(w, k, func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElem(w, strconv.Itoa(val), i != len(t)-1)
			}
		})
	case float64:
		addElem(w, k, strconv.FormatFloat(t, 'f', 10, 64))
	case float32:
		addElem(w, k, strconv.FormatFloat(float64(t), 'f', 10, 32))
	case []float64:
		addElements(w, k, func(w *bytes.Buffer) {
			for i, val := range t {
				addArrayElem(w, strconv.FormatFloat(val, 'f', 10, 64), i != len(t)-1)
			}
		})
	case []float32:
		addElements(w, k, func(w *bytes.Buffer) {
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
