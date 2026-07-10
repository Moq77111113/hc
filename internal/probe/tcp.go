//go:build !hc_slim || hc_tcp

package probe

import (
	"context"
	"net"
	"net/url"
)

func init() {
	register("tcp", tcpProber{})
}

// tcpProber treats a successful TCP connection as healthy.
type tcpProber struct{}

func (tcpProber) Probe(ctx context.Context, target *url.URL) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", target.Host)
	if err != nil {
		return err
	}
	return conn.Close()
}
