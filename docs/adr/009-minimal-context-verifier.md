# ADR-009: Minimal-context verifier

**Status:** Accepted · 2026-11-20 (amended by the v2.1 addendum)

## Context

A verifier that sees the whole page can 'fill in' with general knowledge and approve a question the page doesn't support.

## Decision

A separate request receives **only** the evidence spans (max 3), the stem and the options, never the key. It returns the option it would choose and pass/fail/uncertain: do the spans alone justify the answer, and would another option also be defensible? `uncertain` blocks the question.

## Alternatives considered

- **Verifier with the full page** — easier to pass, weaker guarantee.
- **No verifier** — relies on generation alone.

## Consequences

More AI calls per question (cost is covered by ADR-013). A different model may be used as verifier to reduce correlated errors, decided by the eval.
