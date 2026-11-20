# ADR-001: No answer without proof

**Status:** Accepted · 2026-11-20 (amended by the v2.1 addendum)

## Context

The user is a 10-year-old. A wrong question or a wrong answer key teaches her something false, and she has no way to notice. Language models can produce fluent questions with wrong keys.

## Decision

Every exercise must be backed by proof: **computed** where it can be computed (mode A), **grounded and verified** in the source page (mode B), or **already approved** by a human (mode C). The principle reads: *no answer without proof; residual risk is measured and stated.* It is a quality objective with a measured error rate, not an absolute guarantee.

## Alternatives considered

- **Let the model generate questions and answers freely** — cheap and flexible, but errors can't be detected or measured.
- **Human-written bank only** — safe, but it can't follow her textbook or scale to every subject.

## Consequences

Every mode needs its own proof mechanism (property tests, validation + verifier, review). More work up front, and every quality claim in the README can be backed by a number.
