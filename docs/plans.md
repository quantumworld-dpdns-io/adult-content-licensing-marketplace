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
- [x] Runtime schema bootstrap + `migrate` command
- [x] Durable audit store (memory/postgres) + protected audit endpoint with query filters/pagination
- [x] Request ID middleware + structured access logging
- [x] In-process latency metrics snapshot endpoint (`/v1/metrics/latency`)
- [x] k6 load smoke scaffold for p95 baseline verification
- [x] Multi-phase scaffolds for infra/contracts/rust/julia/quantum/federated/frontend/security/observability

## Next Build Depth

1. Persist audit events to PostgreSQL (not only in-memory stream).
2. Emit metrics to Prometheus/OpenTelemetry instead of in-process snapshot only.
3. Add distributed tracing context propagation and span attributes.
4. Add endpoint-level authorization integration tests against running service.
5. Add load test harness to validate p95 target assumptions.
