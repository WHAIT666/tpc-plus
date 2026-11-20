# ADR-012: Family isolation, roles and transactional idempotency

**Status:** Accepted · 2026-11-20 (amended by the v2.1 addendum)

## Context

Even with one family today, the API must be safe for many. A child's tablet is an easy target for token theft. Flaky networks resend answers.

## Decision

- The family always comes from the token, never from the URL or body. Authorisation tests prove that one family can't read another's data (IDOR tests).
- Roles `child`, `parent`, `operator`: the child answers and sees progress, but can't revoke devices, change parental settings or delete data.
- 1-hour access tokens, rotating refresh tokens with reuse detection (reusing an old one revokes the chain), single-use pairing codes.
- Each answer carries an `Idempotency-Key`. Storing the attempt and updating progress happen in **one transaction**, and a repeat with the same key returns the stored result, even under concurrency.

## Alternatives considered

- **Trust IDs in the request** — the classic IDOR bug.
- **Long-lived tokens** — one leak lasts forever.

## Consequences

More auth code than one user needs. This is accepted on purpose, because it is the evidence of secure API design.
