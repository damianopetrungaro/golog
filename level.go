package golog

import (
	"errors"
	"strings"
)

const (
	_ Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var ErrLevelNotParsed = errors.New("golog: could not parse level")

type Level int

func ParseLevel(s string) (Level, error) {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN", "WARNING":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	}

	return 0, ErrLevelNotParsed
}

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return ""
	}
}
