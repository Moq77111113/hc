package test

import (
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

func TestTCPProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, nginxRequest()), httpPort)

	assertOpenClosed(t, "tcp", addr)
}

// TestTimeout: an unreachable target must fail within the -timeout budget.
func TestTimeout(t *testing.T) {
	began := time.Now()
	code, out := runHC(t, "-timeout", "500ms", "tcp://10.255.255.1:1")

	if code != 1 {
		t.Errorf("exit %d, want 1\n%s", code, out)
	}
	if elapsed := time.Since(began); elapsed >= 2*time.Second {
		t.Errorf("took %s, want < 2s", elapsed)
	}
}
