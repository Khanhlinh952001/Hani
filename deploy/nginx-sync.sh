#!/usr/bin/env bash
set -euo pipefail

SITES_AVAILABLE="/etc/nginx/sites-available"
SITES_ENABLED="/etc/nginx/sites-enabled"

DISABLE=(
  hanconnect.kr.conf
  be.hanconnect.kr.conf
)

KEEP=(
  ai-pulse.kr.conf
  hani.ai-pulse.kr.conf
  behani.ai-pulse.kr.conf
  be.ai-pulse.kr.conf
  memora.ai.kr.conf
  be.memora.ai.kr.conf
)

echo "==> Remove old hanconnect sites"
for f in "${DISABLE[@]}"; do
  sudo rm -f "${SITES_ENABLED}/${f}" "${SITES_AVAILABLE}/${f}"
done

echo "==> Enable sites"
for f in "${KEEP[@]}"; do
  [[ -f "${SITES_AVAILABLE}/${f}" ]] || { echo "MISSING ${SITES_AVAILABLE}/${f}" >&2; exit 1; }
  sudo ln -sf "${SITES_AVAILABLE}/${f}" "${SITES_ENABLED}/${f}"
  echo "  enabled: ${f}"
done

echo "==> Test & reload nginx"
sudo nginx -t
sudo systemctl reload nginx
echo "Done."
echo "  cert Hani FE:  sudo certbot --nginx -d hani.ai-pulse.kr"
echo "  cert Hani BE:  sudo certbot --nginx -d behani.ai-pulse.kr"
