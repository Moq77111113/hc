package main

import (
	"io"
	"os"
	"path/filepath"
)

// install self-copies the running binary to dest at 0755, atomically. 0755 so a
// non-root container can exec it; self-copy because scratch ships no cp.
func install(dest string) error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	src, err := os.Open(self) //nolint:gosec // our own binary
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	tmp, err := os.CreateTemp(filepath.Dir(dest), ".hc-install-*") //nolint:gosec // operator-chosen dir
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }() //nolint:gosec // our own temp file

	if _, err := io.Copy(tmp, src); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o755); err != nil { //nolint:gosec // must be executable
		return err
	}
	return os.Rename(tmp.Name(), dest) //nolint:gosec // operator-chosen install path
}
