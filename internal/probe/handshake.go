//go:build !hc_slim || hc_postgres || hc_redis

package probe

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/url"
)

// maxHandshakeReply bounds the reply read against a peer that never stops sending.
const maxHandshakeReply = 1 << 10

// dial connects to target.Host and arms the connection with ctx's deadline.
func dial(ctx context.Context, target *url.URL) (net.Conn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", target.Host)
	if err != nil {
		return nil, err
	}
	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}
	return conn, nil
}

// handshake dials target, sends payload, and lets accept judge the reply: the
// shared shape of the byte-level probers.
func handshake(ctx context.Context, target *url.URL, payload []byte, accept func(*bufio.Reader) error) error {
	conn, err := dial(ctx, target)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	if _, err := conn.Write(payload); err != nil {
		return err
	}
	return accept(bufio.NewReader(io.LimitReader(conn, maxHandshakeReply)))
}
