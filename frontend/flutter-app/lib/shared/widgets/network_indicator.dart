import 'dart:async';

import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/sync_provider.dart';

/// Network status indicator banner.
///
/// Shows offline/syncing/online/conflict state near the top of the app.
/// When conflicts exist, tapping the banner navigates to /conflict-resolve.
class NetworkIndicator extends ConsumerStatefulWidget {
  const NetworkIndicator({super.key});

  @override
  ConsumerState<NetworkIndicator> createState() => _NetworkIndicatorState();
}

class _NetworkIndicatorState extends ConsumerState<NetworkIndicator>
    with SingleTickerProviderStateMixin {
  _NetworkStatus _status = _NetworkStatus.online;
  late final AnimationController _animController;
  Connectivity? _connectivity;
  StreamSubscription<ConnectivityResult>? _connectivitySub;

  @override
  void initState() {
    super.initState();
    _animController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 300),
    );

    // Linux desktop(WSL 포함)에서는 NetworkManager DBus가 비활성인 경우가 많아
    // connectivity_plus가 비동기 예외를 내므로 연결 감시는 비활성화한다.
    if (!kIsWeb && defaultTargetPlatform == TargetPlatform.linux) {
      debugPrint(
          '[NetworkIndicator] Linux desktop: connectivity probe disabled');
      return;
    }

    _connectivity = Connectivity();
    _connectivitySub = _connectivity?.onConnectivityChanged.listen(
      _onConnectivityChanged,
      onError: (Object error, StackTrace stackTrace) {
        debugPrint('[NetworkIndicator] Connectivity stream unavailable: ' +
            error.toString());
        _updateStatus(_NetworkStatus.online);
      },
    );
    unawaited(_checkConnectivity());
  }

  @override
  void dispose() {
    _connectivitySub?.cancel();
    _animController.dispose();
    super.dispose();
  }

  void _onConnectivityChanged(ConnectivityResult result) {
    _updateStatus(result != ConnectivityResult.none
        ? _NetworkStatus.online
        : _NetworkStatus.offline);
  }

  Future<void> _checkConnectivity() async {
    final connectivity = _connectivity;
    if (connectivity == null) return;

    try {
      final result = await connectivity.checkConnectivity();
      _updateStatus(result != ConnectivityResult.none
          ? _NetworkStatus.online
          : _NetworkStatus.offline);
    } catch (error, stackTrace) {
      debugPrint(
          '[NetworkIndicator] Connectivity probe failed: ' + error.toString());
      _updateStatus(_NetworkStatus.online);
      debugPrintStack(stackTrace: stackTrace);
    }
  }

  void _updateStatus(_NetworkStatus newStatus) {
    if (_status == newStatus) return;
    setState(() => _status = newStatus);
    if (newStatus == _NetworkStatus.online) {
      _animController.reverse();
    } else {
      _animController.forward();
    }
  }

  void setOffline() => _updateStatus(_NetworkStatus.offline);

  void setSyncing() => _updateStatus(_NetworkStatus.syncing);

  @override
  Widget build(BuildContext context) {
    final syncState = ref.watch(syncProvider);
    final hasConflicts = syncState.hasConflicts;

    if (_status == _NetworkStatus.online && !hasConflicts) {
      return const SizedBox.shrink();
    }

    final conflictText =
        '데이터 충돌 ' + syncState.failedCount.toString() + '건 — 탭하여 해결';

    final (color, icon, text) = hasConflicts
        ? (
            Colors.deepOrange,
            Icons.sync_problem,
            conflictText,
          )
        : switch (_status) {
            _NetworkStatus.offline => (
                Colors.red,
                Icons.cloud_off,
                '오프라인 모드 — 데이터가 로컬에 저장됩니다',
              ),
            _NetworkStatus.syncing => (
                Colors.orange,
                Icons.sync,
                '데이터 동기화 중...'
              ),
            _NetworkStatus.online => (Colors.green, Icons.cloud_done, '연결됨'),
          };

    return SizeTransition(
      sizeFactor:
          hasConflicts ? const AlwaysStoppedAnimation(1.0) : _animController,
      child: GestureDetector(
        onTap: hasConflicts ? () => context.push('/conflict-resolve') : null,
        child: Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
          color: color.withOpacity(0.9),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(icon, size: 14, color: Colors.white),
              const SizedBox(width: 8),
              Flexible(
                child: Text(
                  text,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                    fontWeight: FontWeight.w500,
                  ),
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              if (hasConflicts) ...[
                const SizedBox(width: 4),
                const Icon(Icons.chevron_right, size: 14, color: Colors.white),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

enum _NetworkStatus { online, offline, syncing }
