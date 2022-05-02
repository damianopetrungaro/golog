package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/damianopetrungaro/golog/http"
)

func TestRecorderResponseWriter(t *testing.T) {
	data := []byte(`Hello`)
	statusCode := http.StatusCreated
	httpRec := httptest.NewRecorder()

	rec := &RecorderResponseWriter{
		ResponseWriter: httpRec,
	}

	if _, err := rec.Write(data); err != nil {
		t.Fatalf("could not write: %s", err)
	}

	rec.WriteHeader(statusCode)

	if rec.Size != len(data) {
		t.Error("could not match size")
		t.Errorf("got: %d", rec.Size)
		t.Errorf("want: %d", len(data))
	}
	if rec.Status != statusCode {
		t.Error("could not match status")
		t.Errorf("got: %d", rec.Status)
		t.Errorf("want: %d", statusCode)
	}

	rec.Flush()
	if !httpRec.Flushed {
		t.Error("could not match flushed")
	}
}
