import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

import 'package:manpasik/shared/widgets/pass_verification_webview.dart';

/// 본인인증 서비스 인터페이스 (B2)
///
/// PASS / KG이니시스 등 실제 본인인증 SDK 연동 시
/// 이 인터페이스를 구현합니다. SDK 미설치 시 SimulatedIdentityService 사용.
abstract class IdentityVerificationService {
  Future<VerificationResult> requestVerification({
    required String userId,
    VerificationType type = VerificationType.phone,
  });

  Future<VerificationResult> confirmVerification({
    required String requestId,
    required String verificationCode,
  });

  /// 구성 상태에 따라 적절한 구현체를 반환하는 팩토리
  static IdentityVerificationService create() {
    if (PassIdentityService.isConfigured) {
      return PassIdentityService();
    }
    return SimulatedIdentityService();
  }
}

enum VerificationType { phone, iPin }

/// PASS 본인인증 WebView 기반 서비스
///
/// PASS 인증 페이지를 WebView로 호출하여 실제 본인인증을 진행합니다.
/// PASS_MERCHANT_ID 환경변수 미설정 시 시뮬레이션 모드 폴백.
class PassIdentityService implements IdentityVerificationService {
  static const _merchantId = String.fromEnvironment('PASS_MERCHANT_ID');

  /// PASS 키가 설정되어 있는지 확인
  static bool get isConfigured => _merchantId.isNotEmpty;

  /// NavigatorKey 설정 (앱 시작 시 호출)
  static GlobalKey<NavigatorState>? _navigatorKey;
  static void setNavigatorKey(GlobalKey<NavigatorState> key) {
    _navigatorKey = key;
  }

  // 인증 요청 보관 (requestId → userId 매핑)
  final Map<String, String> _pendingRequests = {};

  @override
  Future<VerificationResult> requestVerification({
    required String userId,
    VerificationType type = VerificationType.phone,
  }) async {
    if (!isConfigured) {
      debugPrint('[PassKYC] MERCHANT_ID 미설정 → 시뮬레이션 모드');
      return _simulateRequest(userId, type);
    }

    final context = _navigatorKey?.currentContext;
    if (context == null) {
      debugPrint('[PassKYC] NavigatorContext 없음 → 시뮬레이션 모드');
      return _simulateRequest(userId, type);
    }

    final result = await Navigator.push<PassVerificationResult>(
      context,
      MaterialPageRoute(
        builder: (_) => PassVerificationWebView(
          merchantId: _merchantId,
        ),
      ),
    );

    if (result == null) {
      return VerificationResult(
        success: false,
        requestId: '',
        message: '사용자가 인증을 취소했습니다',
        status: VerificationStatus.failed,
      );
    }

    if (result.success) {
      final requestId = 'pass_${DateTime.now().millisecondsSinceEpoch}';
      return VerificationResult(
        success: true,
        requestId: requestId,
        message: '본인인증이 완료되었습니다',
        status: VerificationStatus.verified,
        verifiedName: result.name,
        verifiedPhone: result.phone,
      );
    }

    return VerificationResult(
      success: false,
      requestId: '',
      message: result.errorMessage ?? '본인인증 실패',
      status: VerificationStatus.failed,
    );
  }

  @override
  Future<VerificationResult> confirmVerification({
    required String requestId,
    required String verificationCode,
  }) async {
    // PASS 방식은 WebView에서 인증이 완료되므로 별도 확인 불필요
    // requestVerification에서 이미 verified 상태로 반환
    if (_pendingRequests.containsKey(requestId)) {
      _pendingRequests.remove(requestId);
      return VerificationResult(
        success: true,
        requestId: requestId,
        message: '본인인증 확인 완료',
        status: VerificationStatus.verified,
      );
    }
    return VerificationResult(
      success: false,
      requestId: requestId,
      message: '유효하지 않은 인증 요청입니다',
      status: VerificationStatus.failed,
    );
  }

  Future<VerificationResult> _simulateRequest(
    String userId,
    VerificationType type,
  ) async {
    debugPrint('[PassKYC:Sim] 본인인증 요청: $userId ($type)');
    await Future.delayed(const Duration(seconds: 1));
    final requestId = 'kyc_${DateTime.now().millisecondsSinceEpoch}';
    _pendingRequests[requestId] = userId;
    return VerificationResult(
      success: true,
      requestId: requestId,
      message: '시뮬레이션 인증 코드: 123456',
      status: VerificationStatus.pending,
    );
  }
}

/// 시뮬레이션 본인인증 서비스
class SimulatedIdentityService implements IdentityVerificationService {
  @override
  Future<VerificationResult> requestVerification({
    required String userId,
    VerificationType type = VerificationType.phone,
  }) async {
    debugPrint('[SimulatedKYC] 본인인증 요청: $userId ($type)');
    await Future.delayed(const Duration(seconds: 1));
    return VerificationResult(
      success: true,
      requestId: 'kyc_${DateTime.now().millisecondsSinceEpoch}',
      message: '시뮬레이션 인증 코드: 123456',
      status: VerificationStatus.pending,
    );
  }

  @override
  Future<VerificationResult> confirmVerification({
    required String requestId,
    required String verificationCode,
  }) async {
    debugPrint('[SimulatedKYC] 인증 확인: $requestId (code=$verificationCode)');
    await Future.delayed(const Duration(milliseconds: 500));

    if (verificationCode == '123456') {
      return VerificationResult(
        success: true,
        requestId: requestId,
        message: '시뮬레이션 본인인증 완료',
        status: VerificationStatus.verified,
        verifiedName: '홍길동',
        verifiedPhone: '010-****-1234',
      );
    }

    return VerificationResult(
      success: false,
      requestId: requestId,
      message: '인증 코드가 일치하지 않습니다',
      status: VerificationStatus.failed,
    );
  }
}

enum VerificationStatus { pending, verified, failed, expired }

class VerificationResult {
  final bool success;
  final String requestId;
  final String message;
  final VerificationStatus status;
  final String? verifiedName;
  final String? verifiedPhone;
  final String? errorCode;

  const VerificationResult({
    required this.success,
    required this.requestId,
    required this.message,
    required this.status,
    this.verifiedName,
    this.verifiedPhone,
    this.errorCode,
  });
}
