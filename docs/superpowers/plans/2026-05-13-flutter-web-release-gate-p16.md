# Flutter Web Release Gate Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Flutter JS web build 성공과 Wasm dry-run compatibility warning을 명확히 분리하는 release gate를 추가한다.

**Architecture:** 현재 릴리스 타깃은 `flutter build web --no-pub`로 생성되는 JS web artifact로 둔다. gate script는 build 성공 marker를 필수로 확인하고, Wasm dry-run warning은 non-blocking warning으로 출력한다. 향후 Wasm을 릴리스 타깃으로 전환할 때 별도 strict gate를 추가한다.

**Tech Stack:** Bash, Flutter web build, shell parser test.

---

### Task 1: Release gate parser contract

**Files:**
- Create: `scripts/flutter_web_release_gate_test.sh`
- Create: `scripts/flutter_web_release_gate.sh`

- [ ] **Step 1: Write the failing parser test**

Create a shell test that calls `scripts/flutter_web_release_gate.sh --policy-check <fixture-log>`.

The test must verify:
- A log containing `✓ Built build/web` and `Wasm dry run findings:` exits 0.
- Output contains `FLUTTER_WEB_RELEASE_GATE_PASS`.
- Output contains `FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN`.
- A log without the build success marker exits non-zero.

- [ ] **Step 2: Run test to verify RED**

Run:
`bash scripts/flutter_web_release_gate_test.sh`

Expected: FAIL because `scripts/flutter_web_release_gate.sh` does not exist yet.

- [ ] **Step 3: Implement release gate script**

Implement:
- `--policy-check <log>` mode for fast parser tests.
- default mode that runs `frontend/flutter-app` `flutter build web --no-pub`, stores a log under `/tmp`, and applies the policy check.
- `FLUTTER_BIN` override with fallback to `/mnt/d/우리집/flutter_cache/flutter/bin/flutter`.

- [ ] **Step 4: Run parser test to verify GREEN**

Run:
`bash scripts/flutter_web_release_gate_test.sh`

Expected: PASS.

### Task 2: Real build and project gates

**Files:**
- Create: `docs/audit/flutter-web-release-gate-p16.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [ ] **Step 1: Run release gate**

Run:
`bash scripts/flutter_web_release_gate.sh`

Expected: PASS with optional `FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN`.

- [ ] **Step 2: Run DataHub focused regression tests**

Run:
`cd frontend/flutter-app && flutter test --no-pub test/features/data_hub/presentation/data_hub_layout_smoke_test.dart test/features/data_hub/presentation/data_hub_evidence_badge_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart test/features/data_hub/domain/data_hub_domain_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`

- [ ] **Step 3: Run common gates**

Run:
`python3 scripts/validate_ssot_constants.py`

Run:
`bash scripts/security_release_gate.sh`

Run:
`bash scripts/assay_evidence_gate.sh`

- [ ] **Step 4: Run script hygiene checks**

Run:
`bash -n scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh`

Run:
`git diff --check -- scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh CHANGELOG.md CONTEXT.md`
