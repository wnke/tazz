# AGENTS

This document describes on how to work on this project.

## 1) Code architecture

Hexagonal architecture (ports and adapters) keeps business logic independent from external systems.

- **Core (application/use cases):** business workflows and rules.
- **Ports (interfaces):** contracts the core uses or exposes.
- **Adapters:** implementations around the core.

Rule of thumb: **business logic should not depend on frameworks, databases, HTTP, or vendor SDKs**.

---

## 2) Project Structure

This project is organized to reflect hexagonal architecture:

```text
cmd/
  tazz/
    main.go                # composition root (wiring + startup)

internal/
  domain/                  # shared contracts between consumers and providers
    core/                  # structs/error codes shared by core <-> driver adapters
    driven/                # structs/error codes shared by core <-> driven adapters
  core/                    # use cases and orchestration
  adapter/
    driver/                # inbound adapters (HTTP, gRPC, CLI, consumers)
    driven/                # outbound adapters (db, cache, external APIs)
```

Import boundaries:

- Driver adapters may import shared structs/error codes from `internal/domain/core`, but must never import from `internal/core`.
- Core may import shared structs/error codes for outbound integrations from `internal/domain/driven`, while driven adapters implement those contracts.

### Driver vs Driven

- **Driver adapter (inbound):** drives the app by invoking use cases (HTTP handlers, gRPC handlers, CLI commands, message consumers).
- **Driven adapter (outbound):** is driven by the core to reach external systems (repositories, cache clients, API clients, publishers, storage).

---

## 3) Interface Placement and Ports

Define interfaces where they are consumed (consumer-owned interfaces), not in a global interface package, as per the golang idiom.

- Core defines outbound ports it needs.
- Driver adapters define service interfaces they call.
- Driven adapters implement core-defined ports.

This keeps interfaces small, focused, and testable.

---

## 4) Dependency Injection and Wiring

Use constructor-based DI and keep wiring in `cmd/.../main.go`.

### Why

- Decouples construction from behavior.
- Improves testability via mocks/fakes.
- Keeps dependencies one-directional.

### Typical Wiring Flow

1. Build infrastructure clients (DB, HTTP, Kafka, etc.).
2. Build driven adapters from those clients.
3. Build core services with required interfaces.
4. Build driver adapters with core services.
5. Register routes/consumers and start runtime.

---

## 5) Constructor and Options Pattern

Use explicit constructors with mandatory deps as positional args and optional deps as functional options.

- `NewX(requiredDeps..., WithLogger(...), WithTimeout(...))`
- Validate mandatory dependencies early.
- Fail fast on invalid configuration.

### Code Generation for Options

Use `go generate` to reduce boilerplate and enforce consistency.

Example in this project: `options-gen` generating `NewOptions` and `WithXxx` helpers.

---

## 6) Testing Strategy

- **Core tests:** mock outbound ports; validate business behavior and edge cases.
- **Driven adapter tests:** verify protocol/SDK translation and error mapping.
- **Driver adapter tests:** verify request/response mapping and status/code behavior.
- Prefer package-local mocks/fakes generated from interfaces.

Goal: fast unit tests for core, targeted integration tests for adapters.

---

## 7) Error Handling Conventions

- Define domain/application errors with business meaning.
- Wrap errors with context at adapter boundaries (`fmt.Errorf(... %w ...)`).
- Map domain errors to transport-layer errors in driver adapters.
- Keep semantics stable across transports.

---

## 8) Go Best Practices in This Project

- Use structured logging via `log/slog`.
- Manage config through env vars aligned with Twelve-Factor principles.
- Keep packages small and focused.
- Favor readable, idiomatic, simple Go.
- Handle concurrency with explicit synchronization and robust error handling.

---

## 9) Feature Delivery Checklist

When adding a new feature:

1. Define or extend domain types and errors.
2. Add/update use-case logic in `internal/core`.
3. Define/refine interfaces at the consumer side.
4. Implement/update driven adapters.
5. Implement/update driver adapters.
6. Wire dependencies in `cmd/.../main.go`.
7. Add tests (core first, adapters second).
8. Regenerate code and run lint/test tasks.

---

## 10) Common Pitfalls (To be avoided!)

- Importing framework/vendor packages directly in core.
- Oversized god interfaces.
- Shared global interface packages.
- Hidden required dependencies in globals/singletons.
- Transport concerns leaking into core business logic.

---

## 11) Minimal Design Rule

If unsure, default to this:

> Keep business logic pure and injectable; keep external concerns in adapters; wire everything in one place.

---

## 12) Quality Gate

All changes must pass the linter and tests (`justfile check`) and the tests before being considered complete.
