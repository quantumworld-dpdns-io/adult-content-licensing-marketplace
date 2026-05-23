# Implementation Tracker (All Phases)

## Execution Mode

Assumed defaults are locked in `docs/spec/assumed-defaults-mvp.md` and treated as implementation constraints.

## Completed Highlights

- [x] Core services: gateway + license-svc baseline
- [x] RBAC/tier enforcement for API (`admin`, `regularuser`, `auditor`; tier gating)
- [x] Signed bearer token auth support (HMAC claims with expiry)
- [x] Dev token issuance endpoint for local integration testing
- [x] Crypto-only payment enforcement on license creation
- [x] Compliance and threat-model policy endpoints
- [x] OpenAPI updated for bearer auth and policy endpoints
- [x] Threat model doc grounded in STRIDE/PASTA/OWASP
- [x] Repository abstraction for license persistence
- [x] Configurable storage backend (`memory` default, `postgres` option)
- [x] SQL repository implementation scaffold for PostgreSQL
- [x] Multi-phase scaffolds for infra/contracts/rust/julia/quantum/federated/frontend/security/observability

## Next Build Depth

1. Wire a concrete Postgres SQL driver and add integration tests against dockerized Postgres.
2. Replace JSON territory storage with normalized writes to `license_territories`.
3. Add migrations runner command and startup migration guard.
4. Add immutable audit event stream for license create/update actions.
5. Add load test harness to validate p95 target assumptions.
