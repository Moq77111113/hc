//go:build !hc_slim

package probe

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// Every registered scheme must be catalogued and vice versa.
func TestCatalogMatchesRegistrations(t *testing.T) {
	want := map[string]bool{}
	for _, n := range SchemeNames() {
		want[n] = true
	}
	for s := range probers {
		if !want[s] {
			t.Errorf("scheme %q registered but missing from Catalog", s)
		}
	}
	for n := range want {
		if _, ok := probers[n]; !ok {
			t.Errorf("scheme %q in Catalog but not registered", n)
		}
	}
}

// Each prober file must carry the canonical build tag, so a slim build selects
// it correctly. Working dir for a package test is the package dir.
func TestProberFilesHaveCanonicalBuildTag(t *testing.T) {
	for _, s := range Catalog {
		path := fmt.Sprintf("%s.go", s.Name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("scheme %q: %v", s.Name, err)
			continue
		}
		want := fmt.Sprintf("//go:build !hc_slim || hc_%s", s.Name)
		first := strings.SplitN(string(data), "\n", 2)[0]
		if first != want {
			t.Errorf("%s first line = %q, want %q", path, first, want)
		}
	}
}
