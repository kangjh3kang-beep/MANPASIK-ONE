import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';

/// Toss Payments 결제 WebView 위젯
///
/// Toss Payments 결제창을 WebView로 호출하고
/// 결제 성공/실패 시 콜백을 반환합니다.
class TossPaymentWebView extends StatefulWidget {
  const TossPaymentWebView({
    super.key,
    required this.clientKey,
    required this.orderId,
    required this.orderName,
    required this.amount,
    this.customerName,
    this.customerEmail,
    this.successUrl = 'manpasik://payment/success',
    this.failUrl = 'manpasik://payment/fail',
  });

  final String clientKey;
  final String orderId;
  final String orderName;
  final int amount;
  final String? customerName;
  final String? customerEmail;
  final String successUrl;
  final String failUrl;

  @override
  State<TossPaymentWebView> createState() => _TossPaymentWebViewState();
}

class _TossPaymentWebViewState extends State<TossPaymentWebView> {
  late final WebViewController _controller;
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _controller = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setNavigationDelegate(
        NavigationDelegate(
          onPageStarted: (_) => setState(() => _isLoading = true),
          onPageFinished: (_) => setState(() => _isLoading = false),
          onNavigationRequest: _handleNavigation,
        ),
      )
      ..loadHtmlString(_buildPaymentHtml());
  }

  NavigationDecision _handleNavigation(NavigationRequest request) {
    final uri = Uri.tryParse(request.url);
    if (uri == null) return NavigationDecision.navigate;

    // 결제 성공 콜백
    if (request.url.startsWith(widget.successUrl) ||
        request.url.contains('payment/success')) {
      final paymentKey = uri.queryParameters['paymentKey'] ?? '';
      final orderId = uri.queryParameters['orderId'] ?? widget.orderId;
      final amount = int.tryParse(uri.queryParameters['amount'] ?? '') ??
          widget.amount;
      Navigator.pop(context, TossPaymentResult(
        success: true,
        paymentKey: paymentKey,
        orderId: orderId,
        amount: amount,
      ));
      return NavigationDecision.prevent;
    }

    // 결제 실패 콜백
    if (request.url.startsWith(widget.failUrl) ||
        request.url.contains('payment/fail')) {
      final errorCode = uri.queryParameters['code'] ?? 'UNKNOWN';
      final errorMessage = uri.queryParameters['message'] ??
          Uri.decodeComponent(uri.queryParameters['msg'] ?? '결제 실패');
      Navigator.pop(context, TossPaymentResult(
        success: false,
        errorCode: errorCode,
        errorMessage: errorMessage,
      ));
      return NavigationDecision.prevent;
    }

    return NavigationDecision.navigate;
  }

  String _buildPaymentHtml() {
    final customerName = widget.customerName ?? '';
    final customerEmail = widget.customerEmail ?? '';

    return '''
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <script src="https://js.tosspayments.com/v1/payment"></script>
  <style>
    body { display: flex; justify-content: center; align-items: center;
           min-height: 100vh; margin: 0; background: #f8f9fa;
           font-family: -apple-system, sans-serif; }
    .loading { text-align: center; color: #666; }
    .loading p { margin-top: 16px; font-size: 16px; }
  </style>
</head>
<body>
  <div class="loading">
    <p>결제창을 불러오는 중...</p>
  </div>
  <script>
    var tossPayments = TossPayments('${widget.clientKey}');
    tossPayments.requestPayment('카드', {
      amount: ${widget.amount},
      orderId: '${widget.orderId}',
      orderName: '${widget.orderName}',
      ${customerName.isNotEmpty ? "customerName: '$customerName'," : ''}
      ${customerEmail.isNotEmpty ? "customerEmail: '$customerEmail'," : ''}
      successUrl: '${widget.successUrl}',
      failUrl: '${widget.failUrl}',
    }).catch(function (error) {
      if (error.code === 'USER_CANCEL') {
        window.location.href = '${widget.failUrl}?code=USER_CANCEL&message=사용자가 결제를 취소했습니다';
      } else {
        window.location.href = '${widget.failUrl}?code=' + error.code + '&message=' + encodeURIComponent(error.message);
      }
    });
  </script>
</body>
</html>
''';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('결제'),
        leading: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () => Navigator.pop(context, null),
        ),
      ),
      body: Stack(
        children: [
          WebViewWidget(controller: _controller),
          if (_isLoading)
            const Center(child: CircularProgressIndicator()),
        ],
      ),
    );
  }
}

/// Toss Payments 결제 결과
class TossPaymentResult {
  final bool success;
  final String? paymentKey;
  final String? orderId;
  final int? amount;
  final String? errorCode;
  final String? errorMessage;

  const TossPaymentResult({
    required this.success,
    this.paymentKey,
    this.orderId,
    this.amount,
    this.errorCode,
    this.errorMessage,
  });
}
