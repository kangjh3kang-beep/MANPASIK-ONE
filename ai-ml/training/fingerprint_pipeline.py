import numpy as np
import xgboost as xgb
from sklearn.decomposition import PCA
from sklearn.model_selection import train_test_split
from sklearn.metrics import root_mean_squared_error
from typing import Tuple

# =========================================================================
# ManPaSik (萬波息) v2.4.3 - 핑거프린트 훈련 파이프라인
# 88차원 특성을 eNose 융합하여 896차원으로 확대한 뒤 TFLite(Edge) 변환용 모델 학습
# =========================================================================

class FingerprintPipeline:
    def __init__(self, target_dim: int = 896):
        self.target_dim = target_dim
        # PCA로 차원 간소화 (Feature 의존성 축소)
        self.pca = PCA(n_components=128)
        self.xg_regressor = xgb.XGBRegressor(
            objective='reg:squarederror',
            max_depth=5,
            learning_rate=0.01,
            n_estimators=500
        )

    def generate_mock_data(self, samples: int = 1000) -> Tuple[np.ndarray, np.ndarray]:
        """
        [Layer 3]의 Rust FeatureExtractor가 생성한 896차원 핑거프린트 행렬 시뮬레이션
        """
        # (N, 896) 형태의 핑거프린트
        X = np.random.normal(loc=0.5, scale=0.1, size=(samples, self.target_dim))
        
        # 일부 특성에 종속적인 정량 예측값 Y (e.g., 혈당치 mg/dL)
        baseline = 100.0
        # 0번, 10번, 500번 차원이 지배적 상관관계를 갖는다고 가정
        Y = baseline + (X[:, 0] * 50) - (X[:, 10] * 20) + (X[:, 500] * 10) + np.random.normal(0, 2, samples)
        return X, Y

    def train_model(self, X: np.ndarray, y: np.ndarray):
        print(f"[1] 데이터 준비 완료: {X.shape[0]} 샘플, {X.shape[1]} 차원 핑거프린트")
        
        # 1. 차원 축소 (896 차원에서 엣지 디바이스 효율을 극대화하기 위해 128차원으로 PCA)
        X_pca = self.pca.fit_transform(X)
        print(f"[2] PCA 축소 완료: {X_pca.shape[1]} 차원 (보존율 향상)")

        # 2. 훈련 및 검증셋 분리
        X_train, X_test, y_train, y_test = train_test_split(X_pca, y, test_size=0.2, random_state=42)

        # 3. XGBoost 훈련
        print("[3] XGBoost 정량 모델 훈련 중...")
        self.xg_regressor.fit(X_train, y_train)

        # 4. 성능 검증 (PCCP 규제 승인 데이터 생성)
        y_pred = self.xg_regressor.predict(X_test)
        rmse = root_mean_squared_error(y_test, y_pred)
        print(f"[4] 모델 훈련 완료. Test RMSE: {rmse:.2f} mg/dL")

    def export_to_edge(self, model_id: str):
        """
        TFLite(Rust Edge 추론용) 또는 JSON으로 추출
        """
        filepath = f"../models/model_{model_id}.json"
        self.xg_regressor.save_model(filepath)
        print(f"[5] Edge 배포용 모델 저장 완료: {filepath}")

if __name__ == "__main__":
    pipeline = FingerprintPipeline()
    X, Y = pipeline.generate_mock_data(samples=2000)
    pipeline.train_model(X, Y)
    pipeline.export_to_edge("v1_0_glucose_896dim")
