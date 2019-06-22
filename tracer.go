package tracer

import (
	"context"
	"path"
	"strconv"
	"strings"
	"time"
)

func Fetch(ctx context.Context) *Trace {
	trace, _ := ctx.Value(key{}).(*Trace)
	return trace
}

func Inject(ctx context.Context, stack []Call) context.Context {
	return context.WithValue(ctx, key{}, &Trace{in: stack, out: make([]Call, 0, len(stack))})
}

type key struct{}

type Trace struct {
	in       []Call
	out      []Call
	allocate int
}

func (trace *Trace) Start() *Trace {
	if trace == nil {
		return nil
	}

	if len(trace.in) == cap(trace.in) {
		trace.allocate++
	}
	trace.in = append(trace.in, Call{caller: Caller(3), start: time.Now()})
	return trace
}

func (trace *Trace) Stop() {
	if trace == nil {
		return
	}

	var (
		call Call
		last = len(trace.in) - 1
	)
	if last < 0 {
		return
	}
	call, trace.in = trace.in[last], trace.in[:last]
	call.stop = time.Now()
	trace.out = append(trace.out, call)
}

func (trace *Trace) String() string {
	if trace == nil {
		return ""
	}

	builder := strings.Builder{}
	for i := len(trace.out) - 1; i >= 0; i-- {
		call := trace.out[i]
		builder.WriteString(path.Base(call.caller.Name))
		builder.WriteString(": ")
		builder.WriteString(call.stop.Sub(call.start).String())
		builder.WriteRune('\n')
	}
	builder.WriteString("allocates: ")
	builder.WriteString(strconv.Itoa(trace.allocate))
	return builder.String()
}

type Call struct {
	caller      CallerInfo
	start, stop time.Time
}
