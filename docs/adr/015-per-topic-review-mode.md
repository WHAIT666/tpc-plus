# ADR-015: Per-topic review mode — conservative, then dynamic

**Status:** Accepted · 2026-11-20 (new in the v2.1 addendum)

## Context

Reviewing every question forever doesn't scale. Reviewing none is unsafe.

## Decision

Per topic, two modes. **Conservative** until the topic meets ADR-010: no mode B question reaches the child without human approval. **Dynamic** afterwards: automatically validated questions may be shown, with the residual risk stated, a "this question is wrong" button and continuous review of a ~10% sample. A topic falls back to conservative if quality drops.

## Alternatives considered

- **Always review everything** — doesn't scale.
- **Never review** — unsafe.

## Consequences

The review CLI and the stratified sampling are part of M5.
