"""
ManPaSik (萬波息) v2.4.3 - 핑거프린트 훈련 파이프라인
=====================================================

차동측정 물리 기반 구조화된 합성 896차원 핑거프린트 데이터 생성,
5-fold 교차검증, XGBoost 회귀 (혈당 정량) + RandomForest 분류 (30종 바이오마커),
ONNX 형식 모델 내보내기.

핑거프린트 구조 (896차원):
  - [0:88]     전기화학 기본 차분신호 (Sdiff = S_det - alpha * S_ref, alpha=0.98)
  - [88:176]   시간 미분 (dS/dt) — 반응 속도
  - [176:264]  주파수 영역 (FFT 크기) — 노이즈/진동 패턴
  - [264:352]  교류 임피던스 (EIS) 실수부
  - [352:440]  교류 임피던스 (EIS) 허수부
  - [440:528]  온도 보정 계수
  - [528:616]  비선형 커널 변환 (RBF)
  - [616:704]  채널 간 상관 행렬 (상삼각 → flatten)
  - [704:792]  이동평균 필터 (3-tap)
  - [792:880]  2차 차분 (d²S/dt²)
  - [880:896]  환경 메타데이터 (온도, 습도, 배터리, 시간 등)
"""

import os
import json
import warnings
from pathlib import Path
from typing import Tuple, Dict, Any

import numpy as np
from sklearn.model_selection import StratifiedKFold, KFold
from sklearn.metrics import (
    root_mean_squared_error,
    accuracy_score,
    confusion_matrix,
    classification_report,
)
from sklearn.ensemble import RandomForestClassifier
import xgboost as xgb

# ONNX 내보내기 (선택 의존성)
try:
    import onnx
    from skl2onnx import convert_sklearn
    from skl2onnx.common.data_types import FloatTensorType

    ONNX_AVAILABLE = True
except ImportError:
    ONNX_AVAILABLE = False
    warnings.warn("skl2onnx 미설치 — ONNX 내보내기 비활성화. pip install skl2onnx onnx")

# SSOT 상수 (CLAUDE.md 베이스라인)
ALPHA_DEFAULT = 0.98
ALPHA_MIN = 0.90
ALPHA_MAX = 1.10
FINGERPRINT_DIM = 896
NUM_BIOMARKER_CLASSES = 30

# 바이오마커 클래스 정의 (29종 + NonTarget)
BIOMARKER_CLASSES = [
    "glucose",
    "hba1c",
    "uric_acid",
    "creatinine",
    "vitamin_d",
    "tsh",
    "cortisol",
    "crp",
    "cholesterol_total",
    "hdl",
    "ldl",
    "triglycerides",
    "alt",
    "ast",
    "bilirubin",
    "albumin",
    "bun",
    "calcium",
    "potassium",
    "sodium",
    "chloride",
    "iron",
    "ferritin",
    "testosterone",
    "estradiol",
    "progesterone",
    "psa",
    "bnp",
    "troponin_i",
    "non_target",
]

# 바이오마커별 정상 농도 범위 (mg/dL 또는 해당 단위)
BIOMARKER_RANGES: Dict[str, Tuple[float, float]] = {
    "glucose": (70.0, 200.0),
    "hba1c": (4.0, 10.0),
    "uric_acid": (2.0, 10.0),
    "creatinine": (0.5, 5.0),
    "vitamin_d": (10.0, 80.0),
    "tsh": (0.1, 10.0),
    "cortisol": (3.0, 30.0),
    "crp": (0.0, 50.0),
}


class FingerprintDataGenerator:
    """차동측정 물리 기반 구조화된 합성 핑거프린트 데이터 생성기.

    랜덤 데이터가 아닌, 전기화학 센서의 물리적 특성을 반영한
    구조화된 패턴을 생성합니다.
    """

    def __init__(self, seed: int = 42):
        self.rng = np.random.default_rng(seed)
        # 각 바이오마커 클래스별 고유 시그니처 벡터 (88차원)
        self._class_signatures = self._generate_class_signatures()

    def _generate_class_signatures(self) -> np.ndarray:
        """바이오마커별 고유 전기화학 시그니처 생성.

        실제 시스템에서는 각 바이오마커가 다른 산화환원 전위와
        전류 패턴을 보이므로, 이를 시뮬레이션합니다.
        """
        rng = np.random.default_rng(0)  # 재현성을 위한 고정 시드
        signatures = np.zeros((NUM_BIOMARKER_CLASSES, 88))

        for cls in range(NUM_BIOMARKER_CLASSES):
            # 각 클래스는 3~8개의 주요 채널에서 높은 반응을 보임
            num_active = rng.integers(3, 9)
            active_channels = rng.choice(88, size=num_active, replace=False)
            amplitudes = rng.uniform(0.5, 2.0, size=num_active)

            for ch, amp in zip(active_channels, amplitudes):
                # 가우시안 프로파일로 인접 채널에도 영향
                for offset in range(-3, 4):
                    idx = (ch + offset) % 88
                    signatures[cls, idx] += amp * np.exp(-0.5 * (offset / 1.5) ** 2)

        return signatures

    def generate_regression_data(
        self, n_samples: int = 2000
    ) -> Tuple[np.ndarray, np.ndarray]:
        """혈당 정량 회귀용 데이터 생성 (896차원 -> 혈당 mg/dL).

        차동측정 공식: Sdiff_n = S_n - alpha * R_n
        """
        X = np.zeros((n_samples, FINGERPRINT_DIM))
        y = np.zeros(n_samples)

        for i in range(n_samples):
            # 혈당 농도 (mg/dL): 실제 분포를 반영한 로그정규 분포
            concentration = self.rng.lognormal(mean=4.7, sigma=0.3)  # ~110 mg/dL 중앙
            concentration = np.clip(concentration, 40, 400)
            y[i] = concentration

            # 1. 기본 차분신호 (0:88) — 농도에 비례하는 전류 반응
            alpha = ALPHA_DEFAULT + self.rng.normal(0, 0.02)
            alpha = np.clip(alpha, ALPHA_MIN, ALPHA_MAX)
            signal = self._generate_electrochemical_signal(concentration, alpha)
            X[i, 0:88] = signal

            # 2. 시간 미분 (88:176)
            X[i, 88:176] = np.gradient(signal)

            # 3. FFT 크기 (176:264)
            fft_vals = np.abs(np.fft.fft(signal))
            X[i, 176:264] = fft_vals

            # 4. EIS 실수부 (264:352)
            X[i, 264:352] = self._generate_eis_real(concentration, signal)

            # 5. EIS 허수부 (352:440)
            X[i, 352:440] = self._generate_eis_imag(concentration, signal)

            # 6. 온도 보정 (440:528)
            temp = self.rng.normal(25.0, 3.0)  # 환경 온도
            X[i, 440:528] = signal * (1.0 + 0.02 * (temp - 25.0))

            # 7. RBF 커널 변환 (528:616)
            gamma = 0.01
            ref_point = self._class_signatures[0]  # glucose 시그니처
            X[i, 528:616] = np.exp(-gamma * (signal - ref_point) ** 2)

            # 8. 채널 간 상관 (616:704) — 인접 채널 공분산
            for j in range(88):
                X[i, 616 + j] = signal[j] * signal[(j + 1) % 88]

            # 9. 이동평균 (704:792)
            X[i, 704:792] = np.convolve(signal, np.ones(3) / 3, mode="same")

            # 10. 2차 차분 (792:880)
            X[i, 792:880] = np.gradient(np.gradient(signal))

            # 11. 환경 메타데이터 (880:896)
            X[i, 880] = temp / 50.0  # 정규화된 온도
            X[i, 881] = self.rng.uniform(0.3, 0.9)  # 습도
            X[i, 882] = self.rng.uniform(0.7, 1.0)  # 배터리
            X[i, 883] = self.rng.uniform(0, 1)  # 시간 (0~24h 정규화)
            X[i, 884] = alpha  # 사용된 alpha 값
            X[i, 885:896] = self.rng.normal(0, 0.01, size=11)  # 패딩

        return X, y

    def generate_classification_data(
        self, n_samples: int = 3000
    ) -> Tuple[np.ndarray, np.ndarray]:
        """30종 바이오마커 분류용 데이터 생성."""
        samples_per_class = n_samples // NUM_BIOMARKER_CLASSES
        X = np.zeros((n_samples, FINGERPRINT_DIM))
        y = np.zeros(n_samples, dtype=int)

        idx = 0
        for cls in range(NUM_BIOMARKER_CLASSES):
            for _ in range(samples_per_class):
                if idx >= n_samples:
                    break

                y[idx] = cls
                signature = self._class_signatures[cls]

                # 농도 변동 (클래스별)
                concentration = self.rng.uniform(50, 200)

                # 기본 신호 = 시그니처 * 농도 스케일 + 노이즈
                signal = signature * (concentration / 100.0) + self.rng.normal(
                    0, 0.05, size=88
                )
                X[idx, 0:88] = signal

                # 나머지 차원도 구조화
                X[idx, 88:176] = np.gradient(signal)
                X[idx, 176:264] = np.abs(np.fft.fft(signal))
                X[idx, 264:352] = signal * self.rng.uniform(0.8, 1.2)
                X[idx, 352:440] = signal * self.rng.uniform(-0.3, 0.3)
                X[idx, 440:528] = signal * (1.0 + self.rng.normal(0, 0.02))
                gamma = 0.01
                X[idx, 528:616] = np.exp(-gamma * (signal - signature) ** 2)
                for j in range(88):
                    X[idx, 616 + j] = signal[j] * signal[(j + 1) % 88]
                X[idx, 704:792] = np.convolve(signal, np.ones(3) / 3, mode="same")
                X[idx, 792:880] = np.gradient(np.gradient(signal))
                X[idx, 880:896] = self.rng.normal(0.5, 0.1, size=16)

                idx += 1

        # 남은 슬롯 채우기
        while idx < n_samples:
            cls = self.rng.integers(0, NUM_BIOMARKER_CLASSES)
            y[idx] = cls
            X[idx] = X[self.rng.integers(0, idx)] + self.rng.normal(0, 0.02, size=FINGERPRINT_DIM)
            idx += 1

        # 셔플
        perm = self.rng.permutation(n_samples)
        return X[perm], y[perm]

    def _generate_electrochemical_signal(
        self, concentration: float, alpha: float
    ) -> np.ndarray:
        """차동측정 기반 전기화학 신호 생성.

        Sdiff_n = S_det_n - alpha * S_ref_n
        """
        channels = np.arange(88, dtype=float)

        # 검출 전극 신호: Michaelis-Menten 동역학 모사
        km = 50.0  # Michaelis 상수
        vmax = concentration / (concentration + km)

        # 채널별 감도 프로파일 (가우시안 혼합)
        s_det = vmax * (
            0.7 * np.exp(-0.5 * ((channels - 20) / 8) ** 2)
            + 0.3 * np.exp(-0.5 * ((channels - 60) / 12) ** 2)
        )
        s_det += self.rng.normal(0, 0.02, size=88)

        # 참조 전극 신호: 농도 독립적 기준선
        s_ref = 0.1 * np.sin(2 * np.pi * channels / 88) + self.rng.normal(0, 0.01, size=88)

        # 차분식 적용
        s_diff = s_det - alpha * s_ref
        return s_diff

    def _generate_eis_real(
        self, concentration: float, base_signal: np.ndarray
    ) -> np.ndarray:
        """전기화학 임피던스 분광법 (EIS) 실수부 생성."""
        freq = np.logspace(0, 5, 88)
        r_sol = 100 + self.rng.normal(0, 5)  # 용액 저항
        r_ct = 1000 / (1 + concentration / 100)  # 전하전달 저항 (농도 의존)
        omega = 2 * np.pi * freq
        cdl = 1e-6  # 이중층 커패시턴스

        # Randles 등가회로 실수부
        z_real = r_sol + r_ct / (1 + (omega * r_ct * cdl) ** 2)
        z_real = z_real / z_real.max()  # 정규화
        z_real += self.rng.normal(0, 0.01, size=88)
        return z_real

    def _generate_eis_imag(
        self, concentration: float, base_signal: np.ndarray
    ) -> np.ndarray:
        """EIS 허수부 생성."""
        freq = np.logspace(0, 5, 88)
        r_ct = 1000 / (1 + concentration / 100)
        omega = 2 * np.pi * freq
        cdl = 1e-6

        z_imag = -omega * r_ct**2 * cdl / (1 + (omega * r_ct * cdl) ** 2)
        z_imag = z_imag / (np.abs(z_imag).max() + 1e-10)
        z_imag += self.rng.normal(0, 0.01, size=88)
        return z_imag


class FingerprintPipeline:
    """핑거프린트 훈련 파이프라인 (XGBoost 회귀 + RandomForest 분류)."""

    def __init__(
        self,
        output_dir: str = None,
        n_folds: int = 5,
        seed: int = 42,
    ):
        if output_dir is None:
            output_dir = str(
                Path(__file__).resolve().parent.parent / "models" / "fingerprint"
            )
        self.output_dir = Path(output_dir)
        self.output_dir.mkdir(parents=True, exist_ok=True)
        self.n_folds = n_folds
        self.seed = seed
        self.generator = FingerprintDataGenerator(seed=seed)

        # 모델 인스턴스
        self.xgb_regressor = xgb.XGBRegressor(
            objective="reg:squarederror",
            max_depth=6,
            learning_rate=0.05,
            n_estimators=300,
            subsample=0.8,
            colsample_bytree=0.8,
            reg_alpha=0.1,
            reg_lambda=1.0,
            random_state=seed,
            n_jobs=-1,
        )
        self.rf_classifier = RandomForestClassifier(
            n_estimators=200,
            max_depth=20,
            min_samples_split=5,
            min_samples_leaf=2,
            max_features="sqrt",
            random_state=seed,
            n_jobs=-1,
        )

    def run(self):
        """전체 파이프라인 실행."""
        print("=" * 72)
        print("ManPaSik Fingerprint Training Pipeline")
        print(f"  Fingerprint dim: {FINGERPRINT_DIM}")
        print(f"  Classes: {NUM_BIOMARKER_CLASSES}")
        print(f"  Cross-validation folds: {self.n_folds}")
        print(f"  Output dir: {self.output_dir}")
        print("=" * 72)

        # 1. 회귀 모델 (혈당 정량)
        print("\n[Phase 1] XGBoost Glucose Quantification")
        print("-" * 48)
        self._train_regressor()

        # 2. 분류 모델 (바이오마커 분류)
        print("\n[Phase 2] RandomForest Biomarker Classification")
        print("-" * 48)
        self._train_classifier()

        # 3. ONNX 내보내기
        print("\n[Phase 3] ONNX Export")
        print("-" * 48)
        self._export_models()

        # 4. 메타데이터 저장
        self._save_metadata()

        print("\n" + "=" * 72)
        print("Pipeline complete.")
        print("=" * 72)

    def _train_regressor(self):
        """XGBoost 회귀 모델 훈련 (5-fold CV)."""
        X, y = self.generator.generate_regression_data(n_samples=2000)
        print(f"  Data: {X.shape[0]} samples, {X.shape[1]} features")
        print(f"  Target range: {y.min():.1f} ~ {y.max():.1f} mg/dL")

        kf = KFold(n_splits=self.n_folds, shuffle=True, random_state=self.seed)
        fold_rmses = []

        for fold, (train_idx, val_idx) in enumerate(kf.split(X), 1):
            X_train, X_val = X[train_idx], X[val_idx]
            y_train, y_val = y[train_idx], y[val_idx]

            model = xgb.XGBRegressor(**self.xgb_regressor.get_params())
            model.fit(
                X_train,
                y_train,
                eval_set=[(X_val, y_val)],
                verbose=False,
            )

            y_pred = model.predict(X_val)
            rmse = root_mean_squared_error(y_val, y_pred)
            fold_rmses.append(rmse)
            print(f"  Fold {fold}/{self.n_folds}: RMSE = {rmse:.2f} mg/dL")

        mean_rmse = np.mean(fold_rmses)
        std_rmse = np.std(fold_rmses)
        print(f"  CV RMSE: {mean_rmse:.2f} +/- {std_rmse:.2f} mg/dL")

        # 전체 데이터로 최종 모델 훈련
        print("  Training final model on full dataset...")
        self.xgb_regressor.fit(X, y, verbose=False)

        # 최종 성능 확인 (전체 데이터 — 과적합 바운드 확인용)
        y_full_pred = self.xgb_regressor.predict(X)
        full_rmse = root_mean_squared_error(y, y_full_pred)
        print(f"  Final model train RMSE: {full_rmse:.2f} mg/dL")

    def _train_classifier(self):
        """RandomForest 분류 모델 훈련 (5-fold Stratified CV)."""
        X, y = self.generator.generate_classification_data(n_samples=3000)
        print(f"  Data: {X.shape[0]} samples, {X.shape[1]} features")
        print(f"  Classes: {len(np.unique(y))}")

        skf = StratifiedKFold(
            n_splits=self.n_folds, shuffle=True, random_state=self.seed
        )
        fold_accs = []

        for fold, (train_idx, val_idx) in enumerate(skf.split(X, y), 1):
            X_train, X_val = X[train_idx], X[val_idx]
            y_train, y_val = y[train_idx], y[val_idx]

            model = RandomForestClassifier(**self.rf_classifier.get_params())
            model.fit(X_train, y_train)

            y_pred = model.predict(X_val)
            acc = accuracy_score(y_val, y_pred)
            fold_accs.append(acc)
            print(f"  Fold {fold}/{self.n_folds}: Accuracy = {acc:.4f}")

        mean_acc = np.mean(fold_accs)
        std_acc = np.std(fold_accs)
        print(f"  CV Accuracy: {mean_acc:.4f} +/- {std_acc:.4f}")

        # 전체 데이터로 최종 모델 훈련
        print("  Training final model on full dataset...")
        self.rf_classifier.fit(X, y)

        # 혼동 행렬 요약
        y_full_pred = self.rf_classifier.predict(X)
        cm = confusion_matrix(y, y_full_pred)
        per_class_acc = cm.diagonal() / cm.sum(axis=1).clip(min=1)
        print(f"  Final model train accuracy: {accuracy_score(y, y_full_pred):.4f}")
        print(f"  Per-class accuracy: min={per_class_acc.min():.3f}, "
              f"max={per_class_acc.max():.3f}, mean={per_class_acc.mean():.3f}")

        # 하위 5개 클래스 출력
        worst_classes = np.argsort(per_class_acc)[:5]
        print("  Lowest performing classes:")
        for ci in worst_classes:
            name = BIOMARKER_CLASSES[ci] if ci < len(BIOMARKER_CLASSES) else f"class_{ci}"
            print(f"    {name}: {per_class_acc[ci]:.3f}")

    def _export_models(self):
        """모델을 ONNX 형식으로 내보내기."""
        # XGBoost JSON 저장 (기본 — 항상 가능)
        xgb_json_path = self.output_dir / "glucose_regressor.json"
        self.xgb_regressor.save_model(str(xgb_json_path))
        print(f"  XGBoost JSON saved: {xgb_json_path}")

        if not ONNX_AVAILABLE:
            print("  ONNX export skipped (skl2onnx not installed)")
            return

        # XGBoost -> ONNX (xgboost 자체 ONNX 미지원이므로 JSON만 저장)
        # RandomForest -> ONNX
        try:
            initial_type = [("float_input", FloatTensorType([None, FINGERPRINT_DIM]))]
            onnx_model = convert_sklearn(
                self.rf_classifier,
                "biomarker_classifier",
                initial_types=initial_type,
            )
            rf_onnx_path = self.output_dir / "biomarker_classifier.onnx"
            onnx.save_model(onnx_model, str(rf_onnx_path))
            print(f"  RF ONNX saved: {rf_onnx_path}")
        except Exception as e:
            print(f"  RF ONNX export failed: {e}")

    def _save_metadata(self):
        """훈련 메타데이터 저장."""
        metadata = {
            "pipeline_version": "2.4.3",
            "fingerprint_dim": FINGERPRINT_DIM,
            "num_classes": NUM_BIOMARKER_CLASSES,
            "biomarker_classes": BIOMARKER_CLASSES,
            "alpha_default": ALPHA_DEFAULT,
            "alpha_range": [ALPHA_MIN, ALPHA_MAX],
            "regressor": {
                "type": "XGBRegressor",
                "target": "glucose_mg_dl",
                "n_estimators": self.xgb_regressor.get_params()["n_estimators"],
                "max_depth": self.xgb_regressor.get_params()["max_depth"],
            },
            "classifier": {
                "type": "RandomForestClassifier",
                "n_estimators": self.rf_classifier.get_params()["n_estimators"],
                "max_depth": self.rf_classifier.get_params()["max_depth"],
            },
            "cv_folds": self.n_folds,
            "seed": self.seed,
        }
        meta_path = self.output_dir / "pipeline_metadata.json"
        with open(meta_path, "w", encoding="utf-8") as f:
            json.dump(metadata, f, indent=2, ensure_ascii=False)
        print(f"  Metadata saved: {meta_path}")


if __name__ == "__main__":
    pipeline = FingerprintPipeline()
    pipeline.run()
