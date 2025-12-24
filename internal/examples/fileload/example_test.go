package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestExamplePlayFileToStdoutSkips(t *testing.T) {
	t.Setenv("GOTRACKER_SKIP_EXAMPLES", "1")
	done := make(chan struct{})
	go func() {
		ExamplePlayFileToStdout()
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatalf("example did not return promptly when skip env set")
	}
}

func TestExampleFileExists(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	path := filepath.Join(dir, "..", "..", "..", "test", "ode_to_protracker.mod")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected test module to exist at %s: %v", path, err)
	}
}
