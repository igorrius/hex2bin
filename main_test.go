package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetOutputFileName_NoConflict(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "file.hex")
	os.WriteFile(input, []byte("dummy"), 0644)
	name, err := GetOutputFileName(input, "hex2bin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Base(name) != "file.bin" {
		t.Errorf("expected file.bin, got %s", filepath.Base(name))
	}
}

func TestGetOutputFileName_OneConflict(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "file.hex")
	os.WriteFile(input, []byte("dummy"), 0644)
	os.WriteFile(filepath.Join(dir, "file.bin"), []byte("bin"), 0644)
	name, err := GetOutputFileName(input, "hex2bin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Base(name) != "file_1.bin" {
		t.Errorf("expected file_1.bin, got %s", filepath.Base(name))
	}
}

func TestGetOutputFileName_MultipleConflicts(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "file.hex")
	os.WriteFile(input, []byte("dummy"), 0644)
	os.WriteFile(filepath.Join(dir, "file.bin"), []byte("bin"), 0644)
	os.WriteFile(filepath.Join(dir, "file_1.bin"), []byte("bin1"), 0644)
	os.WriteFile(filepath.Join(dir, "file_2.bin"), []byte("bin2"), 0644)
	name, err := GetOutputFileName(input, "hex2bin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Base(name) != "file_3.bin" {
		t.Errorf("expected file_3.bin, got %s", filepath.Base(name))
	}
}
