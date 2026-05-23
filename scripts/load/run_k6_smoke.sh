#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

if ! command -v k6 >/dev/null 2>&1; then
  echo "k6 is not installed. Install from https://k6.io/docs/get-started/installation/"
  exit 1
fi

BASE_URL="$BASE_URL" k6 run tests/load/k6_smoke.js
