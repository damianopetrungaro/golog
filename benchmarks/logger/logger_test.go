package logger

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/damianopetrungaro/golog"
)

func BenchmarkLogger(b *testing.B) {

	b.Run("golog", func(b *testing.B) {
		ctx := context.Background()
		writer := golog.NewBufWriter(
			golog.NewJsonEncoder(golog.DefaultJsonConfig()),
			bufio.NewWriter(io.Discard),
			golog.DefaultErrorHandler(),
		)

		logger := golog.New(writer, golog.NewLevelCheckerOption(golog.DEBUG))

		golog.SetLogger(logger)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				golog.With(golog.Fields{
					golog.Int("int", 10),
					golog.Ints("ints", []int{1, 2, 3, 4, 5, 6, 7}),
					golog.String("string", "a string"),
					golog.Strings("strings", []string{"one", "one", "one", "one", "one", "one"}),
					golog.Int("int_2", 10),
					golog.Ints("ints_2", []int{1, 2, 3, 4, 5, 6, 7}),
					golog.String("string_2", "a string"),
					golog.Strings("strings_2", []string{"one", "one", "one", "one", "one", "one"}),
					golog.Int("int_3", 10),
					golog.Ints("ints_3", []int{1, 2, 3, 4, 5, 6, 7}),
					golog.String("string_3", "a string"),
					golog.Strings("strings_3", []string{"one", "one", "one", "one", "one", "one"}),
					golog.Err(fmt.Errorf("an error occurred")),
				}).Debug(ctx, "This is a message")
			}
		})
	})

	b.Run("zap", func(b *testing.B) {
		encoderCfg := zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			NameKey:        "logger",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), &Discarder{}, zap.DebugLevel)
		logger := zap.New(core).WithOptions()

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.With(
					zap.Int("int", 10),
					zap.Ints("ints", []int{1, 2, 3, 4, 5, 6, 7}),
					zap.String("string", "a string"),
					zap.Strings("strings", []string{"one", "one", "one", "one", "one", "one"}),
					zap.Int("int_2", 10),
					zap.Ints("ints_2", []int{1, 2, 3, 4, 5, 6, 7}),
					zap.String("string_2", "a string"),
					zap.Strings("strings_2", []string{"one", "one", "one", "one", "one", "one"}),
					zap.Int("int_3", 10),
					zap.Ints("ints_3", []int{1, 2, 3, 4, 5, 6, 7}),
					zap.String("string_3", "a string"),
					zap.Strings("strings_3", []string{"one", "one", "one", "one", "one", "one"}),
					zap.Error(fmt.Errorf("an error occurred")),
				).Debug("This is a message")
			}
		})
	})

	b.Run("logrus", func(b *testing.B) {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(io.Discard)
		logrus.SetLevel(logrus.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logrus.WithFields(logrus.Fields{
					"int":       10,
					"ints":      []int{1, 2, 3, 4, 5, 6, 7},
					"string":    "a string",
					"strings":   []string{"one", "one", "one", "one", "one", "one"},
					"int_2":     10,
					"ints_2":    []int{1, 2, 3, 4, 5, 6, 7},
					"string_2":  "a string",
					"strings_2": []string{"one", "one", "one", "one", "one", "one"},
					"int_3":     10,
					"ints_3":    []int{1, 2, 3, 4, 5, 6, 7},
					"string_3":  "a string",
					"strings_3": []string{"one", "one", "one", "one", "one", "one"},
					"error":     fmt.Errorf("an error occurred"),
				}).Debug("This is a message")
			}
		})
	})

	b.Run("golog.Check", func(b *testing.B) {
		ctx := context.Background()
		writer := golog.NewBufWriter(
			golog.NewJsonEncoder(golog.DefaultJsonConfig()),
			bufio.NewWriter(io.Discard),
			golog.DefaultErrorHandler(),
		)

		logger := golog.New(writer, golog.NewLevelCheckerOption(golog.WARN))

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if checked, ok := logger.CheckDebug(ctx, "This is a message"); ok {
					checked.Log(golog.Fields{
						golog.Int("int", 10),
						golog.Ints("ints", []int{1, 2, 3, 4, 5, 6, 7}),
						golog.String("string", "a string"),
						golog.Strings("strings", []string{"one", "one", "one", "one", "one", "one"}),
						golog.Int("int_2", 10),
						golog.Ints("ints_2", []int{1, 2, 3, 4, 5, 6, 7}),
						golog.String("string_2", "a string"),
						golog.Strings("strings_2", []string{"one", "one", "one", "one", "one", "one"}),
						golog.Int("int_3", 10),
						golog.Ints("ints_3", []int{1, 2, 3, 4, 5, 6, 7}),
						golog.String("string_3", "a string"),
						golog.Strings("strings_3", []string{"one", "one", "one", "one", "one", "one"}),
						golog.Err(fmt.Errorf("an error occurred")),
					})
				}
			}
		})
	})

	b.Run("zap.Check", func(b *testing.B) {
		encoderCfg := zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			NameKey:        "logger",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		}
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), &Discarder{}, zap.WarnLevel)
		logger := zap.New(core).WithOptions()

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if ce := logger.Check(zap.DebugLevel, "This is a message"); ce != nil {
					ce.Write(
						zap.Int("int", 10),
						zap.Ints("ints", []int{1, 2, 3, 4, 5, 6, 7}),
						zap.String("string", "a string"),
						zap.Strings("strings", []string{"one", "one", "one", "one", "one", "one"}),
						zap.Int("int_2", 10),
						zap.Ints("ints_2", []int{1, 2, 3, 4, 5, 6, 7}),
						zap.String("string_2", "a string"),
						zap.Strings("strings_2", []string{"one", "one", "one", "one", "one", "one"}),
						zap.Int("int_3", 10),
						zap.Ints("ints_3", []int{1, 2, 3, 4, 5, 6, 7}),
						zap.String("string_3", "a string"),
						zap.Strings("strings_3", []string{"one", "one", "one", "one", "one", "one"}),
						zap.Error(fmt.Errorf("an error occurred")),
					)
				}
			}
		})
	})
}
