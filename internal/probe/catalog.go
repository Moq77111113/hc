package probe

import "sort"

// Scheme is one protocol hc can probe. Aliases route alternate URL schemes to
// the same prober. Name matches the prober's register() call and its
// <name>.go filename.
type Scheme struct {
	Name    string
	Aliases []string
}

// Bundle is a slim binary shipping a fixed scheme subset.
type Bundle struct {
	Binary  string
	Schemes []Scheme
}

var (
	HTTP     = Scheme{Name: "http"}
	HTTPS    = Scheme{Name: "https"}
	TCP      = Scheme{Name: "tcp"}
	Postgres = Scheme{Name: "postgres", Aliases: []string{"pg"}}
	MySQL    = Scheme{Name: "mysql"}
	Redis    = Scheme{Name: "redis"}
)

// Catalog is every scheme; the full hc binary ships all of them.
var Catalog = []Scheme{HTTP, HTTPS, TCP, Postgres, MySQL, Redis}

// Bundles are the slim binaries. Schemes reference the vars above, so a typo is
// a compile error rather than a silent bad build.
var Bundles = []Bundle{
	{Binary: "hc-core", Schemes: []Scheme{HTTP, HTTPS, TCP}},
	{Binary: "hc-sql", Schemes: []Scheme{TCP, Postgres, MySQL}},
}

// SchemeNames returns every scheme and alias name, sorted — the full set the
// default binary registers.
func SchemeNames() []string {
	var names []string
	for _, s := range Catalog {
		names = append(names, s.Name)
		names = append(names, s.Aliases...)
	}
	sort.Strings(names)
	return names
}
