#!/usr/bin/env bash
# Build Next.js frontend and (re)start with PM2.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
FE="$ROOT/fe"

cd "$FE"

if command -v pnpm >/dev/null 2>&1; then
  PKG=pnpm
elif command -v npm >/dev/null 2>&1; then
  PKG=npm
else
  echo "Need pnpm or npm" >&2
  exit 1
fi

echo "→ Installing dependencies ($PKG)…"
if [[ -f pnpm-lock.yaml ]] && [[ "$PKG" == pnpm ]]; then
  pnpm install --frozen-lockfile
else
  $PKG install
fi

echo "→ Building frontend…"
$PKG run build

cd "$ROOT"

if pm2 describe hani-fe >/dev/null 2>&1; then
  echo "→ Reloading hani-fe…"
  pm2 reload ecosystem.config.cjs --only hani-fe --update-env
else
  echo "→ Starting hani-fe…"
  pm2 start ecosystem.config.cjs --only hani-fe --env production
fi

pm2 save 2>/dev/null || true
echo "✓ Done. Logs: pm2 logs hani-fe"
