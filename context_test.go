package tracer_test

import (
	"context"
	"testing"

	. "github.com/kamilsk/tracer"
)

func TestContext(t *testing.T) {
	t.Run("fetch nil", func(t *testing.T) {
		if Fetch(context.Background()) != nil {
			t.Error("expected nil")
		}
		Fetch(context.Background()).Start("no panic").Stop()
	})
	t.Run("fetch injected", func(t *testing.T) {
		ctx := Inject(context.Background(), nil)
		if Fetch(ctx) == nil {
			t.Error("unexpected nil")
		}
		Fetch(ctx).Start("allocation").Stop()
	})
}

// BenchmarkContext/injecting-12         	20000000	        72.1 ns/op	      80 B/op	       2 allocs/op
// BenchmarkContext/fetching-12          	300000000	         5.90 ns/op	       0 B/op	       0 allocs/op
// BenchmarkContext/combine-12           	20000000	        80.8 ns/op	      80 B/op	       2 allocs/op
func BenchmarkContext(b *testing.B) {
	b.Run("injecting", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Inject(context.Background(), nil)
		}
	})
	b.Run("fetching", func(b *testing.B) {
		ctx := Inject(context.Background(), nil)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Fetch(ctx)
		}
	})
	b.Run("combine", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx := Inject(context.Background(), nil)
			for i := 0; i < 3; i++ {
				_ = Fetch(ctx)
			}
		}
	})
}
