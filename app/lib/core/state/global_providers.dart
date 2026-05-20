// lib/core/state/global_providers.dart
import 'package:flutter_riverpod/flutter_riverpod.dart';

// =========================================================================
// 마스터플랜 UX v1.4.2 대응 전역 상태 관리 모델(Riverpod)
// =========================================================================

/// 1. 디지털 트윈 상태 (리더기 잔차 평가 기반 게이지 상태)
enum TwinHealthGrade { good, degraded, replace, critical }

final twinStatusProvider = StateProvider<TwinHealthGrade>((ref) {
  // 초기값은 Good, 실제로는 Rust API를 지속적(polling)으로 호출하여 갱신
  return TwinHealthGrade.good;
});

/// 2. H-010 / K-070 등 유기적 환경 구성: 컨텍스트 카드 피드 상태
class ContextCard {
  final String id;
  final String type; // 'insight', 'nudge', 'recommendation'
  final String content;
  final String ctaPath; // 딥링크 목적지 (ex: /measure/progress)
  
  ContextCard(this.id, this.type, this.content, this.ctaPath);
}

final contextCardFeedProvider = FutureProvider<List<ContextCard>>((ref) async {
  // Mock 데이터 반환 (향후 DB 혹은 Rust 엔진 판단에 따름)
  await Future.delayed(const Duration(milliseconds: 500));
  return [
    ContextCard('1', 'nudge', '어제 수면이 6시간 미만이네요. 영양 섭취에 유의하세요.', '/coach/chat'),
    ContextCard('2', 'recommendation', 'AI 추천: 혈당 강하율 15% 상승 효과 제품', '/market/product/123'),
  ];
});

/// 3. 측정 전선 흐름 제어(Measurement State)
enum MeasurementPhase { IDLE, PREFLIGHT_CHECK, SCANNING_NFC, PREPARING, MEASURING, INFERENCE, COMPLETED, ERROR }

class MeasurementState {
  final MeasurementPhase phase;
  final double progressPercent;
  final String? errorMessage;
  
  MeasurementState({required this.phase, this.progressPercent = 0.0, this.errorMessage});
}

final measurementStateProvider = StateProvider<MeasurementState>((ref) {
  return MeasurementState(phase: MeasurementPhase.IDLE);
});
