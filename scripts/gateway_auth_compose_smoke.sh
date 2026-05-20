#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="${ROOT_DIR}/infrastructure/docker/docker-compose.gateway-auth-smoke.yml"
PROJECT_NAME="${MANPASIK_GATEWAY_AUTH_SMOKE_PROJECT:-manpasik-gateway-auth-smoke-$$}"

find_go() {
  if [[ -n "${MANPASIK_GO_BINARY:-}" ]]; then
    printf '%s\n' "${MANPASIK_GO_BINARY}"
    return 0
  fi
  if [[ -x "/home/kangjh3kang/sdk/go-go1.26.2/bin/go" ]]; then
    printf '%s\n' "/home/kangjh3kang/sdk/go-go1.26.2/bin/go"
    return 0
  fi
  command -v go
}

reserve_port() {
  python3 - <<'PY'
import socket

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
    sock.bind(("127.0.0.1", 0))
    print(sock.getsockname()[1])
PY
}

if ! command -v docker >/dev/null 2>&1; then
  echo "GATEWAY_AUTH_COMPOSE_SMOKE_STATUS=blocked_docker_unavailable"
  echo "GATEWAY_AUTH_COMPOSE_SMOKE_REASON=docker_cli_not_found_in_current_wsl_distro"
  exit 2
fi

if ! docker compose version >/dev/null 2>&1; then
  echo "GATEWAY_AUTH_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable"
  exit 2
fi

GO_BIN="$(find_go)"
if [[ -z "${GO_BIN}" ]]; then
  echo "GATEWAY_AUTH_COMPOSE_SMOKE_STATUS=blocked_go_unavailable"
  exit 2
fi

export MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_PORT="${MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_PORT:-$(reserve_port)}"

cleanup() {
  docker compose -p "${PROJECT_NAME}" -f "${COMPOSE_FILE}" down --remove-orphans -v >/dev/null 2>&1 || true
}
trap cleanup EXIT

cd "${ROOT_DIR}"

echo "GATEWAY_AUTH_COMPOSE_SMOKE_PROJECT=${PROJECT_NAME}"
echo "GATEWAY_AUTH_COMPOSE_SMOKE_HTTP_ADDR=127.0.0.1:${MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_PORT}"

docker compose -p "${PROJECT_NAME}" -f "${COMPOSE_FILE}" up -d --build auth-service gateway

MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_ADDR="127.0.0.1:${MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_PORT}" \
  "${GO_BIN}" test -v ./backend/services/gateway/cmd -run TestGatewayAuthExternalEndpointSmoke -count=1

echo "GATEWAY_AUTH_COMPOSE_SMOKE_STATUS=passed"
