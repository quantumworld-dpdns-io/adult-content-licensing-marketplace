# adult-content-licensing-marketplace

Adult content licensing marketplace with territory-aware licensing and AI-training policy controls.

## Locked MVP Defaults (2026-05-23)

- AWS-first architecture
- p95 <= 300ms target for core API
- 500 peak RPS, 20k DAU, 2k peak concurrency
- Crypto-only payments (USDC/USDT/ETH)
- Roles: `admin`, `regularuser`, `auditor`
- Tiers: `alpha`, `beta`, `vip1`, `vip2`, `vip3`
- Compliance emphasis: KYC/AML, age verification, GDPR/CCPA, 2257, DMCA
- Threat modeling references: STRIDE, PASTA, OWASP Top 10

Full spec: `docs/spec/assumed-defaults-mvp.md`

## Quick Start

```bash
make init
make test
make run
```

## Current API

- `GET /healthz`
- `GET /readyz`
- `GET /v1/licenses` (requires `X-Role`)
- `POST /v1/licenses` (requires role/tier and crypto currency)
- `GET /v1/policy/compliance` (admin/auditor)
- `GET /v1/policy/threat-model` (admin/auditor)

OpenAPI: `api/openapi.yaml`

## License

[MIT](LICENSE)
