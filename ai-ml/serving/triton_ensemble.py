import json
import numpy as np

# =========================================================================
# ManPaSik (萬波息) - Triton Inference Server (Ensemble Mock)
# 역할: 여러 AI 모델(보정, 정량, 비표적 이상탐지)을 DAG 형태로 엮어 추론 파이프라인 형성
# =========================================================================

class TritonEnsemblePipeline:
    def __init__(self):
        # 모의 상태: 메모리에 올려둔 3개의 서로 다른 모델 노드
        self.models = {
            "node_1_drift_detector": self._mock_drift_detector,
            "node_2_xgboost_glucose": self._mock_xgboost_glucose,
            "node_3_anomaly_detector": self._mock_anomaly_detector
        }

    def execute_pipeline(self, fingerprint: np.ndarray, env_temp: float) -> dict:
        """
        [Flutter App] -> [Go 백엔드] -> [Triton] 으로 전달된 896차원 핑거프린트 분해 추론
        """
        print(f"[Triton] 앙상블 파이프라인 가동 (입력 차원: {fingerprint.shape})")

        # 1. 이상 탐지 노드 (비표적 물질 감지)
        is_anomaly, anomaly_score = self.models["node_3_anomaly_detector"](fingerprint)
        if is_anomaly:
            print("[Triton] 🚨 치명적 비표적 물질 이상 탐지 -> 추론 중지")
            return {"status": "error", "message": "Unknown substance detected", "anomaly_score": anomaly_score}

        # 2. 드리프트 감지 노드 (온도 기반 보정계수)
        drift_factor = self.models["node_1_drift_detector"](env_temp)
        
        # 3. 메인 타겟(혈당) 정량 추론 모터
        predicted_concentration = self.models["node_2_xgboost_glucose"](fingerprint, drift_factor)
        
        # 4. 종합 응답
        return {
            "status": "success",
            "concentrations": {
                "glucose": round(predicted_concentration, 2)
            },
            "anomaly_score": round(anomaly_score, 4),
            "drift_compensation_applied": True
        }

    # ----- Mock Model Implementations -----
    def _mock_anomaly_detector(self, fingerprint: np.ndarray):
        # Mahalanobis 거리 등 분산 분석 시뮬레이션
        score = float(np.mean(fingerprint) * 0.01)
        # 0.5 이상이면 이상치로 판단
        return (score > 0.5), score

    def _mock_drift_detector(self, temp: float):
        # 25도를 기준으로 한 민감도 변화 시뮬레이션
        offset = (temp - 25.0) * 0.05
        return 1.0 + offset

    def _mock_xgboost_glucose(self, fingerprint: np.ndarray, drift_factor: float):
        # 가상의 정량값 계산 (896차원 벡터 내적 기반)
        base_val = float(np.sum(fingerprint[:10]) * 3.5) + 80.0
        return base_val * drift_factor

if __name__ == "__main__":
    ensemble = TritonEnsemblePipeline()
    dummy_fingerprint = np.random.normal(0.5, 0.1, size=(896,))
    
    result = ensemble.execute_pipeline(dummy_fingerprint, env_temp=30.0)
    print("\n[Triton Pipeline Result]")
    print(json.dumps(result, indent=2, ensure_ascii=False))
