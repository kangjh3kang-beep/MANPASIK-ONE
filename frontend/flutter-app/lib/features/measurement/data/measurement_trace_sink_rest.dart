import 'dart:async';

import 'package:manpasik/core/services/app_logger.dart';
import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/features/measurement/application/measurement_golden_path_orchestrator.dart';

class MeasurementGoldenPathRestTraceSink {
  MeasurementGoldenPathRestTraceSink(
    this._client, {
    this.source = 'flutter',
    this.route = '/measure',
  });

  final ManPaSikRestClient _client;
  final String source;
  final String route;

  void call(MeasurementGoldenPathTraceEvent event) {
    unawaited(_send(event));
  }

  Future<void> _send(MeasurementGoldenPathTraceEvent event) async {
    try {
      await _client.recordMeasurementTraceEvent(
        data: event.toRemoteObservabilityJson(
          source: source,
          route: route,
        ),
      );
    } catch (error) {
      AppLogger.instance.warning(
        'remote_trace_failed phase=${event.phase.name}',
        tag: 'MeasureGoldenPath',
        error: error,
      );
    }
  }
}
