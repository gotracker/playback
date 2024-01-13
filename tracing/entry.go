package tracing

import (
	"fmt"
	"strings"
)

type entry[TPrefix Ticker, TPayload fmt.Stringer] struct {
	prefix    TPrefix
	operation string
	comment   string
	payload   TPayload
}

func (e entry[TPrefix, TPayload]) String() string {
	var chunks []string
	if len(e.operation) > 0 {
		chunks = append(chunks, e.operation)
	}
	if line := fmt.Sprint(e.payload); len(line) > 0 {
		chunks = append(chunks, line)
	}
	if len(e.comment) > 0 {
		chunks = append(chunks, "//", e.comment)
	}
	return strings.Join(chunks, " ")
}

func (e entry[TPrefix, TPayload]) GetTick() Tick {
	return e.prefix.GetTick()
}

func (e entry[TPrefix, TPayload]) Prefix() string {
	return e.prefix.String()
}

///////////////////////////////////////////////////////////

func (t *Tracing) trace(tick Tick, op string) {
	t.traceWithComment(tick, op, "")
}

type emptyPayload struct{}

func (emptyPayload) String() string {
	return ""
}

var empty emptyPayload

func (t *Tracing) traceWithComment(tick Tick, op, comment string) {
	if t.tracingFile == nil {
		return
	}
	traceWithPayload(t, tick, op, comment, empty)
}

func traceWithPayload[TPrefix Ticker, TPayload fmt.Stringer](t *Tracing, prefix TPrefix, op, comment string, payload TPayload) {
	e := entry[TPrefix, TPayload]{
		prefix:    prefix,
		operation: op,
		comment:   comment,
		payload:   payload,
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.updates = append(t.updates, e)
}
