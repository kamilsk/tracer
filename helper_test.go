package tracer_test

import (
	"context"

	. "github.com/kamilsk/tracer"
)

//go:noinline
func callerA() CallerInfo {
	return Caller(2)
}

func callerB() CallerInfo {
	return callerA()
}

func callerC() CallerInfo {
	return func() CallerInfo {
		return Caller(2)
	}()
}

func traceRoot(ctx context.Context) {
	call := Fetch(ctx).Start().Mark("root")
	defer call.Stop()

	call.Checkpoint("checkpointA")
	traceA(ctx)

	call.Checkpoint("checkpointB")
	traceB(ctx)
}

func traceA(ctx context.Context) {
	call := Fetch(ctx).Start().Mark("A")
	defer call.Stop()

	call.Checkpoint("checkpointA1")
	traceA1(ctx)

	call.Checkpoint("checkpointA2")
	traceA2(ctx)
}

func traceA1(ctx context.Context) {
	defer Fetch(ctx).Start().Stop()
	func(ctx context.Context) {
		defer Fetch(ctx).Start().Stop()
	}(ctx)
}

func traceA2(ctx context.Context) {
	defer Fetch(ctx).Start().Stop()
}

func traceB(ctx context.Context) {
	call := Fetch(ctx).Start().Mark("B")
	defer call.Stop()

	call.Checkpoint("checkpointB1")
	traceB1(ctx)

	call.Checkpoint("checkpointB2")
	func(ctx context.Context) {
		defer Fetch(ctx).Start().Stop()
	}(ctx)
}

func traceB1(ctx context.Context) {
	defer Fetch(ctx).Start().Stop()
	func(ctx context.Context) {
		defer Fetch(ctx).Start().Stop()
	}(ctx)
}
