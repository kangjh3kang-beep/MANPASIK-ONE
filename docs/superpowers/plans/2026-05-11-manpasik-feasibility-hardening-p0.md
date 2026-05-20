# ManPaSik Feasibility Hardening P0 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 만파식 범용분석시스템의 실현 가능성을 떨어뜨리는 P0 기술 격차를 먼저 닫아 SSOT, 측정 의미론, mock-free 운영 게이트, 품질 게이트의 기반을 만든다.

**Architecture:** 최신 `ManPaSik_Tech_Spec_v2.4.3.md`를 SSOT로 삼고, 코드·AGENTS·CI가 같은 상수를 검증하게 한다. 측정 결과는 공통 raw/fingerprint 파이프라인과 카트리지별 assay manifest를 분리해, `s_corrected`를 임의 단위로 직접 노출하지 않게 한다.

**Tech Stack:** Python 3 validation script, GitHub Actions, Go measurement-service, Rust manpasik-engine, Markdown audit docs.

---

## File Structure

- Create: `scripts/validate_ssot_constants.py`
  - `ManPaSik_Tech_Spec_v2.4.3.md`, `rust-core/manpasik-engine/src/lib.rs`, `AGENTS.md`의 핵심 상수 일치 여부를 검증한다.
- Modify: `.github/workflows/ci.yml`
  - CI 초반에 SSOT 검증 job을 추가하고 Rust clippy를 blocking gate로 전환한다.
- Modify: `AGENTS.md`
  - `alpha=0.98`, `MAX_CHANNELS=1792`, classifier output 30 등 최신 스펙과 코드 기준으로 낡은 상수를 정리한다.
- Create: `docs/audit/ssot-constants-gate.md`
  - 변경 이유, 검증 범위, 품질 게이트, 잔여 리스크를 기록한다.
- Future Modify: `backend/services/measurement-service/internal/handler/grpc.go`
  - `SCorrected -> PrimaryValue`, `Unit=mg/dL`, `Confidence=0.95` 하드코딩을 assay-aware 계산으로 교체한다.
- Future Create: `backend/shared/assay/registry.go`, `backend/shared/assay/registry_test.go`
  - 카트리지별 analyte, unit, precision policy, confidence policy를 중앙화한다.

## Task 1: SSOT Constants Gate

**Files:**
- Create: `scripts/validate_ssot_constants.py`
- Modify: `.github/workflows/ci.yml`
- Modify: `AGENTS.md`
- Create: `docs/audit/ssot-constants-gate.md`

- [x] **Step 1: Write the validator**

Create `scripts/validate_ssot_constants.py` with these checks:

```python
#!/usr/bin/env python3
from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def read(path: str) -> str:
    return (ROOT / path).read_text(encoding="utf-8")


def find_float(name: str, text: str) -> float:
    pattern = rf"{re.escape(name)}\s*:\s*f64\s*=\s*([0-9]+(?:\.[0-9]+)?)"
    match = re.search(pattern, text)
    if not match:
        raise AssertionError(f"{name} not found")
    return float(match.group(1))


def find_usize(name: str, text: str) -> int:
    pattern = rf"{re.escape(name)}\s*:\s*usize\s*=\s*([0-9]+)"
    match = re.search(pattern, text)
    if not match:
        raise AssertionError(f"{name} not found")
    return int(match.group(1))


def require(condition: bool, message: str, failures: list[str]) -> None:
    if not condition:
        failures.append(message)


def main() -> int:
    spec = read("ManPaSik_Tech_Spec_v2.4.3.md")
    lib = read("rust-core/manpasik-engine/src/lib.rs")
    agents = read("AGENTS.md")

    expected_alpha = find_float("ALPHA_DEFAULT", spec)
    expected_max_channels = find_usize("FINGERPRINT_DIM_MAX", spec)

    failures: list[str] = []
    require(
        find_float("DEFAULT_ALPHA", lib) == expected_alpha,
        f"rust DEFAULT_ALPHA must equal spec ALPHA_DEFAULT {expected_alpha}",
        failures,
    )
    require(
        find_usize("MAX_CHANNELS", lib) == expected_max_channels,
        f"rust MAX_CHANNELS must equal spec FINGERPRINT_DIM_MAX {expected_max_channels}",
        failures,
    )
    require(
        f"기본값 = {expected_alpha:.2f}" in agents,
        f"AGENTS.md must document alpha default {expected_alpha:.2f}",
        failures,
    )
    require(
        f"MAX_CHANNELS = {expected_max_channels}" in agents,
        f"AGENTS.md must document MAX_CHANNELS = {expected_max_channels}",
        failures,
    )
    require("기본값 = 0.95" not in agents, "AGENTS.md still contains obsolete alpha 0.95", failures)
    require("MAX_CHANNELS = 896" not in agents, "AGENTS.md still contains obsolete MAX_CHANNELS = 896", failures)

    if failures:
        for failure in failures:
            print(f"SSOT_CHECK_FAIL: {failure}", file=sys.stderr)
        return 1

    print(
        "SSOT_CHECK_PASS "
        f"alpha={expected_alpha:.2f} "
        f"max_channels={expected_max_channels}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
```

- [x] **Step 2: Run the validator and verify it fails before doc alignment**

Run: `python3 scripts/validate_ssot_constants.py`

Expected before `AGENTS.md` update: FAIL with obsolete alpha/max channel messages.

- [x] **Step 3: Align `AGENTS.md` constants**

Change the differential section to `기본값 = 0.98` and the fingerprint section to `MAX_CHANNELS = 1792`.

- [x] **Step 4: Wire the CI gate**

Add a `ssot-governance` job before language-specific jobs:

```yaml
  ssot-governance:
    name: SSOT Governance
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Validate SSOT constants
        run: python3 scripts/validate_ssot_constants.py
```

Remove `continue-on-error: true` from Rust clippy so lint failures block release.

- [x] **Step 5: Verify**

Run:

```bash
python3 scripts/validate_ssot_constants.py
python3 - <<'PY'
import yaml
from pathlib import Path
yaml.safe_load(Path('.github/workflows/ci.yml').read_text())
print('CI_YAML_PARSE_PASS')
PY
```

Expected: `SSOT_CHECK_PASS alpha=0.98 max_channels=1792` and `CI_YAML_PARSE_PASS`.

## Task 2: Measurement Assay Registry

**Files:**
- Create: `backend/shared/assay/registry.go`
- Create: `backend/shared/assay/registry_test.go`
- Modify: `backend/services/measurement-service/internal/handler/grpc.go`
- Modify: `backend/services/measurement-service/internal/service/measurement_test.go`
- Create: `docs/audit/measurement-assay-semantics.md`

- [x] **Step 1: Write failing tests**

Add tests that assert `Glucose` returns `mg/dL`, `Crp` returns `mg/L`, unknown cartridge types return an explicit error, and confidence is calculated from assay policy rather than a literal `0.95`.

- [x] **Step 2: Implement registry**

Create a small registry keyed by cartridge type with fields: `Code`, `Analyte`, `Unit`, `PrimaryValueSource`, `ConfidenceFloor`, `ConfidenceCeiling`.

- [x] **Step 3: Replace handler hardcoding**

Replace `PrimaryValue: req.GetDifferential().GetSCorrected()`, `Unit: "mg/dL"`, and `Confidence: 0.95` with registry-derived values.

- [x] **Step 4: Verify**

Run: `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/shared/assay ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`

Expected: PASS.

## Task 3: Prod Mock-Free Gate

**Files:**
- Modify: `rust-core/manpasik-engine/src/ai/mod.rs`
- Create: `rust-core/manpasik-engine/tests/ai_prod_mock_gate.rs`
- Modify: `docs/audit/mock-retirement-register.md`
- Create: `docs/audit/ai-prod-mock-free-gate.md`

- [x] **Step 1: Write failing Rust tests**

Add tests that set production mode and assert `predict()` returns an error when no real runtime-backed model is loaded.

- [x] **Step 2: Implement explicit simulation policy**

Add an `InferenceMode` or equivalent flag so simulation is allowed only in tests/demo mode and impossible to use silently in production mode.

- [x] **Step 3: Verify**

Run: `cd rust-core && cargo test -p manpasik-engine ai_prod_mock_gate`

Expected: PASS.

## Task 4: Security Release Gates

**Files:**
- Modify: `.github/workflows/ci.yml`
- Create: `scripts/security_release_gate.sh`
- Create: `docs/audit/security-release-gates.md`

- [x] **Step 1: Add script checks**

The script must fail if production compose/k8s manifests contain default secrets, Elasticsearch security disabled, or `Just Works` BLE pairing is enabled for production profile.

- [x] **Step 2: Wire CI**

Run the script in `ssot-governance` after the constants check.

- [x] **Step 3: Verify**

Run: `bash scripts/security_release_gate.sh`

Expected: PASS for dev-only defaults and FAIL if any production profile carries those defaults.

## Task 5: Documentation And Shared Context

**Files:**
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`
- Create or Modify: task-specific `docs/audit/*.md`

- [x] **Step 1: Add latest changelog entry**

Insert a new Korean entry at the top with status, changed files, decisions, quality gates, and next steps.

- [x] **Step 2: Update context**

Update only the sections affected by P0 hardening: SSOT gate, measurement semantics, mock-free gate, and residual Docker compose smoke blocker.

- [x] **Step 3: Verify changed-file hygiene**

Run: `git diff --check -- scripts/validate_ssot_constants.py .github/workflows/ci.yml AGENTS.md docs/audit/ssot-constants-gate.md CHANGELOG.md CONTEXT.md`

Expected: no output and exit 0.
