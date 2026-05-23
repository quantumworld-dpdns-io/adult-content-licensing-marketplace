# Architecture Overview

## Runtime Core

- `gateway` (Go): API ingress and health/readiness endpoints.
- `license-svc` (Go): baseline license CRUD API.
- PostgreSQL schema in `db/migrations`.

## Platform

- `deploy/docker/docker-compose.yml`: local infra (Postgres/Redis/Kafka + gateway).
- `deploy/nginx/nginx.conf`: reverse proxy + rate limit baseline.
- `infra/terraform` and `infra/pulumi`: IaC scaffolds for environment expansion.

## Advanced Stack Scaffolds

- Smart contracts: `contracts/`
- Rust services/WASM: `rust/`
- Julia analytics: `julia/`
- Quantum experiments: `qasm/`
- Federated learning seed: `federated/`
- MCP/agent integration seed: `mcp/`
- Vault baseline policies/config: `security/vault/`
- Observability baseline: `observability/`
