# M2 · Usable app — specification

**Goal:** the child practises fractions on her own tablet: she pairs the app once, picks a topic, does a session of 5 exercises with a maths keypad and gets feedback in Portuguese. Every answer is graded on the server by the M1 engine. No scanning, no AI, no database yet.

**Done when:**

- [ ] **She completes 3 sessions alone, on different days**, on her tablet. Evidence: the attempt log shows 3 complete sessions from her device on 3 different days, and her parents confirm she did them alone.
- [ ] The app installs and **updates without a cable** (tested before January 2027, when the author moves abroad).
- [ ] The API is served over HTTPS from an EU server; the app only works with a paired device.
- [ ] The session payload never contains answers, keys or misconceptions (tested).
- [ ] CI is green for backend and Android: build, lint, unit tests, contract drift check.
- [ ] Docs: README Evidence row, devlog, screen recording (demo data only) in `docs/media/`, ADRs, portfolio notes ([Definition of done](milestones.md#definition-of-done)).

**Gates before she uses the app** (from M0): written parental consent on file, and the tablet checklist done.

> **How we work:** each step ends with a small commit on a branch and a review before the next step starts. Every finished, green step is pushed.

---

## Decisions for M2

| Decision | Record |
|---|---|
| OpenAPI 3.1 contract is the source of truth; Go server types are generated with `oapi-codegen`, the Kotlin client is generated or checked against it | [ADR-017](adr/017-openapi-contract.md), now Accepted |
| Early EU hosting: one small VPS, a free subdomain, Caddy for HTTPS, the API in a Docker image built by CI | [ADR-021](adr/021-early-eu-hosting-and-device-pairing.md) |
| Device pairing with a single-use code and a per-device token; attempts in an append-only JSONL log | [ADR-021](adr/021-early-eu-hosting-and-device-pairing.md) |
| Exercises are referenced by `skill@version:level:seed`; the server regenerates them to grade, so answers never leave the server | [ADR-022](adr/022-stateless-exercise-references.md) |
| The APK is sideloaded from GitHub Releases and updates itself after the parents allow it once | [ADR-023](adr/023-sideloaded-apk-with-in-app-update.md) |

These pull a minimal version of M3 (device access, attempt log) and M6 (deploy) forward, because the done criterion needs the real tablet at home. M3 still replaces them with rotating tokens, Postgres and transactional idempotency; M6 still adds backups, observability and the threat model ([milestones.md](milestones.md)).

## Layout

```text
api/openapi.yaml                 the contract
backend/
  cmd/api/                       HTTP server
  cmd/admin/                     pairing codes and the session report (run on the server)
  cmd/practice/                  M1 terminal session (now uses internal/feedback)
  internal/feedback/             all Portuguese messages for the child, shared by CLI and API
  internal/httpapi/              generated types + handlers, auth, limits
  internal/devices/              pairing codes and device tokens (hashed), JSON file store
  internal/attemptlog/           append-only JSONL log, dedupe index, report
  Dockerfile                     static binary on a distroless image
deploy/
  compose.yml, Caddyfile         API + Caddy on the VPS
android/                         Kotlin + Compose app (io.github.whait666.tpcplus)
```

## API (v1)

All endpoints except `/healthz` and `/v1/pair` need `Authorization: Bearer <device token>`. Bodies are JSON, at most 4 KB; every request has a timeout.

| Method and path | Purpose |
|---|---|
| `GET /healthz` | Liveness, for Caddy and the deploy job |
| `POST /v1/pair` | `{code}` → `{device_token}`. Codes are single-use and expire after 15 minutes |
| `GET /v1/skills` | Skills with Portuguese titles and their levels |
| `POST /v1/sessions` | `{skill_id, level}` → `{session_id, exercises: [{ref, prompt, kind, options?}]}`, always 5 exercises |
| `POST /v1/sessions/{session_id}/answers` | `{position, ref, answer}` → `{outcome, expected, message, hint?}`. Repeating the same `(session_id, position)` returns the first result |
| `GET /v1/app/latest` | `{version_code, version_name, apk_url, sha256}` for the in-app update |

**Wire rules, each with a test at the API boundary:**

- `answer` is always the text the child entered. For a choice exercise it is the **number shown next to the option, from 1** (`"2"`), exactly what the grader expects.
- An exercise in the session payload has no `answer`, `correct` or `misconception` field.
- `ref` is `skill@version:level:seed`. A ref whose generator version no longer exists gets `410 Gone`, and the app starts a new session.
- `message` and `hint` come from `internal/feedback`, so the CLI and the app say the same thing.
- Errors use one shape: `{error: {code, message}}`, with `message` safe to show to the child.

## Attempt log

One JSON object per line, append-only, on a Docker volume:

```json
{"t":"2026-12-03T17:42:10Z","device":"d_7f3a","session":"s_91c2","position":2,"ref":"mat.fracoes_equivalentes@1:3:812","answer":"6/10","outcome":"correct","ms":8400}
```

No names and no free text beyond what she typed. `cmd/admin report` prints, per device and day, the sessions started and completed. This is the evidence for the done criterion.

## App

| Screen | What it does |
|---|---|
| Pairing | First launch only: enter the 6-character code; the token is stored with the Android Keystore |
| Home | Three topic cards with friendly Portuguese titles, a level picker 1–3 |
| Exercise | One exercise per screen, a progress row of 5 dots, a maths keypad (0–9, `/`, `,`, delete, OK) or 4 big option buttons |
| Feedback | Correct / almost (not simplified) / not yet, the right answer and the misconception hint |
| Summary | Score out of 5 and "Mais uma?" |

- **Look:** it must not look like a worksheet (M0 finding 6): a bright theme, large type and touch targets, small animations on correct answers. Themes, avatar and collectables are backlog, with rewards approved by her parents.
- **Architecture:** single activity, Compose Navigation, ViewModel + `StateFlow` (UDF), Hilt, Retrofit + OkHttp + kotlinx.serialization, DataStore for settings. Room arrives when there is offline data.
- **Network:** a friendly message and a retry button when offline; a retried answer reuses its `(session_id, position)`, so it is never counted twice.
- **Tests:** ViewModel unit tests with a fake API, Compose UI tests for the keypad and option buttons on the JVM (Robolectric), all in CI.

## Toolchain

Checked on 2026-11-20; step 1 confirms them in CI.

| Tool | Version |
|---|---|
| Android Gradle Plugin | 9.5.2 |
| Gradle | 9.9.2 |
| Kotlin | 2.4.21 |
| Compose BOM | 2026.11.01 |
| JDK in CI | 21 (Temurin) |
| oapi-codegen | v2.8.2 |
| openapi-generator | v7.27.0 |

This PC has no Android SDK yet, so Android builds are proven in CI first and in Android Studio once installed.

## Steps, in order

| # | Step | Done when |
|---|---|---|
| 1 | Android skeleton: Gradle project, one Compose screen, Android CI job (build, lint, unit tests). Install a debug APK on her tablet by cable | CI green; the app opens on her tablet; tablet checklist done |
| 2 | Contract: `api/openapi.yaml`, Go code generation, CI check that generated code is up to date | The spec lints; CI fails if the code drifts |
| 3 | Backend API: skills, sessions, answers, shared `internal/feedback`, exercise refs, attempt log with dedupe | Handler tests pass, including "no answers in the payload" and the choice wire format |
| 4 | Device pairing, auth middleware, body and rate limits, `cmd/admin pair` and `report` | Unpaired requests get 401; a code works once |
| 5 | Deploy: Dockerfile, image to GHCR from CI, VPS with Caddy and the DuckDNS subdomain, deploy job | `https://<name>.duckdns.org/healthz` answers from the tablet's browser |
| 6 | App: network layer, pairing screen, token storage | The tablet pairs with the real server |
| 7 | App: home, exercise, keypad, feedback and summary, wired to the API | A full session of 5 on the tablet; ViewModel and UI tests green |
| 8 | App: look and feel from M0 (theme, touch targets, animation) | Screenshots in demo mode; reviewed with her once |
| 9 | Release and update: signed APK on GitHub Releases from a tag, in-app update check | An update installs on her tablet **without a cable**, before January |
| 10 | Field test: first session with the author present (observation notes stay in `private/`), then 3 sessions alone on different days | `cmd/admin report` shows them; parents confirm |
| 11 | Docs: README evidence, devlog, screen recording, portfolio notes | Definition of done ticked |

## The author's to-dos

- Install **Android Studio** as a normal user (not Administrator).
- Create the VPS account (EU region, e.g. Hetzner CX22) with a working payment method, and a **DuckDNS** subdomain; add an SSH key.
- With the tablet: the tablet checklist from the M0 guide (Android/HyperOS version, Google Play services, developer options, "install unknown apps").
- Keep the parents' **written** consent on file before step 10.

## Risks

- **Sideloading rules.** Google is rolling out developer verification for apps installed outside the Play Store. Step 1 proves a sideloaded APK installs on her tablet; if that changes, Play internal testing is the fallback (one-time fee).
- **Xiaomi HyperOS** may need extra settings to install over USB or to keep the app from being killed. The tablet checklist records them.
- **One VPS, no backups yet.** The JSONL log is small; copy it off the server weekly until M6 adds real backups.
- **Kotlin is new.** Steps are small, and each one is reviewed before the next.

## Out of scope for M2

Scanning and OCR (M4), AI (M5), Postgres, rotating refresh tokens and transactional idempotency (M3), placement and adaptive levels (M3), parent views, avatars and rewards, offline practice.
