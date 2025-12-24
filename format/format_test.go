package format

import (
	"bytes"
	"os"
	"testing"
)

func TestLoadFromReaderUnknownFormat(t *testing.T) {
	_, _, err := LoadFromReader("zzz", bytes.NewReader(nil))
	if err == nil || err.Error() != "unsupported format" {
		t.Fatalf("expected unsupported format error, got %v", err)
	}
}

func TestLoadFromReaderNoFormatFallsThrough(t *testing.T) {
	_, _, err := LoadFromReader("", bytes.NewReader([]byte("not a module")))
	if err == nil {
		t.Fatalf("expected error for unsupported data")
	}
	if err.Error() != "unsupported format" {
		t.Fatalf("expected unsupported format error, got %v", err)
	}
}

func TestLoadNonexistentFileReturnsNotExist(t *testing.T) {
	_, _, err := Load("this_file_should_not_exist_12345.s3m")
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected os.IsNotExist, got %v", err)
	}
}

func TestLoadUnsupportedFormatFromExistingFile(t *testing.T) {
	f, err := os.CreateTemp("", "unsupported*.bin")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	name := f.Name()
	_, _ = f.WriteString("not a module")
	_ = f.Close()
	t.Cleanup(func() { _ = os.Remove(name) })

	_, _, err = Load(name)
	if err == nil {
		t.Fatalf("expected unsupported format error")
	}
	if err.Error() != "unsupported format" {
		t.Fatalf("expected unsupported format, got %v", err)
	}
}
