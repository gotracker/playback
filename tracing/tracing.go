package tracing

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/gotracker/playback/index"
)

type Tracing struct {
	tracingFile *os.File
	chMap       map[int]*tracingChannelState
	traces      []tracingMsgFunc
	c           chan func(w io.Writer)
	wg          sync.WaitGroup

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

func (t *Tracing) EnableTracing(filename string) error {
	var err error
	t.tracingFile, err = os.Create(filename)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tracing) Close() {
	if t.c != nil {
		close(t.c)
	}
	if t.tracingFile != nil {
		t.tracingFile.Close()
	}
	t.wg.Wait()
}

func (t *Tracing) OutputTraces() {
	if t.tracingFile == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	var updates []entryIntf
	updates, t.updates = t.updates, nil

	go func() {
		logger := log.New(t.tracingFile, "", 0)
		for _, u := range updates {
			if tick := u.GetTick(); !tick.Equals(t.prevTick) {
				fmt.Fprintln(t.tracingFile)
				t.prevTick = tick
			}

			logger.Println("[" + u.Prefix() + "] " + u.String())
		}
	}()
}

func (t *Tracing) SetTracingTick(order index.Order, row index.Row, tick int) {
	t.mu.Lock()
	t.tick = Tick{
		Order: order,
		Row:   row,
		Tick:  tick,
	}
	t.mu.Unlock()
}

func (t *Tracing) GetTracingTick() Tick {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.tick
}

func (t *Tracing) Trace(op string) {
	t.TraceWithComment(op, "")
}

func (t *Tracing) TraceWithComment(op, comment string) {
	traceWithPayload(t, t.GetTracingTick(), op, comment, empty)
}

func (t *Tracing) TraceValueChange(op string, prev, new any) {
	t.TraceValueChangeWithComment(op, prev, new, "")
}

func (t *Tracing) TraceValueChangeWithComment(op string, prev, new any, comment string) {
	traceWithPayload(t, t.GetTracingTick(), op, comment, valueUpdate{
		old: prev,
		new: new,
	})
}

func (t *Tracing) TraceChannel(ch index.Channel, op string) {
	t.TraceChannelWithComment(ch, op, "")
}

func (t *Tracing) TraceChannelWithComment(ch index.Channel, op, comment string) {
	tc := tickChannel{
		tick: t.GetTracingTick(),
		ch:   ch,
	}
	traceWithPayload(t, tc, op, comment, empty)
}

func (t *Tracing) TraceChannelValueChange(ch index.Channel, op string, prev, new any) {
	t.TraceChannelValueChangeWithComment(ch, op, prev, new, "")
}

func (t *Tracing) TraceChannelValueChangeWithComment(ch index.Channel, op string, prev, new any, comment string) {
	tc := tickChannel{
		tick: t.GetTracingTick(),
		ch:   ch,
	}
	traceWithPayload(t, tc, op, comment, valueUpdate{
		old: prev,
		new: new,
	})
}
