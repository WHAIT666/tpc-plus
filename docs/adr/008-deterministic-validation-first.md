# ADR-008: Deterministic validation first

**Status:** Accepted · 2026-11-20 (amended by the v2.1 addendum)

## Context

Checks that code can do are free, fast and reliable. Model-based checks are not.

## Decision

Before any model-based verification, Go code validates each mode B question: up to 3 **evidence spans** must exist in the page text (position and content), the key must be one of the options, the format must be valid. A **concept fingerprint** (spans + correct answer) is compared with the homework questions and her recent sessions on that topic, to reject near-duplicates. Mode A generators exclude the exact numbers from the homework and combinations she saw recently.

## Alternatives considered

- **Trust the model to validate itself** — correlated errors.

## Consequences

The discard rate at each step becomes a production metric.
