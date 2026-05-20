# Measurement Assay Semantics Gate

## 목적

기존 `StreamMeasurement` 경로는 `s_corrected`를 그대로 `primary_value`로 복사하고, 단위를 `mg/dL`, confidence를 `0.95`로 고정했다. 이 구조는 CRP, HbA1c, TSH, 비타민 등 카트리지별 단위와 성능 정책이 다른 만파식 범용분석시스템의 실현 가능성을 낮춘다.

이번 변경은 측정 결과 의미론을 `backend/shared/assay` registry로 분리해, 세션의 카트리지 타입에 따라 단위와 confidence policy를 결정하게 한다.

## 변경 범위

- `backend/shared/assay/registry.go`
  - `glucose`, `crp`를 포함한 건강 바이오마커 기본 assay 정의를 추가했다.
  - `blood_glucose`, `cartridge-glucose`, `0x01`, `0x0d` 같은 alias를 canonical assay key로 정규화한다.
  - 알 수 없는 카트리지는 `ErrUnknownAssay`로 거부한다.
  - confidence는 고정 literal 대신 raw channel completeness를 assay별 floor/ceiling 사이에 매핑한다.
- `backend/services/measurement-service/internal/handler/grpc.go`
  - handler에서 `PrimaryValue`, `Unit`, `Confidence` 하드코딩을 제거했다.
- `backend/services/measurement-service/internal/service/measurement.go`
  - 세션 메타데이터 backfill 이후 `assay.Evaluate()`로 측정 의미론을 적용한다.
  - upstream에서 이미 `PrimaryValue`, `Unit`, `Confidence`를 명시한 경우 기존 값을 보존하되, 카트리지 타입은 반드시 registry에 존재해야 한다.
- 테스트
  - registry alias, unknown cartridge, confidence completeness 정책을 검증했다.
  - gRPC stream, TCP loopback, helper process smoke가 새 assay-derived confidence 계약을 검증한다.
  - CRP 세션은 `mg/L` 단위로 응답하는지 검증한다.

## 품질 게이트

- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS

## 잔여 리스크

- 현재 registry는 P0 하드코딩 제거용 최소 정의다. 다음 단계에서는 각 assay에 LoB/LoD/LoQ, linearity, interference, reference method, LOINC, UCUM, validation status를 추가해야 한다.
- `PrimaryValueSourceCorrectedSignal`은 아직 모든 assay에 동일하게 적용된다. 실제 인허가 수준에서는 카트리지별 calibration curve와 lot calibration manifest를 통과해야 한다.
- `CustomResearch`와 비표적 카트리지는 운영 진단 경로에서 거부된다. 연구 모드에서만 허용할 별도 policy gate가 필요하다.
