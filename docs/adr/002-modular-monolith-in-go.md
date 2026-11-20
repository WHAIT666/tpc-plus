# ADR-002: Modular monolith in Go

**Status:** Accepted · 2026-11-20

## Context

One developer, one real user and a handful of features that share data (sessions, attempts, questions, budget). The backend must be easy to run, test and deploy.

## Decision

A single Go service with clear internal packages (`internal/engine/exact`, `internal/grading`, `internal/budget`, …), one PostgreSQL database and several small binaries in `cmd/` (API, practice CLI, eval, review).

## Alternatives considered

- **Microservices / Kubernetes** — operational cost with no benefit at this scale.
- **Node.js/TypeScript backend** — already on the CV; Go adds a new, in-demand skill and fits exact arithmetic and CLIs well.

## Consequences

Package boundaries must be kept by discipline (no import cycles, `internal/`). Deployment is one container plus Postgres.
