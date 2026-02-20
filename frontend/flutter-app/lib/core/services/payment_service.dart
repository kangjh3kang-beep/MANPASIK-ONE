import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

import 'package:manpasik/core/services/rest_client.dart';
import 'package:manpasik/shared/widgets/toss_payment_webview.dart';

/// PG 결제 서비스 인터페이스 (B1)
///
/// Toss Payments / NHN KCP 등 실제 PG SDK 연동 시
/// 이 인터페이스를 구현합니다. SDK 미설치 시 SimulatedPaymentService 사용.
abstract class PaymentService {
  Future<PaymentResult> requestPayment({
    required String orderId,
    required int amountKrw,
    required String orderName,
    String? customerName,
    String? customerEmail,
  });

  Future<PaymentResult> confirmPayment({
    required String paymentKey,
    required String orderId,
    required int amountKrw,
  });

  Future<RefundResult> requestRefund({
    required String paymentKey,
    String? reason,
    int? refundAmountKrw,
  });

  /// 구성 상태에 따라 적절한 구현체를 반환하는 팩토리
  static PaymentService create({
    ManPaSikRestClient? restClient,
    String? userId,
  }) {
    if (TossPaymentService.isConfigured && restClient != null) {
      return TossPaymentService(restClient: restClient, userId: userId);
    }
    return SimulatedPaymentService();
  }
}

/// Toss Payments WebView 기반 결제 서비스
///
/// 실제 Toss Payments 결제창(WebView)을 열어 카드 결제를 진행하고
/// 결제 성공 시 REST API로 서버 확인(confirm)을 수행합니다.
/// TOSS_CLIENT_KEY 환경변수 미설정 시 시뮬레이션 모드 폴백.
class TossPaymentService implements PaymentService {
  TossPaymentService({required this.restClient, this.userId});

  final ManPaSikRestClient restClient;
  final String? userId;

  static const _clientKey = String.fromEnvironment('TOSS_CLIENT_KEY');

  /// Toss 키가 설정되어 있는지 확인
  static bool get isConfigured => _clientKey.isNotEmpty;

  /// NavigatorKey 설정 (앱 시작 시 호출)
  static GlobalKey<NavigatorState>? _navigatorKey;
  static void setNavigatorKey(GlobalKey<NavigatorState> key) {
    _navigatorKey = key;
  }

  @override
  Future<PaymentResult> requestPayment({
    required String orderId,
    required int amountKrw,
    required String orderName,
    String? customerName,
    String? customerEmail,
  }) async {
    if (!isConfigured) {
      debugPrint('[TossPG] CLIENT_KEY 미설정 → 시뮬레이션 모드');
      return _simulatePayment(orderId, amountKrw, orderName);
    }

    final context = _navigatorKey?.currentContext;
    if (context == null) {
      debugPrint('[TossPG] NavigatorContext 없음 → 시뮬레이션 모드');
      return _simulatePayment(orderId, amountKrw, orderName);
    }

    final result = await Navigator.push<TossPaymentResult>(
      context,
      MaterialPageRoute(
        builder: (_) => TossPaymentWebView(
          clientKey: _clientKey,
          orderId: orderId,
          orderName: orderName,
          amount: amountKrw,
          customerName: customerName,
          customerEmail: customerEmail,
        ),
      ),
    );

    if (result == null) {
      return PaymentResult(
        success: false,
        paymentKey: '',
        orderId: orderId,
        amountKrw: amountKrw,
        method: '',
        message: '사용자가 결제를 취소했습니다',
        errorCode: 'USER_CANCEL',
      );
    }

    if (result.success && result.paymentKey != null) {
      return confirmPayment(
        paymentKey: result.paymentKey!,
        orderId: orderId,
        amountKrw: amountKrw,
      );
    }

    return PaymentResult(
      success: false,
      paymentKey: '',
      orderId: orderId,
      amountKrw: amountKrw,
      method: '',
      message: result.errorMessage ?? '결제 실패',
      errorCode: result.errorCode,
    );
  }

  @override
  Future<PaymentResult> confirmPayment({
    required String paymentKey,
    required String orderId,
    required int amountKrw,
  }) async {
    try {
      final res = await restClient.confirmPayment(
        paymentKey,
        pgTransactionId: paymentKey,
        pgProvider: 'toss',
      );
      return PaymentResult(
        success: true,
        paymentKey: paymentKey,
        orderId: orderId,
        amountKrw: amountKrw,
        method: res['payment_method'] as String? ?? 'card',
        message: '결제가 완료되었습니다',
      );
    } catch (e) {
      return PaymentResult(
        success: false,
        paymentKey: paymentKey,
        orderId: orderId,
        amountKrw: amountKrw,
        method: '',
        message: '결제 확인 실패: $e',
        errorCode: 'CONFIRM_FAILED',
      );
    }
  }

  @override
  Future<RefundResult> requestRefund({
    required String paymentKey,
    String? reason,
    int? refundAmountKrw,
  }) async {
    try {
      final res = await restClient.refundPayment(
        paymentKey,
        reason: reason,
      );
      return RefundResult(
        success: true,
        refundKey: res['refund_id'] as String? ?? 'ref_$paymentKey',
        message: '환불이 완료되었습니다',
      );
    } catch (e) {
      return RefundResult(
        success: false,
        refundKey: '',
        message: '환불 실패: $e',
        errorCode: 'REFUND_FAILED',
      );
    }
  }

  Future<PaymentResult> _simulatePayment(
    String orderId, int amountKrw, String orderName,
  ) async {
    debugPrint('[TossPG:Sim] 결제 요청: $orderName ($amountKrw원)');
    await Future.delayed(const Duration(seconds: 1));
    return PaymentResult(
      success: true,
      paymentKey: 'sim_${DateTime.now().millisecondsSinceEpoch}',
      orderId: orderId,
      amountKrw: amountKrw,
      method: 'simulated_card',
      message: '시뮬레이션 결제 성공',
    );
  }
}

/// 시뮬레이션 결제 서비스 (SDK 미설치 시 기본 동작)
class SimulatedPaymentService implements PaymentService {
  @override
  Future<PaymentResult> requestPayment({
    required String orderId,
    required int amountKrw,
    required String orderName,
    String? customerName,
    String? customerEmail,
  }) async {
    debugPrint('[SimulatedPG] 결제 요청: $orderName ($amountKrw원)');
    await Future.delayed(const Duration(seconds: 1));
    return PaymentResult(
      success: true,
      paymentKey: 'sim_${DateTime.now().millisecondsSinceEpoch}',
      orderId: orderId,
      amountKrw: amountKrw,
      method: 'simulated_card',
      message: '시뮬레이션 결제 성공',
    );
  }

  @override
  Future<PaymentResult> confirmPayment({
    required String paymentKey,
    required String orderId,
    required int amountKrw,
  }) async {
    debugPrint('[SimulatedPG] 결제 승인: $paymentKey');
    await Future.delayed(const Duration(milliseconds: 500));
    return PaymentResult(
      success: true,
      paymentKey: paymentKey,
      orderId: orderId,
      amountKrw: amountKrw,
      method: 'simulated_card',
      message: '시뮬레이션 결제 승인 완료',
    );
  }

  @override
  Future<RefundResult> requestRefund({
    required String paymentKey,
    String? reason,
    int? refundAmountKrw,
  }) async {
    debugPrint('[SimulatedPG] 환불 요청: $paymentKey');
    await Future.delayed(const Duration(milliseconds: 500));
    return RefundResult(
      success: true,
      refundKey: 'ref_${DateTime.now().millisecondsSinceEpoch}',
      message: '시뮬레이션 환불 완료',
    );
  }
}

class PaymentResult {
  final bool success;
  final String paymentKey;
  final String orderId;
  final int amountKrw;
  final String method;
  final String message;
  final String? errorCode;

  const PaymentResult({
    required this.success,
    required this.paymentKey,
    required this.orderId,
    required this.amountKrw,
    required this.method,
    required this.message,
    this.errorCode,
  });
}

class RefundResult {
  final bool success;
  final String refundKey;
  final String message;
  final String? errorCode;

  const RefundResult({
    required this.success,
    required this.refundKey,
    required this.message,
    this.errorCode,
  });
}
