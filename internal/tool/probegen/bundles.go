package main

import "github.com/Moq77111113/hc/internal/probe"

// Bundle is a slim binary shipping a fixed scheme subset. It lives here, not in
// the probe package, because only this generator consumes it.
type Bundle struct {
	Binary  string
	Schemes []probe.Scheme
}

// Bundles are the slim binaries. Schemes reference the probe scheme vars, so a
// typo is a compile error rather than a silent bad build.
var Bundles = []Bundle{
	{Binary: "hc-core", Schemes: []probe.Scheme{probe.HTTP, probe.HTTPS, probe.TCP}},
	{Binary: "hc-sql", Schemes: []probe.Scheme{probe.TCP, probe.Postgres, probe.MySQL}},
}
