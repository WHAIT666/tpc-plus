# ADR-004: Fixed catalogue of subjects and topics

**Status:** Accepted · 2026-11-20

## Context

Free-text topics can't be evaluated or mapped to generators and progress.

## Decision

`catalog/topics.yaml` is the single source of truth for subjects and topics, with IDs such as `mat.fracoes_equivalentes`, based on the Portuguese national curriculum (Aprendizagens Essenciais). Classification must pick from this list (enums in the schema).

## Alternatives considered

- **Free-text topic from the model** — flexible, but impossible to evaluate.

## Consequences

Classification accuracy becomes measurable (top-1/top-3). New topics need a catalogue change, which is intended.
