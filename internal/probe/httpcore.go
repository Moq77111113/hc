//go:build !hc_slim || hc_http || hc_https

package probe

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// maxStatusBytes caps how far we read looking for the final status line, so a
// server that never sends a newline can't grow the buffer without bound.
const maxStatusBytes = 8 << 10

// probeHTTP sends a minimal HTTP/1.1 GET and treats a 2xx/3xx final status as
// healthy. It reads only the status line and never touches net/http; that
// omission is the ~1.8 MB saving. The caller must give ctx a deadline.
func probeHTTP(ctx context.Context, conn net.Conn, target *url.URL) error {
	if dl, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(dl); err != nil {
			return err
		}
	}

	req := fmt.Sprintf(
		"GET %s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n",
		target.RequestURI(), target.Host,
	)
	if _, err := conn.Write([]byte(req)); err != nil {
		return err
	}

	reader := bufio.NewReader(io.LimitReader(conn, maxStatusBytes))
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if !strings.HasPrefix(line, "HTTP/") {
			continue // header line of a 1xx interim response
		}

		code, err := statusCode(line)
		if err != nil {
			return err
		}
		if code < 200 {
			continue // 1xx interim, the real status is still coming
		}
		if code >= 400 {
			return fmt.Errorf("status %d", code)
		}
		return nil
	}
}

// statusCode extracts the numeric code from a status line like "HTTP/1.1 200 OK".
func statusCode(line string) (int, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, fmt.Errorf("malformed status line %q", strings.TrimSpace(line))
	}
	return strconv.Atoi(fields[1])
}

// hostPort returns target's host:port, applying defaultPort when none is set.
func hostPort(target *url.URL, defaultPort string) string {
	if p := target.Port(); p != "" {
		return net.JoinHostPort(target.Hostname(), p)
	}
	return net.JoinHostPort(target.Hostname(), defaultPort)
}
