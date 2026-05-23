"""
ManPaSik 모델 레지스트리 — MLOps 모델 버전 관리
- 모델 버전 형식: MPS-{ModelType}-v{Major}.{Minor}.{Patch}
- A/B 테스트 지원
- 드리프트 모니터링 연동
- PCCP 변경 범위 추적

사용법:
    registry = ModelRegistry("./models")
    registry.register("calibration", "1.0.0", model_path, metrics)
    active = registry.get_active("calibration")
"""

import json
import os
import hashlib
from dataclasses import dataclass, asdict, field
from datetime import datetime, timezone
from typing import Dict, List, Optional


@dataclass
class ModelVersion:
    """모델 버전 메타데이터"""
    model_type: str
    version: str
    model_id: str  # MPS-{type}-v{version}
    path: str
    format: str  # "onnx", "tflite", "xgboost_json"
    created_at: str
    metrics: Dict[str, float] = field(default_factory=dict)
    training_data_hash: str = ""
    bias_evaluation: Dict[str, float] = field(default_factory=dict)
    pccp_change_type: str = "none"  # "none", "SPS", "ACP"
    status: str = "staged"  # "staged", "active", "deprecated", "rolled_back"
    checksum: str = ""


@dataclass
class DriftAlert:
    """드리프트 경고"""
    model_id: str
    metric: str
    baseline_value: float
    current_value: float
    threshold: float
    severity: str  # "warning", "critical"
    timestamp: str


class ModelRegistry:
    """MLOps 모델 레지스트리"""

    def __init__(self, base_dir: str = "./models"):
        self.base_dir = base_dir
        self.registry_path = os.path.join(base_dir, "registry.json")
        self.versions: Dict[str, List[ModelVersion]] = {}
        self.drift_alerts: List[DriftAlert] = []
        self._load()

    def _load(self):
        """레지스트리 로드"""
        if os.path.exists(self.registry_path):
            with open(self.registry_path, 'r') as f:
                data = json.load(f)
                for model_type, versions in data.get("versions", {}).items():
                    self.versions[model_type] = [ModelVersion(**v) for v in versions]

    def _save(self):
        """레지스트리 저장"""
        os.makedirs(self.base_dir, exist_ok=True)
        data = {
            "versions": {k: [asdict(v) for v in vs] for k, vs in self.versions.items()},
            "updated_at": datetime.now(timezone.utc).isoformat(),
        }
        with open(self.registry_path, 'w') as f:
            json.dump(data, f, indent=2, ensure_ascii=False)

    def register(self, model_type: str, version: str, model_path: str,
                 metrics: Dict[str, float], model_format: str = "onnx",
                 pccp_change_type: str = "none",
                 bias_evaluation: Optional[Dict[str, float]] = None) -> ModelVersion:
        """새 모델 버전 등록"""
        model_id = f"MPS-{model_type}-v{version}"

        # 체크섬
        checksum = ""
        if os.path.exists(model_path):
            with open(model_path, 'rb') as f:
                checksum = hashlib.sha256(f.read()).hexdigest()

        mv = ModelVersion(
            model_type=model_type,
            version=version,
            model_id=model_id,
            path=model_path,
            format=model_format,
            created_at=datetime.now(timezone.utc).isoformat(),
            metrics=metrics,
            bias_evaluation=bias_evaluation or {},
            pccp_change_type=pccp_change_type,
            status="staged",
            checksum=checksum,
        )

        if model_type not in self.versions:
            self.versions[model_type] = []
        self.versions[model_type].append(mv)
        self._save()
        return mv

    def promote(self, model_type: str, version: str) -> ModelVersion:
        """모델을 active로 승격 (이전 active는 deprecated)"""
        versions = self.versions.get(model_type, [])
        target = None
        for v in versions:
            if v.version == version:
                target = v
            elif v.status == "active":
                v.status = "deprecated"

        if target is None:
            raise ValueError(f"버전 {version} 없음")

        target.status = "active"
        self._save()
        return target

    def rollback(self, model_type: str) -> Optional[ModelVersion]:
        """이전 버전으로 롤백 (H6 실패 격리)"""
        versions = self.versions.get(model_type, [])
        deprecated = [v for v in versions if v.status == "deprecated"]
        if not deprecated:
            return None

        # 가장 최근 deprecated를 active로
        latest = sorted(deprecated, key=lambda v: v.created_at, reverse=True)[0]

        for v in versions:
            if v.status == "active":
                v.status = "rolled_back"

        latest.status = "active"
        self._save()
        return latest

    def get_active(self, model_type: str) -> Optional[ModelVersion]:
        """현재 활성 모델 버전 반환"""
        for v in self.versions.get(model_type, []):
            if v.status == "active":
                return v
        return None

    def check_drift(self, model_type: str, current_metrics: Dict[str, float],
                    threshold_pct: float = 5.0) -> List[DriftAlert]:
        """드리프트 모니터링 — 기준 대비 성능 하락 탐지"""
        active = self.get_active(model_type)
        if not active:
            return []

        alerts = []
        for metric, baseline in active.metrics.items():
            current = current_metrics.get(metric)
            if current is None or baseline == 0:
                continue

            drop_pct = ((baseline - current) / baseline) * 100
            if drop_pct > threshold_pct:
                severity = "critical" if drop_pct > threshold_pct * 2 else "warning"
                alert = DriftAlert(
                    model_id=active.model_id,
                    metric=metric,
                    baseline_value=baseline,
                    current_value=current,
                    threshold=threshold_pct,
                    severity=severity,
                    timestamp=datetime.now(timezone.utc).isoformat(),
                )
                alerts.append(alert)
                self.drift_alerts.append(alert)

        return alerts


if __name__ == "__main__":
    registry = ModelRegistry("/tmp/manpasik_models")

    # 모델 등록
    v1 = registry.register(
        "calibration", "1.0.0", "/models/calibration/model.onnx",
        metrics={"rmse": 0.023, "r2": 0.97},
        bias_evaluation={"age_gap": 0.02, "gender_gap": 0.01},
    )
    print(f"등록: {v1.model_id}")

    # 승격
    registry.promote("calibration", "1.0.0")
    active = registry.get_active("calibration")
    print(f"활성: {active.model_id if active else 'none'}")

    # 드리프트 체크
    alerts = registry.check_drift("calibration", {"rmse": 0.035, "r2": 0.92})
    for a in alerts:
        print(f"⚠ 드리프트: {a.metric} {a.baseline_value:.3f} → {a.current_value:.3f} ({a.severity})")
