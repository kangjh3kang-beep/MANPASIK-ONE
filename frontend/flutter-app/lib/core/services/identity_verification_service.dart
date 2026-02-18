import 'package:flutter/foundation.dart';

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
}

enum VerificationType { phone, iPin }

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
