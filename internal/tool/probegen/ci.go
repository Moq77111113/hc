package main

import (
	"fmt"
	"strings"

	"github.com/Moq77111113/hc/internal/probe"
)

const (
	ciMarkerBegin = "          # probes:begin"
	ciMarkerEnd   = "          # probes:end"
)

// renderMatrix builds the slim job's tag list: each scheme alone (isolation)
// then each bundle's full tag string (the shipped artifacts).
func renderMatrix(catalog []probe.Scheme, bundles []Bundle) string {
	var lines []string
	for _, s := range catalog {
		lines = append(lines, fmt.Sprintf(`          - "hc_slim %s"`, schemeTag(s)))
	}
	for _, b := range bundles {
		lines = append(lines, fmt.Sprintf(`          - "%s"`, slimTags(b.Schemes)))
	}
	return strings.Join(lines, "\n")
}

// writeCIMatrix rewrites the marked block in the workflow at path.
func writeCIMatrix(path string) error {
	return replaceMarkedBlock(path, ciMarkerBegin, ciMarkerEnd, renderMatrix(probe.Catalog, Bundles))
}
