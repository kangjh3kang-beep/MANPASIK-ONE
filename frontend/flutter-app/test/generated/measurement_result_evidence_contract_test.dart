import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/generated/manpasik.pb.dart';

void main() {
  test('MeasurementResult carries assay evidence contract fields', () {
    final encoded = (MeasurementResult()
          ..sessionId = 'session-evidence-1'
          ..primaryValue = 9.8
          ..unit = 'mg/dL'
          ..confidence = 0.91
          ..evidenceStatus = 'research_only'
          ..diagnosticReady = false
          ..evidenceGaps.add('clinical_lock_required'))
        .writeToBuffer();

    final decoded = MeasurementResult.fromBuffer(encoded);

    expect(decoded.evidenceStatus, 'research_only');
    expect(decoded.diagnosticReady, isFalse);
    expect(decoded.evidenceGaps, contains('clinical_lock_required'));
  });
}
