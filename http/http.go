package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/damianopetrungaro/golog"
)

var defaultLogHandle LogHandle = func(rec *RecorderResponseWriter, r *http.Request, logger golog.Logger) {
	logger.With(
		golog.String("request_method", r.Method),
		golog.String("request_uri", r.RequestURI),
		golog.Int("response_size", rec.Size),
		golog.Int("response_status", rec.Status),
		golog.Int64("latency", time.Since(rec.StartAt).Nanoseconds()),
	).Info(r.Context(), "request handled")
}

var ErrCouldNotMatchHijacker = errors.New("golog: could not match hijacker")

func DefaultLogHandle() LogHandle {
	return defaultLogHandle
}

type LogHandle func(*RecorderResponseWriter, *http.Request, golog.Logger)

type Handler struct {
	Next      http.Handler
	Logger    golog.Logger
	LogHandle LogHandle
}

func NewHandler(next http.Handler, logger golog.Logger, logHandle LogHandle) *Handler {
	return &Handler{Next: next, Logger: logger, LogHandle: logHandle}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := &RecorderResponseWriter{
		ResponseWriter: w,
		Status:         http.StatusOK,
		StartAt:        time.Now(),
	}

	h.Next.ServeHTTP(rec, r)
	h.LogHandle(rec, r, h.Logger)
}

func Middleware(logger golog.Logger, logHandle LogHandle) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return NewHandler(next, logger, logHandle)
	}
}
