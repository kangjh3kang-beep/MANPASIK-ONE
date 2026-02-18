import 'package:flutter/widgets.dart';

import 'package:manpasik/core/services/app_logger.dart';

/// 앱 라이프사이클 관리자
///
/// 백그라운드/포그라운드 전환 시 리소스 관리:
/// - 백그라운드: 불필요한 연결 해제, 배터리 절약
/// - 포그라운드: 토큰 갱신, 데이터 동기화
class AppLifecycleObserver with WidgetsBindingObserver {
  AppLifecycleObserver({
    this.onResumed,
    this.onPaused,
    this.onDetached,
  });

  final VoidCallback? onResumed;
  final VoidCallback? onPaused;
  final VoidCallback? onDetached;

  final _log = AppLogger.instance;

  /// WidgetsBinding에 옵저버 등록
  void register() {
    WidgetsBinding.instance.addObserver(this);
    _log.info('AppLifecycleObserver registered', tag: 'Lifecycle');
  }

  /// WidgetsBinding에서 옵저버 해제
  void unregister() {
    WidgetsBinding.instance.removeObserver(this);
    _log.info('AppLifecycleObserver unregistered', tag: 'Lifecycle');
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    switch (state) {
      case AppLifecycleState.resumed:
        _log.info('App resumed (foreground)', tag: 'Lifecycle');
        onResumed?.call();
      case AppLifecycleState.paused:
        _log.info('App paused (background)', tag: 'Lifecycle');
        onPaused?.call();
      case AppLifecycleState.detached:
        _log.info('App detached', tag: 'Lifecycle');
        onDetached?.call();
      case AppLifecycleState.inactive:
        _log.debug('App inactive', tag: 'Lifecycle');
      case AppLifecycleState.hidden:
        _log.debug('App hidden', tag: 'Lifecycle');
    }
  }
}
