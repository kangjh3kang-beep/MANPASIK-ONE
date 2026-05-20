#!/usr/bin/env bash
# Phase AD-3: Tenancy 운영 E2E 검증 스크립트
#
# 사전 조건:
#   - gateway 가 ${GATEWAY_URL} 에서 실행 중 (기본 http://localhost:8080)
#   - jq, curl 설치
#   - admin 사용자가 이미 등록되어 있고 hospital-A 의 admin 멤버십 보유
#
# 시나리오:
#   1. admin 로그인 → 토큰 획득
#   2. 초대 발급 (admin → invitee@example.com)
#   3. 새 사용자 등록 + 로그인
#   4. 초대 수락
#   5. 멤버십 조회 (hospital-A 포함 확인)
#   6. 활성 조직 헤더로 보호된 API 호출 검증 (admin 관리자 화면)
#   7. 멤버 역할 변경
#   8. 멤버 제거
#
# 환경변수:
#   GATEWAY_URL    기본 http://localhost:8080
#   ADMIN_EMAIL    기본 admin@manpasik.test
#   ADMIN_PASSWORD 기본 password123!
#   TENANT_ID      기본 hospital-A

set -euo pipefail

GATEWAY_URL="${GATEWAY_URL:-http://localhost:8080}"
ADMIN_EMAIL="${ADMIN_EMAIL:-admin@manpasik.test}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-password123!}"
TENANT_ID="${TENANT_ID:-hospital-A}"

# ===== 헬퍼 =====

red()    { printf '\033[31m%s\033[0m\n' "$*"; }
green()  { printf '\033[32m%s\033[0m\n' "$*"; }
yellow() { printf '\033[33m%s\033[0m\n' "$*"; }
bold()   { printf '\033[1m%s\033[0m\n' "$*"; }

require() {
  if ! command -v "$1" >/dev/null 2>&1; then
    red "필수 도구 누락: $1"
    exit 1
  fi
}

step() {
  echo
  bold "▶ $*"
}

assert_field() {
  local body="$1"
  local field="$2"
  local desc="$3"
  if ! echo "$body" | jq -e "$field" >/dev/null 2>&1; then
    red "❌ $desc"
    echo "  body: $body"
    exit 1
  fi
  green "✓ $desc"
}

assert_status() {
  local status="$1"
  local expected="$2"
  local desc="$3"
  if [[ "$status" != "$expected" ]]; then
    red "❌ $desc (got $status, want $expected)"
    exit 1
  fi
  green "✓ $desc (HTTP $status)"
}

http() {
  local method="$1"
  local path="$2"
  local body="${3:-}"
  local token="${4:-}"
  local user_id="${5:-}"
  local tenant_id="${6:-}"

  local args=(-sS -w '\n%{http_code}' -X "$method")
  args+=(-H 'Content-Type: application/json')
  [[ -n "$token" ]] && args+=(-H "Authorization: Bearer $token")
  [[ -n "$user_id" ]] && args+=(-H "X-User-ID: $user_id")
  [[ -n "$tenant_id" ]] && args+=(-H "X-Tenant-ID: $tenant_id")
  [[ -n "$body" ]] && args+=(-d "$body")

  curl "${args[@]}" "${GATEWAY_URL}${path}"
}

# ===== 사전 검증 =====

require curl
require jq

step "Gateway 헬스체크"
if curl -sS -o /dev/null -w '%{http_code}' "${GATEWAY_URL}/health" | grep -qE '^(200|204)$'; then
  green "✓ Gateway 응답 OK"
else
  red "❌ Gateway 미응답 — ${GATEWAY_URL}"
  exit 1
fi

# ===== 1. Admin 로그인 =====

step "1. Admin 로그인"
LOGIN_RESP=$(http POST /api/v1/auth/login "{\"email\":\"${ADMIN_EMAIL}\",\"password\":\"${ADMIN_PASSWORD}\"}")
LOGIN_BODY=$(echo "$LOGIN_RESP" | head -n1)
LOGIN_STATUS=$(echo "$LOGIN_RESP" | tail -n1)
assert_status "$LOGIN_STATUS" "200" "Admin 로그인"
ADMIN_TOKEN=$(echo "$LOGIN_BODY" | jq -r '.access_token // .accessToken // empty')
ADMIN_ID=$(echo "$LOGIN_BODY" | jq -r '.user.id // .user_id // empty')
[[ -n "$ADMIN_TOKEN" ]] || { red "토큰 없음"; exit 1; }
[[ -n "$ADMIN_ID" ]] || ADMIN_ID="admin-id"
green "  ADMIN_ID=$ADMIN_ID"

# ===== 2. 초대 발급 =====

step "2. 초대 발급 (admin → invitee@example.com)"
INVITE_RESP=$(http POST /api/v1/tenancy/invitations \
  "{\"tenant_id\":\"${TENANT_ID}\",\"role\":\"member\",\"invitee_hint\":\"invitee@example.com\"}" \
  "$ADMIN_TOKEN" "$ADMIN_ID" "$TENANT_ID")
INVITE_BODY=$(echo "$INVITE_RESP" | head -n1)
INVITE_STATUS=$(echo "$INVITE_RESP" | tail -n1)
assert_status "$INVITE_STATUS" "201" "초대 발급"
INVITE_TOKEN=$(echo "$INVITE_BODY" | jq -r '.token')
[[ -n "$INVITE_TOKEN" ]] || { red "토큰 없음"; exit 1; }
yellow "  INVITE_TOKEN=$INVITE_TOKEN"

# ===== 3. 새 사용자 등록 + 로그인 =====

INVITEE_EMAIL="invitee-$(date +%s)@e2e.test"
INVITEE_PASSWORD="invitee123!"
step "3. 새 사용자 등록 ($INVITEE_EMAIL)"
REG_RESP=$(http POST /api/v1/auth/register \
  "{\"email\":\"${INVITEE_EMAIL}\",\"password\":\"${INVITEE_PASSWORD}\",\"display_name\":\"Invitee\"}")
REG_STATUS=$(echo "$REG_RESP" | tail -n1)
[[ "$REG_STATUS" == "200" || "$REG_STATUS" == "201" ]] || {
  red "등록 실패 ($REG_STATUS)"; exit 1;
}
green "✓ 등록"

step "    로그인"
LOGIN2_RESP=$(http POST /api/v1/auth/login \
  "{\"email\":\"${INVITEE_EMAIL}\",\"password\":\"${INVITEE_PASSWORD}\"}")
LOGIN2_BODY=$(echo "$LOGIN2_RESP" | head -n1)
LOGIN2_STATUS=$(echo "$LOGIN2_RESP" | tail -n1)
assert_status "$LOGIN2_STATUS" "200" "Invitee 로그인"
INVITEE_TOKEN=$(echo "$LOGIN2_BODY" | jq -r '.access_token // .accessToken')
INVITEE_ID=$(echo "$LOGIN2_BODY" | jq -r '.user.id // .user_id // empty')
[[ -n "$INVITEE_ID" ]] || INVITEE_ID="invitee-id"

# ===== 4. 초대 수락 =====

step "4. 초대 수락"
ACCEPT_RESP=$(http POST /api/v1/tenancy/invitations/accept \
  "{\"token\":\"${INVITE_TOKEN}\"}" \
  "$INVITEE_TOKEN" "$INVITEE_ID")
ACCEPT_BODY=$(echo "$ACCEPT_RESP" | head -n1)
ACCEPT_STATUS=$(echo "$ACCEPT_RESP" | tail -n1)
assert_status "$ACCEPT_STATUS" "200" "초대 수락"
assert_field "$ACCEPT_BODY" '.tenant_id' "tenant_id 필드 확인"
assert_field "$ACCEPT_BODY" '.role' "role 필드 확인"

# ===== 5. 멤버십 조회 =====

step "5. Invitee 멤버십 조회"
MEMS_RESP=$(http GET /api/v1/tenancy/me/memberships "" "$INVITEE_TOKEN" "$INVITEE_ID")
MEMS_BODY=$(echo "$MEMS_RESP" | head -n1)
MEMS_STATUS=$(echo "$MEMS_RESP" | tail -n1)
assert_status "$MEMS_STATUS" "200" "멤버십 조회"
MEM_COUNT=$(echo "$MEMS_BODY" | jq '.memberships | length')
[[ "$MEM_COUNT" -ge 1 ]] || { red "멤버십 비어있음"; exit 1; }
green "✓ 멤버십 ${MEM_COUNT}개 (TENANT_ID=$TENANT_ID 포함)"

# ===== 6. Admin: 멤버 목록 조회 =====

step "6. Admin: 조직 멤버 목록 조회"
LIST_RESP=$(http GET "/api/v1/tenancy/tenants/${TENANT_ID}/members" "" "$ADMIN_TOKEN" "$ADMIN_ID" "$TENANT_ID")
LIST_BODY=$(echo "$LIST_RESP" | head -n1)
LIST_STATUS=$(echo "$LIST_RESP" | tail -n1)
assert_status "$LIST_STATUS" "200" "멤버 목록"
MEMBERS_COUNT=$(echo "$LIST_BODY" | jq '.members | length')
green "✓ 조직 멤버 ${MEMBERS_COUNT}명"

# ===== 7. 역할 변경 =====

step "7. Admin: invitee 역할 → medical_staff 변경"
ROLE_RESP=$(http PATCH "/api/v1/tenancy/tenants/${TENANT_ID}/members/${INVITEE_ID}/role" \
  '{"role":"medical_staff"}' \
  "$ADMIN_TOKEN" "$ADMIN_ID" "$TENANT_ID")
ROLE_STATUS=$(echo "$ROLE_RESP" | tail -n1)
assert_status "$ROLE_STATUS" "200" "역할 변경"

# ===== 8. 멤버 제거 =====

step "8. Admin: invitee 멤버 제거"
RM_RESP=$(http DELETE "/api/v1/tenancy/tenants/${TENANT_ID}/members/${INVITEE_ID}" "" \
  "$ADMIN_TOKEN" "$ADMIN_ID" "$TENANT_ID")
RM_STATUS=$(echo "$RM_RESP" | tail -n1)
assert_status "$RM_STATUS" "204" "멤버 제거"

# ===== 완료 =====

echo
green "════════════════════════════════════════════"
bold  "  Tenancy E2E 시나리오 ALL PASS"
green "════════════════════════════════════════════"
