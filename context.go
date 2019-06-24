package tracer

import "context"

type key struct{}

// Fetch tries to get the tracer from a context or returns safe nil.
//
//  tracer.Fetch(context.Background()).Start().Stop() // won't panic
//
func Fetch(ctx context.Context) *Trace {
	trace, _ := ctx.Value(key{}).(*Trace)
	return trace
}

// Inject returns a new context with injected into it the tracer.
//
//  func (server *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
//  	req = req.WithContext(tracer.Inject(req.Context(), make([]*trace.Call, 0, 10)))
//  	server.routing.Handle(rw, req)
//  }
//
func Inject(ctx context.Context, stack []*Call) context.Context {
	return context.WithValue(ctx, key{}, &Trace{stack: stack})
}
