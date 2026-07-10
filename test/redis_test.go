package test

import (
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Moq77111113/hc/test/support"
)

const (
	redisImage = "redis:8-alpine"
	redisPort  = "6379"
)

func redisRequest() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        redisImage,
		ExposedPorts: []string{redisPort + "/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
}

// redisAuthRequest requires a password, so an unauthenticated PING gets -NOAUTH.
func redisAuthRequest() testcontainers.ContainerRequest {
	req := redisRequest()
	req.Cmd = []string{"redis-server", "--requirepass", "secret"}
	return req
}

func TestRedisProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := support.Endpoint(t, support.Start(t, redisRequest()), redisPort)

	support.AssertOpenClosed(t, "redis", addr)
}

func TestRedisProbeNoAuthIsHealthy(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := support.Endpoint(t, support.Start(t, redisAuthRequest()), redisPort)

	if code, out := support.RunHC(t, "redis://"+addr); code != 0 {
		t.Errorf("noauth: exit %d, want 0\n%s", code, out)
	}
}
