"""
ManPaSik (萬波息) - Triton Inference Server Ensemble Pipeline
=============================================================

실제 추론 파이프라인:
  1. DriftDetector: 온도 보상 다항식 회귀 기반 캘리브레이션
  2. AnomalyDetector: Isolation Forest 기반 OOD (Out-of-Distribution) 탐지
  3. GlucoseQuantifier: XGBoost 모델 로드 및 혈당 정량
  4. Full pipeline: input -> drift correction -> anomaly check -> quantification -> output

핑거프린트 입력: 896차원 (SSOT FINGERPRINT_DIM)
"""

import json
import os
import warnings
from dataclasses import dataclass, field
from pathlib import Path
from typing import Dict, List, Optional, Tuple

import numpy as np
from sklearn.ensemble import IsolationForest

try:
    import xgboost as xgb

    XGB_AVAILABLE = True
except ImportError:
    XGB_AVAILABLE = False

# SSOT 상수
FINGERPRINT_DIM = 896
ALPHA_DEFAULT = 0.98


@dataclass
class PipelineResult:
    """앙상블 파이프라인 출력 구조체."""

    status: str  # "success", "anomaly_detected", "error"
    glucose_mg_dl: Optional[float] = None
    confidence: float = 0.0
    anomaly_score: float = 0.0
    is_anomaly: bool = False
    drift_factor: float = 1.0
    temperature_c: float = 25.0
    warnings: List[str] = field(default_factory=list)

    def to_dict(self) -> Dict:
        return {
            "status": self.status,
            "concentrations": {
                "glucose": round(self.glucose_mg_dl, 2) if self.glucose_mg_dl else None,
            },
            "confidence": round(self.confidence, 4),
            "anomaly_score": round(self.anomaly_score, 4),
            "is_anomaly": self.is_anomaly,
            "drift_compensation": {
                "factor": round(self.drift_factor, 6),
                "temperature_c": round(self.temperature_c, 1),
                "applied": abs(self.drift_factor - 1.0) > 1e-6,
            },
            "warnings": self.warnings,
        }


class DriftDetector:
    """온도 보상 캘리브레이션 모듈.

    센서 감도는 온도에 따라 비선형으로 변하므로,
    다항식 회귀를 사용하여 보정 계수를 계산합니다.

    보정 모델: drift_factor = c0 + c1*dT + c2*dT^2 + c3*dT^3
    여기서 dT = T - T_ref (T_ref = 25.0 degC)
    """

    def __init__(self, reference_temp: float = 25.0, poly_degree: int = 3):
        self.reference_temp = reference_temp
        self.poly_degree = poly_degree

        # 기본 다항식 계수 (실험적 캘리브레이션 결과 기반)
        # 25도에서 factor=1.0, 온도 변화에 따른 감도 드리프트 모사
        # 이 계수는 실제 시스템에서는 캘리브레이션 데이터로 피팅합니다.
        self._coefficients = np.array([1.0, -0.008, 0.0002, -0.000003])

        # 캘리브레이션 이력
        self._calibration_history: List[Tuple[float, float]] = []

    def fit(self, temperatures: np.ndarray, measured_factors: np.ndarray):
        """캘리브레이션 데이터로 다항식 계수 피팅.

        Parameters
        ----------
        temperatures : 온도 배열 (degC)
        measured_factors : 실측 보정 계수 배열
        """
        dt = temperatures - self.reference_temp
        self._coefficients = np.polyfit(dt, measured_factors, self.poly_degree)[::-1]

    def compute_drift_factor(self, temperature: float) -> float:
        """주어진 온도에서의 드리프트 보정 계수 계산."""
        dt = temperature - self.reference_temp
        factor = sum(
            c * dt**i for i, c in enumerate(self._coefficients)
        )
        # 물리적으로 합리적인 범위로 클램핑 (0.8 ~ 1.2)
        factor = float(np.clip(factor, 0.8, 1.2))
        self._calibration_history.append((temperature, factor))
        return factor

    def compensate_fingerprint(
        self, fingerprint: np.ndarray, temperature: float
    ) -> np.ndarray:
        """핑거프린트에 드리프트 보정 적용.

        전기화학 채널 (0:528)만 보정하고, 메타데이터 채널은 유지합니다.
        """
        factor = self.compute_drift_factor(temperature)
        compensated = fingerprint.copy()

        # 전기화학 신호 채널만 보정 (0:528)
        compensated[:528] = fingerprint[:528] / factor

        return compensated


class AnomalyDetector:
    """Isolation Forest 기반 OOD (Out-of-Distribution) 탐지.

    비표적 물질이나 오염된 샘플을 식별하여
    잘못된 정량 결과가 출력되는 것을 방지합니다.
    """

    def __init__(
        self,
        contamination: float = 0.05,
        n_estimators: int = 200,
        seed: int = 42,
    ):
        self.model = IsolationForest(
            contamination=contamination,
            n_estimators=n_estimators,
            random_state=seed,
            n_jobs=-1,
        )
        self._is_fitted = False
        self._anomaly_threshold = -0.3  # decision_function 임계값

    def fit(self, X_normal: np.ndarray):
        """정상 샘플로 Isolation Forest 학습.

        Parameters
        ----------
        X_normal : 정상(비이상) 핑거프린트 행렬 (N, 896)
        """
        self.model.fit(X_normal)
        self._is_fitted = True

    def detect(self, fingerprint: np.ndarray) -> Tuple[bool, float]:
        """이상 탐지 수행.

        Returns
        -------
        (is_anomaly, anomaly_score)
          - is_anomaly: True이면 OOD 의심
          - anomaly_score: -1(정상) ~ +1(이상) 범위, 높을수록 이상
        """
        if not self._is_fitted:
            # 미학습 상태: 통계 기반 폴백
            return self._statistical_fallback(fingerprint)

        x = fingerprint.reshape(1, -1)
        # decision_function: 음수면 이상, 양수면 정상
        raw_score = float(self.model.decision_function(x)[0])
        # 정규화: -1(이상) ~ +1(정상) -> 0(정상) ~ 1(이상)으로 변환
        anomaly_score = float(np.clip(-raw_score, 0.0, 1.0))
        is_anomaly = raw_score < self._anomaly_threshold

        return is_anomaly, anomaly_score

    def _statistical_fallback(self, fingerprint: np.ndarray) -> Tuple[bool, float]:
        """Isolation Forest 미학습 시 통계 기반 폴백.

        Mahalanobis 거리 근사 (대각 공분산 가정).
        """
        signal = fingerprint[:88]
        mean = np.mean(signal)
        std = np.std(signal) + 1e-10
        # Z-score 기반 이상 판정
        z_scores = np.abs((signal - mean) / std)
        max_z = float(np.max(z_scores))
        anomaly_score = float(np.clip(max_z / 5.0, 0.0, 1.0))  # 5-sigma 정규화
        is_anomaly = max_z > 4.0  # 4-sigma 이상이면 이상
        return is_anomaly, anomaly_score


class GlucoseQuantifier:
    """XGBoost 기반 혈당 정량 추론 모듈.

    학습된 모델을 로드하여 896차원 핑거프린트로부터
    혈당 농도 (mg/dL)를 예측합니다.
    """

    def __init__(self, model_path: Optional[str] = None):
        self.model: Optional[xgb.XGBRegressor] = None
        self._model_path = model_path

        if model_path and os.path.exists(model_path):
            self.load_model(model_path)

    def load_model(self, model_path: str):
        """XGBoost 모델 파일 로드."""
        if not XGB_AVAILABLE:
            warnings.warn("xgboost 미설치 - 정량 추론 불가")
            return

        self.model = xgb.XGBRegressor()
        self.model.load_model(model_path)
        self._model_path = model_path

    def predict(self, fingerprint: np.ndarray) -> Tuple[float, float]:
        """혈당 정량 예측.

        Returns
        -------
        (predicted_mg_dl, confidence)
        """
        if self.model is None:
            # 모델 미로드 시 물리 기반 폴백
            return self._physics_fallback(fingerprint)

        x = fingerprint.reshape(1, -1)
        predicted = float(self.model.predict(x)[0])

        # 신뢰도 추정: 앙상블 내 트리별 예측 분산 기반
        # XGBoost는 직접적인 불확실성 추정을 제공하지 않으므로
        # 입력 신호 품질에 기반한 휴리스틱 적용
        signal_quality = self._assess_signal_quality(fingerprint)
        confidence = float(np.clip(signal_quality, 0.5, 0.99))

        # 물리적 범위 클램핑 (20 ~ 600 mg/dL)
        predicted = float(np.clip(predicted, 20.0, 600.0))

        return predicted, confidence

    def _physics_fallback(self, fingerprint: np.ndarray) -> Tuple[float, float]:
        """모델 미로드 시 물리 기반 폴백 추정.

        Michaelis-Menten 역변환으로 대략적 농도 추정.
        """
        signal = fingerprint[:88]
        peak = float(np.max(signal))
        km = 50.0
        if peak > 0 and peak < 1.0:
            concentration = km * peak / (1.0 - peak + 1e-10)
        else:
            concentration = float(np.mean(signal) * 100.0 + 80.0)
        concentration = float(np.clip(concentration, 20.0, 600.0))
        return concentration, 0.60  # 낮은 신뢰도

    def _assess_signal_quality(self, fingerprint: np.ndarray) -> float:
        """신호 품질 평가 (SNR 기반)."""
        signal = fingerprint[:88]
        noise_est = fingerprint[792:880]  # 2차 차분 = 고주파 노이즈 추정

        signal_power = float(np.mean(signal**2))
        noise_power = float(np.mean(noise_est**2)) + 1e-10

        snr_db = 10 * np.log10(signal_power / noise_power + 1e-10)
        # SNR 20dB이상 -> 0.95, 0dB -> 0.50
        quality = float(np.clip(0.5 + snr_db / 40.0, 0.5, 0.99))
        return quality


class TritonEnsemblePipeline:
    """Triton Inference Server 앙상블 파이프라인.

    DAG 형태로 3개의 모델 노드를 엮어 추론합니다:
      input -> [DriftDetector] -> [AnomalyDetector] -> [GlucoseQuantifier] -> output
    """

    def __init__(self, models_dir: Optional[str] = None):
        if models_dir is None:
            models_dir = str(
                Path(__file__).resolve().parent.parent / "models" / "fingerprint"
            )

        self.drift_detector = DriftDetector()
        self.anomaly_detector = AnomalyDetector()

        # 모델 파일 로드 시도
        glucose_model_path = os.path.join(models_dir, "glucose_regressor.json")
        self.glucose_quantifier = GlucoseQuantifier(
            model_path=glucose_model_path if os.path.exists(glucose_model_path) else None
        )

    def fit_anomaly_detector(self, X_normal: np.ndarray):
        """이상 탐지기 학습 (정상 데이터 필요)."""
        self.anomaly_detector.fit(X_normal)

    def execute_pipeline(
        self,
        fingerprint: np.ndarray,
        env_temp: float = 25.0,
    ) -> PipelineResult:
        """앙상블 추론 파이프라인 실행.

        Parameters
        ----------
        fingerprint : 896차원 핑거프린트 벡터
        env_temp : 환경 온도 (degC)

        Returns
        -------
        PipelineResult with glucose prediction and quality metrics
        """
        result = PipelineResult(status="error", temperature_c=env_temp)

        # 입력 검증
        if fingerprint.shape[0] != FINGERPRINT_DIM:
            result.warnings.append(
                f"Invalid fingerprint dimension: {fingerprint.shape[0]} "
                f"(expected {FINGERPRINT_DIM})"
            )
            return result

        # Step 1: 드리프트 보정
        drift_factor = self.drift_detector.compute_drift_factor(env_temp)
        compensated = self.drift_detector.compensate_fingerprint(fingerprint, env_temp)
        result.drift_factor = drift_factor

        if abs(drift_factor - 1.0) > 0.05:
            result.warnings.append(
                f"Significant drift compensation applied: factor={drift_factor:.4f} "
                f"at {env_temp:.1f} degC"
            )

        # Step 2: 이상 탐지
        is_anomaly, anomaly_score = self.anomaly_detector.detect(compensated)
        result.anomaly_score = anomaly_score
        result.is_anomaly = is_anomaly

        if is_anomaly:
            result.status = "anomaly_detected"
            result.warnings.append(
                f"OOD sample detected (anomaly_score={anomaly_score:.4f}). "
                f"Quantification unreliable."
            )
            # 이상 탐지 시에도 참고용 정량값은 제공 (H6: 실패 격리)
            glucose, confidence = self.glucose_quantifier.predict(compensated)
            result.glucose_mg_dl = glucose
            result.confidence = confidence * 0.3  # 신뢰도 대폭 하향
            return result

        # Step 3: 혈당 정량
        glucose, confidence = self.glucose_quantifier.predict(compensated)
        result.glucose_mg_dl = glucose
        result.confidence = confidence
        result.status = "success"

        # 위험 수준 경고
        if glucose < 70:
            result.warnings.append(
                f"Low glucose warning: {glucose:.1f} mg/dL (hypoglycemia risk)"
            )
        elif glucose > 180:
            result.warnings.append(
                f"High glucose warning: {glucose:.1f} mg/dL (hyperglycemia risk)"
            )

        return result


if __name__ == "__main__":
    print("ManPaSik Triton Ensemble Pipeline Demo")
    print("=" * 50)

    ensemble = TritonEnsemblePipeline()

    # 정상 데이터로 이상 탐지기 학습
    rng = np.random.default_rng(42)
    normal_data = rng.normal(0.5, 0.1, size=(500, FINGERPRINT_DIM))
    ensemble.fit_anomaly_detector(normal_data)

    # 정상 핑거프린트 테스트
    print("\n--- Normal Sample (25 degC) ---")
    normal_fp = rng.normal(0.5, 0.1, size=(FINGERPRINT_DIM,))
    result = ensemble.execute_pipeline(normal_fp, env_temp=25.0)
    print(json.dumps(result.to_dict(), indent=2, ensure_ascii=False))

    # 고온 환경 테스트
    print("\n--- High Temperature (38 degC) ---")
    result_hot = ensemble.execute_pipeline(normal_fp, env_temp=38.0)
    print(json.dumps(result_hot.to_dict(), indent=2, ensure_ascii=False))

    # 이상 핑거프린트 테스트
    print("\n--- Anomalous Sample ---")
    anomaly_fp = rng.uniform(-5, 5, size=(FINGERPRINT_DIM,))  # 비정상 분포
    result_anomaly = ensemble.execute_pipeline(anomaly_fp, env_temp=25.0)
    print(json.dumps(result_anomaly.to_dict(), indent=2, ensure_ascii=False))
