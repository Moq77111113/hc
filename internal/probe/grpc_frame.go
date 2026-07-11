//go:build !hc_slim || hc_grpc

package probe

import (
	"encoding/binary"
	"io"
)

// http2Preface is the client connection preface sent first under prior-knowledge
// h2c, the cleartext HTTP/2 transport gRPC speaks.
const http2Preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

const (
	// frameHeaderLen is the fixed 9-byte header prefixing every HTTP/2 frame.
	frameHeaderLen = 9
	// maxFramePayload bounds a single frame read against a peer that never stops.
	maxFramePayload = 1 << 20
)

// HTTP/2 frame types used by the gRPC health probe.
const (
	frameData     = 0x00
	frameHeaders  = 0x01
	frameRST      = 0x03
	frameSettings = 0x04
	frameGoAway   = 0x07
)

// HTTP/2 frame flags used by the gRPC health probe.
const (
	flagEndStream  = 0x01
	flagAck        = 0x01
	flagEndHeaders = 0x04
)

// frameHeader is the parsed 9-byte prefix of an HTTP/2 frame.
type frameHeader struct {
	length   uint32
	kind     byte
	flags    byte
	streamID uint32
}

// writeFrame writes one HTTP/2 frame: the 9-byte header then payload.
func writeFrame(w io.Writer, kind, flags byte, streamID uint32, payload []byte) error {
	var header [frameHeaderLen]byte
	header[0] = byte(len(payload) >> 16)
	header[1] = byte(len(payload) >> 8)
	header[2] = byte(len(payload))
	header[3] = kind
	header[4] = flags
	binary.BigEndian.PutUint32(header[5:], streamID)
	if _, err := w.Write(header[:]); err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	_, err := w.Write(payload)
	return err
}

// readFrameHeader reads and parses the 9-byte header of the next frame.
func readFrameHeader(r io.Reader) (frameHeader, error) {
	var buf [frameHeaderLen]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return frameHeader{}, err
	}
	return frameHeader{
		length:   uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2]),
		kind:     buf[3],
		flags:    buf[4],
		streamID: binary.BigEndian.Uint32(buf[5:]),
	}, nil
}
