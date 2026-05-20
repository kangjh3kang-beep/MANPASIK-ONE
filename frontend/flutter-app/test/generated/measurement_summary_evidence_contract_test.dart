import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/generated/manpasik.pb.dart';

void main() {
  test('MeasurementSummary preserves evidence fields through protobuf', () {
    final summary = MeasurementSummary(
      sessionId: 'session-history',
      cartridgeType: 'glucose',
      primaryValue: 99.5,
      unit: 'mg/dL',
      evidenceStatus: 'research_only',
      diagnosticReady: false,
      evidenceGaps: ['clinical_lock_required'],
    );

    final decoded = MeasurementSummary.fromBuffer(summary.writeToBuffer());

    expect(decoded.evidenceStatus, 'research_only');
    expect(decoded.diagnosticReady, isFalse);
    expect(decoded.evidenceGaps, contains('clinical_lock_required'));
  });
}
