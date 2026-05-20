#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

failures=()

check_absent() {
  local file="$1"
  local pattern="$2"
  local label="$3"
  if [[ -f "$file" ]] && grep -Eq "$pattern" "$file"; then
    failures+=("$file contains forbidden production pattern: $label")
  fi
}

check_present() {
  local file="$1"
  local pattern="$2"
  local label="$3"
  if [[ ! -f "$file" ]] || ! grep -Eq "$pattern" "$file"; then
    failures+=("$file is missing required production gate: $label")
  fi
}

production_files=(
  "infrastructure/kubernetes/base/config/configmap.yaml"
  "infrastructure/kubernetes/overlays/production/config-patch.yaml"
)

for file in "${production_files[@]}"; do
  check_absent "$file" 'dev-secret-change-in-production|manpasik_dev|minioadmin' "default development secret"
  check_absent "$file" 'KEYCLOAK_ADMIN_PASSWORD:\s*"admin"|KEYCLOAK_ADMIN_PASSWORD=admin' "default Keycloak admin password"
  check_absent "$file" 'xpack\.security\.enabled:\s*"false"|xpack\.security\.enabled=false' "Elasticsearch security disabled"
  check_absent "$file" 'DB_SSLMODE:\s*"disable"|DB_SSLMODE=disable' "database SSL disabled"
done

check_present \
  "infrastructure/kubernetes/base/config/configmap.yaml" \
  'MAX_FINGERPRINT_DIMENSION:\s*"1792"' \
  "MAX_FINGERPRINT_DIMENSION 1792"

check_present \
  "infrastructure/kubernetes/base/config/configmap.yaml" \
  'DEFAULT_ALPHA:\s*"0\.98"' \
  "DEFAULT_ALPHA 0.98"

check_present \
  "infrastructure/kubernetes/overlays/production/config-patch.yaml" \
  'DB_SSLMODE:\s*"require"' \
  "production DB_SSLMODE require"

check_present \
  "infrastructure/kubernetes/overlays/production/config-patch.yaml" \
  'TENANCY_ENFORCED:\s*"true"' \
  "production tenancy enforcement"

if (( ${#failures[@]} > 0 )); then
  for failure in "${failures[@]}"; do
    echo "SECURITY_RELEASE_GATE_FAIL: $failure" >&2
  done
  exit 1
fi

echo "SECURITY_RELEASE_GATE_PASS"
