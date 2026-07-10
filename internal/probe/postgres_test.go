//go:build !hc_slim || hc_postgres

package probe

import (
	"context"
	"io"
	"net"
	"net/url"
	"testing"
)

func TestPostgresProberHealthyOnSSLReply(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		var req [8]byte
		if _, err := io.ReadFull(conn, req[:]); err != nil {
			return
		}
		conn.Write([]byte{'S'})
	}()

	u, _ := url.Parse("postgres://" + ln.Addr().String())
	if err := (postgresProber{}).Probe(context.Background(), u); err != nil {
		t.Fatalf("want healthy, got %v", err)
	}
}

func TestPostgresProberUnhealthyOnUnexpectedReply(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		var req [8]byte
		if _, err := io.ReadFull(conn, req[:]); err != nil {
			return
		}
		conn.Write([]byte{'X'}) // neither 'S' nor 'N'
	}()

	u, _ := url.Parse("postgres://" + ln.Addr().String())
	if err := (postgresProber{}).Probe(context.Background(), u); err == nil {
		t.Fatal("want error on unexpected SSLRequest reply, got nil")
	}
}
