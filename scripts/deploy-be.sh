#!/usr/bin/env bash
# Build Go API and (re)start with PM2.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BE="$ROOT/be"

cd "$BE"

if [[ ! -f .env ]]; then
  echo "Warning: be/.env not found — copy from docs/DEPLOY.md" >&2
fi

echo "→ Building API…"
go build -o bin/api ./cmd/api

cd "$ROOT"

if pm2 describe hani-be >/dev/null 2>&1; then
  echo "→ Reloading hani-be…"
  pm2 reload ecosystem.config.cjs --only hani-be --update-env
else
  echo "→ Starting hani-be…"
  pm2 start ecosystem.config.cjs --only hani-be --env production
fi

pm2 save 2>/dev/null || true
echo "✓ Done. Logs: pm2 logs hani-be"
