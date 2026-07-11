//go:build !hc_slim || hc_amqp

package probe

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
)

func init() { register(AMQP, amqpProber{}) }

const (
	// amqpFrameMethod is the frame type of an AMQP 0-9-1 METHOD frame; a live
	// broker opens negotiation with Connection.Start in one.
	amqpFrameMethod = 0x01
)

// amqpProtocolHeader is the AMQP 0-9-1 protocol header the client sends first.
// A live broker answers with a METHOD frame; a version mismatch echoes its own
// header back instead.
var amqpProtocolHeader = []byte{'A', 'M', 'Q', 'P', 0x00, 0x00, 0x09, 0x01}

// amqpProber proves an AMQP broker is up by sending the 0-9-1 protocol header
// and checking the reply opens a METHOD frame. It does not authenticate or open
// a channel: liveness, not login.
type amqpProber struct{}

func (amqpProber) Probe(ctx context.Context, target *url.URL) error {
	return handshake(ctx, target, amqpProtocolHeader, func(r *bufio.Reader) error {
		frameType, err := r.ReadByte()
		if err != nil {
			return err
		}
		if frameType != amqpFrameMethod {
			return fmt.Errorf("unexpected AMQP frame type %#x", frameType)
		}
		return nil
	})
}
