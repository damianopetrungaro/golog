package opencensus

import (
	"github.com/damianopetrungaro/golog/v2"
)

// NewProductionLogger returns a pre-configured logger for production environment
// adds a decorator for supporting opencensus decorator
func NewProductionLogger(lvl golog.Level) (golog.StdLogger, golog.Flusher) {
	logger, flusher := golog.NewProductionLogger(lvl)
	logger = logger.WithDecorator(decorator)
	return logger, flusher
}

// NewDevelopmentLogger returns a pre-configured logger for development environment
// adds a decorator for supporting opencensus decorator
func NewDevelopmentLogger(lvl golog.Level) (golog.StdLogger, golog.Flusher) {
	logger, flusher := golog.NewDevelopmentLogger(lvl)
	logger = logger.WithDecorator(decorator)
	return logger, flusher
}
