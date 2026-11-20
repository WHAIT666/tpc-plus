# Milestones

Each milestone ends with a checkable outcome. It is **done** only when its documentation is also done (see [Definition of done](#definition-of-done)).

| Milestone | Scope | Done when |
|---|---|---|
| **M0 · Discovery** | Interviews with the child and her parents, school report, written consent, topic choice, photos, tablet check | MVP subjects and topics chosen from evidence; consent on file; tablet verified; at least 10 labelled photos per subject; [m0-discovery.md](m0-discovery.md) written |
| **M1 · Exact engine** | Maths generators, exact grader, answer parser, property tests, fuzzing, terminal CLI. No AI, no database | Thousands of exercises per skill verified by an independent implementation; CI green; CLI demo recorded ([m1-spec.md](m1-spec.md)) |
| **M2 · Usable app** | Kotlin + Compose, manual topic choice, session of 5, maths keypad, simple API; pulled forward from M3/M6: device pairing, an append-only attempt log and an early EU deploy ([ADR-021](adr/021-early-eu-hosting-and-device-pairing.md)) · [m2-spec.md](m2-spec.md) | **She completes 3 sessions alone, on different days** |
| **M3 · Data and security** | Postgres (replacing the M2 JSONL log), families, rotating tokens (replacing M2 device pairing), family isolation tests, idempotency, attempt log, placement check, simple rule, spaced review, fixed check-ins | Cross-family tests pass; every attempt is recorded exactly once |
| **M4 · Scan and classify** | ML Kit scanner + metrics, OCR gate, `/v1/classify` with sufficiency, classification eval | Accuracy and thresholds measured on real photos; manual choice still available |
| **M5 · HGP (B + C)** | Generation, validation, similarity, minimal-context verifier, replenishment, textbook-page request, bank C, human review, reports, budget and ceiling | HGP topics meet the criterion (≥100 reviewed, 0 critical, Wilson ≥ 93%) |
| **M6 · Production** | Hardening the M2 deploy, CI/CD, backups with tested restore, observability, threat model, README with evidence, demo mode, 60-second video | A recruiter can see everything working without her tablet |
| **M7 · Adaptive and expansion** | BKT evaluated offline against the simple rule; Natural Sciences, Portuguese and English, one at a time | Each new subject meets the criterion before it opens |

Phase 2, after M6: a Go systems project (a Redis-compatible server or an interpreter).

## Definition of done

A milestone is not done until all of these are true:

- [ ] Code merged, CI green, no skipped tests.
- [ ] README updated: status, roadmap row, and any **measured** numbers in the Evidence table.
- [ ] A devlog entry in `docs/devlog/`: what was built, what broke, what was learned, numbers.
- [ ] New decisions recorded as ADRs.
- [ ] Screenshots, GIFs or a short video in `docs/media/`, recorded with demo data only.
- [ ] `docs/portfolio-notes.md` updated with CV-ready bullets backed by real numbers.
