# Portfolio notes

A running log of CV and portfolio material. **Rule:** a bullet moves to "Ready" only when every number in it has been measured and can be shown in this repository.

## Ready

*(nothing yet)*

## Drafts — fill the [ ] with real numbers, do not use before they are true

- Built an Android app (Kotlin, Jetpack Compose) and a Go backend that turn a photo of any 5th-grade homework page into new practice on the same subject and topic, designed with and for a real student.
- Three practice modes: exact generators with rational arithmetic for maths, grounded LLM generation with deterministic validation and a minimal-context verifier, and a bank of human-approved questions.
- Exact maths engine: [N] skills, [M] exercises per skill verified against an independent implementation with property tests and fuzzing.
- Question quality gated by human review: [N] questions reviewed, 0 critical errors, [X]% pass rate (95% Wilson lower bound [Y]%).
- Subject and topic classification with LLM structured outputs: [X]% top-1 topic accuracy on a [N]-item labelled eval set; p95 latency [T] s; [C] per session.
- Threat-modelled API: per-family isolation proven by authorisation tests, rotating refresh tokens, idempotent answers and an AI cost circuit breaker; deployed with Docker and GitHub Actions in an EU region.

## Interview stories (STAR)

- **Discovery changed the plan.** The interview and the school report moved the first text subject from Natural Sciences to History & Geography, and showed that the real competitor is the class group chat ([m0-discovery.md](m0-discovery.md)).

## Media to capture

- [ ] M1: terminal GIF of a practice session (`docs/media/m1_cli-session_YYYY-MM-DD.gif`)
- [ ] M2: screen recording of a full session on the tablet, in demo mode
- [ ] M5: screenshot of the review queue with Wilson statistics
- [ ] M6: 60-second product video
