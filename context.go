package tracer

import "context"

type key struct{}

func Fetch(ctx context.Context) *Trace {
	trace, _ := ctx.Value(key{}).(*Trace)
	return trace
}

func Inject(ctx context.Context, stack []*Call) context.Context {
	return context.WithValue(ctx, key{}, &Trace{stack: stack})
}
