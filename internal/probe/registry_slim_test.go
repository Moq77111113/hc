//go:build hc_slim

package probe

import "testing"

// A slim build compiles in only the probers whose build tag was set, so the
// registered set must be a non-empty subset of the catalog. Per-scheme
// exclusion is proven at compile time by the canonical build tags (guarded by
// TestProberFilesHaveCanonicalBuildTag)
func TestSlimBuildRegistersCatalogSubset(t *testing.T) {
	catalogued := make(map[string]bool)
	for _, name := range SchemeNames() {
		catalogued[name] = true
	}

	if len(probers) == 0 {
		t.Fatal("slim build registered no probers")
	}

	for scheme := range probers {
		if !catalogued[scheme] {
			t.Errorf("registered scheme %q is not in the catalog", scheme)
		}
	}
}
