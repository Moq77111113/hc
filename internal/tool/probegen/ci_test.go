package main

import (
	"strings"
	"testing"

	"github.com/Moq77111113/hc/internal/probe"
)

func TestRenderMatrix(t *testing.T) {
	out := renderMatrix(probe.Catalog, Bundles)

	checks := []string{
		`          - "hc_slim hc_http"`,                     // per-scheme
		`          - "hc_slim hc_redis"`,                    // the previously missing one
		`          - "hc_slim hc_mysql"`,                    // the previously missing one
		`          - "hc_slim hc_http hc_https hc_tcp"`,     // per-bundle hc-core
		`          - "hc_slim hc_tcp hc_postgres hc_mysql"`, // per-bundle hc-sql
	}
	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Errorf("matrix missing %q", c)
		}
	}
}
