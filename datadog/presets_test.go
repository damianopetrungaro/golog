package datadog

import (
	"testing"

	"github.com/damianopetrungaro/golog"
)

func Test_NewProductionLogger(t *testing.T) {
	lvl := golog.INFO
	logger, _ := NewProductionLogger(lvl)
	w, ok := logger.Writer.(*golog.BufWriter)
	if !ok {
		t.Fatalf("could not match writer")
	}
	enc, ok := w.Encoder.(golog.JsonEncoder)
	if !ok {
		t.Fatalf("could not match ecoder")
	}
	if enc.Config != DefaultJsonConfig() {
		t.Fatal("could not match config")
	}
}

func Test_NewDevelopmentLogger(t *testing.T) {
	lvl := golog.INFO
	logger, _ := NewDevelopmentLogger(lvl)
	w, ok := logger.Writer.(*golog.BufWriter)
	if !ok {
		t.Fatalf("could not match writer")
	}
	enc, ok := w.Encoder.(golog.TextEncoder)
	if !ok {
		t.Fatalf("could not match ecoder")
	}
	textConfigMatchHelper(t, enc.Config, DefaultTextConfig())
}

func textConfigMatchHelper(t *testing.T, cfg, cfg2 golog.TextConfig) {
	t.Helper()

	if cfg.LevelKeyName != cfg2.LevelKeyName {
		t.Fatal("could not match LevelKeyName")
	}
	if cfg.MessageKeyName != cfg2.MessageKeyName {
		t.Fatal("could not match MessageKeyName")
	}
	if cfg.TimeLayout != cfg2.TimeLayout {
		t.Fatal("could not match TimeLayout")
	}
}
