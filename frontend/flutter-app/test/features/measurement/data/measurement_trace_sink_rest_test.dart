import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/measurement/application/measurement_golden_path_orchestrator.dart';
import 'package:manpasik/features/measurement/data/measurement_trace_sink_rest.dart';

void main() {
  test('MeasurementGoldenPathRestTraceSink forwards PHI-minimized payload',
      () async {
    final received = Completer<Map<String, dynamic>>();
    final server = await HttpServer.bind(InternetAddress.loopbackIPv4, 0);
    addTearDown(() async {
      await server.close(force: true);
    });

    unawaited(() async {
      await for (final request in server) {
        final body = await utf8.decoder.bind(request).join();
        if (!received.isCompleted) {
          received.complete(jsonDecode(body) as Map<String, dynamic>);
        }
        expect(request.method, 'POST');
        expect(request.uri.path, '/api/v1/measurements/trace-events');
        request.response
          ..statusCode = HttpStatus.accepted
          ..headers.contentType = ContentType.json
          ..write(jsonEncode({'accepted': true}));
        await request.response.close();
      }
    }());

    final client = ManPaSikRestClient(
      baseUrl: 'http://${server.address.host}:${server.port}/api/v1',
    );
    final sink = MeasurementGoldenPathRestTraceSink(
      client,
      source: 'flutter-test',
      route: '/measure-test',
    );

    sink(MeasurementGoldenPathTraceEvent(
      phase: MeasurementGoldenPathPhase.serverProcessed,
      elapsedMs: 42,
      occurredAt: DateTime.utc(2026, 5, 1, 12),
      sessionId: 'session-1',
      cartridgeId: 'cart-1',
      engineMode: 'native',
      primaryValue: 95,
      unit: 'mg/dL',
      confidence: 0.96,
    ));

    final payload = await received.future.timeout(const Duration(seconds: 2));
    expect(payload['schema_version'], 'measure_trace.v1');
    expect(payload['source'], 'flutter-test');
    expect(payload['route'], '/measure-test');
    expect(payload['phase'], 'serverProcessed');
    expect(payload['session_id'], 'session-1');
    expect(payload['has_primary_value'], isTrue);
    expect(payload.containsKey('primary_value'), isFalse);
  });
}
