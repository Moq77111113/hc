//go:build !hc_slim || hc_https

package probe

import (
	"context"
	"crypto/tls"
	"net/url"
)

func init() { register(HTTPS, httpsProber{}) }

// defaultHTTPSPort is used when the target URL omits a port.
const defaultHTTPSPort = "443"

// httpsProber probes over TLS. It deliberately skips certificate validation:
// hc checks liveness ("does it answer over TLS?"), not cert validity, so it
// still works against internal/self-signed endpoints. ServerName is set for SNI.
type httpsProber struct{}

func (httpsProber) Probe(ctx context.Context, target *url.URL) error {
	d := tls.Dialer{Config: &tls.Config{
		ServerName:         target.Hostname(),
		InsecureSkipVerify: true, //nolint:gosec // liveness probe, not cert validation (see type doc)
	}}
	conn, err := d.DialContext(ctx, "tcp", hostPort(target, defaultHTTPSPort))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return probeHTTP(ctx, conn, target)
}
