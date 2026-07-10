# hc

**You hardened your image. Your health check died with it.**

You moved to `scratch` / distroless / Wolfi / Chainguard to shrink the attack
surface, and the same move that removed the shell and `curl` broke every
`HEALTHCHECK` you had. `hc` is the one binary you add back: no shell, no curl, no
libc, no config. Copy it in, point it at a target. Its exit code is Docker's
contract: `0` healthy, `1` unhealthy.

[![Go Reference](https://pkg.go.dev/badge/github.com/Moq77111113/hc.svg)](https://pkg.go.dev/github.com/Moq77111113/hc)
[![Release](https://img.shields.io/github/v/release/Moq77111113/hc)](https://github.com/Moq77111113/hc/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## The problem

Minimal images are more secure. They also delete the tools your health checks
lean on: no shell, no curl, no package manager. So the check every tutorial
hands you has nothing to run:

```dockerfile
# no curl in the image, so this never runs
HEALTHCHECK CMD curl -f http://localhost:8080/health || exit 1
```

A secure image that can't tell anyone whether it's alive.

## Quick start

```dockerfile
COPY --from=ghcr.io/moq77111113/hc /hc /hc
HEALTHCHECK CMD ["/hc", "http://localhost:8080/health"]
```

That's the whole integration. No sidecar, no libc, no per-service check script.

## Probes

The URL scheme picks the probe:

```sh
hc http://localhost:8080/health   # healthy on 2xx-3xx
hc tcp://localhost:6379           # healthy if the TCP connection opens
hc postgres://localhost:5432      # healthy if PostgreSQL accepts connections
```

| Scheme            | Healthy when                                         |
| ----------------- | ---------------------------------------------------- |
| `http` / `https`  | response status is `2xx` or `3xx`                    |
| `tcp`             | the TCP connection is established                    |
| `postgres` / `pg` | the server answers the readiness handshake (no auth) |

`https` checks that the service answers over TLS; it does **not** validate the
certificate. hc reports liveness, not cert validity, so it works against internal
and self-signed endpoints.

Set a deadline with `-timeout` (default `3s`):

```sh
hc -timeout 1s http://localhost:8080/health
```

## Health-check an image you don't own

The pain doubles with third-party hardened images like Keycloak, Postgres, or a
DHI base whose Dockerfile you don't control. You can't add a `HEALTHCHECK`
without forking the image and maintaining your own copy forever.

Don't. Inject `hc` through a shared volume: the target image is never modified,
and it all runs non-root.

```yaml
# kubernetes 1.33+: mount hc from its image, read-only
volumes:
  - name: healthz
    image: { reference: ghcr.io/moq77111113/hc-core }
containers:
  - name: app                     # third-party image, untouched
    volumeMounts:
      - { name: healthz, mountPath: /healthz, readOnly: true }
    livenessProbe:
      exec: { command: ["/healthz/hc", "tcp://localhost:8080"] }
```

Older clusters and Docker Compose seed the volume with `hc install /healthz/hc`
(a `scratch` image has no `cp`). Manifests for image-volume, initContainer, and
Compose live in [`deploy/`](deploy/).

## Lean by build

The default binary carries every probe. Need less surface? Pull a variant, and
the probes you don't use aren't compiled in:

| Image     | Probes            | Size     |
| --------- | ----------------- | -------- |
| `hc`      | all               | ~4.5 MB  |
| `hc-core` | http, https, tcp  | ~4.5 MB  |
| `hc-sql`  | tcp, postgres     | ~2.3 MB  |

```dockerfile
COPY --from=ghcr.io/moq77111113/hc-sql /hc /hc
```

## Install

As a build stage (recommended):

```dockerfile
COPY --from=ghcr.io/moq77111113/hc /hc /hc
```

Or grab a [release binary](https://github.com/Moq77111113/hc/releases), or from source:

```sh
go install github.com/Moq77111113/hc@latest
```

## Missing a protocol?

`redis`, `mysql`, and `amqp` are on the list. Need one of those, or a scheme
that isn't here? [Open an issue](https://github.com/Moq77111113/hc/issues), or send a PR.

## Philosophy

`hc` answers one question: "is this endpoint alive?". No monitoring, no retries,
no state, no metrics, no history. Want those? Wrong binary, on purpose.

## License

MIT
