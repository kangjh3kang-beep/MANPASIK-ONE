#!/usr/bin/env bash
# ManPaSik 의존성 보안 스캔 스크립트
# OWASP A06: Vulnerable and Outdated Components
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
REPORT_DIR="$SCRIPT_DIR/reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

mkdir -p "$REPORT_DIR"

echo "=== ManPaSik 의존성 보안 스캔 ==="
echo "시작 시각: $(date)"
echo "프로젝트 루트: $PROJECT_ROOT"
echo ""

TOTAL_VULNS=0
SCAN_RESULTS=""

# ─── Go 모듈 취약점 스캔 ───
scan_go_modules() {
    echo "▶ [1/4] Go 모듈 취약점 스캔..."

    GO_SERVICES=$(find "$PROJECT_ROOT/backend/services" -name "go.mod" -type f 2>/dev/null || true)

    if [ -z "$GO_SERVICES" ]; then
        echo "  ⚠ Go 서비스를 찾을 수 없습니다"
        return
    fi

    for mod in $GO_SERVICES; do
        service_dir=$(dirname "$mod")
        service_name=$(basename "$service_dir")
        echo "  스캔: $service_name"

        if command -v govulncheck &> /dev/null; then
            govulncheck -C "$service_dir" ./... >> "$REPORT_DIR/go_vulns_${TIMESTAMP}.txt" 2>&1 || true
        else
            # govulncheck가 없으면 go list로 대체
            (cd "$service_dir" && go list -m -json all 2>/dev/null | grep -i "vuln\|deprecated" >> "$REPORT_DIR/go_deps_${TIMESTAMP}.txt" 2>&1) || true
        fi
    done

    echo "  ✓ Go 모듈 스캔 완료"
}

# ─── Flutter/Dart 패키지 스캔 ───
scan_flutter_packages() {
    echo "▶ [2/4] Flutter/Dart 패키지 스캔..."

    PUBSPEC_FILES=$(find "$PROJECT_ROOT/frontend" -name "pubspec.yaml" -type f 2>/dev/null || true)

    if [ -z "$PUBSPEC_FILES" ]; then
        echo "  ⚠ Flutter 프로젝트를 찾을 수 없습니다"
        return
    fi

    for pubspec in $PUBSPEC_FILES; do
        app_dir=$(dirname "$pubspec")
        app_name=$(basename "$app_dir")
        echo "  스캔: $app_name"

        # 의존성 목록 추출
        if command -v dart &> /dev/null; then
            (cd "$app_dir" && dart pub outdated --json >> "$REPORT_DIR/dart_outdated_${TIMESTAMP}.json" 2>&1) || true
        fi

        # pubspec.yaml에서 직접 의존성 버전 확인
        echo "  의존성 목록 ($app_name):" >> "$REPORT_DIR/dart_deps_${TIMESTAMP}.txt"
        grep -E '^\s+\w+:' "$pubspec" >> "$REPORT_DIR/dart_deps_${TIMESTAMP}.txt" 2>/dev/null || true
    done

    echo "  ✓ Flutter 패키지 스캔 완료"
}

# ─── Rust crate 스캔 ───
scan_rust_crates() {
    echo "▶ [3/4] Rust crate 취약점 스캔..."

    CARGO_FILES=$(find "$PROJECT_ROOT" -name "Cargo.toml" -maxdepth 3 -type f 2>/dev/null || true)

    if [ -z "$CARGO_FILES" ]; then
        echo "  ⚠ Rust 프로젝트를 찾을 수 없습니다"
        return
    fi

    for cargo in $CARGO_FILES; do
        crate_dir=$(dirname "$cargo")
        crate_name=$(basename "$crate_dir")
        echo "  스캔: $crate_name"

        if command -v cargo-audit &> /dev/null; then
            (cd "$crate_dir" && cargo audit --json >> "$REPORT_DIR/rust_audit_${TIMESTAMP}.json" 2>&1) || true
        else
            echo "  ⚠ cargo-audit 미설치 (cargo install cargo-audit)" >> "$REPORT_DIR/rust_audit_${TIMESTAMP}.txt"
            # Cargo.lock에서 직접 의존성 확인
            if [ -f "$crate_dir/Cargo.lock" ]; then
                grep -E '^\[\[package\]\]' "$crate_dir/Cargo.lock" | wc -l >> "$REPORT_DIR/rust_deps_${TIMESTAMP}.txt"
            fi
        fi
    done

    echo "  ✓ Rust crate 스캔 완료"
}

# ─── Docker 이미지 스캔 ───
scan_docker_images() {
    echo "▶ [4/4] Docker 이미지 보안 스캔..."

    DOCKERFILES=$(find "$PROJECT_ROOT" -name "Dockerfile" -type f 2>/dev/null || true)

    if [ -z "$DOCKERFILES" ]; then
        echo "  ⚠ Dockerfile을 찾을 수 없습니다"
        return
    fi

    for dockerfile in $DOCKERFILES; do
        service_dir=$(dirname "$dockerfile")
        service_name=$(basename "$service_dir")
        echo "  분석: $service_name/Dockerfile"

        # Dockerfile 보안 체크리스트
        {
            echo "=== $service_name ==="

            # root 사용자 실행 확인
            if ! grep -q "USER" "$dockerfile"; then
                echo "  [WARN] USER 지시문 없음 (root로 실행될 수 있음)"
                TOTAL_VULNS=$((TOTAL_VULNS + 1))
            fi

            # latest 태그 사용 확인
            if grep -q ":latest" "$dockerfile"; then
                echo "  [WARN] :latest 태그 사용 (버전 고정 권장)"
                TOTAL_VULNS=$((TOTAL_VULNS + 1))
            fi

            # COPY vs ADD 확인
            if grep -q "^ADD " "$dockerfile"; then
                echo "  [INFO] ADD 대신 COPY 사용 권장"
            fi

            # 민감 파일 복사 확인
            if grep -qE "COPY.*\.(env|key|pem|crt)" "$dockerfile"; then
                echo "  [CRIT] 민감 파일 복사 감지"
                TOTAL_VULNS=$((TOTAL_VULNS + 1))
            fi

            echo ""
        } >> "$REPORT_DIR/docker_scan_${TIMESTAMP}.txt"

        # Trivy 스캔 (설치되어 있는 경우)
        if command -v trivy &> /dev/null; then
            trivy fs "$service_dir" --severity HIGH,CRITICAL --format json \
                >> "$REPORT_DIR/trivy_${service_name}_${TIMESTAMP}.json" 2>&1 || true
        fi
    done

    echo "  ✓ Docker 이미지 스캔 완료"
}

# ─── 실행 ───
scan_go_modules
scan_flutter_packages
scan_rust_crates
scan_docker_images

# ─── 결과 요약 ───
echo ""
echo "=== 스캔 결과 요약 ==="
echo "발견된 경고/취약점: $TOTAL_VULNS"
echo "상세 리포트: $REPORT_DIR/"
echo ""

if [ -d "$REPORT_DIR" ]; then
    echo "생성된 리포트:"
    ls -la "$REPORT_DIR/"*"${TIMESTAMP}"* 2>/dev/null || echo "  (리포트 파일 없음)"
fi

echo ""
echo "=== 권장 보안 도구 ==="
echo "  Go:      govulncheck (go install golang.org/x/vuln/cmd/govulncheck@latest)"
echo "  Rust:    cargo-audit (cargo install cargo-audit)"
echo "  Docker:  trivy (https://aquasecurity.github.io/trivy)"
echo "  통합:    snyk (https://snyk.io)"
echo ""
echo "완료 시각: $(date)"
