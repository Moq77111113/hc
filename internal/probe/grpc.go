//go:build !hc_slim || hc_grpc

package probe

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
)

func init() { register(GRPC, grpcProber{}) }

// maxResponseFrames bounds the read loop against a peer that keeps sending
// control frames without ever answering the health call.
const maxResponseFrames = 32

// grpcProber proves a gRPC server is healthy by calling the standard
// grpc.health.v1.Health/Check RPC over cleartext HTTP/2 and requiring a SERVING
// response. It hand-rolls the HTTP/2 exchange rather than depend on a gRPC stack.
type grpcProber struct{}

func (grpcProber) Probe(ctx context.Context, target *url.URL) error {
	conn, err := dial(ctx, target)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	if err := writeHealthCall(conn, target.Host); err != nil {
		return err
	}

	r := bufio.NewReader(conn)
	for range maxResponseFrames {
		header, err := readFrameHeader(r)
		if err != nil {
			return err
		}
		if header.length > maxFramePayload {
			return fmt.Errorf("oversized HTTP/2 frame: %d bytes", header.length)
		}

		switch header.kind {
		case frameData:
			payload := make([]byte, header.length)
			if _, err := io.ReadFull(r, payload); err != nil {
				return err
			}
			if !grpcServing(payload) {
				return fmt.Errorf("gRPC health status not SERVING")
			}
			return nil
		case frameHeaders:
			if err := skip(r, header.length); err != nil {
				return err
			}
			if header.flags&flagEndStream != 0 {
				return fmt.Errorf("gRPC health returned no SERVING message")
			}
		case frameSettings:
			if header.flags&flagAck == 0 {
				if err := writeFrame(conn, frameSettings, flagAck, 0, nil); err != nil {
					return err
				}
			}
			if err := skip(r, header.length); err != nil {
				return err
			}
		case frameGoAway, frameRST:
			return fmt.Errorf("gRPC connection refused (frame %#x)", header.kind)
		default:
			if err := skip(r, header.length); err != nil {
				return err
			}
		}
	}
	return fmt.Errorf("no gRPC health response")
}

// writeHealthCall sends the client preface, the request HEADERS, and the DATA
// frame carrying the HealthCheckRequest on stream 1, in one write.
func writeHealthCall(w io.Writer, authority string) error {
	buf := []byte(http2Preface)
	buf = appendFrame(buf, frameSettings, 0, 0, nil)
	buf = appendFrame(buf, frameHeaders, flagEndHeaders, 1, hpackHeaders(authority, healthCheckPath))
	buf = appendFrame(buf, frameData, flagEndStream, 1, healthCheckRequest)
	_, err := w.Write(buf)
	return err
}

// skip discards n payload bytes of a frame the probe does not read.
func skip(r io.Reader, n uint32) error {
	_, err := io.CopyN(io.Discard, r, int64(n))
	return err
}
