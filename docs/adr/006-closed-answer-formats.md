# ADR-006: Closed answer formats

**Status:** Accepted · 2026-11-20

## Context

Open answers need subjective grading, which a language model would do inconsistently.

## Decision

Only closed formats in the MVP: numeric answers (integer, fraction, decimal with a comma), multiple choice, true/false and matching. Every answer is graded deterministically against a key.

## Alternatives considered

- **Open answers graded by a model** — richer, but neither verifiable nor reproducible.

## Consequences

Some skills can't be practised yet (written explanations). Grading is exact and testable.
