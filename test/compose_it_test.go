package test

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	composeFile    = "testdata/compose.yaml"
	composeProject = "hc-it-test"
	healthDeadline = 30 * time.Second
	healthPoll     = 500 * time.Millisecond
)

// TestComposeInjectionGoesHealthy mirrors deploy/docker-compose.yaml: hc-install
// seeds a volume, the unmodified nginx app execs the injected hc from its
// HEALTHCHECK, and must reach Docker's "healthy" state.
func TestComposeInjectionGoesHealthy(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)

	t.Cleanup(func() {
		down := exec.Command("docker", "compose", "-f", composeFile, "-p", composeProject, "down", "-v")
		if out, err := down.CombinedOutput(); err != nil {
			t.Logf("compose down: %v\n%s", err, out)
		}
	})

	up := exec.Command("docker", "compose", "-f", composeFile, "-p", composeProject, "up", "-d")
	if out, err := up.CombinedOutput(); err != nil {
		t.Fatalf("compose up: %v\n%s", err, out)
	}

	containerID := composeContainerID(t)
	waitForHealthy(t, containerID)
}

// composeContainerID resolves the running "app" service to its container ID.
func composeContainerID(t *testing.T) string {
	t.Helper()

	cmd := exec.Command("docker", "compose", "-f", composeFile, "-p", composeProject, "ps", "-q", "app")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("compose ps app: %v", err)
	}
	id := strings.TrimSpace(string(out))
	if id == "" {
		t.Fatal("compose ps app: no container ID")
	}
	return id
}

// waitForHealthy polls the container's health status until it reports
// "healthy" or healthDeadline elapses.
func waitForHealthy(t *testing.T, containerID string) {
	t.Helper()

	deadline := time.Now().Add(healthDeadline)
	var lastStatus string
	for time.Now().Before(deadline) {
		cmd := exec.Command("docker", "inspect", "-f", "{{.State.Health.Status}}", containerID)
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("docker inspect health: %v", err)
		}
		lastStatus = strings.TrimSpace(string(out))
		if lastStatus == "healthy" {
			return
		}
		time.Sleep(healthPoll)
	}
	t.Fatalf("app container did not become healthy within %s, last status %q", healthDeadline, lastStatus)
}

// TestDeployComposeConfigValid asserts the shipped example in deploy/ is
// valid Compose, independent of whether its images are pullable.
func TestDeployComposeConfigValid(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)

	cmd := exec.Command("docker", "compose", "-f", "../deploy/docker-compose.yaml", "config")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("docker compose config: %v\n%s", err, out)
	}
}
