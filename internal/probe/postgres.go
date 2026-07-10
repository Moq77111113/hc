//go:build !hc_slim || hc_postgres

package probe

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/url"
)

func init() {
	register("postgres", postgresProber{})
	register("pg", postgresProber{})
}

// pgSSLRequestCode is the magic request code of the PostgreSQL SSLRequest packet.
const pgSSLRequestCode = 80877103

// postgresProber proves a PostgreSQL server is accepting connections without
// credentials. It sends the fixed 8-byte SSLRequest packet; a live server
// always answers with a single byte ('S' or 'N') before authentication, so a
// valid reply means "ready", the same handshake pg_isready relies on.
type postgresProber struct{}

func (postgresProber) Probe(ctx context.Context, target *url.URL) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", target.Host)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			return err
		}
	}

	var pkt [8]byte
	binary.BigEndian.PutUint32(pkt[0:4], uint32(len(pkt)))
	binary.BigEndian.PutUint32(pkt[4:8], pgSSLRequestCode)
	if _, err := conn.Write(pkt[:]); err != nil {
		return err
	}

	var reply [1]byte
	if _, err := io.ReadFull(conn, reply[:]); err != nil {
		return err
	}
	if reply[0] != 'S' && reply[0] != 'N' {
		return fmt.Errorf("unexpected SSLRequest reply %q", reply[0])
	}
	return nil
}
