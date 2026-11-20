# ADR-017: OpenAPI as the API contract

**Status:** Proposed · 2026-11-20

## Context

Two languages (Go and Kotlin) share one API. Without a contract, they drift apart silently and bugs appear at runtime.

## Decision

An OpenAPI spec in the repository is the single source of truth for the API. CI checks the Go handlers against it (contract tests), and the Kotlin client models are generated from it or checked against it.

## Alternatives considered

- **Hand-kept DTOs on both sides** — drift.
- **gRPC** — strong contract, but heavier on Android and less common for a public REST API.

## Consequences

Decide before M2, when the first endpoint exists.
