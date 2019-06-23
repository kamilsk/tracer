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

type Call struct {
	id          string
	labels      []string
	caller      CallerInfo
	start, stop time.Time
	checkpoints []*Checkpoint
	allocates   int
}

func (call *Call) Checkpoint(labels ...string) *Checkpoint {
	if call == nil {
		return nil
	}

	var id string
	if len(labels) > 0 {
		id, labels = labels[0], labels[1:]
	}
	checkpoint := &Checkpoint{id: id, labels: labels, timestamp: time.Now()}
	if len(call.checkpoints) == cap(call.checkpoints) {
		call.allocates++
	}
	call.checkpoints = append(call.checkpoints, checkpoint)

	return checkpoint
}

func (call *Call) Stop() {
	if call == nil {
		return
	}

	call.stop = time.Now()
}

type Checkpoint struct {
	id        string
	labels    []string
	timestamp time.Time
}
