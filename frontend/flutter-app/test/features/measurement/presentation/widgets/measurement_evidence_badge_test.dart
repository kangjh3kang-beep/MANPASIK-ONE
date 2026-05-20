import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/features/measurement/presentation/widgets/measurement_evidence_badge.dart';

void main() {
  testWidgets('research_only evidence badge는 진단 판정 문구 없이 연구용으로 표시한다',
      (tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: Scaffold(
          body: MeasurementEvidenceBadge(
            evidenceStatus: 'research_only',
            diagnosticReady: false,
            evidenceGaps: ['clinical_lock_required'],
          ),
        ),
      ),
    );

    expect(find.text('연구용'), findsOneWidget);
    expect(find.textContaining('정상'), findsNothing);
    expect(find.textContaining('위험'), findsNothing);
    expect(find.textContaining('진단'), findsNothing);
    expect(find.textContaining('확정'), findsNothing);
  });
}
