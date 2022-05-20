package encoder

import (
	"context"
	"fmt"
	"testing"

	"github.com/damianopetrungaro/golog"
)

/**
goarch: amd64
pkg: json_encoder_benchmarks
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkJsonEncoder/encoding/json-12             545086              2307 ns/op            4994 B/op         72 allocs/op
BenchmarkJsonEncoder/gojson-12                    803994              1795 ns/op            3425 B/op         27 allocs/op
BenchmarkJsonEncoder/jsoniter-12                  669554              1691 ns/op            3426 B/op         27 allocs/op
BenchmarkJsonEncoder/golog-12                    2084641               548.0 ns/op          1311 B/op          6 allocs/op
BenchmarkJsonEncoder/golog.map-12                1924228               654.6 ns/op          1248 B/op          6 allocs/op
PASS
ok      json_encoder_benchmarks 7.650s
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
