# ADR-005: Exact arithmetic with big.Rat

**Status:** Accepted · 2026-11-20

## Context

Fractions and decimals in floating point give wrong answers (0.1 + 0.2 ≠ 0.3). A grader that marks a correct answer wrong destroys trust.

## Decision

All maths values are `math/big.Rat`. Decimal answers are parsed into exact rationals, never `float64`.

## Alternatives considered

- **float64 with tolerance** — tolerances are arbitrary and still fail at the edges.
- **int64 numerator/denominator** — workable, but overflow handling is on us.

## Consequences

Slightly more verbose code. Grading is exact, and equivalent forms (2/4 = 1/2) are easy to compare.
