package tracer

import "runtime"

// Caller returns information about a caller at position after the skip steps.
//
//  func StoreToDatabase(data Payload) error {
//  	defer stats.NewTiming().Send(Caller(2).Name)
//
//  	// do something heavy
//  }
//
func Caller(skip int) CallerInfo {
	pc := make([]uintptr, 1)
	runtime.Callers(skip, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return CallerInfo{f.Entry(), f.Name(), file, line}
}

// CallerInfo holds information about a caller.
type CallerInfo struct {
	Entry uintptr // the entry address of the function
	Name  string  // the name of the function
	File  string  // the file name and
	Line  int     // line number of the source code of the function
}
