package tracer_test

import (
	"context"
	"strings"
	"testing"

	. "github.com/kamilsk/tracer"
)

func TestTrace_Start(t *testing.T) {
	(*Trace)(nil).Start().Mark("no panic")
	(&Trace{}).Start().Mark("one allocation")
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

		call := trace.Start().Mark("fn call")
		call.Checkpoint().Mark("checkpoint")
		call.Stop()
		if expected, obtained := "allocates at call stack: 1", trace.String(); !strings.Contains(obtained, expected) {
			t.Errorf("\n expected: %+#v \n obtained: %+#v", expected, obtained)
		}
	})
}

func TestCall_Checkpoint(t *testing.T) {
	(*Call)(nil).Checkpoint().Mark("no panic")
	(&Call{}).Checkpoint().Mark("one allocation")
}

func TestCall_Mark(t *testing.T) {
	(*Call)(nil).Mark("no panic")
	(&Call{}).Mark("by id")
}

func TestCall_Stop(t *testing.T) {
	(*Call)(nil).Mark("no panic").Stop()
	(&Call{}).Mark("success").Stop()
}

func TestCheckpoint_Mark(t *testing.T) {
	(*Checkpoint)(nil).Mark("no panic")
	(&Checkpoint{}).Mark("by id")
}

// BenchmarkTracing/silent-12         	  200000	      7755 ns/op	    1816 B/op	      24 allocs/op
// BenchmarkTracing/full-12           	  200000	      8880 ns/op	    3944 B/op	      45 allocs/op
func BenchmarkTracing(b *testing.B) {
	b.Run("silent", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			traceRoot(Inject(context.Background(), make([]*Call, 0, 9)))
		}
	})
	b.Run("full", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ctx := Inject(context.Background(), make([]*Call, 0, 9))
			traceRoot(ctx)
			_ = Fetch(ctx).String()
		}
	})
}
