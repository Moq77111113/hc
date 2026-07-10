//go:build !hc_slim || hc_redis

package probe

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"strings"
)

func init() { register("redis", redisProber{}) }

const (
	// redisPing is the inline PING command; redisPong is its healthy reply.
	redisPing = "PING\r\n"
	redisPong = "+PONG"
	// redisNoAuth is the error a password-protected server returns to an
	// unauthenticated PING; it still proves the server is alive.
	redisNoAuth = "-NOAUTH"
)

// redisProber sends an inline PING and treats "+PONG" as healthy. A
// password-protected server answers "-NOAUTH", which still proves it is alive,
// so hc reports that as healthy too: liveness, not authorization.
type redisProber struct{}

func (redisProber) Probe(ctx context.Context, target *url.URL) error {
	return handshake(ctx, target, []byte(redisPing), func(r *bufio.Reader) error {
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		if !strings.HasPrefix(line, redisPong) && !strings.HasPrefix(line, redisNoAuth) {
			return fmt.Errorf("unexpected PING reply %q", strings.TrimSpace(line))
		}
		return nil
	})
}
