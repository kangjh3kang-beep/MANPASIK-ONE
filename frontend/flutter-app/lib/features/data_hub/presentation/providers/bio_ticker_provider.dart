import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/features/data_hub/presentation/providers/monitoring_providers.dart';

// ───────────────────────────────────────────────────
// BioTickerProvider - 생체 데이터(BPM) 기반 글로벌 애니메이션 속도 제어
// [PDCA Design] Bio-Rhythm Sync: 사용자 생체 데이터와 UI 맥동 주기를 동기화하여
// 시스템 전체가 살아있는 유기체처럼 느껴지도록 하는 핵심 모듈.
// ───────────────────────────────────────────────────

final bioTickerSpeedProvider = Provider<double>((ref) {
  // 모니터링 대시보드에 선택된 바이오 데이터 로드
  final bioData = ref.watch(selectedBioDataProvider);
  
  // Pulse(심박수) 데이터가 있으면 파싱
  if (bioData.containsKey('Pulse')) {
    final pulseStr = bioData['Pulse']?.toString() ?? '72';
    final pulse = double.tryParse(RegExp(r'[\d.]+').firstMatch(pulseStr)?.group(0) ?? '72') ?? 72.0;
    
    // 심박수 72BPM을 기준(1.0배속)으로 삼고, 비율에 따라 애니메이션 속도 계산
    // 예: 108BPM -> 1.5배속, 48BPM -> 0.66배속
    // (극단적인 상황 방지를 위해 0.5배속 ~ 2.0배속으로 클램프)
    final speed = (pulse / 72.0).clamp(0.5, 2.0);
    return speed;
  }
  
  return 1.0; // 데이터가 없을 경우 기본 속도
});
