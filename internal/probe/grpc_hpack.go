//go:build !hc_slim || hc_grpc

package probe

// HPACK static-table indices (RFC 7541 Appendix A) for the request pseudo- and
// regular headers the gRPC health call needs.
const (
	hpackAuthority   = 1
	hpackMethodPost  = 3
	hpackPath        = 4
	hpackSchemeHTTP  = 6
	hpackContentType = 31
)

// hpackIndexed is an indexed header field: the full name+value pair sits in the
// static table, so a single byte references it (RFC 7541 §6.1).
func hpackIndexed(index byte) byte { return 0x80 | index }

// hpackInt appends value as an HPACK integer with an N-bit prefix (RFC 7541
// §5.1). pattern carries the representation's high bits above the prefix.
func hpackInt(dst []byte, value, prefixBits int, pattern byte) []byte {
	limit := (1 << prefixBits) - 1
	if value < limit {
		return append(dst, pattern|byte(value&0xff))
	}
	dst = append(dst, pattern|byte(limit&0xff))
	value -= limit
	for value >= 128 {
		dst = append(dst, byte((value%128+128)&0xff))
		value /= 128
	}
	return append(dst, byte(value&0xff))
}

// hpackString appends a length-prefixed raw (non-Huffman) string literal.
func hpackString(dst []byte, s string) []byte {
	dst = hpackInt(dst, len(s), 7, 0x00)
	return append(dst, s...)
}

// hpackLiteral appends a literal-without-indexing header (RFC 7541 §6.2.2) whose
// name is the static-table entry nameIndex and whose value is the given string.
func hpackLiteral(dst []byte, nameIndex int, value string) []byte {
	dst = hpackInt(dst, nameIndex, 4, 0x00)
	return hpackString(dst, value)
}

// hpackHeaders encodes the request header block for a gRPC unary call to path on
// authority: the four pseudo-headers plus content-type and te.
func hpackHeaders(authority, path string) []byte {
	var b []byte
	b = append(b, hpackIndexed(hpackMethodPost), hpackIndexed(hpackSchemeHTTP))
	b = hpackLiteral(b, hpackPath, path)
	b = hpackLiteral(b, hpackAuthority, authority)
	b = hpackLiteral(b, hpackContentType, "application/grpc")
	b = hpackInt(b, 0, 4, 0x00) // te: new name (no static index)
	b = hpackString(b, "te")
	b = hpackString(b, "trailers")
	return b
}
