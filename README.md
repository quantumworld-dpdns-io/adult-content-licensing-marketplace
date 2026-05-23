# adult-content-licensing-marketplace

Adult content licensing marketplace with territory-aware licensing and policy constraints.

## Current Status

This repo now includes a runnable Go gateway service and CI baseline.

## Quick Start

```bash
git clone https://github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace.git
cd adult-content-licensing-marketplace

make init
make test
make run
```

Gateway defaults to `:8080` and exposes:

- `GET /healthz`
- `GET /readyz`

## Project Structure

```text
.
├── cmd/gateway/           # Gateway binary entrypoint
├── internal/config/       # Internal configuration loading
├── pkg/httpx/             # Shared HTTP handlers/router
├── docs/                  # Plans and architecture docs
├── tests/                 # Reserved for integration/e2e tests
└── .github/workflows/     # CI
```

## Roadmap

Implementation tracker: `docs/plans.md`

## License

[MIT](LICENSE)
