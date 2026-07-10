//go:build !hc_slim || hc_example

package probe

import (
	"context"
	"net/url"
)

// To add a probe:
//   1. in catalog.go: declare `Example = Scheme{Name: "example"}`, add it to
//      Catalog, and to a Bundle in probegen/bundles.go if it ships slim.
//   2. cp internal/probe/_template.go internal/probe/example.go
//   3. replace every "example"/"Example" with <scheme> (build tag, register, type)
//   4. go generate ./...
func init() { register(Example, exampleProber{}) }

type exampleProber struct{}

// Probe returns nil when the target is healthy.
func (exampleProber) Probe(ctx context.Context, target *url.URL) error {
	return nil
}
