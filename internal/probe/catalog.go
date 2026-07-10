package probe

import "sort"

// Scheme is one protocol hc can probe. Aliases route alternate URL schemes to
// the same prober. Name matches the prober's register() call and its
// <name>.go filename.
type Scheme struct {
	Name    string
	Aliases []string
}

var (
	// HTTP probes http:// targets with a request and checks the status code.
	HTTP = Scheme{Name: "http"}
	// HTTPS probes https:// targets over TLS.
	HTTPS = Scheme{Name: "https"}
	// TCP probes tcp:// targets by establishing a connection.
	TCP = Scheme{Name: "tcp"}
	// Postgres probes postgres:// (alias pg://) targets via the startup handshake.
	Postgres = Scheme{Name: "postgres", Aliases: []string{"pg"}}
	// MySQL probes mysql:// targets via the server handshake packet.
	MySQL = Scheme{Name: "mysql"}
	// Redis probes redis:// targets with an inline PING.
	Redis = Scheme{Name: "redis"}
)

// Catalog is every scheme; the full hc binary ships all of them.
var Catalog = []Scheme{HTTP, HTTPS, TCP, Postgres, MySQL, Redis}

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
