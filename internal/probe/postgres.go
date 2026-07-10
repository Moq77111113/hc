//go:build !hc_slim || hc_postgres

package probe

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net/url"
)

func init() {
	register(Postgres, postgresProber{})
}

// pgSSLRequestCode is the magic request code of the PostgreSQL SSLRequest packet.
const pgSSLRequestCode = 80877103

// pgSSLSupported and pgSSLUnsupported are the two valid single-byte answers to
// an SSLRequest; either one proves the server is up and past the wire handshake.
const (
	pgSSLSupported   = 'S'
	pgSSLUnsupported = 'N'
)

// postgresProber proves a PostgreSQL server is accepting connections without
// credentials. It sends the fixed SSLRequest packet; a live server answers a
// single byte ('S' or 'N') before authentication, the same handshake pg_isready
// relies on.
type postgresProber struct{}

func (postgresProber) Probe(ctx context.Context, target *url.URL) error {
	var pkt [8]byte
	binary.BigEndian.PutUint32(pkt[0:4], uint32(len(pkt)))
	binary.BigEndian.PutUint32(pkt[4:8], pgSSLRequestCode)

	return handshake(ctx, target, pkt[:], func(r *bufio.Reader) error {
		reply, err := r.ReadByte()
		if err != nil {
			return err
		}
		if reply != pgSSLSupported && reply != pgSSLUnsupported {
			return fmt.Errorf("unexpected SSLRequest reply %q", reply)
		}
		return nil
	})
}
