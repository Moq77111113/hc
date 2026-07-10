// Command hc is a dependency-free health probe for minimal container images.
//
// distroless and scratch images ship no shell and no curl, so the usual
// HEALTHCHECK CMD curl ... cannot run. hc is a single static binary you copy
// into any image and point at a target; its exit code follows Docker's
// health contract (0 healthy, 1 unhealthy).
//
//	COPY --from=ghcr.io/moq77111113/hc /hc /hc
//	HEALTHCHECK CMD ["/hc", "http://localhost:8080/health"]
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/Moq77111113/hc/internal/probe"
)

const (
	exitHealthy   = 0
	exitUnhealthy = 1
)

// version is the release identifier, injected at build time via
// -ldflags "-X main.version=<tag>". Unset builds report "dev".
var version = "dev"

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "install" {
		runInstall(os.Args[2:])
		return
	}
	if len(os.Args) >= 2 && os.Args[1] == "version" {
		os.Exit(runVersion(os.Stdout))
	}
	os.Exit(runProbe())
}

// runVersion prints the build version and exits healthy, so `hc version`
// composes with the same exit-code contract as every other invocation.
func runVersion(out io.Writer) int {
	if _, err := fmt.Fprintf(out, "hc %s\n", version); err != nil {
		return exitUnhealthy
	}
	return exitHealthy
}

// runInstall handles the `hc install DEST` subcommand: it self-copies the
// running binary to DEST and exits with the probe-mode exit codes so it
// composes with the same healthy/unhealthy contract callers already expect.
func runInstall(args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: hc install DEST\n")
		os.Exit(exitUnhealthy)
	}
	if err := install(args[0]); err != nil {
		os.Exit(fail("install: %v", err))
	}
	os.Exit(exitHealthy)
}

func runProbe() int {
	timeout := flag.Duration("timeout", 3*time.Second, "max time to wait for a healthy response")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		return exitUnhealthy
	}

	target, err := url.Parse(flag.Arg(0))
	if err != nil || target.Scheme == "" {
		return fail("invalid target %q: need a scheme, e.g. http://host/health", flag.Arg(0))
	}

	prober, ok := probe.Get(target.Scheme)
	if !ok {
		return fail("unsupported scheme %q: have %s", target.Scheme, probe.SupportedSchemes())
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	if err := prober.Probe(ctx, target); err != nil {
		return fail("unhealthy: %v", err)
	}
	return exitHealthy
}

func fail(format string, args ...any) int {
	fmt.Fprintf(os.Stderr, "hc: "+format+"\n", args...)
	return exitUnhealthy
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: hc [-timeout DUR] TARGET\n")
	fmt.Fprintf(os.Stderr, "       hc install DEST\n")
	fmt.Fprintf(os.Stderr, "       hc version\n\n")
	fmt.Fprintf(os.Stderr, "TARGET is a URL whose scheme selects the probe: %s\n\n", probe.SupportedSchemes())
	fmt.Fprintf(os.Stderr, "examples:\n")
	fmt.Fprintf(os.Stderr, "  hc http://localhost:8080/health\n")
	fmt.Fprintf(os.Stderr, "  hc tcp://localhost:6379\n")
	fmt.Fprintf(os.Stderr, "  hc postgres://localhost:5432\n")
	fmt.Fprintf(os.Stderr, "  hc redis://localhost:6379\n")
	fmt.Fprintf(os.Stderr, "  hc install /healthz/hc\n")
}
