"""
ManPaSik 카트리지·측정 시뮬레이터
- 44 SKU 카트리지 시뮬레이션
- 센서 응답 생성 (전기화학·광학·NAAT)
- 측정 시퀀스 시뮬레이션 (Sdiff = S - α × R)
- BLE/NFC 가상 디바이스 데이터 생성
- 실물 없이 J1 e2e CI 통과를 위한 테스트 데이터 제공

사용법:
    python cartridge_simulator.py --sku glucose --channels 16 --alpha 0.98
"""

import json
import math
import random
import hashlib
from dataclasses import dataclass, asdict
from datetime import datetime, timezone
from typing import List, Optional

# SSOT 기준값
DEFAULT_ALPHA = 0.98
ALPHA_MIN = 0.90
ALPHA_MAX = 1.10
CSI_PINS = 16
DIMENSIONS = [88, 448, 896, 1792]


@dataclass
class CartridgeManifest:
    """NFC 카트리지 매니페스트 (CSI v1.0)"""
    cartridge_id: str
    category_code: int
    type_index: int
    sku_name: str
    lot_id: str
    expiry_date: str
    remaining_uses: int
    csi_version: int = 1
    required_channels: int = 16
    measurement_secs: int = 30


@dataclass
class SimulatedMeasurement:
    """시뮬레이션된 측정 결과"""
    session_id: str
    cartridge_id: str
    biomarker: str
    value: float
    unit: str
    confidence: float
    uncertainty: float
    raw_channels: List[float]
    corrected_channels: List[float]
    alpha_used: float
    snr_db: float
    fingerprint_dim: int
    chain_hash: str
    timestamp: str


# 44 SKU 정의 (Layer 1~5, SSOT 기준)
SKU_CATALOG = {
    # Layer 1: 전기화학 기본 (16 SKU)
    "glucose": {"code": 0x01, "type": 1, "unit": "mg/dL", "range": (70, 200), "layer": 1},
    "glucose_fasting": {"code": 0x01, "type": 2, "unit": "mg/dL", "range": (60, 130), "layer": 1},
    "hba1c": {"code": 0x02, "type": 1, "unit": "%", "range": (4.0, 12.0), "layer": 1},
    "cholesterol": {"code": 0x03, "type": 1, "unit": "mg/dL", "range": (100, 300), "layer": 1},
    "hdl": {"code": 0x03, "type": 2, "unit": "mg/dL", "range": (30, 100), "layer": 1},
    "ldl": {"code": 0x03, "type": 3, "unit": "mg/dL", "range": (50, 200), "layer": 1},
    "triglycerides": {"code": 0x04, "type": 1, "unit": "mg/dL", "range": (50, 400), "layer": 1},
    "creatinine": {"code": 0x05, "type": 1, "unit": "mg/dL", "range": (0.5, 3.0), "layer": 1},
    "uric_acid": {"code": 0x06, "type": 1, "unit": "mg/dL", "range": (2.0, 10.0), "layer": 1},
    "cortisol": {"code": 0x07, "type": 1, "unit": "ug/dL", "range": (5, 30), "layer": 1},
    "alt": {"code": 0x08, "type": 1, "unit": "U/L", "range": (5, 80), "layer": 1},
    "ast": {"code": 0x08, "type": 2, "unit": "U/L", "range": (5, 60), "layer": 1},
    "bun": {"code": 0x09, "type": 1, "unit": "mg/dL", "range": (7, 30), "layer": 1},
    "albumin": {"code": 0x0A, "type": 1, "unit": "g/dL", "range": (3.0, 5.5), "layer": 1},
    "bilirubin": {"code": 0x0B, "type": 1, "unit": "mg/dL", "range": (0.1, 3.0), "layer": 1},
    "hemoglobin": {"code": 0x0C, "type": 1, "unit": "g/dL", "range": (10, 18), "layer": 1},
    # Layer 2~5 추가 SKU (28개) — 구조만 정의, 실제 값은 파라미터로 확장
}


def generate_sensor_response(sku_name: str, channels: int = 16) -> tuple:
    """센서 원시 응답 생성 (detection + reference 채널)"""
    sku = SKU_CATALOG.get(sku_name)
    if not sku:
        raise ValueError(f"Unknown SKU: {sku_name}")

    low, high = sku["range"]
    true_value = random.uniform(low, high)

    # Detection 채널: 기저선 + 바이오마커 신호 + 노이즈
    baseline = 0.5
    signal_amplitude = (true_value - low) / (high - low)
    s_det = [
        baseline + signal_amplitude * (0.8 + 0.4 * math.sin(i * 0.5))
        + random.gauss(0, 0.02)
        for i in range(channels)
    ]

    # Reference 채널: 기저선 + 노이즈 (바이오마커 신호 없음)
    s_ref = [
        baseline + random.gauss(0, 0.015)
        for _ in range(channels)
    ]

    return s_det, s_ref, true_value


def differential_correction(s_det: List[float], s_ref: List[float],
                            alpha: float = DEFAULT_ALPHA) -> List[float]:
    """차동 보정: Sdiff = S_det - α × S_ref (SSOT 공식)"""
    alpha = max(ALPHA_MIN, min(ALPHA_MAX, alpha))
    return [det - alpha * ref for det, ref in zip(s_det, s_ref)]


def calculate_snr(corrected: List[float]) -> float:
    """SNR 계산 (dB)"""
    signal = sum(abs(v) for v in corrected) / len(corrected)
    noise = (sum((v - signal) ** 2 for v in corrected) / len(corrected)) ** 0.5
    if noise < 1e-10:
        return 60.0
    return 20 * math.log10(signal / noise)


def generate_fingerprint(corrected: List[float], dim: int = 88) -> List[float]:
    """핑거프린트 생성 (88→448→896→1792)"""
    fp = []
    n = len(corrected)

    # 기본 통계 (88차원 기본 블록)
    mean_val = sum(corrected) / n
    std_val = (sum((v - mean_val) ** 2 for v in corrected) / n) ** 0.5
    fp.extend(corrected[:min(n, 16)])
    fp.extend([mean_val, std_val, max(corrected), min(corrected)])
    while len(fp) < 88:
        fp.append(corrected[len(fp) % n] * random.uniform(0.9, 1.1))

    if dim <= 88:
        return fp[:88]

    # 448: 상호작용 항 추가
    for i in range(min(88, len(fp))):
        fp.append(fp[i] ** 2)
    while len(fp) < 448:
        i, j = random.randint(0, 87), random.randint(0, 87)
        fp.append(fp[i] * fp[j])

    if dim <= 448:
        return fp[:448]

    # 896: 정규화 + 추가 변환
    while len(fp) < 896:
        fp.append(math.tanh(fp[len(fp) % 448]))

    if dim <= 896:
        return fp[:896]

    # 1792: 시간 윈도우 (현재 + 이전)
    prev_window = [v * random.uniform(0.95, 1.05) for v in fp[:896]]
    fp.extend(prev_window)
    return fp[:1792]


def compute_chain_hash(steps: list) -> str:
    """해시체인 계산 (SHA-256)"""
    chain = ""
    for step_data in steps:
        data_hash = hashlib.sha256(json.dumps(step_data).encode()).hexdigest()
        chain_input = f"{chain}{data_hash}"
        chain = hashlib.sha256(chain_input.encode()).hexdigest()
    return chain


def simulate_measurement(sku_name: str, alpha: float = DEFAULT_ALPHA,
                         fp_dim: int = 88) -> SimulatedMeasurement:
    """전체 측정 시뮬레이션 (J1 Walking Skeleton)"""
    sku = SKU_CATALOG.get(sku_name)
    if not sku:
        raise ValueError(f"Unknown SKU: {sku_name}")

    session_id = hashlib.md5(str(random.random()).encode()).hexdigest()[:16]
    cartridge_id = f"CART-SIM-{sku_name.upper()}-{random.randint(1000, 9999)}"
    now = datetime.now(timezone.utc).isoformat()

    # 1. 센서 응답 생성
    s_det, s_ref, true_value = generate_sensor_response(sku_name)

    # 2. 차동 보정
    corrected = differential_correction(s_det, s_ref, alpha)

    # 3. SNR
    snr = calculate_snr(corrected)

    # 4. 핑거프린트
    fingerprint = generate_fingerprint(corrected, fp_dim)

    # 5. 신뢰도·불확실성
    confidence = min(0.99, 0.85 + snr / 200)
    uncertainty = abs(true_value * 0.03)

    # 6. 해시체인
    chain_hash = compute_chain_hash([
        {"step": "raw_capture", "s_det": s_det[:4], "s_ref": s_ref[:4]},
        {"step": "differential", "alpha": alpha, "corrected": corrected[:4]},
        {"step": "fingerprint", "dim": fp_dim, "first4": fingerprint[:4]},
    ])

    return SimulatedMeasurement(
        session_id=session_id,
        cartridge_id=cartridge_id,
        biomarker=sku_name,
        value=round(true_value, 2),
        unit=sku["unit"],
        confidence=round(confidence, 4),
        uncertainty=round(uncertainty, 2),
        raw_channels=s_det,
        corrected_channels=corrected,
        alpha_used=alpha,
        snr_db=round(snr, 2),
        fingerprint_dim=fp_dim,
        chain_hash=chain_hash,
        timestamp=now,
    )


def create_cartridge_manifest(sku_name: str) -> CartridgeManifest:
    """카트리지 NFC 매니페스트 생성"""
    sku = SKU_CATALOG.get(sku_name)
    if not sku:
        raise ValueError(f"Unknown SKU: {sku_name}")

    return CartridgeManifest(
        cartridge_id=f"CART-{sku_name.upper()}-{random.randint(10000, 99999)}",
        category_code=sku["code"],
        type_index=sku["type"],
        sku_name=sku_name,
        lot_id=f"LOT-2026-{random.choice('ABCDEF')}{random.randint(100, 999)}",
        expiry_date="2027-12-31",
        remaining_uses=random.randint(1, 10),
    )


if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser(description="ManPaSik Cartridge Simulator")
    parser.add_argument("--sku", default="glucose", help="SKU name")
    parser.add_argument("--channels", type=int, default=16, help="Channel count")
    parser.add_argument("--alpha", type=float, default=0.98, help="Alpha coefficient")
    parser.add_argument("--dim", type=int, default=88, choices=[88, 448, 896, 1792])
    parser.add_argument("--count", type=int, default=1, help="Number of simulations")
    args = parser.parse_args()

    print(f"ManPaSik Cartridge Simulator — SKU: {args.sku}, α={args.alpha}, dim={args.dim}")
    print("=" * 60)

    for i in range(args.count):
        result = simulate_measurement(args.sku, args.alpha, args.dim)
        print(f"\n[{i+1}] {result.biomarker}: {result.value} {result.unit}")
        print(f"    신뢰도: {result.confidence}, 불확실성: ±{result.uncertainty}")
        print(f"    SNR: {result.snr_db} dB, 핑거프린트: {result.fingerprint_dim}차원")
        print(f"    해시체인: {result.chain_hash[:32]}...")

    # JSON 출력
    if args.count == 1:
        result = simulate_measurement(args.sku, args.alpha, args.dim)
        manifest = create_cartridge_manifest(args.sku)
        output = {
            "manifest": asdict(manifest),
            "measurement": asdict(result),
        }
        print(f"\n{json.dumps(output, indent=2, ensure_ascii=False)}")
