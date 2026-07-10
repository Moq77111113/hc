//go:build hc_slim && hc_redis && !hc_http && !hc_https && !hc_tcp && !hc_postgres

package probe

import "testing"

func TestSlimRedisBuildRegistersOnlyRedis(t *testing.T) {
	if _, ok := probers["redis"]; !ok {
		t.Error("redis must be registered in a slim redis build")
	}
	for _, s := range []string{"http", "https", "tcp", "postgres", "pg"} {
		if _, ok := probers[s]; ok {
			t.Errorf("scheme %q must NOT be registered in a slim redis build", s)
		}
	}
}
