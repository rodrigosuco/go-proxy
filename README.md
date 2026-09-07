# go-proxy

A minimal reverse proxy written from scratch in Go, with a hand-rolled per-IP token-bucket rate limiter. No frameworks, no third-party dependencies — just the standard library.

This is a learning project, built to get hands-on with `net/http/httputil.ReverseProxy` and with real concurrency problems (shared state across goroutines, data races, atomic map operations) rather than just reading about them.

## What it does

- Forwards every incoming HTTP request to a fixed backend target using `httputil.ReverseProxy`.
- Sets `X-Forwarded-*` headers on the outbound request via `ProxyRequest.SetXForwarded`.
- Rate-limits each client by IP address using a token bucket:
  - Each IP gets its own bucket, created on first request and stored in a `sync.Map`.
  - Every request consumes one token; once a bucket is empty, further requests get `429 Too Many Requests` until the bucket refills.
  - Buckets refill to full capacity after a fixed interval (30s) rather than leaking tokens continuously.
- Bucket state is protected per-IP with a `sync.Mutex`, so concurrent requests from the same client don't race on the token count.

## Architecture

```
cmd/
  main.go                        - wires up the target URL, the proxy, and the HTTP server
  middleware/
    ratelimit_middeware.go       - HTTP middleware: resolves client IP, loads/creates its bucket, applies the limit
  ratelimit/
    rate_limiter.go              - the token bucket itself (Bucket struct + CheckLimit/restoreTokens)
```

Request flow: `net/http` server → `RateLimitMiddleware` (checks/updates the caller's bucket) → `httputil.ReverseProxy` (forwards to the backend).

## Running it

Requires Go 1.26+.

```bash
go run ./cmd
```

By default it proxies to `http://localhost:3000/api/v1` and listens on `:8080`. The target is currently hardcoded in `main.go` — there's no config file or env var support yet (see [Limitations](#known-limitations--not-production-ready)).

## Known limitations / not production-ready

This project is intentionally small in scope, and it shows: it's a demonstration of specific Go concepts, not a general-purpose reverse proxy. Known gaps, roughly in order of importance:

- **No eviction of rate-limit state.** The `sync.Map` of buckets grows for every distinct IP the proxy has ever seen and never shrinks — an unbounded memory leak in a long-running process.
- **No server timeouts.** The HTTP server is started with `http.ListenAndServe` and no `ReadTimeout`/`WriteTimeout`/`IdleTimeout`, so a slow or idle client connection can be held open indefinitely.
- **Hardcoded, single backend target.** No config file, flags, or env vars — the destination URL is a literal in `main.go`, and there's no support for routing to multiple backends.
- **No TLS.** Plain HTTP only.
- **Global mutable rate-limit state.** The bucket store is an exported package-level `sync.Map`, not encapsulated behind a constructor — fine for a single-binary demo, but not how you'd want it structured behind a public API.
- **Off-by-one in bucket sizing.** A brand-new bucket currently starts with one fewer token than a bucket that has gone through a restore cycle (14 vs. 15) — a known inconsistency, not (yet) fixed.

## Why build this instead of using nginx/Envoy/Traefik/Caddy?

You generally shouldn't, for a real deployment — those tools have already solved TLS termination, HTTP/2, connection handling edge cases, and service discovery far more thoroughly than a weekend project ever will. This exists to learn the mechanics that those tools hide: how a reverse proxy actually rewrites and forwards a request, and how to keep shared state correct under real concurrent load without reaching for a library to do it for me.

## License

No license has been chosen yet — treat this as source-available for reading/learning until a `LICENSE` file is added.
