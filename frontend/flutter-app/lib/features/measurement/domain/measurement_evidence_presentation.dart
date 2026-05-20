class MeasurementEvidencePresentation {
  const MeasurementEvidencePresentation({
    required this.badgeLabel,
    required this.detailText,
  });

  final String badgeLabel;
  final String detailText;

  factory MeasurementEvidencePresentation.from({
    required String evidenceStatus,
    required bool diagnosticReady,
    required List<String> evidenceGaps,
  }) {
    if (evidenceStatus == 'research_only') {
      return const MeasurementEvidencePresentation(
        badgeLabel: '연구용',
        detailText: '참고용 분석 결과입니다. 의료적 판단에는 별도 검증 자료가 필요합니다.',
      );
    }
    if (evidenceStatus == 'clinical_locked' && diagnosticReady) {
      return const MeasurementEvidencePresentation(
        badgeLabel: '임상 검증',
        detailText: '검증된 분석 기준이 적용된 결과입니다.',
      );
    }
    if (evidenceStatus == 'analytical_locked') {
      return const MeasurementEvidencePresentation(
        badgeLabel: '분석 검증',
        detailText: '분석 성능 기준을 검토 중인 결과입니다.',
      );
    }
    return const MeasurementEvidencePresentation(
      badgeLabel: '검증 확인 중',
      detailText: '검증 상태를 확인하는 중입니다.',
    );
  }
}
