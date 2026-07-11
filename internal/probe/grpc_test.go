//go:build !hc_slim || hc_grpc

package probe

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/url"
	"testing"
)

// grpcMessage wraps a protobuf body in the 5-byte gRPC length prefix.
func grpcMessage(body []byte) []byte {
	msg := []byte{0x00, 0x00, 0x00, 0x00, byte(len(body))}
	return append(msg, body...)
}

// serveGRPC accepts one connection and writes reply, then drains the client so
// the connection stays open until the client is done reading.
func serveGRPC(t *testing.T, reply func(conn net.Conn)) string {
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
		reply(conn)
		_, _ = io.Copy(io.Discard, conn)
	}()

	return ln.Addr().String()
}

func probeGRPC(t *testing.T, addr string) error {
	t.Helper()
	u, _ := url.Parse("grpc://" + addr)
	return (grpcProber{}).Probe(context.Background(), u)
}

func TestGRPCProberHealthyOnServing(t *testing.T) {
	addr := serveGRPC(t, func(conn net.Conn) {
		_ = writeFrame(conn, frameSettings, 0, 0, nil)
		_ = writeFrame(conn, frameData, flagEndStream, 1, grpcMessage([]byte{0x08, grpcStatusServing}))
	})
	if err := probeGRPC(t, addr); err != nil {
		t.Fatalf("want healthy, got %v", err)
	}
}

func TestGRPCProberUnhealthyOnNotServing(t *testing.T) {
	addr := serveGRPC(t, func(conn net.Conn) {
		_ = writeFrame(conn, frameSettings, 0, 0, nil)
		_ = writeFrame(conn, frameData, flagEndStream, 1, grpcMessage([]byte{0x08, 0x02})) // NOT_SERVING
	})
	if err := probeGRPC(t, addr); err == nil {
		t.Fatal("want error on NOT_SERVING, got nil")
	}
}

func TestGRPCProberUnhealthyOnTrailersOnly(t *testing.T) {
	// An error RPC (e.g. no health service) returns trailers only, no DATA frame.
	addr := serveGRPC(t, func(conn net.Conn) {
		_ = writeFrame(conn, frameSettings, 0, 0, nil)
		_ = writeFrame(conn, frameHeaders, flagEndStream|flagEndHeaders, 1, []byte{0x00})
	})
	if err := probeGRPC(t, addr); err == nil {
		t.Fatal("want error on trailers-only response, got nil")
	}
}

func TestGRPCServingParsesStatus(t *testing.T) {
	if !grpcServing(grpcMessage([]byte{0x08, grpcStatusServing})) {
		t.Error("SERVING message must parse as healthy")
	}
	if grpcServing(grpcMessage([]byte{0x08, 0x02})) {
		t.Error("NOT_SERVING message must not parse as healthy")
	}
	if grpcServing([]byte{0x00, 0x00}) {
		t.Error("truncated message must not parse as healthy")
	}
}

func TestHpackHeadersEncodesRequest(t *testing.T) {
	out := hpackHeaders("localhost:50051", healthCheckPath)

	// :method POST and :scheme http are static-table indexed fields.
	if out[0] != 0x83 || out[1] != 0x86 {
		t.Fatalf("want indexed :method/:scheme prefix 0x83 0x86, got %#x %#x", out[0], out[1])
	}
	for _, want := range [][]byte{
		[]byte(healthCheckPath),
		[]byte("localhost:50051"),
		[]byte("application/grpc"),
		[]byte("trailers"),
	} {
		if !bytes.Contains(out, want) {
			t.Errorf("headers missing %q", want)
		}
	}
}
