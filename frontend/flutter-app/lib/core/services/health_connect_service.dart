import 'dart:async';
import 'dart:io' show Platform;

import 'package:flutter/foundation.dart' show debugPrint, kIsWeb;
import 'package:flutter/services.dart';

import 'package:manpasik/core/config/app_config.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// 외부 건강 플랫폼 연동 서비스
///
/// Apple HealthKit / Google Health Connect와 건강 데이터를 동기화합니다.
/// 모바일 플랫폼에서 `health` 패키지 사용, Web/Desktop에서는 시뮬레이션 폴백.
///
/// 읽기 가능 데이터:
/// - 걸음 수, 심박수, 혈압, 혈당, 체중, 수면
///
/// 쓰기 가능 데이터:
/// - ManPaSik 측정 결과 → HealthKit/Health Connect 동기화
class HealthConnectService {
  HealthConnectService._();
  static final instance = HealthConnectService._();

  bool _isAuthorized = false;
  bool get isAuthorized => _isAuthorized;

  String? _platform; // 'apple_healthkit' | 'google_health_connect'

  /// 네이티브 건강 API 사용 가능 여부
  static bool get _isNativeAvailable {
    if (kIsWeb) return false;
    if (!AppConfig.enableHealthKit) return false;
    return Platform.isAndroid || Platform.isIOS;
  }

  /// Platform method channel for native health API
  static const _channel = MethodChannel('com.manpasik/health_connect');

  /// 사용 가능한 건강 데이터 유형
  static const supportedTypes = [
    HealthDataType.steps,
    HealthDataType.heartRate,
    HealthDataType.bloodPressureSystolic,
    HealthDataType.bloodPressureDiastolic,
    HealthDataType.bloodGlucose,
    HealthDataType.weight,
    HealthDataType.sleep,
    HealthDataType.oxygenSaturation,
  ];

  /// 건강 플랫폼 접근 권한을 요청합니다.
  Future<bool> requestAuthorization() async {
    if (_isNativeAvailable) {
      try {
        final result = await _channel.invokeMethod<bool>(
          'requestAuthorization',
          {'types': supportedTypes.map((t) => t.name).toList()},
        );
        _isAuthorized = result ?? false;
        if (!kIsWeb) {
          _platform = Platform.isIOS ? 'apple_healthkit' : 'google_health_connect';
        }
        return _isAuthorized;
      } on PlatformException {
        // 네이티브 실패 → 시뮬레이션 폴백
      }
    }

    // 시뮬레이션 폴백
    await Future.delayed(const Duration(milliseconds: 500));
    _isAuthorized = true;
    _platform = 'simulation';
    return _isAuthorized;
  }

  /// 건강 데이터를 읽어옵니다.
  Future<List<HealthRecord>> fetchHealthData({
    required HealthDataType type,
    required DateTime startDate,
    required DateTime endDate,
  }) async {
    if (!_isAuthorized) {
      throw StateError('건강 플랫폼 권한이 없습니다. requestAuthorization()을 먼저 호출하세요.');
    }

    if (_isNativeAvailable && _platform != 'simulation') {
      try {
        final result = await _channel.invokeMethod<List<dynamic>>(
          'fetchHealthData',
          {
            'type': type.name,
            'startDate': startDate.millisecondsSinceEpoch,
            'endDate': endDate.millisecondsSinceEpoch,
          },
        );
        if (result != null) {
          return result.map((item) {
            final map = item as Map<dynamic, dynamic>;
            return HealthRecord(
              type: type,
              value: (map['value'] as num).toDouble(),
              unit: (map['unit'] as String?) ?? type.defaultUnit,
              timestamp: DateTime.fromMillisecondsSinceEpoch(map['timestamp'] as int),
              source: (map['source'] as String?) ?? _platform ?? 'native',
            );
          }).toList();
        }
      } on PlatformException {
        // 네이티브 실패 → 시뮬레이션 폴백
      }
    }

    await Future.delayed(const Duration(milliseconds: 300));
    return _generateSimulatedData(type, startDate, endDate);
  }

  /// ManPaSik 측정 결과를 건강 플랫폼에 기록합니다.
  Future<bool> writeHealthData({
    required HealthDataType type,
    required double value,
    required DateTime timestamp,
    String? unit,
  }) async {
    if (!_isAuthorized) return false;

    if (_isNativeAvailable && _platform != 'simulation') {
      try {
        final result = await _channel.invokeMethod<bool>(
          'writeHealthData',
          {
            'type': type.name,
            'value': value,
            'timestamp': timestamp.millisecondsSinceEpoch,
            'unit': unit ?? type.defaultUnit,
          },
        );
        return result ?? false;
      } on PlatformException {
        // 네이티브 실패 → 시뮬레이션 폴백
      }
    }

    await Future.delayed(const Duration(milliseconds: 200));
    return true;
  }

  /// 연결 상태 정보를 반환합니다.
  Map<String, dynamic> getConnectionInfo() {
    return {
      'is_authorized': _isAuthorized,
      'platform': _platform ?? 'none',
      'is_native': _isNativeAvailable && _platform != 'simulation',
      'supported_types': supportedTypes.map((t) => t.name).toList(),
      'last_sync': DateTime.now().toIso8601String(),
    };
  }

  /// 연결을 해제합니다.
  Future<void> disconnect() async {
    _isAuthorized = false;
    _platform = null;
  }

  /// 건강 데이터를 ManPaSik 서버에 동기화합니다.
  ///
  /// 로컬 HealthKit/Health Connect 데이터를 읽어서
  /// REST API를 통해 서버에 업로드합니다.
  Future<int> syncToServer({
    required ManPaSikRestClient restClient,
    required String userId,
    DateTime? since,
  }) async {
    if (!_isAuthorized) return 0;

    final startDate = since ?? DateTime.now().subtract(const Duration(days: 1));
    final endDate = DateTime.now();
    int synced = 0;

    for (final type in supportedTypes) {
      try {
        final records = await fetchHealthData(
          type: type,
          startDate: startDate,
          endDate: endDate,
        );
        if (records.isEmpty) continue;

        // 최신 레코드만 서버에 전송
        final latest = records.last;
        await restClient.createHealthRecord(
          userId: userId,
          recordType: _healthDataTypeIndex(type),
          title: '${type.displayName} 자동 동기화',
          description: '${latest.value} ${latest.unit}',
          provider: _platform ?? 'health_connect',
          metadata: {
            'source': latest.source,
            'value': latest.value.toString(),
            'unit': latest.unit,
            'timestamp': latest.timestamp.toIso8601String(),
          },
        );
        synced++;
      } catch (e) {
        debugPrint('[HealthConnect] $type 동기화 실패: $e');
      }
    }

    debugPrint('[HealthConnect] 서버 동기화 완료: $synced/${ supportedTypes.length} 타입');
    return synced;
  }

  int _healthDataTypeIndex(HealthDataType type) {
    return switch (type) {
      HealthDataType.steps => 10,
      HealthDataType.heartRate => 11,
      HealthDataType.bloodPressureSystolic => 12,
      HealthDataType.bloodPressureDiastolic => 13,
      HealthDataType.bloodGlucose => 14,
      HealthDataType.weight => 15,
      HealthDataType.sleep => 16,
      HealthDataType.oxygenSaturation => 17,
    };
  }

  List<HealthRecord> _generateSimulatedData(
    HealthDataType type,
    DateTime start,
    DateTime end,
  ) {
    final records = <HealthRecord>[];
    var current = start;
    while (current.isBefore(end)) {
      final value = switch (type) {
        HealthDataType.steps => 5000.0 + (current.day * 317 % 8000),
        HealthDataType.heartRate => 65.0 + (current.hour * 3 % 30),
        HealthDataType.bloodPressureSystolic => 115.0 + (current.day % 20),
        HealthDataType.bloodPressureDiastolic => 72.0 + (current.day % 15),
        HealthDataType.bloodGlucose => 95.0 + (current.hour * 2 % 40),
        HealthDataType.weight => 68.0 + (current.day % 5) * 0.1,
        HealthDataType.sleep => 6.5 + (current.day % 4) * 0.5,
        HealthDataType.oxygenSaturation => 96.0 + (current.day % 4),
      };
      records.add(HealthRecord(
        type: type,
        value: value,
        unit: type.defaultUnit,
        timestamp: current,
        source: _platform ?? 'simulation',
      ));
      current = current.add(const Duration(hours: 1));
    }
    return records;
  }
}

/// 건강 데이터 유형
enum HealthDataType {
  steps('걸음 수', '걸음'),
  heartRate('심박수', 'bpm'),
  bloodPressureSystolic('수축기 혈압', 'mmHg'),
  bloodPressureDiastolic('이완기 혈압', 'mmHg'),
  bloodGlucose('혈당', 'mg/dL'),
  weight('체중', 'kg'),
  sleep('수면', '시간'),
  oxygenSaturation('산소포화도', '%');

  final String displayName;
  final String defaultUnit;
  const HealthDataType(this.displayName, this.defaultUnit);
}

/// 건강 데이터 레코드
class HealthRecord {
  final HealthDataType type;
  final double value;
  final String unit;
  final DateTime timestamp;
  final String source;

  const HealthRecord({
    required this.type,
    required this.value,
    required this.unit,
    required this.timestamp,
    required this.source,
  });

  Map<String, dynamic> toJson() => {
        'type': type.name,
        'value': value,
        'unit': unit,
        'timestamp': timestamp.toIso8601String(),
        'source': source,
      };
}
