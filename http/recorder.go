package http

import (
	"bufio"
	"net"
	"net/http"
	"time"
)

type RecorderResponseWriter struct {
	Size    int
	Status  int
	StartAt time.Time
	http.ResponseWriter
}

func (rec *RecorderResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := rec.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}

	return nil, nil, ErrCouldNotMatchHijacker
}

func (rec *RecorderResponseWriter) Flush() {
	if f, ok := rec.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (rec *RecorderResponseWriter) Write(bytes []byte) (int, error) {
	size, err := rec.ResponseWriter.Write(bytes)
	rec.Size += size
	return size, err
}

func (rec *RecorderResponseWriter) WriteHeader(statusCode int) {
	rec.Status = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}
