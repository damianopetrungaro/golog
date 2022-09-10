package datadog

import (
	"github.com/damianopetrungaro/golog/v2"
)

// NewProductionLogger returns a pre-configured logger for production environment
// adds a decorator for supporting datadog enoder
func NewProductionLogger(lvl golog.Level) (golog.StdLogger, golog.Flusher) {
	return golog.NewWithEncoder(lvl, NewJsonEncoder())
}

// NewDevelopmentLogger returns a pre-configured logger for development environment
// adds a decorator for supporting datadog enoder
func NewDevelopmentLogger(lvl golog.Level) (golog.StdLogger, golog.Flusher) {
	return golog.NewWithEncoder(lvl, NewTextEncoder())
}
