# ADR-023: Sideloaded APK with in-app update

**Status:** Accepted · 2026-11-20

## Context

The app has one user, on one tablet. After January 2027 the author cannot reach the tablet in person, so installing an update must not need a cable or a computer.

## Decision

- Release APKs are signed in CI from a version tag and published on **GitHub Releases**. The signing key is kept offline by the author and as a CI secret, never in the repository.
- On start, the app asks `GET /v1/app/latest`. When a newer version exists it offers an update, downloads the APK, checks its SHA-256 and hands it to the system installer.
- The parents allow "install unknown apps" for TPC+ once; after that each update is two taps.
- The full update path is tested on her tablet before January 2027.

## Alternatives considered

- **Google Play internal testing** — automatic updates, but a one-time fee, account verification and review delays for a single user. It is the fallback if sideloading becomes harder (Google is rolling out developer verification for sideloaded apps).
- **Update by cable** — impossible after the move.

## Consequences

The release APK is public, so it must contain no secrets: the device token comes from pairing ([ADR-021](021-early-eu-hosting-and-device-pairing.md)). The update check is one more endpoint, and the installer flow needs a permission and a test on the real tablet.
