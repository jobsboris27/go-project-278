#!/usr/bin/env bash
set -euo pipefail

echo "[run.sh] Starting service"

echo "[run.sh] Running DB migrations"
if ! goose -dir ./db/migrations postgres "${DATABASE_URL}" up; then
    echo "[run.sh] Warning: migration failed, continuing anyway"
fi

echo "[run.sh] Starting Caddy"
caddy run --config /etc/caddy/Caddyfile &

echo "[run.sh] Starting Go app"
exec /app/bin/app