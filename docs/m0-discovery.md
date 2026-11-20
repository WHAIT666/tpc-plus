# M0 · Discovery summary

**Date:** November 2026 · **Status:** in progress

This is the public, anonymised summary. Raw notes, photos and school documents stay in a private folder outside the repository.

## Method

- Short informal interview with the child (10 years old, 5th grade, Portugal), with past homework in hand. Questions were about concrete past events, not hypothetical features.
- The first interim school report of the year (grades and the class teacher's comment), shared by the parents.
- Written consent from the parents before any of her data was used.
- Still to do: a think-aloud observation (maths exercises plus a History & Geography page), corrected maths tests, photos for the eval set and the tablet checklist.

## Findings

| # | Finding | Evidence | What it means for the product |
|---|---|---|---|
| 1 | Maths is where she struggles most | School report; the teacher's comment names maths explicitly; she ranks it among her least-liked subjects | Mode A (maths) stays the core of the MVP |
| 2 | History & Geography of Portugal (HGP) is weaker than Natural Sciences | School report; she ranks HGP among her least-liked subjects | HGP replaces Natural Sciences as the first mode B subject ([ADR-016](adr/016-mvp-subjects-maths-and-hgp.md)) |
| 3 | When stuck, she asks the class group chat and receives the answer | Interview | The real competitor is a group chat: instant, social, effortless. Practice after homework tests what copying hides. A "hint, not the answer" mode goes to the backlog |
| 4 | The teacher asks for more regular homework, more autonomy and more attention | School report | Short sessions, a low-effort start, and progress the parents can see |
| 5 | She prefers short exercises | Interview | One question per screen, sessions of 5 |
| 6 | Her favourite apps are Roblox, YouTube and Pinterest; she likes visual arts | Interview | Visual, customisable UI (themes, avatar, collectables), with rewards approved by the parents. It must not look like a worksheet |

## Decisions

- MVP subjects: **Maths + HGP** ([ADR-016](adr/016-mvp-subjects-maths-and-hgp.md)).
- First maths skills: fractions. They are the classic 5th-grade hurdle and the best fit for exact arithmetic. The order is confirmed against the corrected tests and the textbook index.
- Proposed success criteria, to agree with the parents: maths moves from below the pass mark to at least a pass at the next evaluation, and she practises alone about 3 times a week (from the app's logs). This is evidence of progress, not proof of cause.

## Open questions

- Are her maths errors conceptual or caused by distraction? The observation will decide this, and it changes the UX: explanation and graded steps for concept errors, short exercises with self-checks for distraction.
- Tablet: Android/HyperOS version, Google Play services, parental supervision, how the app will be installed.
- Time limits and rewards the parents accept.
