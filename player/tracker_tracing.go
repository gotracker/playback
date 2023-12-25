package player

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"text/tabwriter"

	ansi "github.com/fatih/color"
)

type tracingMsgFunc func() string

type tracingState struct {
	chMap  map[int]*tracingChannelState
	traces []tracingMsgFunc
	c      chan func(w io.Writer)
	wg     sync.WaitGroup
}

type tracingChannelState struct {
	traces []tracingMsgFunc
}

func (t *Tracker) TraceChannel(ch int, msgFunc tracingMsgFunc) {
	if t.tracingFile == nil {
		return
	}

	tc := t.tracingState.chMap[ch]
	if tc == nil {
		tc = &tracingChannelState{}
		t.tracingState.chMap[ch] = tc
	}

	tc.traces = append(tc.traces, msgFunc)
}

func (t *Tracker) TraceTick(msgFunc tracingMsgFunc) {
	if t.tracingFile == nil {
		return
	}

	t.tracingState.traces = append(t.tracingState.traces, msgFunc)
}

type tracingColumn struct {
	heading string
	rows    []string
}

type TracingTable struct {
	cols      []*tracingColumn
	rowColors []*ansi.Color
	name      string
	maxRows   int
}

func NewTracingTable(name string, headers ...string) TracingTable {
	tt := TracingTable{
		name: name,
	}
	for _, h := range headers {
		tt.cols = append(tt.cols, &tracingColumn{
			heading: h,
		})
	}
	return tt
}

func (tt *TracingTable) AddRow(cols ...any) {
	tt.AddRowColor(nil, cols...)
}

func (tt *TracingTable) AddRowColor(c []ansi.Attribute, cols ...any) {
	var cc *ansi.Color
	if len(c) > 0 {
		cc = ansi.Set(c...)
		cc.EnableColor()
	}
	tt.rowColors = append(tt.rowColors, cc)
	for i, col := range cols {
		tc := tt.cols[i]
		tc.rows = append(tc.rows, fmt.Sprint(col))
	}
	tt.maxRows++
}

func (tt TracingTable) Fprintln(w io.Writer, colSep string, withRowNums bool) error {
	head := []string{tt.name}
	for _, c := range tt.cols {
		head = append(head, c.heading)
	}
	if _, err := fmt.Fprintln(w, strings.Join(head, colSep)); err != nil {
		return err
	}

	for r := 0; r < tt.maxRows; r++ {
		numCols := len(tt.cols)
		colStart := 0
		if withRowNums {
			numCols++
			colStart++
		}
		cols := []string{""}
		if withRowNums {
			cols[0] = fmt.Sprintf("[%d]", r+1)
		}
		for _, col := range tt.cols {
			if r >= len(col.rows) {
				return errors.New("not enough rows to satisfy TracingTable writer")
			}
			cols = append(cols, col.rows[r])
		}
		if _, err := fmt.Fprintln(w, strings.Join(cols, colSep)); err != nil {
			return err
		}
	}

	return nil
}

func (tt *TracingTable) WriteOut(w io.Writer) error {
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 1, 1, 1, ' ', 0)
	if err := tt.Fprintln(tw, "\t", true); err != nil {
		return err
	}
	if err := tw.Flush(); err != nil {
		return err
	}

	for r, txt := range strings.Split(sb.String(), "\n") {
		var rc *ansi.Color
		if rcr := r - 1; rcr >= 0 && rcr < len(tt.rowColors) {
			rc = tt.rowColors[rcr]
		}
		if rc != nil {
			txt = rc.Sprint(txt)
		}
		if _, err := fmt.Fprintln(w, txt); err != nil {
			return err
		}
	}
	return nil
}

type Traceable interface {
	OutputTraces(out chan<- func(w io.Writer))
}

func (t *Tracker) OutputTraces() {
	if t.tracingFile != nil && t.Traceable != nil {
		if t.tracingState.c == nil {
			t.tracingState.c = make(chan func(w io.Writer), 1000*1000)
			go func() {
				defer close(t.tracingState.c)
				defer t.tracingFile.Close()

				t.tracingState.wg.Add(1)
				defer t.tracingState.wg.Done()

				for tr := range t.tracingState.c {
					tr(t.tracingFile)
				}
			}()
		}
		t.Traceable.OutputTraces(t.tracingState.c)
	}
}
