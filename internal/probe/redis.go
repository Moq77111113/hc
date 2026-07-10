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

// redisProber sends an inline PING and treats "+PONG" as healthy. A
// password-protected server answers "-NOAUTH", which still proves it is alive,
// so hc reports that as healthy too: liveness, not authorization.
type redisProber struct{}

func (redisProber) Probe(ctx context.Context, target *url.URL) error {
	return handshake(ctx, target, []byte("PING\r\n"), func(r *bufio.Reader) error {
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		if !strings.HasPrefix(line, "+PONG") && !strings.HasPrefix(line, "-NOAUTH") {
			return fmt.Errorf("unexpected PING reply %q", strings.TrimSpace(line))
		}
		return nil
	})
}
