import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:manpasik/features/home/presentation/home_screen.dart';
import 'package:manpasik/features/measurement/domain/measurement_repository.dart';
import 'package:manpasik/shared/providers/ecosystem_providers.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  setUp(() {
    SharedPreferences.setMockInitialValues({});
  });

  testWidgets('HomeScreen은 최근 측정 research_only evidence를 연구용 배지로 표시한다',
      (tester) async {
    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          homeDashboardProvider.overrideWithValue(
            HomeDashboardData(
              healthScore: 92,
              latestMeasurement: MeasurementHistoryItem(
                sessionId: 'session-home',
                cartridgeType: 'glucose',
                primaryValue: 99.5,
                unit: 'mg/dL',
                measuredAt: DateTime.now(),
                evidenceStatus: 'research_only',
                diagnosticReady: false,
                evidenceGaps: const ['clinical_lock_required'],
              ),
            ),
          ),
          globalBadgeProvider.overrideWithValue(const {}),
        ],
        child: const MaterialApp(
          home: HomeScreen(),
        ),
      ),
    );

    await tester.pump();

    expect(find.text('연구용'), findsOneWidget);
    expect(find.textContaining('정상'), findsNothing);
    expect(find.textContaining('위험'), findsNothing);
    expect(find.textContaining('진단'), findsNothing);
    expect(find.textContaining('확정'), findsNothing);
  });
}
