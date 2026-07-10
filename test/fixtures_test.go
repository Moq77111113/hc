package test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// start launches req, registers cleanup, and fails the test on error.
func start(t *testing.T, req testcontainers.ContainerRequest) testcontainers.Container {
	t.Helper()
	ctr, err := testcontainers.GenericContainer(context.Background(),
		testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
	testcontainers.CleanupContainer(t, ctr)
	if err != nil {
		t.Fatalf("start %s: %v", req.Image, err)
	}
	return ctr
}

// endpoint returns the container's "host:port" for an exposed port.
func endpoint(t *testing.T, ctr testcontainers.Container, port string) string {
	t.Helper()
	ep, err := ctr.PortEndpoint(context.Background(), port+"/tcp", "")
	if err != nil {
		t.Fatalf("endpoint: %v", err)
	}
	return ep
}

// onNetwork attaches req to nw under alias, for container-to-container probing.
func onNetwork(req testcontainers.ContainerRequest, nw *testcontainers.DockerNetwork, alias string) testcontainers.ContainerRequest {
	req.Networks = []string{nw.Name}
	req.NetworkAliases = map[string][]string{nw.Name: {alias}}
	return req
}

func pgRequest() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        "postgres:18-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env:          map[string]string{"POSTGRES_PASSWORD": "pw"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
}

func redisRequest() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        "redis:8-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
}

// redisAuthRequest requires a password, so an unauthenticated PING gets -NOAUTH.
func redisAuthRequest() testcontainers.ContainerRequest {
	req := redisRequest()
	req.Cmd = []string{"redis-server", "--requirepass", "secret"}
	return req
}

func nginxRequest() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        "nginx:alpine",
		ExposedPorts: []string{"80/tcp"},
		WaitingFor:   wait.ForHTTP("/"),
	}
}

// nginxTLSRequest serves HTTPS with a fresh self-signed cert.
func nginxTLSRequest(t *testing.T) testcontainers.ContainerRequest {
	t.Helper()
	dir := t.TempDir()
	selfSignedCert(t, dir)
	return testcontainers.ContainerRequest{
		Image:        "nginx:alpine",
		ExposedPorts: []string{"443/tcp"},
		WaitingFor:   wait.ForListeningPort("443/tcp"),
		Files: []testcontainers.ContainerFile{
			{HostFilePath: filepath.Join(dir, "cert.pem"), ContainerFilePath: "/etc/nginx/certs/cert.pem", FileMode: 0o644},
			{HostFilePath: filepath.Join(dir, "key.pem"), ContainerFilePath: "/etc/nginx/certs/key.pem", FileMode: 0o644},
			{Reader: strings.NewReader(nginxTLSConf), ContainerFilePath: "/etc/nginx/conf.d/tls.conf", FileMode: 0o644},
		},
	}
}

const nginxTLSConf = `server {
    listen 443 ssl;
    ssl_certificate     /etc/nginx/certs/cert.pem;
    ssl_certificate_key /etc/nginx/certs/key.pem;
    location / { return 200 "ok"; }
}
`
