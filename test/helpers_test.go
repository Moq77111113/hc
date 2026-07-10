package test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// runHC runs the built binary and returns its exit code and combined output.
func runHC(t *testing.T, args ...string) (int, string) {
	t.Helper()
	cmd := exec.Command(hcBinary, args...)
	out, err := cmd.CombinedOutput()
	if cmd.ProcessState == nil {
		t.Fatalf("hc did not run: %v\n%s", err, out)
	}
	return cmd.ProcessState.ExitCode(), string(out)
}

// assertOpenClosed checks scheme is healthy against addr and unhealthy against
// a closed port on the same host: the shared shape of the connect-based probes.
func assertOpenClosed(t *testing.T, scheme, addr string) {
	t.Helper()
	if code, out := runHC(t, scheme+"://"+addr); code != 0 {
		t.Errorf("open: exit %d, want 0\n%s", code, out)
	}
	if code, out := runHC(t, scheme+"://"+hostOf(t, addr)+":1"); code != 1 {
		t.Errorf("closed: exit %d, want 1\n%s", code, out)
	}
}

// hostOf strips the port from a "host:port" endpoint.
func hostOf(t *testing.T, addr string) string {
	t.Helper()
	h, _, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("split %q: %v", addr, err)
	}
	return h
}

// selfSignedCert writes cert.pem and key.pem into dir, valid for localhost and 127.0.0.1.
func selfSignedCert(t *testing.T, dir string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("generate serial: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	writePEM(t, filepath.Join(dir, "cert.pem"), "CERTIFICATE", der)

	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	writePEM(t, filepath.Join(dir, "key.pem"), "EC PRIVATE KEY", keyBytes)
}

func writePEM(t *testing.T, path, blockType string, bytes []byte) {
	t.Helper()
	f, err := os.Create(path) //nolint:gosec // test fixture under t.TempDir()
	if err != nil {
		t.Fatalf("create %s: %v", path, err)
	}
	defer func() { _ = f.Close() }()

	if err := pem.Encode(f, &pem.Block{Type: blockType, Bytes: bytes}); err != nil {
		t.Fatalf("encode %s: %v", path, err)
	}
}
