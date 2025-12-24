package s3m

import (
	"bytes"
	"testing"
)

func TestLoadFromReaderRejectsInvalidData(t *testing.T) {
	_, err := S3M.LoadFromReader(bytes.NewReader([]byte("bad")), nil)
	if err == nil {
		t.Fatalf("expected error for invalid S3M data")
	}
}
