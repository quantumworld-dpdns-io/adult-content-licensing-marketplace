# Assumed Defaults MVP Spec (Locked)

Date locked: 2026-05-23

## Product Scope

MVP is a production-capable product baseline with:

- Go API services (gateway + license service behavior in gateway)
- Crypto-only purchase policy
- Tiered regular user model
- Admin and third-party auditor access model
- Compliance policy artifacts and enforcement hooks

## Fixed Assumptions

- Cloud strategy: AWS-first for launch architecture. GCP/Azure kept as future portability target.
- Performance target: p95 API latency <= 300ms for core endpoints under baseline load.
- Capacity target: 500 peak RPS, 20k DAU, 2k peak concurrency.
- Payment model: crypto only.
- On-chain licensing: optional in MVP; testnet-compatible configuration.
- Compliance implementation order: age verification + DMCA workflow + KYC-lite checks.

## RBAC and Tier Model

Roles:

- admin
- regularuser
- auditor

Regular user tiers:

- alpha
- beta
- vip1
- vip2
- vip3

Permissions:

- admin: full read/write on licenses and policy endpoints.
- regularuser:
  - alpha/beta: read licenses only.
  - vip1/vip2/vip3: read + create licenses.
- auditor: read-only access + policy read access.

## Threat Model Priorities

Applied methods:

- STRIDE
- PASTA (high-level stage mapping)
- OWASP Top 10 (web/API)

Priority risks (MVP top 3):

1. API abuse
2. Data breach
3. Fraud

## Legal and Content Policy Position

- Allowed content: legal content under applicable US/EU law and platform policy.
- Prohibited content: illegal content under applicable US/EU law and platform policy.
- Enforcement basis: contractual terms + technical controls (logging, flags, policy checks, auditability).

## API Acceptance Criteria (MVP)

- `GET /healthz` and `GET /readyz` return 200.
- `GET /v1/licenses` requires valid role context.
- `POST /v1/licenses` enforces role+tier matrix and crypto-only currency rule.
- `GET /v1/policy/compliance` returns compliance stance and controls.
- `GET /v1/policy/threat-model` returns active threat modeling methods and priorities.
