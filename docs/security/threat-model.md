# Threat Model (MVP)

## Methods

- STRIDE for threat categorization.
- PASTA for process alignment and attack simulation planning.
- OWASP Top 10 2021 for web/API coverage.

## Priority Risks

1. API abuse
2. Data breach
3. Fraud

## Baseline Mitigations

- Role/tier based authorization checks for protected APIs.
- Request identity propagation via headers (temporary until JWT is added).
- Audit-friendly response model for policy endpoints.
- Crypto-only payment guardrail at license creation.
- Input validation and safe defaults on API handlers.

## Residual Risks

- No persistent audit log store yet.
- No WAF enforcement in runtime path yet.
- No production-grade KYC vendor integration yet.
