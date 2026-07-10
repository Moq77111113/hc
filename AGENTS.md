# AGENTS.md (hc)

Instructions for coding agents working in this repo. Read `README.md` first;
it defines what hc is and, just as importantly, what it is not.

## What hc is

A single static Go binary that health-checks one target and exits `0`
(healthy) or `1` (unhealthy), following Docker's HEALTHCHECK contract. It is
copied into hardened/minimal images (Chainguard, Wolfi, distroless, scratch)
that ship no shell or curl:

```dockerfile
COPY --from=ghcr.io/moq77111113/hc /hc /hc
HEALTHCHECK CMD ["/hc", "http://localhost:8080/health"]
```

It can also be injected into a third-party image you don't rebuild, via a shared
volume seeded by `hc install DEST` (see `deploy/`). It stays one-shot either way.

## Scope

In scope:

- Schemes: one connect-level readiness probe per URL scheme. The registry is
  the source of truth — `hc` with no args (or `SupportedSchemes()`) lists the
  live set; don't maintain a roster here.
- Modular builds: `hc_slim` opts out of every probe, a per-scheme `hc_<scheme>`
  tag opts one back in. The `HC_TAGS` build arg drives a local `docker build`;
  the published image variants (`hc` full, plus the slim bundles in
  `.goreleaser.yaml`) come from each build's `-tags`.
- Injection: the `hc install DEST` subcommand plus the `deploy/` manifests.
- The `-timeout` flag (default 3s).

Do NOT add without the owner's say-so. hc stays **one-shot**:

- No central daemon, no long-running or listening process, no socket, no config
  file, no state.
- No CLI flags beyond `-timeout` and the `install` subcommand.

Scope creep is the failure mode this project exists to resist. If a change adds
surface area, stop and ask.

## Architecture

- `main.go`: the CLI surface. Dispatch `install`, else parse the target, pick a
  prober by URL scheme, apply the timeout, map the result to an exit code. No
  protocol logic here.
- `install.go`: the `hc install DEST` self-copy (a scratch image has no `cp`, so
  injection seeds a shared volume this way).
- `internal/probe/`: the probers.
  - `registry.go`: the `Prober` interface, the scheme→prober map, `register`,
    `Get`, `SupportedSchemes`. The only extension seam.
  - `<scheme>.go`: one file per prober, build-tagged and self-registering in
    `init()`. One nameable protocol per file.
  - Shared probe machinery lives beside them: `httpcore.go` (the hand-rolled
    HTTP/1.1 status probe behind `http`/`https`, no `net/http`) and
    `handshake.go` (the `dial` + `handshake` helper the byte-level probers reuse
    — dial, send a payload, judge the reply).

### Adding a protocol, the only sanctioned extension

1. Add `internal/probe/<scheme>.go`: a `Prober` struct with one
   `Probe(ctx, *url.URL) error` method (return `nil` when healthy), a
   `//go:build !hc_slim || hc_<scheme>` tag, and `func init() { register("<scheme>", …) }`.
2. Add `hc_<scheme>` to any bundle build (`hc-core`, `hc-sql`) in
   `.goreleaser.yaml` that should include it.

That is the entire change. If a protocol needs more, the design is wrong. Raise
it, don't work around it.

## Coding conventions

- The scheme→prober map encodes the decision; callers never choose a probe by
  hand. Keep that shape. Narrow types, imposed flow.
- Guard clauses early, flat returns, no deep nesting.
- Doc comment on every exported symbol; that's the contract. Terse: what, when,
  why, not prose.
- WHY-comments only for the non-obvious (the postgres SSLRequest trick, the
  hand-rolled HTTP dropping net/http, the https skip-verify). No WHAT-comments.
- No inline comments explaining what a line does. If a line needs one to be
  understood, rename or extract until it doesn't; refactor over annotate. The
  only comments inside a function body are the rare WHY above and `//nolint`
  directives.
- No new runtime dependencies. `CGO_ENABLED=0` must keep building on `scratch`.
  The standard library is the ceiling.

## Build & check

```sh
go vet ./...
go test ./...                                        # default build: every probe
go test -tags "hc_slim hc_tcp" ./...                 # a slim build compiles + passes
CGO_ENABLED=0 go build -ldflags="-s -w" -o /tmp/hc . # must build with CGO off
docker build -t hc:test .                            # must produce a scratch image
docker run --rm hc:test                              # prints usage, exits 1
```

## Git

- Branch pattern `{type}/{short-description}`; never work on `main`.
- Commit pattern `{type}({scope}): {short-description}`; short commits, no body,
  no self-attribution.
- Never commit documentation (README, VISION, this file); the owner does that.
- Never use force flags to bypass `.gitignore`.
