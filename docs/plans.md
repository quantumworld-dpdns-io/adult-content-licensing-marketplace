# Implementation Plan (Grounded to Current Repo)

This repository started as an empty scaffold while `docs/plans.md` listed a 138-item long-term roadmap.
This file now tracks what has actually been implemented and what is next.

## Scope Decision

- Implemented now: foundational platform work that makes the repo buildable, testable, and CI-ready.
- Deferred: large multi-language / multi-cloud / blockchain / ML integrations until the core services and data model stabilize.

## Completed in this implementation pass

### Phase A: Foundation

- [x] Add Go module (`go.mod`)
- [x] Add Makefile task runner (`init`, `fmt`, `lint`, `test`, `build`, `run`, `clean`)
- [x] Create backend-oriented directory structure:
  - `cmd/gateway`
  - `internal/config`
  - `pkg/httpx`
- [x] Add GitHub Actions CI for formatting, vet, tests, and build

### Phase B: First Running Service (Gateway)

- [x] Implement HTTP gateway entrypoint with graceful shutdown
- [x] Implement health endpoints:
  - `GET /healthz`
  - `GET /readyz`
- [x] Implement environment-based configuration (`PORT`, default `8080`)

### Phase C: Test Baseline

- [x] Add unit tests for config loading
- [x] Add handler tests for health/readiness endpoints

## Backlog (next priority order)

1. Add PostgreSQL schema + migrations for users/licenses/territories.
2. Add `license-svc` with CRUD endpoints and validation.
3. Introduce structured logging + request IDs + basic metrics.
4. Add OpenAPI 3.1 spec for gateway and services.
5. Add Dockerfiles + local `docker compose` stack.
6. Add auth baseline (JWT) and role model (creator/buyer/admin).
7. Add contract integration spike behind an interface (no chain lock-in yet).

## Out of Scope for this pass

The original long-horizon roadmap items (full Terraform/Pulumi estate, Solidity suite, Rust/Julia/Quantum/Federated stack, and all AI agent integrations) remain strategic backlog and are intentionally not claimed as implemented in this repository yet.
