package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/damianopetrungaro/golog"
	. "github.com/damianopetrungaro/golog/http"
)

func TestHandler_ServeHTTP(t *testing.T) {
	fld := golog.Bool("injected", true)
	logger := golog.StdLogger{}.With(fld)
	var success bool
	var loggerMatched bool
	wantContentType := "no-content-type"

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		success = true
	})

	h := Middleware(logger, func(w *RecorderResponseWriter, r *http.Request, l golog.Logger) {
		if l.(golog.StdLogger).Fields[0] != fld {
			t.Error("could not match logger")
		}

		if ct := r.Header.Get("Content-Type"); ct != wantContentType {
			t.Error("could not match logger")
		}

		loggerMatched = true

	})

	srv := httptest.NewServer(h(next))
	if _, err := srv.Client().Post(srv.URL, wantContentType, http.NoBody); err != nil {
		t.Errorf("could not do request: %s", err)
	}

	if !success {
		t.Error("could not match handler call")
	}

	if !loggerMatched {
		t.Error("could not match logger")
	}
}

func TestDefaultLogHandle(t *testing.T) {
	logHandle := DefaultLogHandle()

	rec := &RecorderResponseWriter{
		Size:    1000,
		Status:  http.StatusOK,
		StartAt: time.Now(),
	}

	r := &http.Request{
		Method:     http.MethodPost,
		RequestURI: "https://www.goggle.com",
	}

	w := &inmem{}
	logger := golog.StdLogger{Writer: w}
	logHandle(rec, r, logger)

	for _, e := range w.entry.Fields() {
		switch e.Key() {
		case "request_method":
			if r.Method != e.Value() {
				t.Fatal("could not match request_method")
			}
		case "request_uri":
			if r.RequestURI != e.Value() {
				t.Fatal("could not match request_uri")
			}
		case "response_size":
			if rec.Size != e.Value() {
				t.Fatal("could not match response_size")
			}
		case "response_status":
			if rec.Status != e.Value() {
				t.Fatal("could not match response_status")
			}
		case "latency":
			if 100_000 <= e.Value().(int64) {
				t.Fatal("could not match latency")
			}
		}
	}
}

type inmem struct {
	entry golog.StdEntry
}

func (i *inmem) WriteEntry(e golog.Entry) {
	i.entry = e.(golog.StdEntry)
}
func (*inmem) Write([]byte) (int, error) {
	return 0, nil
}
