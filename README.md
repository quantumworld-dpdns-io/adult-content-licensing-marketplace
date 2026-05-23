# adult-content-licensing-marketplace

Adult content licensing marketplace with territory-aware licensing and AI-training policy controls.

## Current State

This repo now implements all roadmap phases at baseline depth, with runnable Go services and scaffolded multi-stack integrations.

## Quick Start

```bash
make init
make test
make run
```

Gateway endpoints:

- `GET /healthz`
- `GET /readyz`
- `GET /v1/licenses`
- `POST /v1/licenses`

## Key Paths

- `cmd/gateway`, `cmd/license-svc` — runnable Go services
- `db/migrations` — PostgreSQL schema migrations
- `api/openapi.yaml` — OpenAPI 3.1 baseline
- `deploy/docker` — local docker compose stack
- `infra/terraform`, `infra/pulumi` — IaC scaffolds
- `contracts` — Solidity baseline
- `rust`, `julia`, `qasm`, `federated` — advanced-stack scaffolds
- `security/vault` — Vault policy/config baseline
- `tests/robot` — Robot framework scaffold
- `observability` — Grafana + OTel baseline

## Plan Tracker

See `docs/plans.md`.

## License

[MIT](LICENSE)
