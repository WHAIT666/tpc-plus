# ADR-022: Stateless exercise references

**Status:** Accepted · 2026-11-20

## Context

The app shows exercises and the server grades the answers. The server needs to know which exercise an answer belongs to, without sending the answer to the tablet and without a database (M3).

## Decision

Each exercise is identified by a reference `skill@version:level:seed`, for example `mat.fracoes_equivalentes@1:3:812`. M1 generators are deterministic, so the server regenerates the exercise from its reference to grade an answer.

- The session payload carries the reference, the prompt, the kind and the option texts, never the answer, `correct` flags or misconceptions. A test checks this.
- A reference whose generator version no longer exists gets `410 Gone`; the app starts a new session.
- A session is a random ID; the server picks a random seed and uses seed, seed+1, … for its 5 exercises, as the CLI does.

## Alternatives considered

- **Send the answers to the app and grade there** — grading logic in two languages, and the answer is one network inspector away.
- **Store each session's exercises on the server** — needs storage that M2 does not have, and adds nothing determinism does not already give.

## Consequences

Grading is reproducible from the log alone. Bumping a generator version invalidates sessions in progress, which is acceptable for 5-exercise sessions. Anyone could compute answers from a reference, which does not matter for a practice app with no stakes.
