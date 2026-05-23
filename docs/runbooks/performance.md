# Performance Runbook

## Goal

Validate baseline p95 latency objective (`<= 300ms`) under a lightweight smoke workload.

## Preconditions

- Gateway running on `http://localhost:8080`
- `k6` installed locally

## Run

```bash
bash scripts/load/run_k6_smoke.sh
```

## Notes

- This is a smoke test, not a full capacity benchmark.
- Expand with authenticated and write-heavy scenarios next.
