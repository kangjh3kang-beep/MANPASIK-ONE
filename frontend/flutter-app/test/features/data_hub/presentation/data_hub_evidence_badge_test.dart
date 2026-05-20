import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/features/data_hub/domain/data_hub_repository.dart';
import 'package:manpasik/features/data_hub/presentation/data_hub_screen.dart';
import 'package:manpasik/features/data_hub/presentation/providers/data_hub_providers.dart';

void main() {
  testWidgets('DataHubScreen은 research_only summary를 연구용 배지로 표시한다',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(390, 844));
    addTearDown(() => tester.binding.setSurfaceSize(null));

    final now = DateTime(2026, 5, 13, 10, 0);
    const summary = BiomarkerSummary(
      biomarkerType: '혈당',
      displayName: '혈당',
      unit: 'mg/dL',
      latestValue: 99.5,
      averageValue: 98.2,
      minValue: 93,
      maxValue: 104,
      referenceMin: 70,
      referenceMax: 140,
      totalMeasurements: 3,
      trend: 'stable',
      latestEvidenceStatus: 'research_only',
      latestDiagnosticReady: false,
      latestEvidenceGaps: ['clinical_lock_required'],
    );

    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          dataHubSummariesProvider.overrideWith(
            (ref) => Stream<List<BiomarkerSummary>>.value([summary]),
          ),
          dataHubTrendProvider.overrideWith(
            (ref, params) async => [
              TrendDataPoint(
                timestamp: now,
                value: 99.5,
                unit: 'mg/dL',
                biomarkerType: params.metric,
                isWithinRange: true,
                evidenceStatus: 'research_only',
                diagnosticReady: false,
                evidenceGaps: const ['clinical_lock_required'],
              ),
            ],
          ),
        ],
        child: const MaterialApp(
          home: DataHubScreen(),
        ),
      ),
    );

    await tester.pump();
    await tester.pump(const Duration(milliseconds: 50));

    expect(find.text('연구용'), findsAtLeastNWidgets(1));
    expect(find.textContaining('정상'), findsNothing);
    expect(find.textContaining('위험'), findsNothing);
    expect(find.textContaining('진단'), findsNothing);
    expect(find.textContaining('확정'), findsNothing);
  });
}
