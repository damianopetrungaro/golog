package datadog

import (
	"time"

	"github.com/damianopetrungaro/golog"
)

var (
	defaultJsonConfig = golog.JsonConfig{
		LevelKeyName:   "status",
		MessageKeyName: "message",
		TimeLayout:     time.RFC3339,
	}

	defaultTextConfig = golog.TextConfig{
		LevelKeyName:   "status",
		MessageKeyName: "message",
		TimeLayout:     time.RFC3339,
	}
)

// NewTimestampDecorator returns a TimestampDecorator with the given field name and layout
func NewTimestampDecorator() golog.TimestampDecorator {
	return golog.TimestampDecorator{TimestampFieldName: "time", TimestampLayout: time.RFC3339}
}

// DefaultJsonConfig returns a default JsonConfig
func DefaultJsonConfig() golog.JsonConfig {
	return defaultJsonConfig
}

// NewJsonEncoder returns a JsonEncoder
func NewJsonEncoder(cfg golog.JsonConfig) golog.JsonEncoder {
	return golog.JsonEncoder{Config: cfg}
}

// DefaultTextConfig returns a default TextConfig
func DefaultTextConfig() golog.TextConfig {
	return defaultTextConfig
}

// NewTextEncoder returns a TextEncoder
func NewTextEncoder(cfg golog.TextConfig) golog.TextEncoder {
	return golog.TextEncoder{Config: cfg}
}
