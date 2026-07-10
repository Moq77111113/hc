//go:build !hc_slim || hc_mysql

package probe

import (
	"context"
	"net"
	"net/url"
	"testing"
)

func serveMySQLGreeting(t *testing.T, greeting []byte) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Write(greeting)
	}()

	return ln.Addr().String()
}

func TestMySQLProberHealthyOnHandshakeV10(t *testing.T) {
	greeting := []byte{0x0a, 0x00, 0x00, 0x00, mysqlProtocolV10}
	u, _ := url.Parse("mysql://" + serveMySQLGreeting(t, greeting))
	if err := (mysqlProber{}).Probe(context.Background(), u); err != nil {
		t.Fatalf("want healthy, got %v", err)
	}
}

func TestMySQLProberUnhealthyOnErrPacket(t *testing.T) {
	// 0xff opens an ERR packet, not a handshake: a server refusing the wire.
	greeting := []byte{0x10, 0x00, 0x00, 0x00, 0xff}
	u, _ := url.Parse("mysql://" + serveMySQLGreeting(t, greeting))
	if err := (mysqlProber{}).Probe(context.Background(), u); err == nil {
		t.Fatal("want error on ERR packet, got nil")
	}
}
