package test

import (
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

func TestRedisProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, redisRequest()), redisPort)

	assertOpenClosed(t, "redis", addr)
}

func TestRedisProbeNoAuthIsHealthy(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := endpoint(t, start(t, redisAuthRequest()), redisPort)

	if code, out := runHC(t, "redis://"+addr); code != 0 {
		t.Errorf("noauth: exit %d, want 0\n%s", code, out)
	}
}
