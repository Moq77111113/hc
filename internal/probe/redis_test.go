//go:build !hc_slim || hc_redis

package probe

import (
	"bufio"
	"context"
	"net"
	"net/url"
	"testing"
)

func serveRedisReply(t *testing.T, reply string) string {
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
		if _, err := bufio.NewReader(conn).ReadString('\n'); err != nil {
			return
		}
		_, _ = conn.Write([]byte(reply))
	}()

	return ln.Addr().String()
}

func TestRedisProberHealthyOnPong(t *testing.T) {
	u, _ := url.Parse("redis://" + serveRedisReply(t, "+PONG\r\n"))
	if err := (redisProber{}).Probe(context.Background(), u); err != nil {
		t.Fatalf("want healthy, got %v", err)
	}
}

func TestRedisProberHealthyOnNoAuth(t *testing.T) {
	u, _ := url.Parse("redis://" + serveRedisReply(t, "-NOAUTH Authentication required.\r\n"))
	if err := (redisProber{}).Probe(context.Background(), u); err != nil {
		t.Fatalf("want healthy on NOAUTH, got %v", err)
	}
}

func TestRedisProberUnhealthyOnGarbage(t *testing.T) {
	u, _ := url.Parse("redis://" + serveRedisReply(t, "not resp\r\n"))
	if err := (redisProber{}).Probe(context.Background(), u); err == nil {
		t.Fatal("want error on non-RESP reply, got nil")
	}
}
