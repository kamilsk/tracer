package tracer

import "runtime"

// Caller returns information about the current caller.
//
//  func StoreToDatabase(data Payload) error {
//  	timer := stats.NewTiming()
//  	defer timer.Send(Caller(2).Name)
//
//  	// do something heavy
//  }
//
func Caller(skip int) CallerInfo {
	pc := make([]uintptr, 1)
	runtime.Callers(skip, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return CallerInfo{f.Name(), file, line}
}

// CallerInfo holds information about a caller.
type CallerInfo struct {
	Name string
	File string
	Line int
}
