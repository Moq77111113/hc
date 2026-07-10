//go:build !hc_slim || hc_http

package probe

import (
	"context"
	"net"
	"net/url"
)

func init() { register("http", httpProber{}) }

// httpProber issues a minimal HTTP/1.1 GET over a plain TCP connection and
// treats any 2xx/3xx status as healthy.
type httpProber struct{}

func (httpProber) Probe(ctx context.Context, target *url.URL) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", hostPort(target, "80"))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return probeHTTP(ctx, conn, target)
}
