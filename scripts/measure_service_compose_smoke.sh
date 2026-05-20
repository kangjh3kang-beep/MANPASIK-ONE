#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="${ROOT_DIR}/infrastructure/docker/docker-compose.measurement-smoke.yml"
PROJECT_NAME="${MANPASIK_MEASUREMENT_SMOKE_PROJECT:-manpasik-measure-smoke-$$}"

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
  echo "MEASUREMENT_COMPOSE_SMOKE_STATUS=blocked_docker_unavailable"
  echo "MEASUREMENT_COMPOSE_SMOKE_REASON=docker_cli_not_found_in_current_wsl_distro"
  exit 2
fi

if ! docker compose version >/dev/null 2>&1; then
  echo "MEASUREMENT_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable"
  exit 2
fi

GO_BIN="$(find_go)"
if [[ -z "${GO_BIN}" ]]; then
  echo "MEASUREMENT_COMPOSE_SMOKE_STATUS=blocked_go_unavailable"
  exit 2
fi

export MANPASIK_MEASUREMENT_SMOKE_GRPC_PORT="${MANPASIK_MEASUREMENT_SMOKE_GRPC_PORT:-$(reserve_port)}"
export MANPASIK_MEASUREMENT_SMOKE_HTTP_PORT="${MANPASIK_MEASUREMENT_SMOKE_HTTP_PORT:-$(reserve_port)}"

cleanup() {
  docker compose -p "${PROJECT_NAME}" -f "${COMPOSE_FILE}" down --remove-orphans -v >/dev/null 2>&1 || true
}
trap cleanup EXIT

cd "${ROOT_DIR}"

echo "MEASUREMENT_COMPOSE_SMOKE_PROJECT=${PROJECT_NAME}"
echo "MEASUREMENT_COMPOSE_SMOKE_GRPC_ADDR=127.0.0.1:${MANPASIK_MEASUREMENT_SMOKE_GRPC_PORT}"
echo "MEASUREMENT_COMPOSE_SMOKE_HTTP_ADDR=127.0.0.1:${MANPASIK_MEASUREMENT_SMOKE_HTTP_PORT}"

docker compose -p "${PROJECT_NAME}" -f "${COMPOSE_FILE}" up -d --build measurement-service

MANPASIK_MEASUREMENT_SERVICE_SMOKE_ADDR="127.0.0.1:${MANPASIK_MEASUREMENT_SMOKE_GRPC_PORT}" \
MANPASIK_MEASUREMENT_SERVICE_SMOKE_HTTP_ADDR="127.0.0.1:${MANPASIK_MEASUREMENT_SMOKE_HTTP_PORT}" \
MANPASIK_MEASUREMENT_SERVICE_SMOKE_VERSION="compose-smoke" \
  "${GO_BIN}" test -v ./backend/services/measurement-service/cmd -run TestMeasurementServiceExternalEndpointSmoke -count=1

echo "MEASUREMENT_COMPOSE_SMOKE_STATUS=passed"
