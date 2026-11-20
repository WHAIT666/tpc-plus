# ADR-018: Privacy-safe crash reporting on Android

**Status:** Proposed · 2026-11-20

## Context

If the app crashes on her tablet, nobody finds out why. A 10-year-old won't report it.

## Decision

Add crash reporting with no personal data: no user IDs beyond a random install ID, no screen content, no OCR text. Choose the tool (e.g. Sentry or Firebase Crashlytics) by how well it can be configured to strip PII.

## Alternatives considered

- **No crash reporting** — blind to failures in the real world.

## Consequences

Must be covered by the parents' consent and the privacy section. Decide in M2.
