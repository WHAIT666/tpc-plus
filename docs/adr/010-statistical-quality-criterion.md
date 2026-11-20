# ADR-010: Statistical quality criterion per topic

**Status:** Accepted · 2026-11-20 (amended by the v2.1 addendum)

## Context

"95% correct on 50 questions" says little: with 50/50 the 95% Wilson lower bound is only 92.9%.

## Decision

Human review with severity (critical, major, minor). A topic meets the criterion with **at least 100 reviewed questions, 0 critical errors and a Wilson lower bound of at least 93%** on the pass rate. Sampling is stratified by topic, question type, level and page, with a minimum per topic. A weak topic doesn't open because the subject's average is good.

## Alternatives considered

- **Fixed percentage on a small sample** — false confidence.
- **Criterion per subject only** — hides weak topics.

## Consequences

Reference points: 100/100 → 96.3%; 99/100 → 94.6%; 30/30 → 88.7%. Topics with few reviews stay conservative (ADR-015).
