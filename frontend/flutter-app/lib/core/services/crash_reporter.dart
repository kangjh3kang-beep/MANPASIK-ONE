import 'dart:async';

import 'package:flutter/foundation.dart';

/// 크래시 리포터 추상 인터페이스
///
/// 프로덕션 환경에서는 Sentry/Crashlytics 구현체로 교체.
/// 현재는 콘솔 출력 기반 ConsoleCrashReporter 제공.
abstract class CrashReporter {
  /// 비치명적 에러 리포트
  Future<void> reportError(dynamic error, StackTrace? stackTrace, {String? context});

  /// 치명적 에러(앱 크래시) 리포트
  Future<void> reportFatalError(dynamic error, StackTrace? stackTrace);

  /// 사용자 식별 정보 설정 (로그인 후 호출)
  void setUser({required String userId, String? email});

  /// 사용자 식별 정보 초기화 (로그아웃 시 호출)
  void clearUser();

  /// 커스텀 태그/컨텍스트 추가
  void setTag(String key, String value);
}

/// 콘솔 기반 크래시 리포터 (개발/테스트용)
///
/// Sentry/Crashlytics 미연동 시 기본 동작.
/// 릴리스 빌드에서도 동작하며, 추후 실 리포터로 교체 가능.
class ConsoleCrashReporter implements CrashReporter {
  final Map<String, String> _tags = {};
  String? _userId;

  @override
  Future<void> reportError(dynamic error, StackTrace? stackTrace, {String? context}) async {
    final buffer = StringBuffer('[CrashReporter] Non-fatal error');
    if (context != null) buffer.write(' ($context)');
    if (_userId != null) buffer.write(' [user=$_userId]');
    for (final entry in _tags.entries) {
      buffer.write(' [${entry.key}=${entry.value}]');
    }
    buffer.write(': $error');
    debugPrint(buffer.toString());
    if (stackTrace != null && kDebugMode) {
      debugPrint(stackTrace.toString());
    }
  }

  @override
  Future<void> reportFatalError(dynamic error, StackTrace? stackTrace) async {
    debugPrint('[CrashReporter] FATAL: $error');
    if (stackTrace != null) {
      debugPrint(stackTrace.toString());
    }
  }

  @override
  void setUser({required String userId, String? email}) {
    _userId = userId;
    debugPrint('[CrashReporter] User set: $userId');
  }

  @override
  void clearUser() {
    _userId = null;
    debugPrint('[CrashReporter] User cleared');
  }

  @override
  void setTag(String key, String value) {
    _tags[key] = value;
  }
}
