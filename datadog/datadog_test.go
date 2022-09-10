package datadog_test

import (
	"testing"
	"time"

	. "github.com/damianopetrungaro/golog/v2/datadog"
)

func TestNewTimestampDecorator(t *testing.T) {
	d := NewTimestampDecorator()
	if d.TimestampLayout != time.RFC3339 {
		t.Error("could not match expected layout")
		t.Errorf("got: %s", d.TimestampLayout)
		t.Errorf("want: %s", time.RFC3339)
	}

	if d.TimestampFieldName != "time" {
		t.Error("could not match expected field name")
		t.Errorf("got: %s", d.TimestampFieldName)
		t.Errorf("want: time")
	}
}

func TestDefaultJsonConfig(t *testing.T) {
	cfg := DefaultJsonConfig()
	if cfg.MessageKeyName != "message" {
		t.Error("could not match message key")
		t.Errorf("got: %s", cfg.MessageKeyName)
		t.Errorf("want: message")
	}
	if cfg.LevelKeyName != "status" {
		t.Error("could not match level key")
		t.Errorf("got: %s", cfg.LevelKeyName)
		t.Errorf("want: status")
	}
	if cfg.TimeLayout != time.RFC3339 {
		t.Error("could not match time layout")
		t.Errorf("got: %s", cfg.TimeLayout)
		t.Errorf("want: %s", time.RFC3339)
	}
}

func TestNewJsonEncoder(t *testing.T) {
	cfg := DefaultJsonConfig()
	enc := NewJsonEncoder()
	if enc.Config != cfg {
		t.Error("could not match config")
		t.Errorf("got: %v", enc.Config)
		t.Errorf("want: %v", cfg)
	}
}

func TestDefaultTextConfig(t *testing.T) {
	cfg := DefaultTextConfig()
	if cfg.MessageKeyName != "message" {
		t.Error("could not match message key")
		t.Errorf("got: %s", cfg.MessageKeyName)
		t.Errorf("want: message")
	}
	if cfg.LevelKeyName != "status" {
		t.Error("could not match level key")
		t.Errorf("got: %s", cfg.LevelKeyName)
		t.Errorf("want: status")
	}
	if cfg.TimeLayout != time.RFC3339 {
		t.Error("could not match time layout")
		t.Errorf("got: %s", cfg.TimeLayout)
		t.Errorf("want: %s", time.RFC3339)
	}
}

func TestNewTextEncoder(t *testing.T) {
	cfg := DefaultTextConfig()
	enc := NewTextEncoder()
	if enc.Config.MessageKeyName != cfg.MessageKeyName &&
		enc.Config.LevelKeyName != cfg.LevelKeyName &&
		enc.Config.TimeLayout != cfg.TimeLayout {
		t.Error("could not match config")
		t.Errorf("got: %v", enc.Config)
		t.Errorf("want: %v", cfg)
	}
}
