package test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Moq77111113/hc/test/support"
)

const (
	nginxImage = "nginx:alpine"
	httpPort   = "80"
	httpsPort  = "443"
)

func nginxRequest() testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        nginxImage,
		ExposedPorts: []string{httpPort + "/tcp"},
		WaitingFor:   wait.ForHTTP("/"),
	}
}

func nginxTLSRequest(t *testing.T) testcontainers.ContainerRequest {
	t.Helper()
	dir := t.TempDir()
	support.SelfSignedCert(t, dir)
	return testcontainers.ContainerRequest{
		Image:        nginxImage,
		ExposedPorts: []string{httpsPort + "/tcp"},
		WaitingFor:   wait.ForListeningPort(httpsPort + "/tcp"),
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

func TestHTTPProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := support.Endpoint(t, support.Start(t, nginxRequest()), httpPort)

	if code, out := support.RunHC(t, "http://"+addr+"/"); code != 0 {
		t.Errorf("http: exit %d, want 0\n%s", code, out)
	}
}

func TestHTTPSProbe(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	addr := support.Endpoint(t, support.Start(t, nginxTLSRequest(t)), httpsPort)

	if code, out := support.RunHC(t, "https://"+addr+"/"); code != 0 {
		t.Errorf("https: exit %d, want 0\n%s", code, out)
	}
}
