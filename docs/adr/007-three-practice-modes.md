# ADR-007: Three practice modes and asking for the textbook page

**Status:** Accepted · 2026-11-20

## Context

Not every photographed page has enough content to ground questions. A maths worksheet or a page with only questions isn't a source.

## Decision

Three modes: **A** exact generators (maths); **B** grounded generation from a page with enough content; **C** a bank of human-approved questions. If a non-maths page has only questions, the app asks for the textbook page. If nothing works, it says so honestly.

## Alternatives considered

- **Always generate from the photo** — produces ungrounded questions.

## Consequences

More paths to build and test. Every session can say which mode produced it.
