# Dart Proto Generation Preflight 게이트

**작성일**: 2026-05-01  
**작성자**: Codex  
**범위**: 전체 Dart generated proto 전환 전 안전성 확인

## 목적

native Measure gRPC stream 직결을 위해 수동 Dart 스텁에 필요한 surface는 보강했다. 다음 위험은 전체 `protoc --dart_out=grpc` 산출물을 기존 수동 스텁 위에 그대로 덮어쓸 때 발생할 수 있는 대규모 API/의존성 충돌이다. 이번 게이트는 실제 파일을 덮어쓰지 않고 임시 디렉터리에만 생성한 뒤, 전환 가능성과 차단 사유를 명확히 보고한다.

## 구현 결과

- `frontend/flutter-app/scripts/check_proto_generation_preflight.sh`를 추가했다.
- 스크립트는 `protoc`과 `protoc-gen-dart` 존재를 확인한다.
- 임시 디렉터리에 `backend/shared/proto/manpasik.proto` Dart gRPC 산출물을 생성한다.
- 생성 산출물에 `MeasurementData`와 `streamMeasurement`가 존재하는지 확인한다.
- 최신 protoc plugin 산출물의 Timestamp import가 현재 pub cache와 호환되는지 확인해 전체 치환 가능성을 분리 보고한다.

## 현재 판정

- `StreamMeasurement` 생성 surface: 사용 가능
- 전체 generated Dart 즉시 치환: 보류
- 보류 사유: 현재 `protobuf 3.1.0` lock/pub cache는 최신 protoc plugin이 생성하는 well-known Timestamp import 경로와 맞지 않는다.
- 격리 compile gate: `docs/audit/dart-proto-full-generated-compile-gate.md`에서 완료
- 격리 compile toolchain: `grpc 5.1.0`, `protobuf 6.0.0`
- 앱 lock 정렬 및 checked-in official generated 교체: `docs/audit/dart-proto-official-generated-replacement.md`에서 완료
- 현재 preflight 판정: `PROTO_PREFLIGHT_FULL_REPLACEMENT=ready_for_compile_gate`

## 품질 게이트

```bash
cd /home/kangjh3kang/Manpasik/frontend/flutter-app
bash scripts/check_proto_generation_preflight.sh
```

예상 상태:

- `PROTO_PREFLIGHT_STATUS=generated_streammeasurement_available`
- `PROTO_PREFLIGHT_FULL_REPLACEMENT=blocked_timestamp_import_incompatible_with_current_pub_cache`
- `PROTO_PREFLIGHT_LOCKED_PROTOBUF_VERSION=3.1.0`

## 다음 단계

- 전체 generated Dart compile gate와 checked-in official generated 교체는 완료했다.
- 다음 단계는 Auth login 사용자 식별자 보강 또는 Docker Compose container smoke runtime PASS 확보다.
- official generated output은 계속 wire golden contract test로 보호한다.
