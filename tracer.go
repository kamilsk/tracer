package tracer

import (
	"path"
	"strconv"
	"strings"
	"time"
)

type Trace struct {
	stack     []*Call
	allocates int
}

func (trace *Trace) Start() *Call {
	if trace == nil {
		return nil
	}

	if len(trace.stack) == cap(trace.stack) {
		trace.allocates++
	}
	call := &Call{caller: Caller(3), start: time.Now()}
	trace.stack = append(trace.stack, call)
	return call
}

func (trace *Trace) String() string {
	if trace == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString("allocates at call stack: ")
	builder.WriteString(strconv.Itoa(trace.allocates))
	builder.WriteString(", detailed call stack:\n")
	for _, call := range trace.stack {
		builder.WriteRune('\t')
		builder.WriteString(path.Base(call.caller.Name))
		builder.WriteString(": ")
		builder.WriteString(call.stop.Sub(call.start).String())
		builder.WriteString(", allocates: ")
		builder.WriteString(strconv.Itoa(call.allocates))
		builder.WriteRune('\n')

		prev := call.start
		for _, checkpoint := range call.checkpoints {
			builder.WriteString("\t\t")
			builder.WriteString(checkpoint.ID())
			builder.WriteString(": ")
			builder.WriteString(checkpoint.timestamp.Sub(prev).String())
			builder.WriteRune('\n')
			prev = checkpoint.timestamp
		}
	}

	return builder.String()
}

type Call struct {
	caller      CallerInfo
	start, stop time.Time
	id          string
	checkpoints []*Checkpoint
	allocates   int
}

func (call *Call) Checkpoint() *Checkpoint {
	if call == nil {
		return nil
	}

	checkpoint := &Checkpoint{timestamp: time.Now()}
	if len(call.checkpoints) == cap(call.checkpoints) {
		call.allocates++
	}
	call.checkpoints = append(call.checkpoints, checkpoint)

	return checkpoint
}

func (call *Call) ID() string {
	return ""
}

func (call *Call) Mark(id string) *Call {
	if call == nil {
		return nil
	}

	call.id = id
	return call
}

func (call *Call) Stop() {
	if call == nil {
		return
	}

	call.stop = time.Now()
}

type Checkpoint struct {
	id        string
	timestamp time.Time
}

func (checkpoint *Checkpoint) ID() string {
	if checkpoint.id == "" {
		return "checkpoint"
	}
	return checkpoint.id
}

func (checkpoint *Checkpoint) Mark(id string) {
	if checkpoint == nil {
		return
	}

	checkpoint.id = id
}
