package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallCopiesExecutableAt0755(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "hc")
	if err := install(dest); err != nil {
		t.Fatalf("install: %v", err)
	}

	info, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Errorf("mode = %v, want -rwxr-xr-x", info.Mode().Perm())
	}

	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	want, _ := os.ReadFile(self)
	got, _ := os.ReadFile(dest)
	if !bytes.Equal(want, got) {
		t.Error("copied bytes differ from the running executable")
	}
}

func TestInstallErrorsWhenParentMissing(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "nope", "hc")
	if err := install(dest); err == nil {
		t.Fatal("want error when parent dir is missing, got nil")
	}
}

func TestInstallErrorsWhenDestIsDir(t *testing.T) {
	if err := install(t.TempDir()); err == nil {
		t.Fatal("want error when dest is a directory, got nil")
	}
}
