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

func TestRegisterBindsSchemeAndAliases(t *testing.T) {
	pg, okPg := Get("postgres")
	alias, okAlias := Get("pg")
	if !okPg || !okAlias {
		t.Fatal("postgres and its alias pg must both be registered")
	}
	if pg != alias {
		t.Error("postgres and pg must resolve to the same prober")
	}
}
