package tracing

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gotracker/playback/index"
)

type tracerFile struct {
	file   *os.File
	chMap  map[int]*tracingChannelState
	traces []tracingMsgFunc
	c      chan func(w io.Writer)
	wg     sync.WaitGroup

	tick     Tick
	updates  []entryIntf
	prevTick Tick
	mu       sync.RWMutex
}

type entryIntf interface {
	GetTick() Tick
	Prefix() string
	fmt.Stringer
}

type tracingMsgFunc func() string

type tracingChannelState struct {
	traces []tracingMsgFunc
}

func (t *tracerFile) Close() error {
	if t.c != nil {
		close(t.c)
	}
	if t.file != nil {
		if err := t.file.Close(); err != nil {
			return err
		}
	}
	t.wg.Wait()
	return nil
}

func (t *tracerFile) OutputTraces() {
	if t.file == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	var updates []entryIntf
	updates, t.updates = t.updates, nil

	go func() {
		logger := log.New(t.file, "", 0)
		for _, u := range updates {
			if tick := u.GetTick(); !tick.Equals(t.prevTick) {
				fmt.Fprintln(t.file)
				t.prevTick = tick
			}

			logger.Println("[" + u.Prefix() + "] " + u.String())
		}
	}()
}

func (t *tracerFile) SetTracingTick(order index.Order, row index.Row, tick int) {
	t.mu.Lock()
	t.tick = Tick{
		Order: order,
		Row:   row,
		Tick:  tick,
	}
	t.mu.Unlock()
}

func (t *tracerFile) GetTracingTick() Tick {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.tick
}

func (t *tracerFile) Trace(op string) {
	t.TraceWithComment(op, "")
}

func (t *tracerFile) TraceWithComment(op, commentFmt string, commentParams ...any) {
	traceWithPayload(t, t.GetTracingTick(), op, fmt.Sprintf(commentFmt, commentParams...), empty)
}

func (t *tracerFile) TraceValueChange(op string, prev, new any) {
	t.TraceValueChangeWithComment(op, prev, new, "")
}

func (t *tracerFile) TraceValueChangeWithComment(op string, prev, new any, commentFmt string, commentParams ...any) {
	traceWithPayload(t, t.GetTracingTick(), op, fmt.Sprintf(commentFmt, commentParams...), valueUpdate{
		old: prev,
		new: new,
	})
}

func (t *tracerFile) TraceChannel(ch index.Channel, op string) {
	t.TraceChannelWithComment(ch, op, "")
}

func (t *tracerFile) TraceChannelWithComment(ch index.Channel, op, commentFmt string, commentParams ...any) {
	tc := tickChannel{
		tick: t.GetTracingTick(),
		ch:   ch,
	}
	traceWithPayload(t, tc, op, fmt.Sprintf(commentFmt, commentParams...), empty)
}

func (t *tracerFile) TraceChannelValueChange(ch index.Channel, op string, prev, new any) {
	t.TraceChannelValueChangeWithComment(ch, op, prev, new, "")
}

func (t *tracerFile) TraceChannelValueChangeWithComment(ch index.Channel, op string, prev, new any, commentFmt string, commentParams ...any) {
	tc := tickChannel{
		tick: t.GetTracingTick(),
		ch:   ch,
	}
	traceWithPayload(t, tc, op, fmt.Sprintf(commentFmt, commentParams...), valueUpdate{
		old: prev,
		new: new,
	})
}
