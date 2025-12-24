package tracing

import (
	"fmt"
	"strings"
)

type entry struct {
	prefix    Ticker
	operation string
	comment   string
	payload   fmt.Stringer
}

func (e entry) String() string {
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

func (e entry) GetTick() Tick {
	return e.prefix.GetTick()
}

func (e entry) Prefix() string {
	return e.prefix.String()
}

///////////////////////////////////////////////////////////

func (t *tracerFile) trace(tick Tick, op string) {
	t.traceWithComment(tick, op, "")
}

type emptyPayload struct{}

func (emptyPayload) String() string {
	return ""
}

var empty emptyPayload

func (t *tracerFile) traceWithComment(tick Tick, op, comment string) {
	if t.file == nil {
		return
	}
	traceWithPayload(t, tick, op, comment, empty)
}

func traceWithPayload(t *tracerFile, prefix Ticker, op, comment string, payload fmt.Stringer) {
	e := entry{
		prefix:    prefix,
		operation: op,
		comment:   comment,
		payload:   payload,
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.updates = append(t.updates, e)
}
