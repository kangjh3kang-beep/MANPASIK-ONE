import 'dart:async';
import 'dart:math';

import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:logger/logger.dart';

import 'package:manpasik/core/config/app_config.dart';
import 'package:manpasik/core/network/offline_queue.dart';

/// 동기화 상태
enum SyncStatus { idle, syncing, error, offline }

/// 동기화 상태 데이터
class SyncState {
  final SyncStatus status;
  final int pendingCount;
  final int syncedCount;
  final int failedCount;
  final String? lastError;
  final DateTime? lastSyncAt;
  final bool hasConflicts;

  const SyncState({
    this.status = SyncStatus.idle,
    this.pendingCount = 0,
    this.syncedCount = 0,
    this.failedCount = 0,
    this.lastError,
    this.lastSyncAt,
    this.hasConflicts = false,
  });

  SyncState copyWith({
    SyncStatus? status,
    int? pendingCount,
    int? syncedCount,
    int? failedCount,
    String? lastError,
    DateTime? lastSyncAt,
    bool? hasConflicts,
  }) =>
      SyncState(
        status: status ?? this.status,
        pendingCount: pendingCount ?? this.pendingCount,
        syncedCount: syncedCount ?? this.syncedCount,
        failedCount: failedCount ?? this.failedCount,
        lastError: lastError ?? this.lastError,
        lastSyncAt: lastSyncAt ?? this.lastSyncAt,
        hasConflicts: hasConflicts ?? this.hasConflicts,
      );
}

/// 자동 동기화 프로바이더
///
/// connectivity_plus로 네트워크 상태 감시 →
/// 재연결 시 OfflineQueue 일괄 전송
final syncProvider =
    StateNotifierProvider<SyncNotifier, SyncState>((ref) => SyncNotifier());

class SyncNotifier extends StateNotifier<SyncState> {
  SyncNotifier() : super(const SyncState()) {
    _init();
  }

  final _log = Logger(printer: PrettyPrinter(methodCount: 0));
  StreamSubscription<ConnectivityResult>? _connectivitySub;
  Timer? _periodicSync;
  final _dio = Dio(BaseOptions(
    baseUrl: AppConfig.baseUrl,
    connectTimeout: const Duration(seconds: 10),
    receiveTimeout: const Duration(seconds: 30),
  ));

  static const _maxRetries = 3;
  int _consecutiveFailures = 0;

  void _init() {
    _connectivitySub = Connectivity().onConnectivityChanged.listen(_onConnectivityChanged);
    _schedulePeriodicSync();
    _updatePendingCount();
  }

  /// Exponential backoff 기반 주기적 동기화 스케줄링
  void _schedulePeriodicSync() {
    _periodicSync?.cancel();
    // 기본 5분, 실패 시 최대 30분까지 증가
    final backoffMinutes = min(5 * pow(2, _consecutiveFailures), 30).toInt();
    _log.d('다음 동기화: ${backoffMinutes}분 후');
    _periodicSync = Timer.periodic(Duration(minutes: backoffMinutes), (_) => syncAll());
  }

  void _onConnectivityChanged(ConnectivityResult result) {
    if (result == ConnectivityResult.none) {
      state = state.copyWith(status: SyncStatus.offline);
      _log.w('네트워크 끊김 — 오프라인 모드');
    } else {
      _log.i('네트워크 복구 — 동기화 시작');
      syncAll();
    }
  }

  void _updatePendingCount() {
    state = state.copyWith(pendingCount: OfflineQueue.instance.pendingCount);
  }

  /// 전체 오프라인 큐 동기화
  Future<void> syncAll() async {
    final queue = OfflineQueue.instance;
    if (queue.pendingCount == 0) {
      state = state.copyWith(status: SyncStatus.idle, pendingCount: 0);
      return;
    }

    state = state.copyWith(status: SyncStatus.syncing);
    var synced = 0;
    var failed = 0;

    final keys = queue.keys;
    for (final key in keys) {
      final request = queue.getByKey(key);
      if (request == null) {
        await queue.remove(key);
        continue;
      }

      try {
        await _sendRequest(request);
        await queue.remove(key);
        synced++;
        state = state.copyWith(
          syncedCount: state.syncedCount + 1,
          pendingCount: queue.pendingCount,
        );
      } catch (e) {
        failed++;
        if (request.retryCount >= _maxRetries) {
          await queue.remove(key);
          _log.e('최대 재시도 초과, 삭제: ${request.method} ${request.path}');
        }
        _log.w('동기화 실패 (${request.retryCount + 1}/$_maxRetries): ${request.path}');
      }
    }

    state = state.copyWith(
      status: failed > 0 ? SyncStatus.error : SyncStatus.idle,
      syncedCount: state.syncedCount,
      failedCount: state.failedCount + failed,
      lastSyncAt: DateTime.now(),
      pendingCount: queue.pendingCount,
    );

    _log.i('동기화 완료: 성공=$synced, 실패=$failed, 대기=${queue.pendingCount}');

    // Exponential backoff 갱신
    if (failed > 0) {
      _consecutiveFailures = min(_consecutiveFailures + 1, 4);
      _schedulePeriodicSync();
    } else if (_consecutiveFailures > 0) {
      _consecutiveFailures = 0;
      _schedulePeriodicSync();
    }

    // 충돌 감지 시 hasConflicts 플래그 설정
    if (failed > 0 && queue.pendingCount > 0) {
      state = state.copyWith(hasConflicts: true);
    }
  }

  /// 충돌 존재 여부 확인
  bool get hasConflicts => state.hasConflicts;

  /// 충돌 해결 완료 처리
  void clearConflicts() {
    state = state.copyWith(hasConflicts: false);
  }

  Future<void> _sendRequest(OfflineRequest request) async {
    final options = Options(headers: request.headers);
    switch (request.method.toUpperCase()) {
      case 'POST':
        await _dio.post(request.path, data: request.body, options: options);
      case 'PUT':
        await _dio.put(request.path, data: request.body, options: options);
      case 'PATCH':
        await _dio.patch(request.path, data: request.body, options: options);
      case 'DELETE':
        await _dio.delete(request.path, data: request.body, options: options);
      default:
        await _dio.get(request.path, options: options);
    }
  }

  /// Dio 인스턴스에 인증 토큰 설정
  void setAuthToken(String token) {
    _dio.options.headers['Authorization'] = 'Bearer $token';
  }

  @override
  void dispose() {
    _connectivitySub?.cancel();
    _periodicSync?.cancel();
    super.dispose();
  }
}
