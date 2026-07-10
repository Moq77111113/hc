//go:build !hc_slim

package probe

import "testing"

func TestDefaultBuildRegistersAllSchemes(t *testing.T) {
	const want = "http, https, mysql, pg, postgres, redis, tcp"
	if got := SupportedSchemes(); got != want {
		t.Fatalf("SupportedSchemes() = %q, want %q", got, want)
	}
	for _, s := range []string{"http", "https", "tcp", "postgres", "pg", "redis", "mysql"} {
		if _, ok := probers[s]; !ok {
			t.Errorf("scheme %q not registered", s)
		}
	}
}
