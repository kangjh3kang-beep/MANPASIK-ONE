import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/features/data_hub/domain/data_hub_repository.dart';
import 'package:manpasik/features/data_hub/presentation/data_hub_screen.dart';
import 'package:manpasik/features/data_hub/presentation/providers/data_hub_providers.dart';

void main() {
  testWidgets('DataHubScreen은 좁은 모바일 폭에서도 evidence header overflow가 없다',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(320, 720));
    addTearDown(() => tester.binding.setSurfaceSize(null));

    const metric = '초장문바이오마커추적항목';
    final now = DateTime(2026, 5, 13, 10, 0);
    const summary = BiomarkerSummary(
      biomarkerType: metric,
      displayName: metric,
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
          dataHubSelectedMetricProvider.overrideWith((ref) => metric),
          dataHubSummariesProvider.overrideWith(
            (ref) => Stream<List<BiomarkerSummary>>.value([summary]),
          ),
          dataHubTrendProvider.overrideWith(
            (ref, params) async => [
              TrendDataPoint(
                timestamp: now.subtract(const Duration(days: 2)),
                value: 101.0,
                unit: 'mg/dL',
                biomarkerType: params.metric,
                isWithinRange: true,
                evidenceStatus: 'research_only',
                diagnosticReady: false,
                evidenceGaps: const ['clinical_lock_required'],
              ),
              TrendDataPoint(
                timestamp: now.subtract(const Duration(days: 1)),
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
    expect(tester.takeException(), isNull);
  });
}
