package test

import (
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

func TestPostgresProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, pgRequest()), "5432")

	if code, out := runHC(t, "postgres://"+addr); code != 0 {
		t.Errorf("open: exit %d, want 0\n%s", code, out)
	}
	if code, out := runHC(t, "postgres://"+hostOf(t, addr)+":1"); code != 1 {
		t.Errorf("closed: exit %d, want 1\n%s", code, out)
	}
}

func TestTCPProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, nginxRequest()), "80")

	if code, out := runHC(t, "tcp://"+addr); code != 0 {
		t.Errorf("open: exit %d, want 0\n%s", code, out)
	}
	if code, out := runHC(t, "tcp://"+hostOf(t, addr)+":1"); code != 1 {
		t.Errorf("closed: exit %d, want 1\n%s", code, out)
	}
}

func TestHTTPProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, nginxRequest()), "80")

	if code, out := runHC(t, "http://"+addr+"/"); code != 0 {
		t.Errorf("http: exit %d, want 0\n%s", code, out)
	}
}

func TestHTTPSProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, nginxTLSRequest(t)), "443")

	if code, out := runHC(t, "https://"+addr+"/"); code != 0 {
		t.Errorf("https: exit %d, want 0\n%s", code, out)
	}
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
