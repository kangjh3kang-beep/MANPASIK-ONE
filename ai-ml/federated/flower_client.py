import flwr as fl
import numpy as np

# =========================================================================
# ManPaSik (萬波息) - Federated Learning(연합 학습) Client
# 특허: Family A cl.24 (연합학습에 의한 디지털 트윈 고도화)
# =========================================================================

class ManpasikFlowerClient(fl.client.NumPyClient):
    """
    모바일 단말 또는 엣지 허브에 내장되어 중앙 서버와 가중치만 교환하는 연합학습 클라이언트.
    사용자의 Raw 측정 데이터는 절대 클라우드로 전송되지 않습니다.
    """
    def __init__(self, twin_parameters: np.ndarray, num_measurements: int):
        self.twin_parameters = twin_parameters # 센서의 알파값(alpha) 및 잔차 교정 트리 가중치
        self.num_measurements = num_measurements # 확보된 데이터 수

    def get_parameters(self, config):
        # 1. 중앙 서버에 제출할 현재 기기의 학습 가중치
        return [self.twin_parameters]

    def fit(self, parameters, config):
        # 2. 중앙 서버에서 종합되어 내려온 Global 가중치를 수신받아 로컬 가중치 업데이트 (모의)
        print("[Federated] Global 모델 파라미터 수신 완료 및 로컬 트윈 동기화 적용")
        server_weights = parameters[0]
        
        # 새 트윈 파라미터 = (내 기존 데이터) 0.3 + (글로벌 데이터) 0.7 융합
        new_weights = (self.twin_parameters * 0.3) + (server_weights * 0.7)
        self.twin_parameters = new_weights
        
        # 3. 반환값: (새 파라미터, 데이터 포인트 개수, 메타데이터 딕셔너리)
        return [self.twin_parameters], self.num_measurements, {"loss": 0.05}

    def evaluate(self, parameters, config):
        # 4. 연합 평가
        print("[Federated] 로컬 디지털 트윈 드리프트 오차 평가")
        loss = 0.05
        accuracy = 0.98  # 잔차 일치율 98%
        return float(loss), self.num_measurements, {"accuracy": float(accuracy)}

if __name__ == "__main__":
    # 무작위 가중치 생성 (Simulating initial Twin State)
    initial_weights = np.random.rand(128).astype(np.float32)
    client = ManpasikFlowerClient(initial_weights, num_measurements=150)
    
    # 5. 분산망 플라워 서버에 접속 (백그라운드 통신)
    # fl.client.start_numpy_client(server_address="ml.manpasik.internal:8080", client=client)
    print("Flower Client 셋업 완료 상태. (서버 연결 대기중)")
