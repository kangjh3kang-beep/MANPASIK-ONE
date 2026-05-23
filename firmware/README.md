# ManPaSik MMUP-OS — 리더기 임베디드 운영체제

> **에픽** E-FW · **MCU** STM32F405 · **표준** IEC 62304 Class B/C · ISO 14971
> **범위** 리더기에 탑재되는 RTOS 기반 임베디드 펌웨어 (범용 커널이 아님)

## 아키텍처

```
firmware/
├── mmup_os/        — RTOS 태스크 스케줄러, 전원관리, 워치독
├── drivers/        — AFE 9블록 드라이버 (embedded-hal trait 추상화)
├── measure_seq/    — 측정 시퀀서 (급전→인큐→읽기→Sdiff 1차)
├── secure/         — 보안부팅 (ECDSA), A/B OTA, 자가진단
└── rust_embedded/  — no_std Rust 안전중요 측정 로직
```

## RTOS 선정 (근거 기반)

| 후보 | 라이선스 | IEC 62304 | 장점 | 단점 |
|------|---------|-----------|------|------|
| **FreeRTOS** | MIT | 자체 인증 필요 | STM32 벤더 지원, 풍부한 생태계 | 인증 비용 |
| **SafeRTOS** | 상용 | 사전인증 | IEC 61508 SIL 3 인증 | 라이선스 비용 |
| **ThreadX (Azure RTOS)** | 상용 | 사전인증 | IEC 62304 Class C, FDA 실적 | MS 의존 |
| **Zephyr** | Apache 2.0 | 자체 인증 필요 | 보안 스택, BLE 내장 | 리소스 무거움 |

**권고**: 프로토타입은 FreeRTOS, 양산은 SafeRTOS/ThreadX. 최종 선정은 HW SSOT + 사람 승인.

## 하네스 엔지니어링 준수

- **H1**: 모든 드라이버는 embedded-hal trait 뒤에 추상화 (구체 IC 비참조)
- **H2**: CSI v1.0 14신호핀 + 2마운트 = 16핀 계약 준수, v2.0 어댑터 의무
- **H4**: Stage-1(전기화학) → Stage-2(광학) → Stage-3(NAAT) 순차 확장
- **H5**: BLE GATT UUID + NFC 매니페스트 버전화 계약
- **H6**: 측정 시퀀서에서 채널별 실패 격리, Result<T,E> 필수

## 빌드 (크로스 컴파일)

```bash
# Rust embedded (ARM Cortex-M4F)
cd rust_embedded
cargo build --target thumbv7em-none-eabihf --release
```
