import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/generated/manpasik.pb.dart';

const _measurementDataGoldenHex = '0a0973657373696f6e2d31'
    '1218000000000000f03f00000000000000400000000000000840'
    '1a24'
    '090000000000005940'
    '110000000000001440'
    '19666666666666ee3f'
    '210000000000d05740'
    '220f'
    '0d0000c441'
    '1500003442'
    '1d9a99ca42';

void main() {
  test('MeasurementData official Dart proto matches Go proto wire contract', () {
    final frame = MeasurementData(
      sessionId: 'session-1',
      rawChannels: [1, 2, 3],
      differential: DifferentialCorrection(
        sDet: 100,
        sRef: 5,
        alpha: 0.95,
        sCorrected: 95.25,
      ),
      envMeta: EnvironmentMeta(
        tempC: 24.5,
        humidityPct: 45,
        pressureKpa: 101.3,
      ),
    );

    expect(_hex(frame.writeToBuffer()), _measurementDataGoldenHex);

    final decoded =
        MeasurementData.fromBuffer(_bytes(_measurementDataGoldenHex));
    expect(decoded.sessionId, 'session-1');
    expect(decoded.rawChannels, [1, 2, 3]);
    expect(decoded.differential.sDet, 100);
    expect(decoded.differential.sRef, 5);
    expect(decoded.differential.alpha, 0.95);
    expect(decoded.differential.sCorrected, 95.25);
    expect(decoded.envMeta.tempC, 24.5);
    expect(decoded.envMeta.humidityPct, 45);
    expect(decoded.envMeta.pressureKpa, closeTo(101.3, 0.0001));
  });
}

String _hex(List<int> bytes) {
  return bytes.map((byte) => byte.toRadixString(16).padLeft(2, '0')).join();
}

List<int> _bytes(String hex) {
  final bytes = <int>[];
  for (var i = 0; i < hex.length; i += 2) {
    bytes.add(int.parse(hex.substring(i, i + 2), radix: 16));
  }
  return bytes;
}
