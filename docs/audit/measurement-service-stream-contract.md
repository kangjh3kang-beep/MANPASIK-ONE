# Measurement Service Stream Contract Audit

**Date**: 2026-05-01  
**Owner**: Codex  
**Scope**: Measure golden path H1 backend stream persistence contract

## 목적

Gateway `/api/v1/measurements/process` 계약 테스트는 REST 요청이 gRPC `StreamMeasurement`로 전달되는 것을 보장한다. 이번 단계는 그 다음 모세혈관인 measurement-service gRPC 핸들러에서 수신 프레임이 실제 서비스 저장 경로로 내려가고, Timescale 계층용 측정 데이터와 Milvus 계층용 fingerprint vector가 함께 생성되는지 고정한다.

## 고정한 계약

`MeasurementHandler.StreamMeasurement`는 각 수신 프레임에 대해 아래를 수행해야 한다.

1. `session_id`가 없으면 `InvalidArgument`로 거부한다.
2. `raw_channels`를 `[]float32` fingerprint vector로 변환한다.
3. differential 보정값을 `service.MeasurementData`에 전달한다.
4. 활성 세션의 `device_id`, `user_id`, `cartridge_id`를 저장 데이터에 보강한다.
5. measurement repository `Store`와 vector repository `StoreFingerprint`를 호출한다.
6. 처리 결과를 `MeasurementResult`로 스트림에 반환하고 `processed_at`을 채운다.

## 추가 검증

파일: `backend/services/measurement-service/internal/handler/grpc_stream_test.go`

- `TestStreamMeasurementStoresMeasurementAndFingerprint`
  - 활성 세션 기반으로 스트림 프레임 1건을 주입한다.
  - 저장된 measurement의 session/device/user/cartridge metadata 보강을 검증한다.
  - differential corrected value, primary value, unit, confidence, environment metadata를 검증한다.
  - `raw_channels`와 fingerprint vector가 저장소와 응답에 동일하게 반영되는지 검증한다.
  - `processed_at` 응답 필드가 nil이 아닌지 검증한다.
- `TestStreamMeasurementRequiresSessionID`
  - 빈 `session_id` 프레임이 저장소 호출 없이 `InvalidArgument`로 차단되는지 검증한다.

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik
export PATH=/home/kangjh3kang/sdk/go-go1.26.2/bin:/usr/local/bin:/usr/bin:/bin
gofmt -w backend/services/measurement-service/internal/handler/grpc_stream_test.go
go test ./backend/services/measurement-service/... ./backend/services/gateway/...
```

결과: PASS

## 다음 단계

- Flutter Measure trace event를 gateway 또는 audit-service remote observability sink로 전달한다.
- 이후 H2에서 mock/stub retirement register의 Measure 관련 항목을 실제 native Rust/DB wiring 상태와 대조한다.
