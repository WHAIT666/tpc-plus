# ADR-021: Early EU hosting and device pairing for M2

**Status:** Accepted · 2026-11-20

## Context

M2 is done only when the child completes 3 sessions alone on her own tablet at home, in Portugal. Every answer is graded on the server ([architecture](../architecture.md)), so the API must be reachable from the internet in M2, not only in M6. The author moves abroad in January 2027 and must be able to run it remotely. Real authentication (rotating tokens, roles) and Postgres belong to M3 ([ADR-012](012-family-isolation-roles-idempotency.md)).

## Decision

- **Hosting:** one small VPS in an EU region (e.g. Hetzner CX22). The API runs as a Docker image built by CI and published to GitHub Container Registry; Caddy terminates HTTPS with automatic certificates for a free **DuckDNS** subdomain. A CI job deploys over SSH.
- **Device access:** `cmd/admin pair` (run on the server) prints a single-use code that expires in 15 minutes. The app exchanges it once for a random per-device token, sent as a Bearer token. The server stores only a hash of each token. A device can be revoked by deleting its entry.
- **Data:** attempts go to an append-only JSONL log on a Docker volume, with no personal data. `cmd/admin report` reads it.
- **Limits:** 4 KB bodies, request timeouts and a per-device rate limit from day one.

## Alternatives considered

- **Server on the author's PC with a tunnel** — free, but the tablet only works while the PC is on, and not at all after the move.
- **Serverless (Cloud Run)** — almost free, but the filesystem is not persistent, so the log would need another service.
- **Run the engine on the tablet** — no server needed, but grading would live in two places, against the architecture.
- **Token baked into the APK** — the APK is public on GitHub Releases, so the token would be too.

## Consequences

Running cost of about €4–5 a month from M2. The free subdomain is less polished than a real domain and can be swapped later without app changes if the app reads the base URL from its build config. M3 replaces pairing tokens with rotating access and refresh tokens, and the JSONL log with Postgres; M6 adds backups, observability and the threat model. Until then the log is copied off the server weekly.
