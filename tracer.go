package tracer

import (
	"path"
	"strconv"
	"strings"
	"time"
)

// Trace holds information about a current execution flow.
type Trace struct {
	stack     []*Call
	allocates int
}

// Start creates a call entry and marks its start time.
//
//  func Do(ctx context.Context) {
//  	call := tracer.Fetch(ctx).Start("id", "labelX", "labelY")
//  	defer call.Stop()
//  }
//
func (trace *Trace) Start(labels ...string) *Call {
	if trace == nil {
		return nil
	}

	var id string
	if len(labels) > 0 {
		id, labels = labels[0], labels[1:]
	}
	call := &Call{id: id, labels: labels, caller: Caller(3), start: time.Now()}
	if len(trace.stack) == cap(trace.stack) {
		trace.allocates++
	}
	trace.stack = append(trace.stack, call)
	return call
}

// String returns a string representation of the current execution flow.
func (trace *Trace) String() string {
	if trace == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString("allocates at call stack: ")
	builder.WriteString(strconv.Itoa(trace.allocates))
	builder.WriteString(", detailed call stack:")
	if len(trace.stack) == 0 {
		builder.WriteString(" ~")
	}
	for _, call := range trace.stack {
		builder.WriteString("\n\tcall ")
		builder.WriteString(path.Base(call.caller.Name))
		if call.id != "" {
			builder.WriteString(" [")
			builder.WriteString(call.id)
			builder.WriteRune(']')
		}
		builder.WriteString(": ")
		builder.WriteString(call.stop.Sub(call.start).String())
		builder.WriteString(", allocates: ")
		builder.WriteString(strconv.Itoa(call.allocates))

		prev := call.start
		for _, checkpoint := range call.checkpoints {
			builder.WriteString("\n\t\tcheckpoint")
			if checkpoint.id != "" {
				builder.WriteString(" [")
				builder.WriteString(checkpoint.id)
				builder.WriteRune(']')
			}
			builder.WriteString(": ")
			builder.WriteString(checkpoint.timestamp.Sub(prev).String())
			prev = checkpoint.timestamp
		}
	}

	return builder.String()
}

// Call holds information about a current function call.
type Call struct {
	id          string
	labels      []string
	caller      CallerInfo
	start, stop time.Time
	checkpoints []Checkpoint
	allocates   int
}

// Checkpoint stores timestamp of a current execution position of the current call.
//
//  func Do(ctx context.Context) {
//  	call := tracer.Fetch(ctx).Start()
//  	defer call.Stop()
//  	...
//  	call.Checkpoint()
//  	...
//  	call.Checkpoint("id", "labelX", "labelY")
//  	...
//  }
//
func (call *Call) Checkpoint(labels ...string) {
	if call == nil {
		return
	}

	var id string
	if len(labels) > 0 {
		id, labels = labels[0], labels[1:]
	}
	checkpoint := Checkpoint{id: id, labels: labels, timestamp: time.Now()}
	if len(call.checkpoints) == cap(call.checkpoints) {
		call.allocates++
	}
	call.checkpoints = append(call.checkpoints, checkpoint)
}

// Stop marks the end time of the current call.
//
//  func Do(ctx context.Context) {
//  	defer tracer.Fetch(ctx).Start().Stop()
//  	...
//  }
//
func (call *Call) Stop() {
	if call == nil {
		return
	}

	call.stop = time.Now()
}

// Checkpoint holds information about a current execution position of a current call.
type Checkpoint struct {
	id        string
	labels    []string
	timestamp time.Time
}
