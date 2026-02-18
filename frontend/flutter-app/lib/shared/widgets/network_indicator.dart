import 'dart:async';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/sync_provider.dart';

/// 네트워크 상태 인디케이터
///
/// 앱 상단에 배치하여 오프라인/동기화/온라인/충돌 상태를 표시합니다.
/// ScaffoldWithBottomNav 등의 상위 위젯에 삽입하여 사용합니다.
/// 충돌 존재 시 탭하여 `/conflict-resolve` 화면으로 이동합니다.
class NetworkIndicator extends ConsumerStatefulWidget {
  const NetworkIndicator({super.key});

  @override
  ConsumerState<NetworkIndicator> createState() => _NetworkIndicatorState();
}

class _NetworkIndicatorState extends ConsumerState<NetworkIndicator> with SingleTickerProviderStateMixin {
  _NetworkStatus _status = _NetworkStatus.online;
  late AnimationController _animController;
  StreamSubscription<ConnectivityResult>? _connectivitySub;

  @override
  void initState() {
    super.initState();
    _animController = AnimationController(vsync: this, duration: const Duration(milliseconds: 300));
    _connectivitySub = Connectivity().onConnectivityChanged.listen(_onConnectivityChanged);
    _checkConnectivity();
  }

  @override
  void dispose() {
    _connectivitySub?.cancel();
    _animController.dispose();
    super.dispose();
  }

  void _onConnectivityChanged(ConnectivityResult result) {
    _updateStatus(result != ConnectivityResult.none ? _NetworkStatus.online : _NetworkStatus.offline);
  }

  Future<void> _checkConnectivity() async {
    final result = await Connectivity().checkConnectivity();
    _updateStatus(result != ConnectivityResult.none ? _NetworkStatus.online : _NetworkStatus.offline);
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

  /// 외부에서 오프라인 상태를 설정할 때 사용
  void setOffline() => _updateStatus(_NetworkStatus.offline);

  /// 외부에서 동기화 상태를 설정할 때 사용
  void setSyncing() => _updateStatus(_NetworkStatus.syncing);

  @override
  Widget build(BuildContext context) {
    final syncState = ref.watch(syncProvider);
    final hasConflicts = syncState.hasConflicts;

    // 충돌이 있으면 온라인이어도 배너 표시
    if (_status == _NetworkStatus.online && !hasConflicts) {
      return const SizedBox.shrink();
    }

    final (color, icon, text) = hasConflicts
        ? (Colors.deepOrange, Icons.sync_problem, '데이터 충돌 ${syncState.failedCount}건 — 탭하여 해결')
        : switch (_status) {
            _NetworkStatus.offline => (Colors.red, Icons.cloud_off, '오프라인 모드 — 데이터가 로컬에 저장됩니다'),
            _NetworkStatus.syncing => (Colors.orange, Icons.sync, '데이터 동기화 중...'),
            _NetworkStatus.online => (Colors.green, Icons.cloud_done, '연결됨'),
          };

    return SizeTransition(
      sizeFactor: hasConflicts ? const AlwaysStoppedAnimation(1.0) : _animController,
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
                  style: const TextStyle(color: Colors.white, fontSize: 12, fontWeight: FontWeight.w500),
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
