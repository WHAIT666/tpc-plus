# ADR-013: AI cost ceiling

**Status:** Accepted · 2026-11-20 (amended by the v2.1 addendum)

## Context

AI cost must never surprise. Concurrent requests can overshoot a naive limit.

## Decision

A monthly ceiling (default €15). The daily quota per family is **derived from the monthly budget**: at ~$0.155 per mode B session, about 3 sessions a day. Before each AI call the budget is **reserved atomically**, then settled with the real cost from `usage`. When the ceiling is reached, mode B refuses and modes A and C keep working.

## Alternatives considered

- **No limit** — unacceptable.
- **Check-then-spend without a reservation** — races overshoot the ceiling.

## Consequences

Every call is recorded in `ai_usage`. An operator view shows the month's spend. A test proves A and C still work at the ceiling.
