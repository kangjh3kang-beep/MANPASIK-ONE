import 'dart:convert';
import 'dart:io';

import 'package:crypto/crypto.dart';
import 'package:dio/dio.dart';
import 'package:dio/io.dart';
import 'package:flutter/foundation.dart' show kIsWeb, kReleaseMode;

import 'package:manpasik/core/config/app_config.dart';

/// SSL Certificate Pinning 설정
///
/// 프로덕션 빌드에서 중간자 공격(MITM)을 방지합니다.
/// AppConfig에서 환경별 인증서 핀과 허용 호스트를 가져옵니다.
/// 디버그 모드에서는 자동으로 비활성화됩니다.
class SslPinning {
  SslPinning._();

  /// Dio에 SSL Pinning을 적용합니다.
  ///
  /// Web 플랫폼 및 디버그 모드에서는 건너뜁니다.
  /// 프로덕션에서는 인증서 SHA-256 핀을 검증합니다.
  static void apply(Dio dio) {
    if (kIsWeb) return;
    if (!kReleaseMode) return;
    if (!AppConfig.sslPinningEnabled) return;

    final httpClientAdapter = dio.httpClientAdapter;
    if (httpClientAdapter is IOHttpClientAdapter) {
      httpClientAdapter.createHttpClient = () {
        final client = HttpClient();
        client.badCertificateCallback = (X509Certificate cert, String host, int port) {
          // 1. 허용된 호스트 확인
          if (!_isAllowedHost(host)) return false;

          // 2. 인증서 SHA-256 핀 검증
          final pins = AppConfig.certificatePins;
          if (pins.isEmpty) return true; // 핀 미설정 시 호스트 검증만

          return _verifyCertificatePin(cert, pins);
        };
        return client;
      };
    }
  }

  /// 인증서의 DER 인코딩 SHA-256 해시가 핀 목록에 포함되는지 검증
  static bool _verifyCertificatePin(X509Certificate cert, List<String> pins) {
    try {
      final certDer = cert.der;
      final certHash = sha256.convert(certDer);
      final certPin = 'sha256/${base64.encode(certHash.bytes)}';

      for (final pin in pins) {
        if (pin == certPin) return true;
      }
      return false;
    } catch (_) {
      return false;
    }
  }

  /// AppConfig 기반 허용 호스트 확인
  static bool _isAllowedHost(String host) {
    return AppConfig.allowedHosts.contains(host);
  }

  /// 인증서 핀 목록 조회
  static List<String> get pinnedCertificates =>
      List.unmodifiable(AppConfig.certificatePins);
}
