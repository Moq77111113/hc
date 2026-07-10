//go:build !hc_slim || hc_mysql

package probe

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
)

func init() { register(MySQL, mysqlProber{}) }

const (
	// mysqlPacketHeaderLen is the 4-byte length+sequence header prefixing every
	// MySQL protocol packet; the greeting payload starts right after it.
	mysqlPacketHeaderLen = 4
	// mysqlProtocolV10 is the version byte that opens a HandshakeV10 greeting,
	// the first packet a live MySQL or MariaDB server sends on connect.
	mysqlProtocolV10 = 0x0a
)

// mysqlProber proves a MySQL server is up by reading the greeting it sends on
// connect. The server speaks first: hc reads the packet header and the leading
// protocol-version byte and treats HandshakeV10 as healthy without
// authenticating — liveness, not login.
type mysqlProber struct{}

func (mysqlProber) Probe(ctx context.Context, target *url.URL) error {
	return handshake(ctx, target, nil, func(r *bufio.Reader) error {
		var header [mysqlPacketHeaderLen]byte
		if _, err := io.ReadFull(r, header[:]); err != nil {
			return err
		}
		version, err := r.ReadByte()
		if err != nil {
			return err
		}
		if version != mysqlProtocolV10 {
			return fmt.Errorf("unexpected handshake protocol version %#x", version)
		}
		return nil
	})
}
