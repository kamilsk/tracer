package tracer_test

import (
	"strings"
	"testing"

	. "github.com/kamilsk/tracer"
)

func TestCaller(t *testing.T) {
	tests := []struct {
		name     string
		caller   func() CallerInfo
		expected []string
	}{
		{"direct caller", callerA, []string{"github.com/kamilsk/tracer_test.callerA"}},
		{"chain caller", callerB, []string{"github.com/kamilsk/tracer_test.callerA"}},
		{"lambda caller", callerC, []string{
			"github.com/kamilsk/tracer_test.callerC",
			"github.com/kamilsk/tracer_test.callerC.func1", // Go 1.10, 1.11 - https://golang.org/doc/go1.12#runtime
		}},
	}
	for _, test := range tests {
		tc := test
		t.Run(test.name, func(t *testing.T) {
			var found bool
			obtained := tc.caller().Name
			for _, expected := range tc.expected {
				if expected == obtained {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("\n expected: %+#v \n obtained: %+#v", strings.Join(tc.expected, " or "), obtained)
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
