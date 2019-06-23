package tracer

import (
	"path"
	"strconv"
	"strings"
	"time"
)

type Trace struct {
	stack     []Call
	allocates int
}

func (trace *Trace) Start(labels ...string) ptr {
	if trace == nil {
		return ptr{}
	}

	var (
		id   string
		tags []string
	)
	if len(labels) > 0 {
		id = labels[0]
		tags = labels[1:]
	}
	call := Call{caller: Caller(3), start: time.Now(), id: id, tags: tags}
	if len(trace.stack) == cap(trace.stack) {
		trace.allocates++
	}
	trace.stack = append(trace.stack, call)
	return ptr{trace, len(trace.stack) - 1}
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
	caller      CallerInfo
	start, stop time.Time
	id          string
	tags        []string
	checkpoints []Checkpoint
	allocates   int
}

type Checkpoint struct {
	id        string
	tags      []string
	timestamp time.Time
}

type ptr struct {
	*Trace
	int
}

func (call ptr) Checkpoint(labels ...string) {
	var (
		id   string
		tags []string
	)
	if len(labels) > 0 {
		id = labels[0]
		tags = labels[1:]
	}
	checkpoint := Checkpoint{id: id, tags: tags, timestamp: time.Now()}
	if len(call.stack[call.int].checkpoints) == cap(call.stack[call.int].checkpoints) {
		call.stack[call.int].allocates++
	}
	call.stack[call.int].checkpoints = append(call.stack[call.int].checkpoints, checkpoint)
}

func (call ptr) Stop() {
	if call.Trace == nil {
		return
	}

	call.stack[call.int].stop = time.Now()
}
