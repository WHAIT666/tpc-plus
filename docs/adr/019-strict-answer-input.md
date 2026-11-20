# ADR-019: Strict answer input rules

**Status:** Accepted · 2026-11-20

## Context

The child types answers on a maths keypad. The grader must accept what she naturally types and reject everything else without panicking ([ADR-006](006-closed-answer-formats.md)). Go's `big.Rat.SetString` is the obvious parser, but it accepts far more than a 5th grader ever means: signs, exponents (`1e5`), hex and octal prefixes (`010/3` is read as 8/3), underscores, `.5` and `5.`.

## Decision

Answers are validated against our own rules first and only then converted with `SetFrac64`; `SetString` is never used.

- Accepted: whole numbers, fractions `a/b` (spaces around the slash are fine) and decimals with a comma or a point (`0,75`, `0.75`).
- ASCII digits only; at most 32 bytes of input and 9 digits per number.
- No leading zeros, except a lone `0` and the digits after the decimal mark (`0,05`).
- Rejected with their own reason: negative numbers, mixed numbers (`1 1/2`), zero denominators.
- A whole number is valid for integer, fraction and decimal exercises; no other form crosses kinds (`0,75` is not accepted for a fraction exercise).
- `1.000` is read as 1: the point is a decimal mark, not a thousands separator.
- The parser reports whether a typed fraction was in lowest terms, because `big.Rat` reduces it immediately.

This refines ADR-006, which listed only the decimal comma.

## Alternatives considered

- **`big.Rat.SetString` with a few extra checks** — every extra syntax it accepts would need its own check, and a missed one becomes a wrong grade.
- **Accept any equivalent form across kinds** — friendlier, but the child would never practise writing the form the exercise asks for.

## Consequences

Every rejection maps to a specific hint for the child. The rules are table-tested row by row and fuzzed. Mixed numbers and thousands separators need a new decision when a skill requires them.
