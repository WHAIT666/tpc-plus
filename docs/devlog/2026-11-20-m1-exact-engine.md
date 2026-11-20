# 2026-11-20 · M1: the exact maths engine

## What was built

- **Answer parser** (`internal/grading`): reads integers, fractions and decimals (comma or point) the way a 5th grader types them, and rejects everything else with a specific reason. It never calls `big.Rat.SetString`; every input is checked against our own rules first ([ADR-019](../adr/019-strict-answer-input.md)).
- **Exact grader:** correct, wrong, *not simplified* or invalid input, plus the expected answer formatted for display. "Not simplified" only applies when the value is right.
- **Generators** (`internal/engine/exact`): 3 skills × 3 levels, all fractions, deterministic from `(version, seed, level)` with an explicitly seeded PCG.
  - `mat.fracoes_equivalentes`: complete, simplify, and a multiple-choice level whose wrong options come from real misconceptions ([ADR-020](../adr/020-misconception-distractors.md)).
  - `mat.fracoes_somar_subtrair`: same denominator, one a multiple of the other, unrelated denominators (irreducible answer).
  - `mat.fracao_de_quantidade`: unit fractions, any fraction, and four short word problems.
- **Terminal CLI** (`cmd/practice`): a session of N exercises with feedback in Portuguese, a hint after invalid input, the misconception behind a wrong option, a score, and the exact command to replay the session.
- **CI** (GitHub Actions): `go vet`, golangci-lint, `go test -race` and a 10 s fuzz run on every backend change.

## Numbers

| What | Value |
|---|---|
| Exercises verified against an independent oracle | 90,000 (3 skills × 3 levels × 10,000 seeds), 0 failures, 0.4 s |
| Fuzzing of the answer parser | 25.1 million inputs in 60 s, 0 crashes |
| Golden files | 9 (seeds 1–20 per skill and level) |
| Tests | 221 passing |
| Code | ~1,300 lines of Go, ~1,500 lines of tests |

## How the tests earn trust

- **Independent oracle.** The property tests live in package `exact_test`, so the compiler stops them from reusing the generators' helpers. They read the numbers back from the prompt the child sees and recompute the answer with `int64` and their own `gcd`, never `big.Rat`.
- **Ranges are reached, not just respected.** Every allowed denominator, multiplier, denominator pair, operation, correct-option position and story must appear at least once in 10,000 seeds.
- **Mutation-checked.** Each property test was shown to fail on a deliberate bug: an off-by-one answer, a narrowed denominator range, a missing guard against equal distractors, a negative difference, the part answered instead of the remainder, and a widened story range.
- **Fuzzing** checks that the parser never panics, only returns known errors, never accepts a negative or over-long input, and that formatting an accepted value and parsing it again gives the same value.

## What broke, and what was learned

- **`big.Rat.SetString` is far too permissive for a child's keypad.** It accepts `-3/4`, `1e5`, `0x10`, `1_000`, `.5`, and reads `010/3` as 8/3 (octal). The fix was to validate with our own rules and only then build the value with `SetFrac64`.
- **`big.Rat` always reduces**, so after parsing `6/8` there is no way to know it was not simplified. The parser now returns `InLowestTerms` alongside the value.
- **Two misconceptions can give the same value.** For 1/4 with ×2, "multiply only the numerator" gives 2/4 and "add 2 to both terms" gives 3/6: both equal 1/2. The generator falls back to "multiply only the denominator", which can never collide. The property tests caught this class of bug by checking that all four options have different values.
- **Linting Portuguese.** `misspell` flagged *denominador* as a typo of *denominator*. It now checks comments only (English); strings for the child are Portuguese.
- **CI did not start at first:** GitHub refused to run Actions on the private repository because of a billing problem on the account. The repository is now public, where Actions are free.
- **Locally, `go test -cover` fails** on this Go install (`internal/coverage is not in std`); plain tests and CI are unaffected. A clean Go reinstall should fix it.

## Decisions

- [ADR-019](../adr/019-strict-answer-input.md): strict answer input rules.
- [ADR-020](../adr/020-misconception-distractors.md): multiple-choice distractors come from named misconceptions.
- Choices are typed as the number shown on screen, from 1; the answer key stores a 0-based index.
- A session with seed *s* uses seeds *s*, *s*+1, …; sessions with adjacent seeds overlap, which is acceptable with random seeds up to 1,000,000.

## Media

- [`docs/media/m1_cli-session_2026-11-20.gif`](../media/m1_cli-session_2026-11-20.gif): a 3-exercise session at level 3, rendered from the real output of the CLI with scripted answers.

## Next

- M2: the Kotlin + Compose app with a 5-exercise session, using this engine behind a small API.
- Finish M0: observation, corrected maths tests (they may reorder the skills), photos and the tablet checklist.
