# ADR-020: Distractors come from named misconceptions

**Status:** Accepted · 2026-11-20

## Context

A multiple-choice question is only as good as its wrong options. Random wrong numbers are easy to dismiss and teach nothing; a wrong option the child actually picks tells us which mistake she is making.

## Decision

Every wrong option in a generated multiple-choice exercise is built from a named misconception and records it:

| Misconception | Example for 3/5 |
|---|---|
| `add_same_to_both_terms` | 5/7 |
| `multiply_one_term` | 6/5 or 3/10 |
| `swap_terms` | 5/3 |

The four options always have different values, and none is equivalent to the answer. When two misconceptions would give the same value, the generator uses a variant that cannot collide. Property tests rebuild every distractor from its rule.

## Alternatives considered

- **Random nearby values** — simpler, but feedback could only say "wrong".
- **Store only the option text** — the mistake could be guessed later, but not reliably.

## Consequences

Feedback can explain the mistake ("Trocaste o numerador com o denominador."), and later the attempt log can show which misconceptions keep coming back. Each new choice skill needs its own list of misconceptions.
