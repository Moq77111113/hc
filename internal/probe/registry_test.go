//go:build !hc_slim

package probe

import (
	"strings"
	"testing"
)

func TestDefaultBuildRegistersCatalog(t *testing.T) {
	names := SchemeNames()
	want := strings.Join(names, ", ")
	if got := SupportedSchemes(); got != want {
		t.Fatalf("SupportedSchemes() = %q, want %q", got, want)
	}
	for _, s := range names {
		if _, ok := probers[s]; !ok {
			t.Errorf("scheme %q not registered", s)
		}
	}
}
