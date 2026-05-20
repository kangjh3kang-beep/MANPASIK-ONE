# DataHub Evidence Metadata P9 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans and superpowers:test-driven-development. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** DataHub trend/summary models가 measurement history evidence fields를 잃지 않도록 보존한다.

**Architecture:** `TrendDataPoint`에 evidence metadata를 직접 추가하고, `BiomarkerSummary`에는 latest measurement evidence metadata를 추가한다. REST repository는 measurement history response의 snake_case/camelCase evidence fields를 decode해 trend points와 summaries로 전달한다.

**Tech Stack:** Flutter domain tests, Flutter REST repository tests, targeted Flutter analyze, existing SSOT/security/assay gates.

---

## Stage Gate Rules

- RED tests first.
- Legacy/default evidence status는 `unknown`, `false`, `[]`다.
- `research_only`는 `diagnosticReady=false`로 보존한다.
- DataHub UI badge 표시는 다음 단계로 분리한다.

## File Structure

- Modify: `frontend/flutter-app/lib/features/data_hub/domain/data_hub_repository.dart`
  - `TrendDataPoint` evidence fields 추가.
  - `BiomarkerSummary` latest evidence fields 추가.
- Modify: `frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart`
  - REST history evidence decode helper 추가.
  - `getTrendData`, `getBiomarkerSummary`, `getAllBiomarkerSummaries`에 evidence metadata 전달.
- Modify: `frontend/flutter-app/test/features/data_hub/domain/data_hub_domain_test.dart`
  - evidence defaults/explicit fields 테스트 추가.
- Modify: `frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart`
  - REST history snake_case/camelCase evidence mapping 테스트 추가.
- Create: `docs/audit/datahub-evidence-metadata-p9.md`
- Modify: `CHANGELOG.md`, `CONTEXT.md`

## Task 1: Domain Evidence Metadata

- [x] **Step 1: Write failing domain tests**

Require default and explicit evidence fields on `TrendDataPoint` and `BiomarkerSummary`.

- [x] **Step 2: Verify RED**

Run domain test and confirm missing fields failure.

- [x] **Step 3: Add domain fields**

Add fields with conservative defaults.

- [x] **Step 4: Verify GREEN**

Run domain test.

## Task 2: REST Repository Evidence Mapping

- [x] **Step 1: Write failing REST tests**

Use local HTTP server responses to assert:
- `getTrendData` maps snake_case evidence fields.
- `getTrendData` maps legacy camelCase evidence fields.
- `getBiomarkerSummary` exposes latest evidence status.

- [x] **Step 2: Verify RED**

Run DataHub REST test and confirm missing mapping/fields failure.

- [x] **Step 3: Implement REST mapping**

Add small decode helpers and propagate evidence metadata.

- [x] **Step 4: Verify GREEN**

Run DataHub REST test.

## Task 3: Audit and Final Gates

- [x] **Step 1: Record audit and logs**

Document TDD record, self-review, quality gates, and next-stage guidance.

- [x] **Step 2: Run final gates**

Run:

```bash
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/data_hub/domain/data_hub_domain_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart
cd frontend/flutter-app && /mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos lib/features/data_hub/domain/data_hub_repository.dart lib/features/data_hub/data/data_hub_repository_rest.dart test/features/data_hub/domain/data_hub_domain_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart
```

Expected: PASS.
