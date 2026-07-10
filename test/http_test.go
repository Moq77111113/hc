package test

import (
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

func TestHTTPProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, nginxRequest()), httpPort)

	if code, out := runHC(t, "http://"+addr+"/"); code != 0 {
		t.Errorf("http: exit %d, want 0\n%s", code, out)
	}
}

func TestHTTPSProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, nginxTLSRequest(t)), httpsPort)

	if code, out := runHC(t, "https://"+addr+"/"); code != 0 {
		t.Errorf("https: exit %d, want 0\n%s", code, out)
	}
}
