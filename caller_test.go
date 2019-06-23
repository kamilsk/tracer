package tracer_test

import (
	"reflect"
	"testing"

	. "github.com/kamilsk/tracer"
)

func TestCaller(t *testing.T) {
	tests := []struct {
		name     string
		caller   func() CallerInfo
		expected string
	}{
		{"direct caller", callerA, "github.com/kamilsk/tracer_test.callerA"},
		{"chain caller", callerB, "github.com/kamilsk/tracer_test.callerA"},
		{"lambda caller", callerC, "github.com/kamilsk/tracer_test.callerC"},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			if expected, obtained := tc.expected, tc.caller().Name; !reflect.DeepEqual(tc.expected, obtained) {
				t.Errorf("\n expected: %+#v \n obtained: %+#v", expected, obtained)
			}
		})
	}
}

// BenchmarkCaller/direct_caller-12         	 5000000	       308 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCaller/chain_caller-12          	 5000000	       272 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCaller/lambda_caller-12         	 5000000	       365 ns/op	       0 B/op	       0 allocs/op
func BenchmarkCaller(b *testing.B) {
	benchmarks := []struct {
		name   string
		caller func() CallerInfo
	}{
		{"direct caller", callerA},
		{"chain caller", callerB},
		{"lambda caller", callerC},
	}
	for _, bm := range benchmarks {
		tc := bm
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = tc.caller()
			}
		})
	}
}
