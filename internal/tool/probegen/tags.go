package main

import (
	"strings"

	"github.com/Moq77111113/hc/internal/probe"
)

// schemeTag is the build tag that compiles one scheme into a slim binary.
func schemeTag(s probe.Scheme) string { return "hc_" + s.Name }

// slimTags is the space-joined -tags string selecting exactly these schemes.
func slimTags(schemes []probe.Scheme) string {
	tags := make([]string, 0, len(schemes)+1)
	tags = append(tags, "hc_slim")
	for _, s := range schemes {
		tags = append(tags, schemeTag(s))
	}
	return strings.Join(tags, " ")
}
