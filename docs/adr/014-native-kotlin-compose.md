# ADR-014: Native Kotlin + Jetpack Compose

**Status:** Accepted · 2026-11-20

## Context

The client runs on one Android tablet (Xiaomi Redmi Pad 2 Pro). The author's goal includes market signal.

## Decision

A native Android app in Kotlin with Jetpack Compose, UDF (ViewModel + StateFlow), Hilt, Room, Retrofit, ML Kit.

## Alternatives considered

- **Expo / React Native** — reuses React skills, but native Android has a larger job pool in the target market and direct ML Kit access.

## Consequences

A new language to learn. M1 is backend-only, so Kotlin arrives in M2 on a working engine.
