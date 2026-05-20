import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/measurement/domain/measurement_evidence_presentation.dart';

void main() {
  group('MeasurementEvidencePresentation', () {
    test('research_only는 진단 판정처럼 보이지 않는 참고용 문구를 반환한다', () {
      final copy = MeasurementEvidencePresentation.from(
        evidenceStatus: 'research_only',
        diagnosticReady: false,
        evidenceGaps: const ['clinical_lock_required'],
      );

      expect(copy.badgeLabel, '연구용');
      expect(copy.detailText, contains('참고용'));
      expect(copy.detailText, isNot(contains('정상')));
      expect(copy.detailText, isNot(contains('위험')));
      expect(copy.detailText, isNot(contains('확정')));
    });

    test('unknown 상태는 검증 확인 중으로 표시한다', () {
      final copy = MeasurementEvidencePresentation.from(
        evidenceStatus: 'unknown',
        diagnosticReady: false,
        evidenceGaps: const [],
      );

      expect(copy.badgeLabel, '검증 확인 중');
      expect(copy.detailText, contains('검증 상태'));
    });
  });
}
