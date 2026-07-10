//go:build hc_slim && hc_mysql && !hc_http && !hc_https && !hc_tcp && !hc_postgres && !hc_redis

package probe

import "testing"

func TestSlimMySQLBuildRegistersOnlyMySQL(t *testing.T) {
	if _, ok := probers["mysql"]; !ok {
		t.Error("mysql must be registered in a slim mysql build")
	}
	for _, s := range []string{"http", "https", "tcp", "postgres", "pg", "redis"} {
		if _, ok := probers[s]; ok {
			t.Errorf("scheme %q must NOT be registered in a slim mysql build", s)
		}
	}
}
