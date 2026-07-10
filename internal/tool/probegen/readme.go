package main

import (
	"fmt"
	"strings"
)

const (
	readmeMarkerBegin = "<!-- probes:begin -->"
	readmeMarkerEnd   = "<!-- probes:end -->"
)

// renderTable builds the bundle reference table. Scheme names use each Scheme's
// primary Name in declared order (aliases stay implementation detail).
func renderTable(bundles []Bundle) string {
	var b strings.Builder
	b.WriteString("| Binary | Schemes |\n")
	b.WriteString("|---|---|\n")
	for _, bd := range bundles {
		names := make([]string, len(bd.Schemes))
		for i, s := range bd.Schemes {
			names[i] = s.Name
		}
		fmt.Fprintf(&b, "| `%s` | %s |\n", bd.Binary, strings.Join(names, ", "))
	}
	return strings.TrimRight(b.String(), "\n")
}

// writeReadme rewrites the marked table block at path.
func writeReadme(path string) error {
	return replaceMarkedBlock(path, readmeMarkerBegin, readmeMarkerEnd, renderTable(Bundles))
}
