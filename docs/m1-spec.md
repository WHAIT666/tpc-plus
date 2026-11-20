# M1 · Exact maths engine — specification

**Goal:** a Go library and a terminal CLI that generate maths exercises with **computed, exact answers**, and grade answers typed the way a Portuguese 5th grader types them. No AI, no database, no app yet.

**Done when:**

- [x] 3 skills × 3 levels implemented (fractions first, see [Skills](#skills)).
- [x] Property tests verify at least 10,000 seeds per skill and level against an **independent implementation** (90,000+ exercises). This number goes in the README.
- [x] Golden tests pin the output of fixed seeds.
- [x] The answer parser is fuzzed and never panics.
- [x] `go run ./cmd/practice` runs a 5-exercise session in the terminal.
- [x] CI is green: vet, lint, tests with `-race`, a short fuzz run.
- [x] Docs: README Evidence row, devlog entry, CLI GIF in `docs/media/`, ADRs for any new decisions ([Definition of done](milestones.md#definition-of-done)).

> **How we work:** each step ends with a small commit on a branch and a review (correctness, tests against this spec, naming, Go idioms) before the next step starts. Results: [M1 devlog](devlog/2026-11-20-m1-exact-engine.md).

---

## Layout

```text
backend/
  go.mod                         module github.com/WHAIT666/tpc-plus/backend
  internal/grading/              answer parsing and exact grading
  internal/engine/exact/         generators, one file per skill, plus a registry
  cmd/practice/                  terminal practice session
  testdata/golden/               pinned generator output
```

## Core concepts

These describe behaviour, not code. Choosing the exact Go types is part of the exercise.

**Exercise:** what a generator produces.

| Field | Meaning |
|---|---|
| Skill ID | e.g. `mat.fracoes_equivalentes` (from the catalogue, [ADR-004](adr/004-fixed-catalogue-of-topics.md)) |
| Generator version | An integer, bumped whenever the output for a given seed changes |
| Seed, level | Together with the version, they reproduce the exercise exactly |
| Prompt | Text in European Portuguese, short, for a 10-year-old |
| Answer kind | `integer`, `fraction`, `decimal` or `choice` |
| Answer | The canonical answer: a `*big.Rat` for numbers, an option index for choices |
| Options | For `choice` only. Distractors come from **real misconceptions**, not random numbers |
| Requires irreducible | If true, an equivalent but unsimplified answer is not accepted |

**Generator:** has a skill ID and a version, and turns `(seed, level)` into an exercise. It must be **deterministic**: same version, seed and level give the same exercise, on any machine. Use `math/rand/v2` with an explicitly seeded PCG source, never the global random functions.

**Grade result:** correct, or one of the reasons *wrong*, *not simplified* or *invalid input*, plus the expected answer formatted for display. "Not simplified" is its own reason because the child deserves a different message ("certo, mas simplifica!").

## Answer parsing rules

The child types on a maths keypad. The parser must accept what she would naturally type and reject everything else **without panicking**.

| Input | Kind | Result |
|---|---|---|
| `3/4`, ` 3 / 4 ` | fraction | 3/4 |
| `6/8` | fraction | equivalent to 3/4. Accepted unless the exercise requires irreducible, then *not simplified* |
| `2` | fraction | 2/1 (an integer is a valid fraction) |
| `0,75` and `0.75` | decimal | 3/4 exactly (comma is the Portuguese decimal mark) |
| `1/0` | any | invalid: zero denominator |
| `-3/4` | any | invalid: no negative numbers in 5th grade |
| `1 1/2` | any | invalid in M1 (mixed numbers come later) |
| `1.000,5`, `3//4`, `abc`, empty | any | invalid |
| More than 32 characters or more than 9 digits in a number | any | invalid (input limits are a security habit) |

## Skills

The first three are fractions: the classic 5th-grade difficulty and a good fit for `big.Rat`. Skill 4 is chosen after you see her corrected tests ([m0-discovery.md](m0-discovery.md)).

### 1 · `mat.fracoes_equivalentes` — equivalent fractions

| Level | Exercise | Kind | Parameters |
|---|---|---|---|
| 1 | Complete: `2/3 = ?/12` | integer | denominator 2–10, multiplier 2–5 |
| 2 | Simplify `18/24` to its irreducible form | fraction, irreducible | irreducible base with denominator 2–12, multiplier 2–6 |
| 3 | Which fraction is equivalent to `3/5`? (4 options) | choice | distractors from misconceptions: add the same number to both terms (`5/7`), multiply only one term (`6/5`), swap them (`5/3`) |

### 2 · `mat.fracoes_somar_subtrair` — adding and subtracting fractions

| Level | Exercise | Kind | Parameters |
|---|---|---|---|
| 1 | Same denominator: `2/7 + 3/7` | fraction | denominator 3–12, result between 0 and 2 |
| 2 | One denominator is a multiple of the other: `1/4 + 3/8` | fraction | |
| 3 | Different denominators: `2/3 − 1/4` | fraction, irreducible | least common multiple ≤ 24; result never negative |

### 3 · `mat.fracao_de_quantidade` — fraction of a quantity

| Level | Exercise | Kind | Parameters |
|---|---|---|---|
| 1 | `1/4 de 20` | integer | unit fractions, whole-number results |
| 2 | `3/5 de 35` | integer | whole-number results |
| 3 | Short word problem: "A Rita tinha 24 cromos e deu 3/8 ao irmão. Com quantos ficou?" | integer | 3–4 templates with different names and objects; whole-number results |

## Tests

| Type | What it proves | How |
|---|---|---|
| **Table tests** | Parser and grader follow the rules above | Every row of the parsing table is a test case, plus your own edge cases |
| **Property tests** | Generators are correct for every seed, not just the ones you tried | For seeds 1–10,000 at every level: no error; the answer recomputed by an **independent oracle** (written in the test file with `int64` and your own `gcd`, **not** `big.Rat`) matches; the invariants hold (parameter ranges, whole-number results where promised, irreducible when required, distractors different from the answer and from each other) |
| **Golden tests** | Output doesn't change by accident | Seeds 1–20 per skill and level stored as JSON in `testdata/golden/`; a `-update` flag rewrites them. Any intentional change bumps the generator version |
| **Fuzzing** | The parser survives anything | `FuzzParse`: it never panics, and when it accepts an input, formatting the result and parsing again gives the same value. Locally `-fuzztime=60s`, in CI `10s` |
| **Benchmarks** (optional) | Generation is fast enough to never matter | `BenchmarkGenerate` per skill |

## CLI

```text
go run ./cmd/practice -skill mat.fracoes_equivalentes -level 2 -n 5 -seed 42
```

It shows one exercise at a time, reads the answer, gives immediate feedback in Portuguese and ends with a score. Without `-seed`, it uses a random seed and prints it, so any session can be replayed. Record a GIF of a session for `docs/media/` (e.g. with `vhs` or `asciinema`).

## CI

`.github/workflows/backend.yml`, triggered on pushes and pull requests that touch `backend/`:

1. `actions/setup-go` using the version in `go.mod`
2. `go vet ./...`
3. `golangci-lint` (official action)
4. `go test -race -count=1 ./...`
5. `go test -run=^$ -fuzz=FuzzParse -fuzztime=10s ./internal/grading`

## Steps, in order

| # | Step | Done when | Learn first |
|---|---|---|---|
| 1 | `go mod init`, folder layout, first commit on branch `m1/parser` | `go build ./...` passes | A Tour of Go (basics, methods, interfaces, errors) |
| 2 | Answer parser, test-first: write the table tests, then make them pass | All parsing rules covered | *Learn Go with Tests* (Hello World, Integers, Arrays, Structs); `math/big` docs for `Rat` |
| 3 | Grader: correct / wrong / not simplified / invalid | Table tests pass | |
| 4 | Fuzz the parser | Runs 60 s locally without failures | Go fuzzing tutorial (go.dev/doc/tutorial/fuzz) |
| 5 | Generator interface, registry and `mat.fracoes_equivalentes` levels 1–2 | Golden tests for levels 1–2 | `math/rand/v2` |
| 6 | Oracle and property tests for skill 1 | 10,000 seeds × levels pass | |
| 7 | Level 3 with misconception distractors | Property tests cover the distractors | |
| 8 | Skills 2 and 3, each with golden and property tests | All green | |
| 9 | CLI | A 5-exercise session works end to end | `flag`, `bufio` |
| 10 | CI workflow | Green on GitHub | GitHub Actions basics |
| 11 | Docs: README evidence, devlog, GIF, portfolio notes | Definition of done ticked | |

## Conventions

- **Commits:** [Conventional Commits](https://www.conventionalcommits.org/) (`feat(grading): parse decimal answers with a comma`), small and focused, one step per pull request.
- **Language:** code, comments, docs and commits in English. Exercise text in European Portuguese.
- **Errors:** return errors with context (`fmt.Errorf("parse %q: %w", in, err)`); no panics outside `main` and tests.

## Out of scope for M1

HTTP API, database, AI, the Android app, excluding the numbers in her homework (that needs scanned input and arrives in M4), mixed numbers, decimals as a separate skill.
