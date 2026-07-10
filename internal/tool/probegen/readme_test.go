package main

import (
	"strings"
	"testing"
)

func TestRenderTable(t *testing.T) {
	out := renderTable(Bundles)

	checks := []string{
		"| Binary | Schemes |",
		"| `hc-core` | http, https, tcp |",
		"| `hc-sql` | tcp, postgres, mysql |",
	}
	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Errorf("table missing %q\n---\n%s", c, out)
		}
	}
}
