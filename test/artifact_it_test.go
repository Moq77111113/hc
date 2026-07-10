package test

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Moq77111113/hc/test/support"
)

// TestScratchImageProbes runs the shipped scratch image against real services
// over a network: it proves the artifact probes tcp and, with no CA bundle, https.
func TestScratchImageProbes(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)

	nw, err := network.New(context.Background())
	testcontainers.CleanupNetwork(t, nw)
	if err != nil {
		t.Fatalf("network: %v", err)
	}

	support.Start(t, support.OnNetwork(pgRequest(), nw, "pg"))
	support.Start(t, support.OnNetwork(nginxTLSRequest(t), nw, "web"))

	if code := runHCImage(t, nw.Name, "tcp://pg:5432"); code != 0 {
		t.Errorf("image tcp probe: exit %d, want 0", code)
	}
	if code := runHCImage(t, nw.Name, "https://web:443"); code != 0 {
		t.Errorf("image https probe (no CA bundle): exit %d, want 0", code)
	}
}

// runHCImage runs the scratch image as a one-shot container on nw and returns its exit code.
func runHCImage(t *testing.T, nw, target string) int {
	t.Helper()
	c := support.Start(t, testcontainers.ContainerRequest{
		Image:      support.Image,
		Networks:   []string{nw},
		Cmd:        []string{target},
		WaitingFor: wait.ForExit(),
	})
	state, err := c.State(context.Background())
	if err != nil {
		t.Fatalf("state: %v", err)
	}
	return state.ExitCode
}

// TestInstallInContainer proves `hc install` self-copies a working 0755 binary.
func TestInstallInContainer(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	outDir := t.TempDir()

	c := support.Start(t, testcontainers.ContainerRequest{
		Image:      support.Image,
		Cmd:        []string{"install", "/out/hc"},
		WaitingFor: wait.ForExit(),
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.Binds = []string{outDir + ":/out"}
		},
	})
	state, err := c.State(context.Background())
	if err != nil {
		t.Fatalf("state: %v", err)
	}
	if state.ExitCode != 0 {
		t.Fatalf("install: exit %d, want 0", state.ExitCode)
	}

	installed := filepath.Join(outDir, "hc")
	info, err := os.Stat(installed)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o755 {
		t.Errorf("mode %v, want 0755", perm)
	}

	out, runErr := exec.Command(installed).CombinedOutput() //nolint:gosec // installed test binary
	var exitErr *exec.ExitError
	if !errors.As(runErr, &exitErr) {
		t.Fatalf("run installed: %v\n%s", runErr, out)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("installed no-arg: exit %d, want 1\n%s", exitErr.ExitCode(), out)
	}
}
