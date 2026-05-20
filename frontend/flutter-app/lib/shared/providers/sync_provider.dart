import 'dart:async';
import 'dart:math';

import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:logger/logger.dart';

import 'package:manpasik/core/config/app_config.dart';
import 'package:manpasik/core/network/offline_queue.dart';

/// Sync status.
enum SyncStatus { idle, syncing, error, offline }

/// Sync state data.
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

/// Auto sync provider.
///
/// Watches connectivity and flushes OfflineQueue on reconnection.
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
    if (!kIsWeb && defaultTargetPlatform == TargetPlatform.linux) {
      _log.w(
          'Linux desktop: connectivity listener disabled (NetworkManager unavailable)');
      _schedulePeriodicSync();
      _updatePendingCount();
      return;
    }

    _connectivitySub = Connectivity().onConnectivityChanged.listen(
      _onConnectivityChanged,
      onError: (Object error, StackTrace stackTrace) {
        _log.w('Connectivity listener unavailable: $error');
      },
    );
    _schedulePeriodicSync();
    _updatePendingCount();
  }

  /// Schedules periodic sync with exponential backoff.
  void _schedulePeriodicSync() {
    _periodicSync?.cancel();
    final backoffMinutes = min(5 * pow(2, _consecutiveFailures), 30).toInt();
    _log.d('Next sync in ${backoffMinutes}m');
    _periodicSync =
        Timer.periodic(Duration(minutes: backoffMinutes), (_) => syncAll());
  }

  void _onConnectivityChanged(ConnectivityResult result) {
    if (result == ConnectivityResult.none) {
      state = state.copyWith(status: SyncStatus.offline);
      _log.w('Network offline - entering offline mode');
    } else {
      _log.i('Network restored - starting sync');
      syncAll();
    }
  }

  void _updatePendingCount() {
    state = state.copyWith(pendingCount: OfflineQueue.instance.pendingCount);
  }

  /// Syncs all pending offline queue requests.
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
          _log.e(
              'Max retries exceeded, dropping request: ${request.method} ${request.path}');
        }
        _log.w(
            'Sync failed (${request.retryCount + 1}/$_maxRetries): ${request.path}');
      }
    }

    state = state.copyWith(
      status: failed > 0 ? SyncStatus.error : SyncStatus.idle,
      syncedCount: state.syncedCount,
      failedCount: state.failedCount + failed,
      lastSyncAt: DateTime.now(),
      pendingCount: queue.pendingCount,
    );

    _log.i(
        'Sync done: ok=$synced, fail=$failed, pending=${queue.pendingCount}');

    if (failed > 0) {
      _consecutiveFailures = min(_consecutiveFailures + 1, 4);
      _schedulePeriodicSync();
    } else if (_consecutiveFailures > 0) {
      _consecutiveFailures = 0;
      _schedulePeriodicSync();
    }

    if (failed > 0 && queue.pendingCount > 0) {
      state = state.copyWith(hasConflicts: true);
    }
  }

  bool get hasConflicts => state.hasConflicts;

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
