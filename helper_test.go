package tracer_test

import . "github.com/kamilsk/tracer"

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
