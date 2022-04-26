package encoder

import (
	"context"
	"fmt"
	"testing"

	"github.com/damianopetrungaro/golog"
)

/**
OUTCOME:

goos: darwin
goarch: amd64
pkg: json_encoder_benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkJsonEncoder/encoding/json-12             526494              2233 ns/op            4998 B/op         70 allocs/op
BenchmarkJsonEncoder/gojson-12                    783475              1624 ns/op            3427 B/op         27 allocs/op
BenchmarkJsonEncoder/jsoniter-12                  783744              1667 ns/op            3428 B/op         27 allocs/op
BenchmarkJsonEncoder/golog-12                    2231253               528.3 ns/op          1311 B/op          6 allocs/op
PASS
ok      json_encoder_benchmarks 6.689s

*/

var entry = golog.NewStdEntry(context.Background(), golog.INFO, "This is a log message", golog.Fields{
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

func BenchmarkJsonEncoder(b *testing.B) {
	b.Run("encoding/json", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			enc := &JsonEncoder{}
			b.ResetTimer()
			for pb.Next() {
				if _, err := enc.StdLinEncode(entry); err != nil {
					b.Fatal("encode failed")
				}
			}
		})
	})

	b.Run("gojson", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			enc := &JsonEncoder{}
			b.ResetTimer()
			for pb.Next() {
				if _, err := enc.GoJsonEncode(entry); err != nil {
					b.Fatal("encode failed")
				}
			}
		})
	})

	b.Run("jsoniter", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			enc := &JsonEncoder{}
			b.ResetTimer()
			for pb.Next() {
				if _, err := enc.GoJsonEncode(entry); err != nil {
					b.Fatal("encode failed")
				}
			}
		})
	})

	b.Run("golog", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			enc := &JsonEncoder{}
			b.ResetTimer()
			for pb.Next() {
				if _, err := enc.ManualEncode(entry); err != nil {
					b.Fatal("encode failed")
				}
			}
		})
	})

	b.Run("golog.map", func(b *testing.B) {
		enc := &JsonEncoder{}
		mapFields := map[string]any{}
		for _, f := range entry.Fields() {
			mapFields[f.Key()] = f.Value()
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if _, err := enc.ManualMapEncode(entry.Lvl, entry.Msg, mapFields); err != nil {
					b.Fatal("encode failed")
				}
			}
		})
	})
}
