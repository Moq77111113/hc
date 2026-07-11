//go:build !hc_slim || hc_grpc

package probe

import "encoding/binary"

const (
	// healthCheckPath is the HTTP/2 path of the standard gRPC health RPC.
	healthCheckPath = "/grpc.health.v1.Health/Check"
	// grpcStatusServing is HealthCheckResponse.ServingStatus SERVING.
	grpcStatusServing = 0x01
	// grpcMessagePrefixLen is the gRPC length-prefix framing each message:
	// one compression flag byte then a 4-byte big-endian length.
	grpcMessagePrefixLen = 5
	// healthStatusField is the protobuf field number of HealthCheckResponse.status.
	healthStatusField = 1
	protoWireVarint   = 0x00
)

// healthCheckRequest is the framed HealthCheckRequest for the default service
// (""): a gRPC message wrapping an empty protobuf body.
var healthCheckRequest = []byte{0x00, 0x00, 0x00, 0x00, 0x00}

// grpcServing reports whether a gRPC DATA payload carries a HealthCheckResponse
// with status SERVING. It strips the message prefix and reads the status field
// without decoding the rest of the message.
func grpcServing(data []byte) bool {
	if len(data) < grpcMessagePrefixLen {
		return false
	}
	size := binary.BigEndian.Uint32(data[1:grpcMessagePrefixLen])
	body := data[grpcMessagePrefixLen:]
	if int(size) > len(body) {
		return false
	}
	body = body[:size]

	if len(body) < 2 {
		return false
	}
	tag := body[0]
	if tag>>3 != healthStatusField || tag&0x07 != protoWireVarint {
		return false
	}
	return body[1] == grpcStatusServing
}
