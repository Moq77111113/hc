// Package probe selects and runs the readiness Prober for a target URL scheme.
// Each prober registers itself from init(); build tags decide which probers
// compile in, so a slim build carries only the schemes it was built with.
package probe

import (
	"context"
	"net/url"
	"sort"
	"strings"
)

// Prober checks the readiness of a single target and returns nil when healthy.
type Prober interface {
	Probe(ctx context.Context, target *url.URL) error
}

// probers maps a URL scheme to the Prober that speaks it. Each prober file
// registers itself from init(); build tags select which files compile, so a
// slim build carries only the probers it was built with.
var probers = map[string]Prober{}

// register binds a scheme to its Prober. Called from each prober's init().
func register(scheme string, p Prober) {
	probers[scheme] = p
}

// Get returns the Prober registered for scheme, if any.
func Get(scheme string) (Prober, bool) {
	p, ok := probers[scheme]
	return p, ok
}

// SupportedSchemes lists the registered schemes, sorted, for help and errors.
func SupportedSchemes() string {
	schemes := make([]string, 0, len(probers))
	for s := range probers {
		schemes = append(schemes, s)
	}
	sort.Strings(schemes)
	return strings.Join(schemes, ", ")
}
