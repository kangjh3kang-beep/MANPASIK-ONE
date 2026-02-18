import 'package:flutter/foundation.dart';

/// 로그 레벨
enum LogLevel { debug, info, warning, error }

/// 앱 전역 로거
///
/// 레벨별 로깅을 제공하며, 릴리스 빌드에서는 warning 이상만 출력.
/// CrashReporter와 연동하여 error 레벨은 자동 리포트.
class AppLogger {
  AppLogger._();
  static final AppLogger instance = AppLogger._();

  LogLevel _minLevel = kDebugMode ? LogLevel.debug : LogLevel.warning;

  /// 최소 로그 레벨 변경
  void setMinLevel(LogLevel level) => _minLevel = level;

  void debug(String message, {String? tag}) =>
      _log(LogLevel.debug, message, tag: tag);

  void info(String message, {String? tag}) =>
      _log(LogLevel.info, message, tag: tag);

  void warning(String message, {String? tag, dynamic error}) =>
      _log(LogLevel.warning, message, tag: tag, error: error);

  void error(String message, {String? tag, dynamic error, StackTrace? stackTrace}) =>
      _log(LogLevel.error, message, tag: tag, error: error, stackTrace: stackTrace);

  void _log(
    LogLevel level,
    String message, {
    String? tag,
    dynamic error,
    StackTrace? stackTrace,
  }) {
    if (level.index < _minLevel.index) return;

    final prefix = _levelPrefix(level);
    final tagStr = tag != null ? '[$tag] ' : '';
    final line = '$prefix $tagStr$message';

    debugPrint(line);
    if (error != null) debugPrint('  Error: $error');
    if (stackTrace != null && kDebugMode) debugPrint(stackTrace.toString());
  }

  String _levelPrefix(LogLevel level) {
    switch (level) {
      case LogLevel.debug:
        return '[D]';
      case LogLevel.info:
        return '[I]';
      case LogLevel.warning:
        return '[W]';
      case LogLevel.error:
        return '[E]';
    }
  }
}
