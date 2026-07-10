package main

import (
	"strings"
	"testing"

	"github.com/Moq77111113/hc/internal/probe"
)

func TestRenderMatrix(t *testing.T) {
	out := renderMatrix(probe.Catalog, Bundles)

	checks := []string{
		`          - "hc_slim hc_http"`,
		`          - "hc_slim hc_redis"`,
		`          - "hc_slim hc_mysql"`,
		`          - "hc_slim hc_http hc_https hc_tcp"`,
		`          - "hc_slim hc_tcp hc_postgres hc_mysql"`,
	}
	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Errorf("matrix missing %q", c)
		}
	}
}
