import 'package:flutter/foundation.dart';

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
