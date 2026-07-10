//go:build !hc_slim || hc_tcp

package probe

import (
	"context"
	"net"
	"net/url"
	"testing"
)

func TestTCPProberHealthyWhenListening(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	u, _ := url.Parse("tcp://" + ln.Addr().String())
	if err := (tcpProber{}).Probe(context.Background(), u); err != nil {
		t.Fatalf("want healthy, got %v", err)
	}
}

func TestTCPProberUnhealthyWhenClosed(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close() // free the port so the dial has nothing to reach

	u, _ := url.Parse("tcp://" + addr)
	if err := (tcpProber{}).Probe(context.Background(), u); err == nil {
		t.Fatal("want error when port closed, got nil")
	}
}
