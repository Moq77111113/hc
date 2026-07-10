package main

import (
	"fmt"
	"os"
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
	data, err := os.ReadFile(path) //nolint:gosec // caller-controlled generator path, not user input
	if err != nil {
		return err
	}
	out, err := replaceBlock(string(data), ciMarkerBegin, ciMarkerEnd, renderMatrix(probe.Catalog, Bundles))
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(out), 0o644) //nolint:gosec // generated workflow, committed to git, world-readable is fine
}
