//go:build !hc_slim || hc_redis

package probe

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

func init() {
	register("redis", redisProber{})
}

// maxRedisReplyBytes bounds the PING reply read against a server that never
// sends a newline.
const maxRedisReplyBytes = 1 << 10

// redisProber sends an inline PING and treats "+PONG" as healthy. A
// password-protected server answers "-NOAUTH", which still proves it is alive,
// so hc reports that as healthy too: liveness, not authorization.
type redisProber struct{}

func (redisProber) Probe(ctx context.Context, target *url.URL) error {
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

	if _, err := conn.Write([]byte("PING\r\n")); err != nil {
		return err
	}

	line, err := bufio.NewReader(io.LimitReader(conn, maxRedisReplyBytes)).ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.HasPrefix(line, "+PONG") && !strings.HasPrefix(line, "-NOAUTH") {
		return fmt.Errorf("unexpected PING reply %q", strings.TrimSpace(line))
	}
	return nil
}
