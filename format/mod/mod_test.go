package mod

import (
	"bytes"
	"testing"
)

func TestLoadFromReaderRejectsInvalidData(t *testing.T) {
	_, err := MOD.LoadFromReader(bytes.NewReader([]byte("bad")), nil)
	if err == nil {
		t.Fatalf("expected error for invalid MOD data")
	}
}
