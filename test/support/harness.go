package support

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Binary is the path to the built host hc binary; Image is the scratch image
// tag. Both are set by Build and consumed by the scenario tests.
var (
	Binary string
	Image  = "hc-it:latest"
)

// Build compiles the host binary and the scratch image once. TestMain calls it
// with the test directory as the working directory, so ".." is the repo root.
func Build() error {
	bin, err := buildHostBinary()
	if err != nil {
		return err
	}
	Binary = bin
	return buildScratchImage()
}

func buildHostBinary() (string, error) {
	dir, err := os.MkdirTemp("", "hc-it-bin-*")
	if err != nil {
		return "", err
	}

	bin := filepath.Join(dir, "hc")
	cmd := exec.Command("go", "build", "-o", bin, ".") //nolint:gosec // building our own binary from the repo root
	cmd.Dir = ".."
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build: %w\n%s", err, out)
	}
	return bin, nil
}

func buildScratchImage() error {
	cmd := exec.Command("docker", "build", "-t", Image, "..")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker build: %w\n%s", err, out)
	}
	return nil
}

// RunHC runs the built binary and returns its exit code and combined output.
func RunHC(t *testing.T, args ...string) (int, string) {
	t.Helper()
	cmd := exec.Command(Binary, args...) //nolint:gosec // running the hc binary we just built
	out, err := cmd.CombinedOutput()
	if cmd.ProcessState == nil {
		t.Fatalf("hc did not run: %v\n%s", err, out)
	}
	return cmd.ProcessState.ExitCode(), string(out)
}

// AssertOpenClosed checks scheme is healthy against addr and unhealthy against
// a closed port on the same host: the shared shape of the connect-based probes.
func AssertOpenClosed(t *testing.T, scheme, addr string) {
	t.Helper()
	if code, out := RunHC(t, scheme+"://"+addr); code != 0 {
		t.Errorf("open: exit %d, want 0\n%s", code, out)
	}
	if code, out := RunHC(t, scheme+"://"+HostOf(t, addr)+":1"); code != 1 {
		t.Errorf("closed: exit %d, want 1\n%s", code, out)
	}
}

// HostOf strips the port from a "host:port" endpoint.
func HostOf(t *testing.T, addr string) string {
	t.Helper()
	h, _, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("split %q: %v", addr, err)
	}
	return h
}
