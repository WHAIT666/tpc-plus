# 2026-11-20 · Discovery and repository setup

## What happened

- Ran the first part of M0: an informal interview with the child, plus her first interim school report, shared by her parents.
- Parents gave consent.
- Set up the repository: README as a product page, architecture, milestones, 18 ADRs (15 from the technical plan v2.1 and its addendum, 3 new) and the M1 specification.

## What changed in the plan

- **Natural Sciences out, HGP in**, as the first mode B subject. The evidence: her school report and her own ranking ([ADR-016](../adr/016-mvp-subjects-maths-and-hgp.md)).
- **M1 starts before M0 is finished.** The maths engine needs no data about her; the corrected tests only change which skill comes first.
- Two proposals from an architecture review: an OpenAPI contract between Go and Kotlin ([ADR-017](../adr/017-openapi-contract.md)) and privacy-safe crash reporting ([ADR-018](../adr/018-privacy-safe-crash-reporting.md)).

## What I learned

- The most useful answer came from asking what she does when she is stuck, not what she would like. She sends the question to the class group chat.
- One of my questions was leading ("so they give you the answer and you don't learn, right?"). Next time: ask whether she could solve a similar one alone, or just give her one.

## Next

- M1 step 1: Go module, answer parser and its tests ([m1-spec.md](../m1-spec.md)).
- Finish M0: observation, corrected maths tests, photos, tablet checklist.
