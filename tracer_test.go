package tracer_test

import (
	"context"
	"strings"
	"testing"

	. "github.com/kamilsk/tracer"
)

func TestTrace_Start(t *testing.T) {
	(*Trace)(nil).Start("no panic")
	(&Trace{}).Start("one allocation")
}

func TestTrace_String(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		if expected, obtained := "", (*Trace)(nil).String(); expected != obtained {
			t.Errorf("\n expected: %+#v \n obtained: %+#v", expected, obtained)
		}
	})
	t.Run("allocations", func(t *testing.T) {
		trace := &Trace{}
		if expected, obtained := "allocates at call stack: 0, detailed call stack: ~",
			trace.String(); !strings.Contains(obtained, expected) {
			t.Errorf("\n expected: %+#v \n obtained: %+#v", expected, obtained)
		}

		call := trace.Start("fn call")
		call.Checkpoint("checkpoint")
		call.Stop()
		if expected, obtained := "allocates at call stack: 1", trace.String(); !strings.Contains(obtained, expected) {
			t.Errorf("\n expected: %+#v \n obtained: %+#v", expected, obtained)
		}
	})
}

func TestCall_Checkpoint(t *testing.T) {
	(*Call)(nil).Checkpoint("no panic")
	(&Call{}).Checkpoint("one allocation")
}

func TestCall_Stop(t *testing.T) {
	(*Call)(nil).Stop()
	(&Call{}).Stop()
}

// BenchmarkTracing/overhead-12         	 2000000	       616 ns/op	     144 B/op	       9 allocs/op
// BenchmarkTracing/silent-12           	  200000	      7478 ns/op	    2320 B/op	      27 allocs/op
// BenchmarkTracing/active-12           	  200000	      9221 ns/op	    4448 B/op	      48 allocs/op
func BenchmarkTracing(b *testing.B) {
	b.Run("overhead", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			traceRoot(context.Background())
		}
	})
	b.Run("silent", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			traceRoot(Inject(context.Background(), make([]*Call, 0, 9)))
		}
	})
	b.Run("active", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx := Inject(context.Background(), make([]*Call, 0, 9))
			traceRoot(ctx)
			_ = Fetch(ctx).String()
		}
	})
}
