//go:build !hc_slim || hc_amqp

package probe

import (
	"context"
	"net"
	"net/url"
	"testing"
)

func serveAMQPReply(t *testing.T, reply []byte) string {
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
		_, _ = conn.Write(reply)
	}()

	return ln.Addr().String()
}

func TestAMQPProberHealthyOnMethodFrame(t *testing.T) {
	// A live broker opens negotiation with Connection.Start in a METHOD frame.
	reply := []byte{amqpFrameMethod, 0x00, 0x00}
	u, _ := url.Parse("amqp://" + serveAMQPReply(t, reply))
	if err := (amqpProber{}).Probe(context.Background(), u); err != nil {
		t.Fatalf("want healthy, got %v", err)
	}
}

func TestAMQPProberUnhealthyOnProtocolReject(t *testing.T) {
	// A version mismatch makes the broker echo a protocol header ("AMQP…") and
	// close, instead of answering with a METHOD frame.
	reply := []byte("AMQP\x00\x00\x09\x01")
	u, _ := url.Parse("amqp://" + serveAMQPReply(t, reply))
	if err := (amqpProber{}).Probe(context.Background(), u); err == nil {
		t.Fatal("want error on protocol-header reply, got nil")
	}
}
