package xm

import (
	"bytes"
	"testing"
)

func TestLoadFromReaderRejectsInvalidData(t *testing.T) {
	_, err := XM.LoadFromReader(bytes.NewReader([]byte("bad")), nil)
	if err == nil {
		t.Fatalf("expected error for invalid XM data")
	}
}
