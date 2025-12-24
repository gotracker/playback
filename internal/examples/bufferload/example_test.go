package main

import (
	"testing"
	"time"
)

func TestExamplePlayBufferToStdoutSkips(t *testing.T) {
	t.Setenv("GOTRACKER_SKIP_EXAMPLES", "1")
	done := make(chan struct{})
	go func() {
		ExamplePlayBufferToStdout()
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatalf("example did not return promptly when skip env set")
	}
}

func TestModfileDataPresent(t *testing.T) {
	if len(modfile) == 0 {
		t.Fatalf("modfile bytes should not be empty")
	}
	if modfile[0] != '1' {
		t.Fatalf("unexpected modfile header")
	}
}
