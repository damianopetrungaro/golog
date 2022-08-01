package datadog

import (
	"strconv"
)

const (
	TraceIDKey = "dd.trace_id"
	SpanIDKey  = "dd.span_id"
)

func TraceIDFromHexFormat(id string) string {
	return convertHexIDToUint(id)
}

func SpanIDFromHexFormat(id string) string {
	return convertHexIDToUint(id)
}

func convertHexIDToUint(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
