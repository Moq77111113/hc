//go:build !hc_slim || hc_http || hc_https

package probe

import (
	"context"
	"net"
	"net/url"
	"testing"
)

func TestProbeHTTPHealthyOn2xx(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	go func() {
		defer server.Close()
		buf := make([]byte, 512)
		server.Read(buf)
		server.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"))
	}()

	u, _ := url.Parse("http://x/health")
	if err := probeHTTP(context.Background(), client, u); err != nil {
		t.Fatalf("want healthy, got %v", err)
	}
}

func TestProbeHTTPSkips1xxInterim(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	go func() {
		defer server.Close()
		buf := make([]byte, 512)
		server.Read(buf)
		server.Write([]byte("HTTP/1.1 100 Continue\r\n\r\nHTTP/1.1 200 OK\r\n\r\n"))
	}()

	u, _ := url.Parse("http://x/health")
	if err := probeHTTP(context.Background(), client, u); err != nil {
		t.Fatalf("want healthy after 1xx interim, got %v", err)
	}
}

func TestProbeHTTPUnhealthyOn5xx(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	go func() {
		defer server.Close()
		buf := make([]byte, 512)
		server.Read(buf)
		server.Write([]byte("HTTP/1.1 503 Service Unavailable\r\n\r\n"))
	}()

	u, _ := url.Parse("http://x/health")
	if err := probeHTTP(context.Background(), client, u); err == nil {
		t.Fatal("want error on 503, got nil")
	}
}

func TestStatusCode(t *testing.T) {
	cases := []struct {
		line string
		code int
		ok   bool
	}{
		{"HTTP/1.1 200 OK\r\n", 200, true},
		{"HTTP/1.1 301 Moved Permanently\r\n", 301, true},
		{"HTTP/1.1 500 Internal Server Error\r\n", 500, true},
		{"garbage\r\n", 0, false},
		{"", 0, false},
	}
	for _, c := range cases {
		code, err := statusCode(c.line)
		if c.ok && (err != nil || code != c.code) {
			t.Errorf("statusCode(%q) = %d,%v; want %d,nil", c.line, code, err, c.code)
		}
		if !c.ok && err == nil {
			t.Errorf("statusCode(%q) = %d,nil; want error", c.line, code)
		}
	}
}

func TestHostPort(t *testing.T) {
	cases := []struct{ raw, def, want string }{
		{"http://host/x", "80", "host:80"},
		{"http://host:8080/x", "80", "host:8080"},
		{"https://host/x", "443", "host:443"},
		{"http://[::1]/x", "80", "[::1]:80"},
		{"http://[::1]:9000/x", "80", "[::1]:9000"},
	}
	for _, c := range cases {
		u, _ := url.Parse(c.raw)
		if got := hostPort(u, c.def); got != c.want {
			t.Errorf("hostPort(%q, %q) = %q; want %q", c.raw, c.def, got, c.want)
		}
	}
}
