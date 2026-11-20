# ADR-017: OpenAPI as the API contract

**Status:** Accepted · 2026-11-20 (proposed earlier the same day; accepted with the M2 plan)

## Context

Two languages (Go and Kotlin) share one API. Without a contract, they drift apart silently and bugs appear at runtime.

## Decision

`api/openapi.yaml` (OpenAPI 3.1) is the single source of truth for the API.

- **Go:** request and response types and the server interface are generated with `oapi-codegen` (standard `net/http`, strict server). Handlers implement the generated interface, so a handler that does not match the spec does not compile.
- **Kotlin:** the client models are generated from the spec with `openapi-generator` (Retrofit + kotlinx.serialization). If the generator does not fit the Android build, hand-written models are checked against the spec's examples in a unit test instead.
- **CI:** the spec is linted, and CI fails if regenerating the Go code changes anything (`go generate` then `git diff --exit-code`).

## Alternatives considered

- **Hand-kept DTOs on both sides** — drift.
- **gRPC** — strong contract, but heavier on Android and less common for a public REST API.
- **Validate responses against the spec only in tests** — catches drift later than a compile error.

## Consequences

One more generated layer and a third-party code generator on each side. Every API change starts in the spec, which also documents the API for free.
