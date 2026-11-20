# Portfolio notes

A running log of CV and portfolio material. **Rule:** a bullet moves to "Ready" only when every number in it has been measured and can be shown in this repository.

## Ready

- Built an exact maths engine in Go for a 5th-grade practice app: 3 fraction skills × 3 levels, with exact rational arithmetic (`math/big.Rat`), deterministic seeded generation and golden tests.
- Verified 90,000 generated exercises against an independent `int64` oracle with property tests (0 failures, 0.4 s), and confirmed the tests catch bugs with deliberate mutations.
- Hardened the answer parser for a child's keypad with table tests and native Go fuzzing: 25.1 million inputs in 60 s, 0 crashes; fuzzing runs on every CI build with `go vet`, golangci-lint and the race detector.
- Designed multiple-choice wrong answers from named misconceptions, so feedback explains the child's mistake instead of just saying "wrong".

## Drafts — fill the [ ] with real numbers, do not use before they are true

- Built an Android app (Kotlin, Jetpack Compose) and a Go backend that turn a photo of any 5th-grade homework page into new practice on the same subject and topic, designed with and for a real student.
- Three practice modes: exact generators with rational arithmetic for maths, grounded LLM generation with deterministic validation and a minimal-context verifier, and a bank of human-approved questions.
- Question quality gated by human review: [N] questions reviewed, 0 critical errors, [X]% pass rate (95% Wilson lower bound [Y]%).
- Subject and topic classification with LLM structured outputs: [X]% top-1 topic accuracy on a [N]-item labelled eval set; p95 latency [T] s; [C] per session.
- Threat-modelled API: per-family isolation proven by authorisation tests, rotating refresh tokens, idempotent answers and an AI cost circuit breaker; deployed with Docker and GitHub Actions in an EU region.

## Interview stories (STAR)

- **The standard library was too permissive.** `big.Rat.SetString` accepts `-3/4`, `1e5`, `0x10` and reads `010/3` as 8/3. A grader built on it would mark nonsense correct. Wrote strict validation first and fuzzed it (25.1M inputs, 0 crashes) ([ADR-019](adr/019-strict-answer-input.md)).
- **Tests that prove they can fail.** Property tests check 90,000 exercises against an oracle that reads the prompt back with `int64` maths; each one was shown to fail on a deliberate bug. They also caught a real edge case: two misconceptions giving the same value (2/4 and 3/6 for 1/4) ([ADR-020](adr/020-misconception-distractors.md)).
- **Discovery changed the plan.** The interview and the school report moved the first text subject from Natural Sciences to History & Geography, and showed that the real competitor is the class group chat ([m0-discovery.md](m0-discovery.md)).

## Media to capture

- [x] M1: terminal GIF of a practice session ([`docs/media/m1_cli-session_2026-11-20.gif`](media/m1_cli-session_2026-11-20.gif))
- [ ] M2: screen recording of a full session on the tablet, in demo mode
- [ ] M5: screenshot of the review queue with Wilson statistics
- [ ] M6: 60-second product video
