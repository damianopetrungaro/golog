package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	h := NewHandler(next, logger, func(w *RecorderResponseWriter, r *http.Request, l golog.Logger) {
		if l.(golog.StdLogger).Fields[0] != fld {
			t.Error("could not match logger")
		}

		if ct := r.Header.Get("Content-Type"); ct != wantContentType {
			t.Error("could not match logger")
		}

		loggerMatched = true
	})

	srv := httptest.NewServer(h)
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
