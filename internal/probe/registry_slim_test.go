//go:build hc_slim && hc_tcp && !hc_http && !hc_https && !hc_postgres

package probe

import "testing"

func TestSlimTCPBuildRegistersOnlyTCP(t *testing.T) {
	if _, ok := probers["tcp"]; !ok {
		t.Error("tcp must be registered in a slim tcp build")
	}
	for _, s := range []string{"http", "https", "postgres", "pg"} {
		if _, ok := probers[s]; ok {
			t.Errorf("scheme %q must NOT be registered in a slim tcp build", s)
		}
	}
}
