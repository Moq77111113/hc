package test

import (
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"github.com/Moq77111113/hc/test/support"
)

func TestTCPProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := support.Endpoint(t, support.Start(t, nginxRequest()), httpPort)

	support.AssertOpenClosed(t, "tcp", addr)
}

func TestTimeout(t *testing.T) {
	began := time.Now()
	code, out := support.RunHC(t, "-timeout", "500ms", "tcp://10.255.255.1:1")

	if code != 1 {
		t.Errorf("exit %d, want 1\n%s", code, out)
	}
	if elapsed := time.Since(began); elapsed >= 2*time.Second {
		t.Errorf("took %s, want < 2s", elapsed)
	}
}
