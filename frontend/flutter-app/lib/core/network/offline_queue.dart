import 'dart:convert';

import 'package:hive_flutter/hive_flutter.dart';
import 'package:logger/logger.dart';

/// Hive 기반 오프라인 요청 큐
///
/// 네트워크 불가 시 REST 요청을 로컬에 저장하고,
/// 재연결 시 순서대로 재전송합니다.
class OfflineQueue {
  OfflineQueue._();
  static final OfflineQueue _instance = OfflineQueue._();
  static OfflineQueue get instance => _instance;

  static const _boxName = 'offline_queue';
  final _log = Logger(printer: PrettyPrinter(methodCount: 0));

  Box<String>? _box;

  /// Hive 박스 초기화 (앱 시작 시 1회 호출)
  Future<void> init() async {
    _box ??= await Hive.openBox<String>(_boxName);
    _log.d('OfflineQueue 초기화 완료 (${_box!.length}건 대기)');
  }

  /// 오프라인 요청 추가
  Future<void> enqueue(OfflineRequest request) async {
    final box = _box;
    if (box == null) return;
    final key = '${DateTime.now().microsecondsSinceEpoch}_${request.method}_${request.path}';
    await box.put(key, jsonEncode(request.toJson()));
    _log.d('오프라인 큐 추가: ${request.method} ${request.path} (총 ${box.length}건)');
  }

  /// 대기 중인 요청 수
  int get pendingCount => _box?.length ?? 0;

  /// 모든 대기 요청 가져오기 (FIFO 순서)
  List<OfflineRequest> getAll() {
    final box = _box;
    if (box == null) return [];
    return box.values
        .map((json) {
          try {
            return OfflineRequest.fromJson(jsonDecode(json) as Map<String, dynamic>);
          } catch (_) {
            return null;
          }
        })
        .whereType<OfflineRequest>()
        .toList();
  }

  /// 특정 요청 제거 (전송 성공 후)
  Future<void> remove(String key) async {
    await _box?.delete(key);
  }

  /// 전체 큐 비우기
  Future<void> clear() async {
    await _box?.clear();
    _log.d('오프라인 큐 초기화 완료');
  }

  /// 큐의 모든 키 가져오기
  List<String> get keys => _box?.keys.cast<String>().toList() ?? [];

  /// 특정 키의 요청 가져오기
  OfflineRequest? getByKey(String key) {
    final json = _box?.get(key);
    if (json == null) return null;
    try {
      return OfflineRequest.fromJson(jsonDecode(json) as Map<String, dynamic>);
    } catch (_) {
      return null;
    }
  }
}

/// 오프라인 큐에 저장되는 요청 데이터
class OfflineRequest {
  final String method;
  final String path;
  final Map<String, dynamic>? body;
  final Map<String, String>? headers;
  final DateTime createdAt;
  final int retryCount;

  const OfflineRequest({
    required this.method,
    required this.path,
    this.body,
    this.headers,
    required this.createdAt,
    this.retryCount = 0,
  });

  Map<String, dynamic> toJson() => {
        'method': method,
        'path': path,
        if (body != null) 'body': body,
        if (headers != null) 'headers': headers,
        'created_at': createdAt.toIso8601String(),
        'retry_count': retryCount,
      };

  factory OfflineRequest.fromJson(Map<String, dynamic> json) => OfflineRequest(
        method: json['method'] as String,
        path: json['path'] as String,
        body: json['body'] as Map<String, dynamic>?,
        headers: (json['headers'] as Map<String, dynamic>?)?.cast<String, String>(),
        createdAt: DateTime.parse(json['created_at'] as String),
        retryCount: (json['retry_count'] as int?) ?? 0,
      );

  OfflineRequest copyWithRetry() => OfflineRequest(
        method: method,
        path: path,
        body: body,
        headers: headers,
        createdAt: createdAt,
        retryCount: retryCount + 1,
      );
}
