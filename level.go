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

const (
	COLOUR_RESET  = "\033[0m"
	COLOUR_RED    = "\033[31m"
	COLOUR_REDBG  = "\033[41m"
	COLOUR_GREEN  = "\033[32m"
	COLOUR_YELLOW = "\033[33m"
	COLOUR_BLUE   = "\033[34m"
	COLOUR_WHITE  = "\033[97m"
)

// ErrLevelNotParsed is an error returned when a given string can't be parsed as a log Level
var ErrLevelNotParsed = errors.New("golog: could not parse level")

// Level is a log severity level
type Level int

// ParseLevel returns a Level given a string, returns an error in case the string is not a recognized one
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

// String returns a string format of a log Level
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
